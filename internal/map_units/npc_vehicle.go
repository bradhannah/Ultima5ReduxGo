package map_units

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

type VehicleDetails struct {
	VehicleType   references2.VehicleType
	SkiffQuantity int

	currentDirection  references2.Direction
	previousDirection references2.Direction
}

func (v *VehicleDetails) SetPartyVehicleDirection(direction references2.Direction) {
	switch v.VehicleType {
	case references2.HorseVehicle, references2.CarpetVehicle:
		if direction == references2.Up || direction == references2.Down {
			return
		}
	case references2.FrigateVehicle, references2.SkiffVehicle, references2.NoPartyVehicle:
	}
	v.previousDirection = v.currentDirection
	v.currentDirection = direction
}

func (v *VehicleDetails) GetSpriteIndex() indexes.SpriteIndex {
	return v.VehicleType.GetSpriteByDirection(v.previousDirection, v.currentDirection)
}

func (v *VehicleDetails) DoesMoveResultInMovement(newDirection references2.Direction) bool {
	if v.VehicleType != references2.FrigateVehicle {
		return true
	}
	if v.currentDirection == newDirection {
		return true
	}
	return false
}
