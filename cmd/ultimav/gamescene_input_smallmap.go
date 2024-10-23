package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func (g *GameScene) smallMapInputHandler(key ebiten.Key) {
	switch key {
	case ebiten.KeyBackquote:
		g.bShowDebugConsole = !g.bShowDebugConsole
	case ebiten.KeyEnter:
		g.addRowStr("Enter")
	case ebiten.KeyUp:
		g.handleMovement(game_state.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(game_state.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(game_state.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(game_state.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyK:

	case ebiten.KeyX:
		g.gameState.Location = references.Britannia_Underworld
		g.gameState.Floor = 0
		g.gameState.Position = g.gameState.LastLargeMapPosition
	case ebiten.KeyE:
		g.addRowStr(fmt.Sprintf("Enter what?"))
	case ebiten.KeyO:
		g.debugConsole.Output.AddRowStr("Open")
		g.addRowStr("Open-")
		g.gameState.SecondaryKeyState = game_state.OpenDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyJ:
		g.debugMessage = "Jimmy"
		g.addRowStr("Jimmy-")
		g.gameState.SecondaryKeyState = game_state.JimmyDoorDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	}

	// only process end of turn if the turn is actually done.
	if g.gameState.SecondaryKeyState == game_state.PrimaryInput {
		g.gameState.SmallMapProcessEndOfTurn()
	}
}

func (g *GameScene) smallMapHandleSecondaryInput() {
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
