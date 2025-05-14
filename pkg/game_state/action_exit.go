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
	ExittedVehicle   references.PartyVehicle
	ResultingVehicle references.PartyVehicle
}

func (g *GameState) ExitVehicle() VehicleExitResults {
	if g.PartyVehicle == references.NoPartyVehicle {
		return VehicleExitResults{
			ExittedVehicle:   references.NoPartyVehicle,
			ResultingVehicle: references.NoPartyVehicle,
		}
	}
	exittedVehicle := g.PartyVehicle
	g.PartyVehicle = references.NoPartyVehicle
	return VehicleExitResults{
		ExittedVehicle:   exittedVehicle,
		ResultingVehicle: references.NoPartyVehicle,
	}
}
