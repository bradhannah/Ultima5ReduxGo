package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *GameScene) largeMapInputHandler(key ebiten.Key) {
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
		g.secondaryKeyState = KlimbDirectionInput
	case ebiten.KeyX:
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

}
