package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameScene) largeMapInputHandler(key ebiten.Key) {
	// if ebiten.IsKeyPressed(ebiten.KeyControl) {
	// 	if ebiten.IsKeyPressed(ebiten.KeyX) {

	// 		return
	// 	}
	// }

	switch key {
	case ebiten.KeyEscape:
		g.DoEscapeMenu()

		return
	case ebiten.KeySpace:
		g.addRowStr("Pass")
		// g.gameState.FinishTurn()
	case ebiten.KeyBackquote:
		g.toggleDebug()
		return
	case ebiten.KeyEnter:
		g.addRowStr("Enter")
	case ebiten.KeyUp:
		g.handleMovement(references2.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(references2.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(references2.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(references2.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyB:
		g.actionBoard()
	case ebiten.KeyK:
		g.addRowStr("Klimb-")
		g.secondaryKeyState = KlimbDirectionInput
	case ebiten.KeyG:
		g.addRowStr("Get what?")
	case ebiten.KeyX:
		g.actionExit()
	case ebiten.KeyP:
		g.addRowStr("Push what?")
	case ebiten.KeyL:
		g.addRowStr("Look-")
		g.secondaryKeyState = LookDirectionInput
	case ebiten.KeyE:
		g.debugMessage = "Enter a place"
		newLocation := g.gameReferences.LocationReferences.WorldLocations.GetLocationByPosition(
			g.gameState.MapState.PlayerLocation.Position)

		if newLocation != references2.EmptyLocation {
			slr := g.gameReferences.LocationReferences.GetLocationReference(newLocation)
			g.gameState.EnterBuilding(slr)
			g.addRowStr(g.gameReferences.LocationReferences.GetLocationReference(newLocation).EnteringText)
		} else {
			g.addRowStr("Enter what?")
		}

	case ebiten.KeyO:
		g.addRowStr("Open-")
		g.appendToCurrentRowStr("Cannot")
	case ebiten.KeyJ:
		g.addRowStr("Jimmy-")
		g.appendToCurrentRowStr("Cannot")
	case ebiten.KeyI:
		g.debugMessage = "Ignite Torch"
		g.addRowStr("Ignite Torch!")

		if !g.gameState.IgniteTorch() {
			g.addRowStr("None owned!")
		}
	case ebiten.KeyT:
		g.addRowStr("Talk to who?")
	default:
		return
	}

	// only process end of turn if the turn is actually done.
	if g.secondaryKeyState == PrimaryInput {
		g.gameState.FinishTurn()
	}
}

func (g *GameScene) largeMapHandleSecondaryInput() {
	switch g.secondaryKeyState {
	case KlimbDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.secondaryKeyState = PrimaryInput
		}
	case LookDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			newPosition := getCurrentPressedArrowKeyAsDirection().GetNewPositionInDirection(&g.gameState.MapState.PlayerLocation.Position)
			topTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(newPosition)
			g.addRowStr(fmt.Sprintf("Thou dost see %s", g.gameReferences.LookReferences.GetTileLookDescription(topTile.Index)))
			g.secondaryKeyState = PrimaryInput
		}
	default:
		// better safe than sorry
		g.secondaryKeyState = PrimaryInput
	}
}
