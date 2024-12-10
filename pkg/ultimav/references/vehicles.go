package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"

const (
	NoPartyVehicle PartyVehicle = iota
	CarpetVehicle
	HorseVehicle
	SkiffVehicle
	FrigateVehicle
)

func (v PartyVehicle) GetSpriteByDirection(previousDirection Direction, direction Direction) indexes.SpriteIndex {
	switch v {
	case CarpetVehicle:
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

func (v PartyVehicle) RequiresNewSprite(currentDirection Direction, newDirection Direction) bool {
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

func (v PartyVehicle) GetMovementPrefix() string {
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
