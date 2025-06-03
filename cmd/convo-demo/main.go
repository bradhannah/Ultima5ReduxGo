// cmd/convo-demo/main.go
package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

/* ------------------------------------------------------------------------- */
/* 1. Minimal stubs so Conversation logic compiles                           */
/* ------------------------------------------------------------------------- */

// GameState only needs Avatar name + karma in the demo.
//type GameState struct {
//	AvatarsName     string
//	Karma           int
//	TurnsSinceStart int
//}

// NPCState keeps “has met avatar” and the TalkScript.
type NPCState struct {
	HasMetAvatar bool
	Script       *references.TalkScript
}

// TurnResults in the real game aggregates side‑effects; here we log them.
//type TurnResults []string
//
///* ------------------------------------------------------------------------- */
///* 2. Constants extracted from magic numbers                                 */
///* ------------------------------------------------------------------------- */
//
//const (
//	// gold extortion ranges come from OddsAndLogic in C#
//	minEraGuardGold = 20
//	maxEraGuardGold = 80
//)

//// guardExtortionAmount mimics OddsAndLogic.GetGuardExtortionAmount(…)
//func guardExtortionAmount() int {
//	return rand.Intn(maxEraGuardGold-minEraGuardGold+1) + minEraGuardGold
//}

/* ------------------------------------------------------------------------- */
/* 3. Basic Conversation engine (pared‑down)                                 */
/* ------------------------------------------------------------------------- */

//type Conversation struct {
//	game  *game_state.GameState
//	npc   *NPCState
//	out   chan references.ScriptItem // to UI
//	in    chan string                // from UI
//	rune  bool                       // rune mode toggle
//	ended bool
//}
//
//// NewConversation wires everything.
//func NewConversation(gs *game_state.GameState, ns *NPCState) *Conversation {
//	return &Conversation{
//		game: gs,
//		npc:  ns,
//		out:  make(chan references.ScriptItem),
//		in:   make(chan string),
//	}
//}
//
//func (c *Conversation) Out() <-chan references.ScriptItem { return c.out }
//func (c *Conversation) In() chan<- string                 { return c.in }
//
//// run starts the dialogue interpreter (single goroutine).
//func (c *Conversation) run() {
//	defer close(c.out)
//
//	script := c.npc.Script
//	if script == nil {
//		log.Fatalf("NPC has no script")
//		return
//	}
//
//	// very simple order: description → greeting only
//	desc, ok := script.GetScriptLine(references.TalkScriptConstantsDescription)
//	if !ok {
//		log.Fatalf("NPC has no description: %v", c.npc.Script)
//	}
//	greeting, ok := script.GetScriptLine(references.TalkScriptConstantsGreeting)
//
//	if !ok {
//		log.Fatalf("NPC has no greeting: %v", c.npc.Script)
//	}
//
//	lines := []references.ScriptLine{desc, greeting}
//
//	for _, line := range lines {
//		for _, item := range line {
//			switch item.Cmd {
//
//			case references.PlainString:
//				c.emitStr(item.Str)
//
//			case references.AvatarsName:
//				c.emitStr(string(c.game.PartyState.Characters[0].Name[:]))
//
//			case references.Pause:
//				time.Sleep(500 * time.Millisecond)
//
//			case references.AskName:
//				c.emitStr("What is thy name, traveller?")
//				playerName := <-c.in
//				c.npc.HasMetAvatar = strings.EqualFold(playerName, string(c.game.PartyState.Characters[0].Name[:]))
//				// c.game.AvatarsName)
//
//			case references.GoldPrompt:
//				g := guardExtortionAmount()
//				c.emitStr(fmt.Sprintf("[Guard demands %d gold]\n", g))
//
//			case references.EndConversation:
//				c.emitStr("\nFarewell.")
//				c.ended = true
//			}
//		}
//		if c.ended {
//			return
//		}
//	}
//}
//
//// helper
//func (c *Conversation) emitStr(s string) {
//	c.out <- references.ScriptItem{Cmd: references.PlainString, Str: s}
//}

/* ------------------------------------------------------------------------- */
/* 4. Tiny CLI driver                                                         */
/* ------------------------------------------------------------------------- */

//go:embed britain2_SAVED.GAM
var saveFile []byte

const (
	xTilesVisibleOnGameScreen = 19
	yTilesVisibleOnGameScreen = 13
)

func main() {

	cfg := config.NewUltimaVConfiguration()

	//var err error
	baseGameReferences, err := references.NewGameReferences(cfg)

	if err != nil {
		log.Fatalf("Couldn't load game references: %v\n", err)
	}

	baseState := game_state.NewGameStateFromLegacySaveBytes(
		saveFile, cfg, baseGameReferences,
		xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen,
	)

	//dataOvl := references.NewDataOvl(cfg)
	//var err error
	//talkReferences := references.NewTalkReferences(cfg, dataOvl)

	//if err != nil {
	//	return
	//}

	//rand.Seed(time.Now().UnixNano())

	// 1) build word dict and parse ONE NPC blob (hard‑coded file)
	//tdict := references.NewWordDict(myCompressedWordList()) // stubbed list
	//raws, _ := references.LoadFile("data/CASTLE.TLK")       // adjust path
	//blob := raws[1]                                         // NPC #1
	//script, _ := references.ParseNPCBlob(blob, tdict)
	//
	//// 2) wire GameState + NPCState
	//gs := &GameState{AvatarsName: "Avatar", Karma: 0}
	talkScripts := baseState.GameReferences.TalkReferences.GetTalkScript(references.Castle)
	talkScript := talkScripts[1]

	ns := &NPCState{HasMetAvatar: false, Script: talkScript}

	// 3) start Conversation
	convo := party_state.NewConversation(1,
		&baseState.PartyState,
		ns.Script,
	)

	//baseState, ns)
	convo.Start()
	//go convo.run()

	scan := bufio.NewScanner(os.Stdin)
	for item := range convo.Out() {
		fmt.Print(item.Str)
		if strings.HasSuffix(item.Str, "> ") { // our engine prompts with "> "
			if scan.Scan() {
				convo.In() <- scan.Text()
			}
		}
	}

	// 4) main loop: print out items, feed user input when asked
	//scan := bufio.NewScanner(os.Stdin)
	//for item := range convo.Out() {
	//	switch item.Cmd {
	//	case references.PlainString:
	//		fmt.Print(item.Str)
	//	case references.PromptUserForInput_UserInterest,
	//		references.PromptUserForInput_NPCQuestion:
	//		fmt.Print("> ")
	//		if scan.Scan() {
	//			convo.In() <- scan.Text()
	//		}
	//	}
	//}
}

/* ------------------------------------------------------------------------- */
/* 5. Tiny stub compressed‑word list for demo                                 */
/* ------------------------------------------------------------------------- */

func myCompressedWordList() []string {
	return []string{
		"—unused—", "the", "and", "of", "to", "is", "in", "that", // etc…
	}
}

/* Notes:
   - references.TalkScriptConstantsDescription/Greeting are placeholders –
     export constants or helper getters in your AST package.
   - The demo only understands PlainString, AvatarsName, Pause, AskName,
     Gold, EndConversation.  Add cases incrementally as you migrate logic.
   - real rune conversion is skipped; toggle simply ignored.
   - TurnResults side‑effects are logged directly inside the switch when you
     add the opcodes later.
*/
