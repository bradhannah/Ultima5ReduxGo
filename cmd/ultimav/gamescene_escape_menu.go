package main

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ui/widgets"

func (g *GameScene) DoEscapeMenu() {
	bl := widgets.NewButtonListModal(
		func() {
			g.dialogStack.PopModalDialog()
			g.keyboard.SetForceWaitAnyKey(250)
		},
		g.keyboard)

	g.dialogStack.PushModalDialog(bl)
}
