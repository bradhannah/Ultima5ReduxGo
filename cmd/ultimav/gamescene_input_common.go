package main

import "fmt"

func (g *GameScene) isDirectionKeyValidAndOutput() bool {
	arrowKey := getArrowKeyPressed()
	bIsArrowKeyPressed := arrowKey != nil

	if !bIsArrowKeyPressed {
		return false
	}
	if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
		return false
	}
	g.appendDirectionToOutput()
	return true
}

func (g *GameScene) commonMapLookSecondary() {
	newPosition := getCurrentPressedArrowKeyAsDirection().GetNewPositionInDirection(&g.gameState.Position)
	topTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(newPosition)
	g.addRowStr(fmt.Sprintf("Thou dost see %s", g.gameReferences.LookReferences.GetTileLookDescription(topTile.Index)))

}
