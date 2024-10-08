package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"

type DoorOpenState int

const (
	OpenDoorLocked DoorOpenState = iota
	//Unlocked
	//
	OpenDoorLockedMagical
	OpenDoorNotADoor
	OpenDoorOpened
)

func (g *GameState) OpenDoor(direction Direction) DoorOpenState {
	const defaultTurnsForDoorOpen = 2
	mapType := GetMapTypeByLocation(g.Location)

	newPosition := direction.GetNewPositionInDirection(&g.Position)
	targetTile := g.LayeredMaps.LayeredMaps[SmallMap].GetTileTopMapOnlyTile(newPosition)

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

	g.LayeredMaps.LayeredMaps[mapType].SetTile(MapOverrideLayer, newPosition, indexes.BrickFloor)

	if g.openDoorPos != nil {
		g.LayeredMaps.LayeredMaps[mapType].UnSetTile(MapOverrideLayer, g.openDoorPos)
	}
	g.openDoorPos = newPosition
	g.openDoorTurns = defaultTurnsForDoorOpen
	return OpenDoorOpened
}
