package main

import (
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
)

func (g *GameScene) createTextCommandExitBuilding() *grammar.TextCommand {
	yesNo := []string{"yes", "no"}

	return grammar.NewTextCommand(
		[]grammar.Match{
			grammar.MatchStringList{
				Strings:              yesNo,
				Description:          "Yes or No",
				CaseSensitive:        false,
				SingleCharacterInput: true,
			},
		},
		func(s string, command *grammar.TextCommand) {
			ib := g.dialogStack.GetOrAssertTopInputBox()
			outputStr := strings.ToLower(ib.GetText())

			if outputStr == "yes" {
				g.gameState.ExitSmallMap()
				g.dialogStack.PopModalDialog()

			} else if outputStr == "no" {
				// we just go back to what we were doing
				g.dialogStack.PopModalDialog()
			}

		})
}
