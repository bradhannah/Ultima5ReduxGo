package game_state

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

func (g *GameState) OpenDoor(direction references.Direction) DoorOpenState {
	const defaultTurnsForDoorOpen = 2
	//  mapType := GetMapTypeByLocation(g.Location)

	newPosition := direction.GetNewPositionInDirection(&g.Position)
	theMap := g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor)
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
