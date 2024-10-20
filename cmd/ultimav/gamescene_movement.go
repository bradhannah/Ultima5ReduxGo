package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func isArrowKeyPressed() bool {
	return getArrowKeyPressed() != nil
	//return ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyRight)
}

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

func getCurrentPressedArrowKeyAsDirection() game_state.Direction {
	arrowKey := getArrowKeyPressed()
	if arrowKey == nil {
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
	g.addRowStr(fmt.Sprintf("> %s", directionStr))

	isPassable := func(pos *references.Position) bool {
		topTile := g.gameState.LayeredMaps.LayeredMaps[g.gameState.GetMapType()].GetTopTile(pos)
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

func (g *GameScene) handleSecondaryInput() {

	arrowKey := getArrowKeyPressed()
	bIsArrowKeyPressed := arrowKey != nil

	switch g.gameState.SecondaryKeyState {
	case game_state.JimmyDoorDirectionInput:
		if g.gameState.Provisions.QtyKeys <= 0 {
			g.addRowStr("No Keys!")
			g.gameState.SecondaryKeyState = game_state.PrimaryInput
			g.keyboard.SetLastKeyPressedNow()
			return
		}

		if !bIsArrowKeyPressed {
			return
		}
		if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
			return
		}

		g.appendToCurrentRowStr(getCurrentPressedArrowKeyAsDirection().GetDirectionCompassName())

		jimmyResult := g.gameState.JimmyDoor(getCurrentPressedArrowKeyAsDirection(), &g.gameState.Characters[0])

		switch jimmyResult {
		case game_state.JimmyUnlocked:
			g.addRowStr("Unlocked!")
		case game_state.JimmyNotADoor:
			g.addRowStr("Not lock!")
		case game_state.JimmyBrokenPick, game_state.JimmyLockedMagical:
			g.addRowStr("Key broke!")

		default:
			panic("unhandled default case")
		}

		g.gameState.SecondaryKeyState = game_state.PrimaryInput

	case game_state.OpenDirectionInput:

		if !bIsArrowKeyPressed {
			return
		}
		if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
			return
		}

		g.appendToCurrentRowStr(getCurrentPressedArrowKeyAsDirection().GetDirectionCompassName())

		switch g.gameState.OpenDoor(getCurrentPressedArrowKeyAsDirection()) {
		case game_state.OpenDoorNotADoor:
			g.addRowStr("Nothing to open!")
		case game_state.OpenDoorLocked:
			g.addRowStr("Locked!")
		case game_state.OpenDoorLockedMagical:
			g.addRowStr("Magically Locked!")
		case game_state.OpenDoorOpened:
			g.addRowStr("Opened!")
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
