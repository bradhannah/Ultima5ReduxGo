package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

func (g *GameState) DebugQuickExitSmallMap() {
	g.Location = references.Britannia_Underworld
	g.Floor = g.LastLargeMapFloor
	g.Position = g.LastLargeMapPosition
	g.CurrentNPCAIController = g.GetCurrentLargeMapNPCAIController()
	g.UpdateLargeMap()
}

type VehicleExitResults struct {
	ExittedVehicle   *NPCVehicle
	ResultingVehicle *NPCVehicle
}

func (g *GameState) ExitVehicle() VehicleExitResults {
	if g.PartyVehicle.VehicleType == references.NoPartyVehicle {
		return VehicleExitResults{
			ExittedVehicle:   nil,
			ResultingVehicle: nil,
		}
	}

	if g.PartyVehicle.VehicleType == references.FrigateVehicle {
		// check for skiffs
	}

	if references.VehicleType(g.PartyVehicle.SkiffQuantity) == references.NoPartyVehicle {
	}
	// exittedVehicle := g.PartyVehicle
	// g.PartyVehicle = references.NoPartyVehicle
	return VehicleExitResults{
		ExittedVehicle:   nil,
		ResultingVehicle: nil,
	}
}
