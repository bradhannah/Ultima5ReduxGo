package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func getArrowKeyPressed() *ebiten.Key {
	keyPressed := ebiten.KeyF24
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

func getCurrentPressedArrowKeyAsDirection() references2.Direction {
	arrowKey := getArrowKeyPressed()
	if arrowKey == nil {
		return references2.NoneDirection
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		return references2.Up
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		return references2.Down
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		return references2.Left
	}
	return references2.Right
}

func (g *GameScene) moveToNewPositionByDirection(direction references2.Direction) {
	// TODO: dear future Brad - you will need to change this big time when dungeons and combat come in
	bLargeMap := g.gameState.MapState.PlayerLocation.Location.GetMapType() == references2.LargeMapType
	switch direction {
	case references2.Up:
		g.gameState.MapState.PlayerLocation.Position.GoUp(bLargeMap)
	case references2.Down:
		g.gameState.MapState.PlayerLocation.Position.GoDown(bLargeMap)
	case references2.Left:
		g.gameState.MapState.PlayerLocation.Position.GoLeft(bLargeMap)
	case references2.Right:
		g.gameState.MapState.PlayerLocation.Position.GoRight(bLargeMap)
	case references2.NoneDirection:
	}
	if g.gameState.PartyVehicle.GetVehicleDetails().VehicleType != references2.NoPartyVehicle {
		g.gameState.PartyVehicle.SetPos(g.gameState.MapState.PlayerLocation.Position)
		g.gameState.PartyVehicle.NPCReference.Schedule.OverrideAllPositions(byte(g.gameState.MapState.PlayerLocation.Position.X), byte(g.gameState.MapState.PlayerLocation.Position.Y))
	}
}

func (g *GameScene) checkAndAutoKlimbStairs(position *references2.Position) bool {
	floorKlimbOffset := g.gameState.MapState.LayeredMaps.GetSmallMapFloorKlimbOffset(*position, g.gameState.MapState.PlayerLocation.Floor)
	if floorKlimbOffset != 0 {
		// are we on stairs? we need to change floors
		g.gameState.MapState.PlayerLocation.Floor += references2.FloorNumber(floorKlimbOffset)
		g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
		return true
	}
	return false
}

// handleMovement
// Handles the movement of the player in the game scene
func (g *GameScene) handleMovement(directionStr string, key ebiten.Key) {
	g.debugMessage = directionStr

	g.addRowStr(fmt.Sprintf("> %s%s", g.gameState.PartyVehicle.GetVehicleDetails().VehicleType.GetMovementPrefix(), directionStr))
	// moveResultsInMovement
	direction := map_state.GetKeyAsDirection(key)
	newPosition := direction.GetNewPositionInDirection(&g.gameState.MapState.PlayerLocation.Position)
	mapType := g.gameState.MapState.PlayerLocation.Location.GetMapType()

	if mapType == references2.LargeMapType {
		if g.gameState.PartyVehicle.GetVehicleDetails().VehicleType == references2.FrigateVehicle && !g.gameState.PartyVehicle.GetVehicleDetails().DoesMoveResultInMovement(direction) {
			g.gameState.PartyVehicle.GetVehicleDetails().SetPartyVehicleDirection(direction)
			g.output.AddRowStr(fmt.Sprintf("Head %s", direction.GetDirectionCompassName()))
			return
		}

		newPosition = newPosition.GetWrapped(references2.XLargeMapTiles, references2.YLargeMapTiles)

	}
	g.gameState.PartyVehicle.GetVehicleDetails().SetPartyVehicleDirection(direction)

	if mapType == references2.SmallMapType && g.gameState.IsOutOfBounds(*newPosition) {
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
		avatarTopTile := g.gameState.GetCurrentLayeredMapAvatarTopTile()
		extraMovementStr := avatarTopTile.GetExtraMovementString()
		if extraMovementStr != "" {
			g.output.AddRowStr(extraMovementStr)
		}

		g.output.AddRowStr(fmt.Sprintf("X: %d, Y: %d", newPosition.X, newPosition.Y))
	} else {
		g.addRowStr("Blocked!")
	}

	if g.gameState.MapState.PlayerLocation.Location.GetMapType() == references2.SmallMapType {
		g.checkAndAutoKlimbStairs(newPosition)
	}
}
