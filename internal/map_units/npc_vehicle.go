package map_units

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

type VehicleDetails struct {
	VehicleType   references.VehicleType
	skiffQuantity int

	currentDirection  references.Direction
	previousDirection references.Direction
}

func NewVehicleDetails(vehicleType references.VehicleType) VehicleDetails {
	return VehicleDetails{
		VehicleType:       vehicleType,
		skiffQuantity:     0,
		currentDirection:  references.Right,
		previousDirection: references.Right,
	}
}

func (v *VehicleDetails) SetPartyVehicleDirection(direction references.Direction) {
	switch v.VehicleType { //nolint:exhaustive
	case references.HorseVehicle, references.CarpetVehicle:
		if direction == references.Up || direction == references.Down {
			return
		}
	case references.FrigateVehicle, references.SkiffVehicle, references.NoPartyVehicle:
	default:
		panic("unhandled default case")
	}

	v.previousDirection = v.currentDirection
	v.currentDirection = direction
}

func (v *VehicleDetails) GetBoardedSpriteIndex() indexes.SpriteIndex {
	return v.VehicleType.GetBoardedSpriteByDirection(v.previousDirection, v.currentDirection)
}

func (v *VehicleDetails) GetUnBoardedSpriteIndex() indexes.SpriteIndex {
	return v.VehicleType.GetUnBoardedSpriteByDirection(v.currentDirection)
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

// SetSkiffQuantity sets the skiff quantity to the given quantity.
func (v *VehicleDetails) SetSkiffQuantity(quantity int) {
	v.skiffQuantity = quantity
}

// DecrementSkiffQuantity decrements the skiff quantity by one.
func (v *VehicleDetails) DecrementSkiffQuantity() {
	if v.skiffQuantity == 0 {
		log.Fatal("skiff quantity is 0, cannot decrement it.")
	}

	v.skiffQuantity--
}

// HasAtLeastOneSkiff returns true if the vehicle has at least one skiff, false otherwise.
func (v *VehicleDetails) HasAtLeastOneSkiff() bool {
	return v.skiffQuantity > 0
}

func (v *VehicleDetails) IncrementSkiffQuantity() {
	v.skiffQuantity++
}
