package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"strings"
)

func (g *GameScene) createTextCommandExitBuilding() *grammar.TextCommand {
	yesNo := []string{"yes", "no"}

	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchStringList{
			Strings:       yesNo,
			Description:   "Yes or No",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			outputStr := strings.ToLower(g.inputBox.GetText())
			fmt.Sprintf(outputStr)
			//locationStr := command.GetIndexAsString(1, outputStr)
			//oof := d.gameScene.gameReferences.LocationReferences.GetLocationByName(locationStr)
			//
			//d.dumpQuickState(oof.FriendlyLocationName)
			//d.gameScene.gameState.EnterBuilding(oof, d.gameScene.gameReferences.TileReferences)

		})
}
