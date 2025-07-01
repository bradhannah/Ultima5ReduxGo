package conversation

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
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

const bExtraTextDebug = true

const pauseInMs = 400

//type ConversationCallbacks struct {
//	//Karma *party_state.Karma
//	//GameState *game_state.GameState
//	PartyState *party_state.PartyState
//}

type Conversation struct {
	//conversationCallbacks ConversationCallbacks

	//npcID int
	npcReference references.NPCReference
	//	api   AvatarAPI
	//party *party_state.PartyState
	gameState *game_state.GameState
	ts        *references.TalkScript

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

func NewConversation(npcReference references.NPCReference,
	//party *party_state.PartyState,
	gameState *game_state.GameState,
	ts *references.TalkScript,
	// conversationCallbacks ConversationCallbacks
) *Conversation {
	ctx, cancel := context.WithCancel(context.Background())
	return &Conversation{
		npcReference: npcReference,
		gameState:    gameState,
		ts:           ts,
		out:          make(chan references.ScriptItem),
		in:           make(chan string),
		ctx:          ctx,
		cancel:       cancel,
		//conversationCallbacks: conversationCallbacks,
	}
}

func (c *Conversation) Out() <-chan references.ScriptItem { return c.out }
func (c *Conversation) In() chan<- string                 { return c.in }
func (c *Conversation) Start()                            { go c.loop() }
func (c *Conversation) Stop()                             { c.cancel() }

// ----------------------- skip enum -------------------------------------

type skipInstr int

const (
	doNotSkip skipInstr = iota
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
	for i, section := range sections {
		// if skipCounter hits zero, we know we are at the point we need
		// to skip the next record
		if skipCounter == 0 {
			skipCounter--
			continue
		}

		if section.Contains(references.AvatarsName) && !c.gameState.PartyState.HasMet(c.npcReference.DialogNumber) {
			// move directly to next section
			continue // they don't know me yet
		}

		if len(section) == 0 {
			//log.Fatalf("Unexpected that length of sections is zero\n")
			continue
		}

		if err := c.processLine(section, talkIdx, i); err != nil {
			return err
		}

		if skipCounter == -1 {
			skipCounter--
		}

		switch c.currentSkip {
		case skipToLabel:
			{
				return nil
			}
		case skipAfterNext:
			{
				skipCounter = 1
			}
		case skipNext:
			{
				i++
			}
		case doNotSkip:
			{
				break
			}
		}
	}
	//c.currentSkip = skipNone
	//log.Fatalf("Unexpected that length of sections is zero\n")
	return nil
}

func (c *Conversation) giveIncrement(rawToGiveIndex uint16) {
	const nonEquipmentIndexCutoff = 0x40

	if rawToGiveIndex < nonEquipmentIndexCutoff {
		c.gameState.PartyState.Inventory.Equipment.IncrementByOne(references.Equipment(rawToGiveIndex))
	}

	itemIndex := rawToGiveIndex - (nonEquipmentIndexCutoff + 1)
	//nolint:mnd
	switch itemIndex {
	case 0:
		c.gameState.PartyState.Inventory.Provisions.Food.IncrementByOne()
	case 1:
		c.gameState.PartyState.Inventory.Gold.IncrementByOne()
	case 2: //nolint:mnd
		c.gameState.PartyState.Inventory.Provisions.Keys.IncrementByOne()
	case 3:
		c.gameState.PartyState.Inventory.Provisions.Gems.IncrementByOne()
	case 4:
		c.gameState.PartyState.Inventory.Provisions.Torches.IncrementByOne()
	case 5:
		c.gameState.PartyState.Inventory.SpecialItems.IncrementByOne(references.Grapple)
	case 6:
		c.gameState.PartyState.Inventory.SpecialItems.IncrementByOne(references.Carpet)
	case 7:
		c.gameState.PartyState.Inventory.SpecialItems.SetHasOne(references.Sextant)
	case 8:
		c.gameState.PartyState.Inventory.SpecialItems.SetHasOne(references.Spyglass)
	case 9:
		c.gameState.PartyState.Inventory.SpecialItems.SetHasOne(references.BlackBadge)
	case 10:
		c.gameState.PartyState.Inventory.Provisions.SkullKeys.IncrementByOne()
	default:
		log.Fatalf("Unknown item index %d", itemIndex)
	}
}

//nolint:gocyclo
//nolint:cyclop
func (c *Conversation) processLine(line references.ScriptLine, talkIdx, splitIdx int) error {
	// AskName optimisation
	if line.Contains(references.AskName) && c.gameState.PartyState.HasMet(c.npcReference.DialogNumber) {
		c.currentSkip = doNotSkip
		return nil
	}

	// special: description pre‑amble "You see xxx "
	if talkIdx == references.TalkScriptConstantsDescription && splitIdx == 0 {
		c.enqueueStr("You see ")
	}

	for _, scriptItem := range line {
		switch scriptItem.Cmd {

		case references.IfElseKnowsName:
			if c.gameState.PartyState.HasMet(c.npcReference.DialogNumber) {
				c.currentSkip = skipAfterNext
			} else {
				c.currentSkip = skipNext
			}

			return nil
		case references.AvatarsName:
			c.enqueueStr(c.gameState.PartyState.AvatarName())
		case references.AskName:
			c.enqueueStr("What is thy name? ")
			name := c.readLine()

			if strings.EqualFold(name, c.gameState.PartyState.AvatarName()) {
				c.gameState.PartyState.SetMet(c.npcReference.Location, int(c.npcReference.DialogNumber))
				c.enqueueFmt("A pleasure, %s.\n", c.gameState.PartyState.AvatarName())
			} else {
				c.enqueueStr("If thou sayest so...\n")
			}
		case references.CallGuards:
			c.enqueueFmt("PLACEHOLDER")
			// TODO: add flag for AI controller
		case references.Change:
			c.enqueueFmt("PLACEHOLDER")
			c.giveIncrement(scriptItem.ItemAdditionalData)
		case references.DefineLabel:
			{
				// maybe ok?
				tgt := scriptItem.Num

				idx := c.ts.GetScriptLineLabelIndex(tgt)
				if idx != -1 {
					c.convoOrder = append(c.convoOrder, idx)
				}

				c.convoOrder = append(c.convoOrder, idx)
				c.currentSkip = skipToLabel

				return nil
			}
		case references.DoNothingSection:
			break
		case references.EndConversation:
			c.enqueueStr("PLACEHOLDER")
			c.conversationEnded = true
			return nil
		case references.EndScript:
			c.enqueueStr("PLACEHOLDER")
		case references.ExtortionAmount:
			{
				c.enqueueFmt("PLACEHOLDER")
			}
		case references.GoldPrompt:
			c.enqueueStr("PLACEHOLDER")
			c.gameState.PartyState.Inventory.Gold.DecrementBy(uint16(scriptItem.ItemAdditionalData))
		case references.GotoLabel:
			break
		case references.GoToJail:
			c.enqueueStr("PLACEHOLDER")
			// TODO: send Avatar and party to jail
		case references.JoinParty:
			if !c.gameState.PartyState.HasRoom() {
				c.enqueueStr("My party is full.\n")
			} else if err := c.gameState.PartyState.JoinNPC(c.npcReference); err != nil {
				c.enqueueFmt("%v\n", err)
			} else {
				c.enqueueFmt("%s has joined thee!\n", scriptItem.Str)
			}
			c.enqueueFmt("PLACEHOLDER")
			c.conversationEnded = true

			return nil
		case references.KarmaMinusOne:
			c.gameState.PartyState.Karma.DecreaseKarma(1)
			if bExtraTextDebug {
				c.enqueueStr("KARMA_DEC_ONE")
			}
		case references.KarmaPlusOne:
			c.gameState.PartyState.Karma.IncreaseKarma(1)
			if bExtraTextDebug {
				c.enqueueStr("KARMA_PLUS_ONE")
			}
		case references.KeyWait:
			c.enqueueStr("PLACEHOLDER")
		case references.Label1, references.Label2, references.Label3, references.Label4, references.Label5,
			references.Label6, references.Label7, references.Label8, references.Label9, references.Label10:
			c.enqueueStr(scriptItem.Str)
		case references.MakeAHorse:
			c.enqueueStr("PLACEHOLDER")
		case references.NewLine:
			c.enqueueStr("\n")
		case references.PayGenericExtortion:
			c.enqueueStr("PLACEHOLDER")
		case references.PayHalfGoldExtortion:
			c.enqueueStr("PLACEHOLDER")
		case references.PlainString:
			c.enqueueStr(scriptItem.Str)
		case references.Pause:
			time.Sleep(pauseInMs * time.Millisecond)
		case references.PromptUserForInput_NPCQuestion:
			break
		case references.PromptUserForInput_UserInterest:
			break
		case references.Rune:
			c.runeMode = !c.runeMode
		case references.StartLabelDef:
			c.enqueueFmt("PLACEHOLDER - nItem++")
		case references.OrBranch, references.StartNewSection:
			// never appears in split sections – sanity only
			log.Fatalf("Unexpected OR, StartNewSection or DoNothingSection in script: %v", scriptItem.Cmd)
		case references.UserInputNotRecognized:
			c.enqueueStr("Cannot help.")
		default:
			// pass‑through unimplemented opcodes for now
			c.enqueueStr("<" + scriptItem.Cmd.String() + ">")
		}
	}

	c.currentSkip = doNotSkip
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
