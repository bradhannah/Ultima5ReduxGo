package references

import (
	"fmt"
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

type VehicleType int

const (
	NoPartyVehicle VehicleType = iota
	CarpetVehicle
	HorseVehicle
	SkiffVehicle
	FrigateVehicle
	NPC
)

func (v VehicleType) GetSpriteByDirection(previousDirection, direction Direction) indexes.SpriteIndex {
	switch v {
	case CarpetVehicle:
		// note previousDirection exists only for horses and carpet because they do not have vertical
		// tiles, so we keep the previous left/right direction
		if direction == Left || (previousDirection == Left && (direction == Up || direction == Down)) {
			return indexes.AvatarRidingCarpetLeft
		}
		return indexes.AvatarRidingCarpetRight
	case HorseVehicle:
		if direction == Left || (previousDirection == Left && (direction == Up || direction == Down)) {
			return indexes.AvatarRidingHorseLeft
		}
		return indexes.AvatarRidingHorseRight
	case FrigateVehicle:
		switch direction {
		case Left:
			return indexes.FrigateLeftFurled
		case Right:
			return indexes.FrigateRightFurled
		case Up:
			return indexes.FrigateUpFurled
		case Down:
			return indexes.FrigateDownFurled
		default:
			return indexes.FrigateRightFurled
		}
	case SkiffVehicle:
		switch direction {
		case Left:
			return indexes.SkiffLeft
		case Right:
			return indexes.SkiffRight
		case Up:
			return indexes.SkiffUp
		case Down:
			return indexes.SkiffDown
		default:
			return indexes.SkiffLeft
		}
	default:
		return indexes.Avatar_KeyIndex
	}
}

func (v VehicleType) RequiresNewSprite(currentDirection, newDirection Direction) bool {
	switch v {
	case CarpetVehicle, HorseVehicle:
		if currentDirection == Up || currentDirection == Down || currentDirection == newDirection {
			return false
		}
		return true
	case SkiffVehicle:
		return currentDirection != newDirection
	case FrigateVehicle:
		return currentDirection != newDirection
	case NoPartyVehicle:
		return false
	}
	panic("Unexpected: must have valid vehicle")
}

func (v VehicleType) GetMovementPrefix() string {
	switch v {
	case CarpetVehicle:
		return "Fly "
	case HorseVehicle:
		return "Ride "
	case SkiffVehicle:
		return "Row "
	case FrigateVehicle:
		return "Row "
	case NoPartyVehicle:
		return ""
	default:
		return ""
	}
}

func (v VehicleType) GetLowerCaseName() string {
	switch v {
	case NoPartyVehicle:
		return "none"
	case FrigateVehicle:
		return "frigate"
	case CarpetVehicle:
		return "carpet"
	case SkiffVehicle:
		return "skiff"
	case HorseVehicle:
		return "horse"
	default:
		log.Fatalf("Bad vehicle type %d", v)
	}
	return ""
}

func (v VehicleType) GetExitString() string {
	if v == NoPartyVehicle {
		return "X-it what?"
	}
	return fmt.Sprintf("Xit- %s!", v.GetLowerCaseName())
}

func (v VehicleType) GetBoardString() string {
	if v == NoPartyVehicle {
		return "X-it what?"
	}
	return fmt.Sprintf("Board %s!", v.GetLowerCaseName())
}

func GetVehicleFromSpriteIndex(s indexes.SpriteIndex) VehicleType {
	switch s {
	case indexes.Carpet2_MagicCarpet:
		return CarpetVehicle
	case indexes.HorseRight, indexes.HorseLeft:
		return HorseVehicle
	case indexes.FrigateDownFurled, indexes.FrigateUpFurled, indexes.FrigateLeftFurled, indexes.FrigateRightFurled:
		return FrigateVehicle
	case indexes.SkiffLeft, indexes.SkiffRight, indexes.SkiffUp, indexes.SkiffDown:
		return SkiffVehicle
	default:
		return NoPartyVehicle
	}
}
