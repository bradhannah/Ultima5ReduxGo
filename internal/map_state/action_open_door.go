package map_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type DoorOpenState int

const (
	OpenDoorLocked DoorOpenState = iota
	OpenDoorLockedMagical
	OpenDoorNotADoor
	OpenDoorOpened
)

func (g *MapState) OpenDoor(direction references.Direction) DoorOpenState {
	const defaultTurnsForDoorOpen = 2

	newPosition := direction.GetNewPositionInDirection(&g.PlayerLocation.Position)
	theMap := g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.PlayerLocation.Floor)
	targetTile := theMap.GetTileTopMapOnlyTile(newPosition)

	switch targetTile.Index {
	case indexes.MagicLockDoor, indexes.MagicLockDoorWithView:
		return OpenDoorLockedMagical
	case indexes.LockedDoor, indexes.LockedDoorView:
		return OpenDoorLocked
	case indexes.RegularDoor, indexes.RegularDoorView:
		break
	default:
		return OpenDoorNotADoor
	}

	theMap.SetTileByLayer(MapOverrideLayer, newPosition, indexes.BrickFloor)

	if g.openDoorPos != nil {
		theMap.UnSetTileByLayer(MapOverrideLayer, g.openDoorPos)
	}
	g.openDoorPos = newPosition
	g.openDoorTurns = defaultTurnsForDoorOpen
	return OpenDoorOpened
}
