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
	ExittedVehicle   *NPCFriendly
	ResultingVehicle *NPCFriendly
}

func (g *GameState) ExitVehicle() VehicleExitResults {
	vehicleType := g.PartyVehicle.GetVehicleDetails().VehicleType

	if vehicleType == references.NoPartyVehicle {
		return VehicleExitResults{
			ExittedVehicle:   nil,
			ResultingVehicle: nil,
		}
	}

	var ver VehicleExitResults

	if vehicleType == references.FrigateVehicle {
		// check for skiffs
		if g.PartyVehicle.GetVehicleDetails().SkiffQuantity <= 0 {
			// boo no skiff
			ver = VehicleExitResults{
				ExittedVehicle:   &g.PartyVehicle,
				ResultingVehicle: nil,
			}
			return ver
		}

		skiff := NewNPCFriendlyVehiceNewRef(references.SkiffVehicle, g.Position, g.Floor)
		g.PartyVehicle = *skiff

		// yay, exit on skiff
		return VehicleExitResults{
			ExittedVehicle:   &g.PartyVehicle,
			ResultingVehicle: skiff,
		}
	}

	prevVehicle := g.PartyVehicle

	return VehicleExitResults{
		ExittedVehicle:   &prevVehicle,
		ResultingVehicle: nil,
	}
}
