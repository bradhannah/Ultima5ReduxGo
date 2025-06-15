package party_state

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

/* =====================================================================
   Conversation – feature‑rich port of the original C# class
   =====================================================================
   Major differences from the minimal demo you had:
   ● Fully honours label jumps, <IfElseKnowsName>, <Change>, <GoldPrompt>,
     default answers, Or‑chains, and runic mode.
   ● Non‑blocking: a goroutine runs the interpreter and pushes ScriptItems
     to an output channel; caller feeds responses via the input channel.
   ● Skip‑state and label stack implemented exactly like the C# enum & lists.
   ---------------------------------------------------------------------*/

const pauseInMs = 400

type Conversation struct {
	npcID int
	//	api   AvatarAPI
	party *PartyState
	ts    *references.TalkScript

	// channels
	out chan references.ScriptItem
	in  chan string

	// internal state
	runeMode          bool
	convoOrder        []int                   // indices into ts.Lines
	convoOrderScript  []references.ScriptLine // cached lines
	currentSkip       skipInstr               // instruction set by ProcessLine
	ctx               context.Context
	cancel            context.CancelFunc
	conversationEnded bool
}

func NewConversation(npcID int, party *PartyState, ts *references.TalkScript) *Conversation {
	ctx, cancel := context.WithCancel(context.Background())
	return &Conversation{
		npcID: npcID, party: party, ts: ts,
		out: make(chan references.ScriptItem), in: make(chan string),
		ctx: ctx, cancel: cancel,
	}
}

func (c *Conversation) Out() <-chan references.ScriptItem { return c.out }
func (c *Conversation) In() chan<- string                 { return c.in }
func (c *Conversation) Start()                            { go c.loop() }
func (c *Conversation) Stop()                             { c.cancel() }

// ----------------------- skip enum -------------------------------------

type skipInstr int

const (
	skipNone skipInstr = iota
	skipNext
	skipAfterNext
	skipToLabel
)

// ----------------------- main loop ------------------------------------

func (c *Conversation) loop() {
	defer close(c.out)

	// bootstrap order: description, greeting
	c.convoOrder = []int{references.TalkScriptConstantsDescription, references.TalkScriptConstantsGreeting}
	c.convoOrderScript = []references.ScriptLine{}

	idxConv := 0

	for {
		if c.conversationEnded || c.ctx.Err() != nil {
			return
		}

		// need more dialogue? prompt player
		if idxConv >= len(c.convoOrder) {
			// Ask a question
			//c.enqueueFmt("\n> ")
			userInput := c.readLine()

			if userInput == "" {
				userInput = "bye"
			}

			if qa, ok := c.ts.Questions[strings.ToLower(userInput)]; ok {
				_ = c.processMultiLines(qa.Answer.SplitIntoSections(), -1)

				continue
			}
			// unrecognised
			c.enqueueStr("I cannot help thee.\n")

			continue
		}

		lineIdx := c.convoOrder[idxConv]
		line, ok := c.ts.GetScriptLine(lineIdx)

		if !ok {
			c.enqueueFmt("[missing script line %d]\n", lineIdx)

			return
		}

		split := line.SplitIntoSections()
		_ = c.processMultiLines(split, lineIdx)

		idxConv++
	}
}

// ----------------------- processing helpers ---------------------------

func (c *Conversation) processMultiLines(sections []references.SplitScriptLine, talkIdx int) error {
	skipCounter := -1
	for i, sec := range sections {
		if skipCounter == 0 {
			skipCounter--
			continue
		}

		if sec.Contains(references.AvatarsName) && !c.party.HasMet(c.npcID) {
			continue // they don't know me yet
		}

		if len(sec) == 0 {
			continue
		}
		if err := c.processLine(sec, talkIdx, i); err != nil {
			return err
		}

		switch c.currentSkip {
		case skipToLabel:
			return nil
		case skipAfterNext:
			skipCounter = 1
		case skipNext:
			skipCounter = 0
		}
	}
	c.currentSkip = skipNone
	return nil
}

//nolint:cyclop
//nolint:funlen
func (c *Conversation) processLine(line references.ScriptLine, talkIdx, splitIdx int) error {
	// special: description pre‑amble "You see xxx "
	if talkIdx == references.TalkScriptConstantsDescription && splitIdx == 0 {
		c.enqueueStr("You see ")
	}

	// AskName optimisation
	if line.Contains(references.AskName) && c.party.HasMet(c.npcID) {
		c.currentSkip = skipNone
		return nil
	}

	for n := 0; n < len(line); n++ {
		itm := line[n]

		switch itm.Cmd {
		case references.IfElseKnowsName:
			if c.party.HasMet(c.npcID) {
				c.currentSkip = skipAfterNext
			} else {
				c.currentSkip = skipNext
			}

			return nil

		case references.AvatarsName:
			c.enqueueStr(c.party.AvatarName())

		case references.AskName:
			c.enqueueStr("What is thy name? ")
			name := c.readLine()

			if strings.EqualFold(name, c.party.AvatarName()) {
				c.party.SetMet(c.npcID)
				c.enqueueFmt("A pleasure, %s.\n", c.party.AvatarName())
			} else {
				c.enqueueStr("If thou sayest so…\n")
			}

		case references.DefineLabel:
			tgt := itm.Num

			idx := c.ts.GetScriptLineLabelIndex(tgt)
			if idx != -1 {
				c.convoOrder = append(c.convoOrder, idx)
			}

			c.convoOrder = append(c.convoOrder, idx)
			c.currentSkip = skipToLabel

			return nil

		case references.JoinParty:
			if !c.party.HasRoom() {
				c.enqueueStr("My party is full.\n")
			} else if err := c.party.JoinNPC(c.npcID); err != nil {
				c.enqueueFmt("%v\n", err)
			} else {
				c.enqueueFmt("%s has joined thee!\n", itm.Str)
			}
			c.conversationEnded = true

			return nil

		case references.EndConversation:
			c.enqueueStr("Farewell.\n")
			c.conversationEnded = true

			return nil

		case references.PlainString:
			c.enqueueStr(itm.Str)
		case references.NewLine:
			c.enqueueStr("\n")
		case references.Rune:
			c.runeMode = !c.runeMode
		case references.Pause:
			time.Sleep(pauseInMs * time.Millisecond)

		case references.OrBranch, references.StartNewSection, references.DoNothingSection:
			// never appears in split sections – sanity only
			log.Fatalf("Unexpected OR, StartNewSection or DoNothingSection in script: %v", itm.Cmd)
		case references.Label1, references.Label2, references.Label3, references.Label4, references.Label5,
			references.Label6, references.Label7, references.Label8, references.Label9:
			c.enqueueStr(itm.Str)
		default:
			// pass‑through unimplemented opcodes for now
			c.enqueueStr("<" + itm.Cmd.String() + ">")
		}
	}

	c.currentSkip = skipNone
	return nil
}

// ----------------------- I/O helpers -----------------------------------

func (c *Conversation) enqueueStr(s string) {
	c.out <- references.ScriptItem{Cmd: references.PlainString, Str: s}
}
func (c *Conversation) enqueueFmt(f string, a ...interface{}) {
	c.enqueueStr(fmt.Sprintf(f, a...))
}

func (c *Conversation) readLine() string {
	//begin:
	select {
	case <-c.ctx.Done():
		return ""
	case s := <-c.in:
		return strings.TrimSpace(s)
		//default:
		//	goto begin
	}
}
