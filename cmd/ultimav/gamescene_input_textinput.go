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
				g.gameState.DebugQuickExitSmallMap()
				g.dialogStack.PopModalDialog()
			} else if outputStr == "no" {
				g.dialogStack.PopModalDialog()
			}
			g.keyboard.SetForceWaitAnyKey(500)
		})
}
