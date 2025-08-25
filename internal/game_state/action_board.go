package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// ActionBoard - handles boarding vehicles
// Based on Commands.md Board section
func (g *GameState) ActionBoard() bool {
	// Not allowed in dungeons
	if g.MapState.PlayerLocation.Location.GetMapType() == references.DungeonMapType {
		g.SystemCallbacks.Message.AddRowStr("Not here!")
		return false
	}

	// Look for vehicle at current position
	vehicle := g.CurrentNPCAIController.GetNpcs().GetVehicleAtPositionOrNil(
		g.MapState.PlayerLocation.Position)

	if vehicle == nil {
		g.SystemCallbacks.Message.AddRowStr("Board what?")
		return false
	}

	// Check if already on a vehicle (must be on foot to board)
	if g.PartyVehicle.GetVehicleDetails().VehicleType != references.NoPartyVehicle {
		g.SystemCallbacks.Message.AddRowStr("On foot!")
		return false
	}

	vehicleType := vehicle.GetVehicleDetails().VehicleType

	// Print boarding message based on vehicle type
	switch vehicleType {
	case references.HorseVehicle:
		// TODO: Check if NPC refuses mount in towns (npc_refuses_mount check)
		g.SystemCallbacks.Message.AddRowStr("Board horse!")

	case references.CarpetVehicle:
		g.SystemCallbacks.Message.AddRowStr("Board carpet!")

	case references.SkiffVehicle:
		g.SystemCallbacks.Message.AddRowStr("Board skiff!")

	case references.FrigateVehicle:
		g.SystemCallbacks.Message.AddRowStr("Board Ship!")

	default:
		g.SystemCallbacks.Message.AddRowStr("Board what?")
		return false
	}

	// Attempt to board the vehicle
	result := g.BoardVehicle(*vehicle)

	switch result {
	case BoardVehicleResultSuccess:
		// Success message already printed above
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true

	case BoardVehicleResultNoVehicle:
		g.SystemCallbacks.Message.AddRowStr("Board what?")
		return false

	case BoardVehicleResultNoSkiffs:
		// Ship boarded but warning about no skiffs
		g.SystemCallbacks.Message.AddRowStr("WARNING: NO SKIFFS ON BOARD!")
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true
	}

	return false
}
