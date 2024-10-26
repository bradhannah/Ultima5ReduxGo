package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func getArrowKeyPressed() *ebiten.Key {
	var keyPressed ebiten.Key = ebiten.KeyF24
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		keyPressed = ebiten.KeyLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		keyPressed = ebiten.KeyRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		keyPressed = ebiten.KeyUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		keyPressed = ebiten.KeyDown
	}
	if keyPressed == ebiten.KeyF24 {
		return nil
	}
	return &keyPressed
}

func getCurrentPressedArrowKeyAsDirection() references.Direction {
	arrowKey := getArrowKeyPressed()
	if arrowKey == nil {
		return references.NoneDirection
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		return references.Up
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		return references.Down
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		return references.Left
	}
	return references.Right
}

func (g *GameScene) moveToNewPositionByDirection(direction references.Direction) {
	// TODO: dear future Brad - you will need to change this big time when dungeons and combat come in
	bLargeMap := g.gameState.Location.GetMapType() == references.LargeMapType
	switch direction {
	case references.Up:
		g.gameState.Position.GoUp(bLargeMap)
	case references.Down:
		g.gameState.Position.GoDown(bLargeMap)
	case references.Left:
		g.gameState.Position.GoLeft(bLargeMap)
	case references.Right:
		g.gameState.Position.GoRight(bLargeMap)
	case references.NoneDirection:
	}
}

func (g *GameScene) handleMovement(directionStr string, key ebiten.Key) {
	g.debugMessage = directionStr
	g.addRowStr(fmt.Sprintf("> %s", directionStr))

	isPassable := func(pos *references.Position) bool {
		theMap := g.gameState.LayeredMaps.GetLayeredMap(g.gameState.Location.GetMapType(), g.gameState.Floor)
		topTile := theMap.GetTopTile(pos)
		return topTile.IsPassable(g.gameState.PartyVehicle)
	}

	direction := game_state.GetKeyAsDirection(key)
	newPosition := direction.GetNewPositionInDirection(&g.gameState.Position)
	if isPassable(newPosition) || g.gameState.DebugOptions.FreeMove {
		g.moveToNewPositionByDirection(direction)
	} else {
		g.addRowStr("Blocked!")
	}
}
