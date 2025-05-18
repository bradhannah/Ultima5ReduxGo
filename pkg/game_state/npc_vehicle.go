package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type VehicleDetails struct {
	VehicleType   references.VehicleType
	SkiffQuantity int

	currentDirection  references.Direction
	previousDirection references.Direction
}

// func NewNPCVehicle(vehicleType references.VehicleType, pos references.Position, floor references.FloorNumber) NPCFriendly {
// 	return *NewNPCFriendlyVehicle()
// 	return *NewNPCFriendly(

// 	)
// 	VehicleDetails{
// 		VehicleType:   vehicleType,
// 		SkiffQuantity: 0,
// 	}
// }

func (n *VehicleDetails) SetPartyVehicleDirection(direction references.Direction) {
	switch n.VehicleType {
	case references.HorseVehicle, references.CarpetVehicle:
		if direction == references.Up || direction == references.Down {
			return
		}
	case references.FrigateVehicle, references.SkiffVehicle, references.NoPartyVehicle:
	}
	n.previousDirection = n.currentDirection
	n.currentDirection = direction
}

func (n *VehicleDetails) GetSpriteIndex() indexes.SpriteIndex {
	return n.VehicleType.GetSpriteByDirection(n.previousDirection, n.currentDirection)
}

func (g *VehicleDetails) DoesMoveResultInMovement(newDirection references.Direction) bool {
	if g.VehicleType != references.FrigateVehicle {
		return true
	}
	if g.currentDirection == newDirection {
		return true
	}
	return false
}
