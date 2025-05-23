package map_units

import (
	"log"

	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

type VehicleDetails struct {
	VehicleType   references2.VehicleType
	skiffQuantity int

	currentDirection  references2.Direction
	previousDirection references2.Direction
}

func NewVehicleDetails(vehicleType references2.VehicleType) VehicleDetails {
	return VehicleDetails{
		VehicleType:       vehicleType,
		skiffQuantity:     0,
		currentDirection:  references2.Right,
		previousDirection: references2.Right,
	}
}

func (v *VehicleDetails) SetPartyVehicleDirection(direction references2.Direction) {
	switch v.VehicleType { //nolint:exhaustive
	case references2.HorseVehicle, references2.CarpetVehicle:
		if direction == references2.Up || direction == references2.Down {
			return
		}
	case references2.FrigateVehicle, references2.SkiffVehicle, references2.NoPartyVehicle:
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

func (v *VehicleDetails) DoesMoveResultInMovement(newDirection references2.Direction) bool {
	if v.VehicleType != references2.FrigateVehicle {
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
