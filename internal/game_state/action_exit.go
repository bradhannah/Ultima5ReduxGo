package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) DebugQuickExitSmallMap() {
	g.MapState.PlayerLocation.Location = references2.Britannia_Underworld
	g.MapState.PlayerLocation.Floor = g.LastLargeMapFloor
	g.MapState.PlayerLocation.Position = g.LastLargeMapPosition
	g.CurrentNPCAIController = g.GetCurrentLargeMapNPCAIController()
	g.MapState.UpdateLargeMap()
}

// ExitVehicle - exits the vehicle the player is currently in
// returns the previously boarded vehicle - or nil if none was found
func (g *GameState) ExitVehicle() *map_units.NPCFriendly {
	vehicleType := g.PartyVehicle.GetVehicleDetails().VehicleType

	if vehicleType == references2.NoPartyVehicle {
		return nil
	}

	if vehicleType == references2.FrigateVehicle {
		frigate := g.PartyVehicle
		g.CurrentNPCAIController.GetNpcs().AddVehicle(g.PartyVehicle)
		// check for skiffs
		if g.PartyVehicle.GetVehicleDetails().SkiffQuantity <= 0 {
			g.PartyVehicle = map_units.NewNPCFriendlyVehiceNoVehicle()
			return &frigate
		}

		skiff := map_units.NewNPCFriendlyVehiceNewRef(references2.SkiffVehicle, g.MapState.PlayerLocation.Position, g.MapState.PlayerLocation.Floor)
		g.PartyVehicle = *skiff

		// yay, exit on skiff
		return &frigate
	}

	prevVehicle := g.PartyVehicle
	g.PartyVehicle = map_units.NewNPCFriendlyVehiceNoVehicle()

	// set the vehicle to now use your existing position as it's normal position
	prevVehicle.SetPos(g.MapState.PlayerLocation.Position)
	prevVehicle.NPCReference.Schedule.OverrideAllPositions(byte(g.MapState.PlayerLocation.Position.X), byte(g.MapState.PlayerLocation.Position.Y))
	prevVehicle.SetFloor(g.MapState.PlayerLocation.Floor)

	g.CurrentNPCAIController.GetNpcs().AddVehicle(prevVehicle)

	return &prevVehicle
}
