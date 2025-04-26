package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func getArrowKeyPressed() *ebiten.Key {
	var keyPressed = ebiten.KeyF24
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

func (g *GameScene) checkAndAutoKlimbStairs(position *references.Position) bool {
	floorKlimbOffset := g.gameState.LayeredMaps.GetSmallMapFloorKlimbOffset(*position, g.gameState.Floor)
	if floorKlimbOffset != 0 {
		// are we on stairs? we need to change floors
		g.gameState.Floor += references.FloorNumber(floorKlimbOffset)
		g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
		return true
	}
	return false
}

func (g *GameScene) handleMovement(directionStr string, key ebiten.Key) {
	g.debugMessage = directionStr

	g.addRowStr(fmt.Sprintf("> %s%s", g.gameState.PartyVehicle.GetMovementPrefix(), directionStr))
	// moveResultsInMovement
	direction := game_state.GetKeyAsDirection(key)
	newPosition := direction.GetNewPositionInDirection(&g.gameState.Position)
	mapType := g.gameState.Location.GetMapType()
	if mapType == references.LargeMapType {
		if g.gameState.PartyVehicle == references.FrigateVehicle && !g.gameState.DoesMoveResultInMovement(direction) {
			g.gameState.SetPartyVehicleDirection(direction)
			g.output.AddRowStr(fmt.Sprintf("Head %s", direction.GetDirectionCompassName()))
			return
		}

		newPosition = newPosition.GetWrapped(references.XLargeMapTiles, references.YLargeMapTiles)

	}
	g.gameState.SetPartyVehicleDirection(direction)

	if mapType == references.SmallMapType && g.gameState.IsOutOfBounds(*newPosition) {
		g.dialogStack.DoModalInputBox(
			"Dost thou wish to leave?",
			g.createTextCommandExitBuilding(),
			g.keyboard)

		g.addRowStr("OUT OF BOUNDS")
		return
	}

	if g.gameState.IsPassable(newPosition) || g.gameState.DebugOptions.FreeMove {
		g.moveToNewPositionByDirection(direction)
		g.debugConsole.Output.AddRowStr(fmt.Sprintf("X: %d, Y: %d", newPosition.X, newPosition.Y))
		extraMovementStr := g.gameState.GetExtraMovementString()
		if extraMovementStr != "" {
			g.output.AddRowStr(extraMovementStr)
		}

		g.output.AddRowStr(fmt.Sprintf("X: %d, Y: %d", newPosition.X, newPosition.Y))
	} else {
		g.addRowStr("Blocked!")
	}

	if g.gameState.Location.GetMapType() == references.SmallMapType {
		g.checkAndAutoKlimbStairs(newPosition)
	}

}
