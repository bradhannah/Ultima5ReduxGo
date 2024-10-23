package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

// Update method for the GameScene
func (g *GameScene) Update(_ *Game) error {
	if g.bShowDebugConsole {
		// if debug console is showing, then we eat all the input
		g.debugConsole.update()
		return nil
	}

	if g.gameState.SecondaryKeyState != game_state.PrimaryInput {
		g.handleSecondaryInput()
		return nil
	}

	// Handle gameplay logic here
	boundKey := g.keyboard.GetBoundKeyPressed(&boundKeysGame)
	if boundKey == nil {
		// allow quick tapping
		g.keyboard.SetAllowKeyPressImmediately()
		return nil
	}

	if !g.keyboard.TryToRegisterKeyPress(*boundKey) {
		return nil
	}

	switch *boundKey {
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
	case ebiten.KeyX:
		g.gameState.Location = references.Britannia_Underworld
		g.gameState.Floor = 0
		g.gameState.Position = g.gameState.LastLargeMapPosition
	case ebiten.KeyE:
		g.debugMessage = "Enter a place"
		newLocation := g.gameReferences.SingleMapReferences.WorldLocations.GetLocationByPosition(g.gameState.Position)

		if newLocation != references.EmptyLocation {
			maps := g.gameReferences.SingleMapReferences.GetLocationReference(newLocation)
			g.gameState.EnterBuilding(maps, g.gameReferences.TileReferences)
			g.addRowStr(fmt.Sprintf("%s",
				g.gameReferences.SingleMapReferences.GetLocationReference(newLocation).EnteringText))

		}

	case ebiten.KeyO:
		g.debugConsole.Output.AddRowStr("Open")
		g.addRowStr("Open-")
		if g.gameState.Location == references.Britannia_Underworld {
			g.appendToCurrentRowStr("Cannot")
			return nil
		}
		g.gameState.SecondaryKeyState = game_state.OpenDirectionInput
		// we don't want the delay, it feels unnatural
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyJ:
		g.debugMessage = "Jimmy"
		g.addRowStr("Jimmy-")
		if g.gameState.Location == references.Britannia_Underworld {
			g.appendToCurrentRowStr("Cannot")
			return nil
		}

		g.gameState.SecondaryKeyState = game_state.JimmyDoorDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	}

	// only process end of turn if the turn is actually done.
	if g.gameState.SecondaryKeyState == game_state.PrimaryInput {
		g.gameState.ProcessEndOfTurn()
	}
	return nil
}
