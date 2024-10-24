package main

func (g *GameScene) appendDirectionToOutput() {
	g.appendToCurrentRowStr(getCurrentPressedArrowKeyAsDirection().GetDirectionCompassName())
}
