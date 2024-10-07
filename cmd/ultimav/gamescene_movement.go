package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func isArrowKeyPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight)
}

func getCurrentPressedArrowKeyAsDirection() game_state.Direction {
	if !isArrowKeyPressed() {
		return game_state.None
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		return game_state.Up
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		return game_state.Down
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		return game_state.Left
	}
	return game_state.Right
}

func (g *GameScene) moveToNewPositionByDirection(direction game_state.Direction) {
	bLargeMap := g.gameState.Location == references.Britannia_Underworld
	switch direction {
	case game_state.Up:
		g.gameState.Position.GoUp(bLargeMap)
	case game_state.Down:
		g.gameState.Position.GoDown(bLargeMap)
	case game_state.Left:
		g.gameState.Position.GoLeft(bLargeMap)
	case game_state.Right:
		g.gameState.Position.GoRight(bLargeMap)
	case game_state.None:
	}
}

func (g *GameScene) handleMovement(directionStr string, key ebiten.Key) {
	g.debugMessage = directionStr
	g.output.AddToContinuousOutput(fmt.Sprintf("> %s", directionStr))

	isPassable := func(pos *references.Position) bool {
		topTile := g.gameState.LayeredMaps.LayeredMaps[g.gameState.GetMapType()].GetTopTile(pos)
		return topTile.IsPassable(g.gameState.PartyVehicle)
	}

	direction := game_state.GetKeyAsDirection(key)
	newPosition := direction.GetNewPositionInDirection(&g.gameState.Position)
	if isPassable(newPosition) {
		g.moveToNewPositionByDirection(direction)
	} else {
		g.output.AddToContinuousOutput("Blocked!")
	}
}

func (g *GameScene) handleSecondaryInput() {
	if !g.keyboard.TryToRegisterKeyPress() {
		return
	}

	switch g.gameState.SecondaryKeyState {
	case game_state.OpenDirectionInput:
		if !isArrowKeyPressed() {
			return
		}
		g.output.AddToContinuousOutput("K")

		g.gameState.OpenDoor(getCurrentPressedArrowKeyAsDirection())

		g.gameState.SecondaryKeyState = game_state.PrimaryInput
	default:
		panic("unhandled default case")
	}
}
