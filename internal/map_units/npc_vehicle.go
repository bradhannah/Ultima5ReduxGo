package map_units

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

func (v *VehicleDetails) SetPartyVehicleDirection(direction references.Direction) {
	switch v.VehicleType {
	case references.HorseVehicle, references.CarpetVehicle:
		if direction == references.Up || direction == references.Down {
			return
		}
	case references.FrigateVehicle, references.SkiffVehicle, references.NoPartyVehicle:
	}
	v.previousDirection = v.currentDirection
	v.currentDirection = direction
}

func (v *VehicleDetails) GetSpriteIndex() indexes.SpriteIndex {
	return v.VehicleType.GetSpriteByDirection(v.previousDirection, v.currentDirection)
}

func (v *VehicleDetails) DoesMoveResultInMovement(newDirection references.Direction) bool {
	if v.VehicleType != references.FrigateVehicle {
		return true
	}
	if v.currentDirection == newDirection {
		return true
	}
	return false
}
