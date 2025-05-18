package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

func (g *GameState) DebugQuickExitSmallMap() {
	g.Location = references.Britannia_Underworld
	g.Floor = g.LastLargeMapFloor
	g.Position = g.LastLargeMapPosition
	g.CurrentNPCAIController = g.GetCurrentLargeMapNPCAIController()
	g.UpdateLargeMap()
}

// type VehicleExitResults struct {
// 	ExittedVehicle   *NPCFriendly
// 	ResultingVehicle *NPCFriendly
// }

// ExitVehicle
// returns the previously boarded vehicle - or nil if none was found
func (g *GameState) ExitVehicle() *NPCFriendly {
	vehicleType := g.PartyVehicle.GetVehicleDetails().VehicleType

	if vehicleType == references.NoPartyVehicle {
		return nil
	}

	if vehicleType == references.FrigateVehicle {
		frigate := g.PartyVehicle
		g.CurrentNPCAIController.GetNpcs().AddVehicle(g.PartyVehicle)
		// check for skiffs
		if g.PartyVehicle.GetVehicleDetails().SkiffQuantity <= 0 {
			g.PartyVehicle = NewNPCFriendlyVehiceNoVehicle()
			return &frigate
		}

		skiff := NewNPCFriendlyVehiceNewRef(references.SkiffVehicle, g.Position, g.Floor)
		g.PartyVehicle = *skiff

		// yay, exit on skiff
		return &frigate
	}

	prevVehicle := g.PartyVehicle
	g.PartyVehicle = NewNPCFriendlyVehiceNoVehicle()

	// set the vehicle to now use your existing position as it's normal position
	prevVehicle.SetPos(g.Position)
	prevVehicle.NPCReference.Schedule.OverrideAllPositions(byte(g.Position.X), byte(g.Position.Y))
	prevVehicle.SetFloor(g.Floor)

	g.CurrentNPCAIController.GetNpcs().AddVehicle(prevVehicle)

	return &prevVehicle
}

///func (g *GameState) setNo
