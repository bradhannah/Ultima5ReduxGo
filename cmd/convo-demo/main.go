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

//go:embed britain2_SAVED.GAM
var saveFile []byte

// NPCState keeps “has met avatar” and the TalkScript.
type NPCState struct {
	HasMetAvatar bool
	Script       *references.TalkScript
}

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
	npcId := 3

	talkScript := baseState.GameReferences.TalkReferences.GetTalkScriptByNpcIndex(references.Castle, 2)

	ns := &NPCState{HasMetAvatar: false, Script: talkScript}

	// 3) start Conversation
	convo := party_state.NewConversation(npcId,
		&baseState.PartyState,
		ns.Script,
	)

	convo.Start()

	scan := bufio.NewScanner(os.Stdin)
	for item := range convo.Out() {
		fmt.Print(item.Str)
		if strings.HasSuffix(item.Str, "> ") { // our engine prompts with "> "
			if scan.Scan() {
				convo.In() <- scan.Text()
			}
		}
	}

}
