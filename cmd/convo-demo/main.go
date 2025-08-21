// cmd/convo-demo/main.go
package main

import (
	"bufio"
	_ "embed"
	"log"
	"os"
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/conversation"
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
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
	npcID := 3
	talkScriptIndex := 2
	talkScript := baseState.GameReferences.TalkReferences.GetTalkScriptByNpcIndex(references.Castle, talkScriptIndex)

	npcState := &NPCState{HasMetAvatar: false, Script: talkScript}

	npcReferences := baseState.GameReferences.NPCReferences.GetNPCReferencesByLocation(references.Britain)
	npcReference := (*npcReferences)[npcID]

	// Pass the correct arguments to NewConversation
	convo := conversation.NewConversation(npcReference,
		baseState,
		npcState.Script,
	)

	convo.Start()
	sleepDuration := 10 * time.Second
	time.Sleep(sleepDuration)

	scanCh := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			scanCh <- scanner.Text() // blocks, but that’s fine in this goroutine
		}
		close(scanCh) // on EOF / error
	}()

	for {
		select {
		case v := <-convo.Out():
			{
				log.Print(v.Str)
			}
		case line, ok := <-scanCh:
			if !ok {
				continue
			}
			convo.In() <- line
		}
	}

	//for item := range convo.Out() {
	//	fmt.Print(item.Str)
	//	if strings.HasSuffix(item.Str, "> ") { // our engine prompts with "> "
	//		if scan.Scan() {
	//			convo.In() <- scan.Text()
	//		}
	//	}
	//}

}
