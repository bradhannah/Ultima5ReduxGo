package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
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

	bIsArrowKeyPressed := isArrowKeyPressed()

	switch g.gameState.SecondaryKeyState {
	case game_state.OpenDirectionInput:
		if !bIsArrowKeyPressed {
			return
		}
		// TODO: don't love that this has to be in every single function - but may be required
		// if it is not in each individually then stray keystrokes (like Action keys like 'O' )
		// can reset the timer and delay the operation
		if !g.keyboard.TryToRegisterKeyPress() {
			return
		}

		g.output.AppendToOutput(getCurrentPressedArrowKeyAsDirection().GetDirectionCompassName())

		switch g.gameState.OpenDoor(getCurrentPressedArrowKeyAsDirection()) {
		case game_state.NotADoor:
			g.output.AddToContinuousOutput("Nothing to open!")
		case game_state.Locked:
			g.output.AddToContinuousOutput("Locked!")
		case game_state.LockedMagical:
			g.output.AddToContinuousOutput("Magically Locked!")
		case game_state.Opened:
			g.output.AddToContinuousOutput("Opened!")
		default:
			log.Fatal("Unrecognized door open state")
		}

		g.gameState.SecondaryKeyState = game_state.PrimaryInput
	default:
		panic("unhandled default case")
	}

	if bIsArrowKeyPressed {
		g.keyboard.SetLastKeyPressedNow()
	}
}
