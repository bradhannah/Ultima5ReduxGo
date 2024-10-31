package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *GameScene) largeMapInputHandler(key ebiten.Key) {
	switch key {
	case ebiten.KeySpace:
		g.addRowStr("Pass")
		g.gameState.DateTime.Advance(game_state.DefaultLargeMapMinutesPerTurn)
	case ebiten.KeyBackquote:
		g.bShowDebugConsole = !g.bShowDebugConsole
	case ebiten.KeyEnter:
		g.addRowStr("Enter")
	case ebiten.KeyUp:
		g.handleMovement(references.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(references.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(references.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(references.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyK:
		g.addRowStr("Klimb-")
		g.secondaryKeyState = KlimbDirectionInput
	case ebiten.KeyG:
		g.addRowStr("Get what?")
	case ebiten.KeyX:
		g.addRowStr("X-it what?")
	case ebiten.KeyP:
		g.addRowStr("Push what?")
	case ebiten.KeyL:
		g.addRowStr("Look-")
		g.secondaryKeyState = LookDirectionInput
	case ebiten.KeyE:
		g.debugMessage = "Enter a place"
		newLocation := g.gameReferences.LocationReferences.WorldLocations.GetLocationByPosition(g.gameState.Position)

		if newLocation != references.EmptyLocation {
			maps := g.gameReferences.LocationReferences.GetLocationReference(newLocation)
			g.gameState.EnterBuilding(maps, g.gameReferences.TileReferences)
			g.addRowStr(fmt.Sprintf("%s",
				g.gameReferences.LocationReferences.GetLocationReference(newLocation).EnteringText))
		} else {
			g.addRowStr(fmt.Sprintf("Enter what?"))
		}

	case ebiten.KeyO:
		g.addRowStr("Open-")
		g.appendToCurrentRowStr("Cannot")
	case ebiten.KeyJ:
		g.addRowStr("Jimmy-")
		g.appendToCurrentRowStr("Cannot")
	}

	// only process end of turn if the turn is actually done.
	if g.secondaryKeyState == PrimaryInput {
		g.gameState.LargeMapProcessEndOfTurn()
	}
}

func (g *GameScene) largeMapHandleSecondaryInput() {

	arrowKey := getArrowKeyPressed()
	bIsArrowKeyPressed := arrowKey != nil

	switch g.secondaryKeyState {
	case KlimbDirectionInput:
		if !bIsArrowKeyPressed {
			return
		}
		if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
			return
		}
		g.appendDirectionToOutput()

		g.secondaryKeyState = PrimaryInput
	case LookDirectionInput:
		if !bIsArrowKeyPressed {
			return
		}
		if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
			return
		}
		g.appendDirectionToOutput()

		newPosition := getCurrentPressedArrowKeyAsDirection().GetNewPositionInDirection(&g.gameState.Position)
		topTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(newPosition)
		g.addRowStr(g.gameReferences.LookReferences.GetTileLookDescription(topTile.Index))

		g.secondaryKeyState = PrimaryInput
	default:
		// better safe than sorry
		g.secondaryKeyState = PrimaryInput
	}
}
