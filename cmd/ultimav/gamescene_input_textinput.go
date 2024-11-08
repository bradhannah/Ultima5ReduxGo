package main

import (
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

			if outputStr == "yes" {
				g.gameState.ExitSmallMap()
				g.bShowInputBox = false
				g.inputBox = nil
			}
		})
}
