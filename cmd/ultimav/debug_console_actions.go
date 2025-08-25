package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameScene) debugMoveOnMap(position references2.Position) {
	g.gameState.MapState.PlayerLocation.Position = position
}

func (g *GameScene) debugFloorUp() bool {
	if g.gameState.MapState.PlayerLocation.Location == references2.Britannia_Underworld {
		if g.gameState.MapState.PlayerLocation.Floor == -1 {
			g.gameState.MapState.PlayerLocation.Floor = 0
			return true
		}
		return false
	}

	singleMap := g.GetCurrentLocationReference()

	_, maxFloor := singleMap.GetFloorMinMax()
	if g.gameState.MapState.PlayerLocation.Floor+1 > maxFloor {
		return false
	}
	g.gameState.MapState.PlayerLocation.Floor++
	g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)

	return true
}

func (g *GameScene) debugFloorDown() bool {
	if g.gameState.MapState.PlayerLocation.Location == references2.Britannia_Underworld {
		if g.gameState.MapState.PlayerLocation.Floor == 0 {
			g.gameState.MapState.PlayerLocation.Floor = -1
			return true
		}
		return false
	}

	singleMap := g.GetCurrentLocationReference()
	minFloor, _ := singleMap.GetFloorMinMax()
	if g.gameState.MapState.PlayerLocation.Floor-1 < minFloor {
		return false
	}
	g.gameState.MapState.PlayerLocation.Floor--
	g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
	return true
}

func (g *GameScene) debugFloorY(nFloor references2.FloorNumber) bool {
	if g.gameState.MapState.PlayerLocation.Location == references2.Britannia_Underworld {
		if nFloor < -1 {
			g.gameState.MapState.PlayerLocation.Floor = -1
		} else if nFloor > 0 {
			g.gameState.MapState.PlayerLocation.Floor = 0
		} else {
			g.gameState.MapState.PlayerLocation.Floor = nFloor
		}
		return true
	}
	singleMap := g.GetCurrentLocationReference()
	minFloor, maxFloor := singleMap.GetFloorMinMax()
	if nFloor < minFloor || nFloor > maxFloor {
		// bad
		return false
	}
	g.gameState.MapState.PlayerLocation.Floor = nFloor
	return true
}

// debugBoardVehicle creates and boards a vehicle at the current player position
// Bypasses all normal validation checks (dungeon restrictions, existing vehicles, etc.)
func (g *GameScene) debugBoardVehicle(vehicleType references2.VehicleType) bool {
	// Exit any current vehicle first
	if g.gameState.PartyVehicle.GetVehicleDetails().VehicleType != references2.NoPartyVehicle {
		g.gameState.ExitVehicle()
	}

	// Create a new vehicle at the current player position
	currentPos := g.gameState.MapState.PlayerLocation.Position
	currentFloor := g.gameState.MapState.PlayerLocation.Floor

	// Create NPCReference for the vehicle
	npcRef := references2.NewNPCReferenceForVehicle(vehicleType, currentPos, currentFloor)

	// Create the vehicle
	vehicle := map_units.NewNPCFriendlyVehicle(vehicleType, *npcRef)

	// Override schedule to keep it at current position
	vehicle.NPCReference.Schedule.OverrideAllPositions(byte(currentPos.X), byte(currentPos.Y))

	// Set initial skiff quantity for frigates
	if vehicleType == references2.FrigateVehicle {
		vehicle.GetVehicleDetails().SetSkiffQuantity(1)
	}

	// Set the position to current player position
	vehicle.SetPos(currentPos)

	// Initialize vehicle direction to Right (default facing)
	vehicle.GetVehicleDetails().SetPartyVehicleDirection(references2.Right)

	// Force board the vehicle (bypass normal validation)
	g.gameState.PartyVehicle = *vehicle

	return true
}
