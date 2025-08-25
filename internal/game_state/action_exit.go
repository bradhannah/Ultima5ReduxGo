package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) DebugQuickExitSmallMap() {
	g.MapState.PlayerLocation.Location = references.Britannia_Underworld
	g.MapState.PlayerLocation.Floor = g.LastLargeMapFloor
	g.MapState.PlayerLocation.Position = g.LastLargeMapPosition
	g.CurrentNPCAIController = g.GetCurrentLargeMapNPCAIController()
	g.MapState.UpdateLargeMap()
}

// ExitVehicle - exits the vehicle the player is currently in
// returns the previously boarded vehicle - or nil if none was found
func (g *GameState) ExitVehicle() *map_units.NPCFriendly {
	exittingVehicleType := g.PartyVehicle.GetVehicleDetails().VehicleType

	if exittingVehicleType == references.NoPartyVehicle {
		return nil
	}

	if exittingVehicleType == references.SkiffVehicle {
		mapVehicle := g.CurrentNPCAIController.GetNpcs().GetVehicleAtPositionOrNil(g.MapState.PlayerLocation.Position)
		if mapVehicle != nil {
			if mapVehicle.GetVehicleDetails().VehicleType == references.FrigateVehicle {
				return nil
			}
		}
	}

	if exittingVehicleType == references.FrigateVehicle {
		frigate := g.PartyVehicle

		// if we have no skiffs, then we are on the boat without a paddle
		if !g.PartyVehicle.GetVehicleDetails().HasAtLeastOneSkiff() {
			g.PartyVehicle = map_units.NewNPCFriendlyVehiceNoVehicle()

			g.CurrentNPCAIController.GetNpcs().AddVehicle(frigate)

			return &frigate
		}

		frigate.GetVehicleDetails().DecrementSkiffQuantity()
		g.CurrentNPCAIController.GetNpcs().AddVehicle(frigate)

		skiff := map_units.NewNPCFriendlyVehiceNewRef(references.SkiffVehicle,
			g.MapState.PlayerLocation.Position,
			g.MapState.PlayerLocation.Floor)
		skiff.SetPos(g.MapState.PlayerLocation.Position)
		g.PartyVehicle = *skiff

		// yay, exit on skiff
		return &frigate
	}

	prevVehicle := g.PartyVehicle
	g.PartyVehicle = map_units.NewNPCFriendlyVehiceNoVehicle()

	// set the vehicle to now use your existing position as it's a normal position
	prevVehicle.SetPos(g.MapState.PlayerLocation.Position)
	prevVehicle.NPCReference.Schedule.OverrideAllPositions(
		byte(g.MapState.PlayerLocation.Position.X), byte(g.MapState.PlayerLocation.Position.Y))
	prevVehicle.SetFloor(g.MapState.PlayerLocation.Floor)

	g.CurrentNPCAIController.GetNpcs().AddVehicle(prevVehicle)

	return &prevVehicle
}

// ActionExit - handles exiting/dismounting from vehicles
// Based on Commands.md patterns and vehicle system
func (g *GameState) ActionExit() bool {
	// Check if currently on/in a vehicle
	currentVehicleType := g.PartyVehicle.GetVehicleDetails().VehicleType
	if currentVehicleType == references.NoPartyVehicle {
		g.SystemCallbacks.Message.AddRowStr("X-it what?")
		return false
	}

	// Attempt to exit the current vehicle
	exittedVehicle := g.ExitVehicle()

	if exittedVehicle == nil {
		g.SystemCallbacks.Message.AddRowStr("X-it what?")
		return false
	}

	// Print appropriate exit message based on vehicle type
	switch exittedVehicle.GetVehicleDetails().VehicleType {
	case references.HorseVehicle:
		g.SystemCallbacks.Message.AddRowStr("Xit- horse!")
	case references.CarpetVehicle:
		g.SystemCallbacks.Message.AddRowStr("Xit- carpet!")
	case references.SkiffVehicle:
		g.SystemCallbacks.Message.AddRowStr("Xit- skiff!")
	case references.FrigateVehicle:
		g.SystemCallbacks.Message.AddRowStr("Xit- frigate!")
	default:
		g.SystemCallbacks.Message.AddRowStr("X-it what?")
		return false
	}

	g.SystemCallbacks.Flow.AdvanceTime(1)
	return true
}
