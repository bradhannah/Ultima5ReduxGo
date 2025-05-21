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

func (m *MapState) OpenDoor(direction references.Direction) DoorOpenState {
	const defaultTurnsForDoorOpen = 2

	newPosition := direction.GetNewPositionInDirection(&m.PlayerLocation.Position)
	theMap := m.LayeredMaps.GetLayeredMap(references.SmallMapType, m.PlayerLocation.Floor)
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

	if m.openDoorPos != nil {
		theMap.UnSetTileByLayer(MapOverrideLayer, m.openDoorPos)
	}
	m.openDoorPos = newPosition
	m.openDoorTurns = defaultTurnsForDoorOpen
	return OpenDoorOpened
}
