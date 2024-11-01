package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

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

func (g *GameScene) commonMapLookSecondary(direction references.Direction) {
	newPosition := direction.GetNewPositionInDirection(&g.gameState.Position)
	topTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(newPosition)
	g.addRowStr(fmt.Sprintf("Thou dost see %s", g.gameReferences.LookReferences.GetTileLookDescription(topTile.Index)))

}
