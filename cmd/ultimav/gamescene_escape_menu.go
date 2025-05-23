package main

import (
	"os"

	"github.com/bradhannah/Ultima5ReduxGo/internal/ui/widgets"
)

func (g *GameScene) DoEscapeMenu() {
	const forceWaitTimeMs = 250
	bl := widgets.NewButtonListModal(
		"Ultima V Redux",
		func() {
			g.dialogStack.PopModalDialog()
			g.keyboard.SetForceWaitAnyKey(forceWaitTimeMs)
		},
		g.keyboard,
		&gameScreenPercents)

	g.dialogStack.PushModalDialog(bl)
	bl.AddButton("Back to Main Menu", func() { g.dialogStack.PopModalDialog() })
	bl.AddButton("Quick Save", func() { g.dialogStack.PopModalDialog() })
	bl.AddButton("Save", func() { g.dialogStack.PopModalDialog() })
	bl.AddButton("Load", func() { g.dialogStack.PopModalDialog() })
	bl.AddButton("Configure", func() { g.dialogStack.PopModalDialog() })
	bl.AddButton("Exit to DOS", func() {
		g.dialogStack.PopModalDialog()

		// todo: popping up a save menu may be added in the future
		os.Exit(0)
	})
}

// func (g *GameScene) getPercentBasedPlacementForEscapeMenu(nIndex int) sprites.PercentBasedPlacement {
// 	return sprites.PercentBasedPlacement{
// 		StartPercentX: .45,
// 		EndPercentX:   .55,
// 		StartPercentY: startEscapeMenuPercentY * (.05 * float64(nIndex)),
// 		EndPercentY:   startEscapeMenuPercentY + .05,
// 	}
// }
