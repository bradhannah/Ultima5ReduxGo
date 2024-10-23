package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"golang.org/x/exp/rand"
)

type JimmyDoorState int

const (
	JimmyDoorNoLockPicks JimmyDoorState = iota
	JimmyUnlocked
	JimmyLockedMagical
	JimmyNotADoor
	JimmyBrokenPick
)

func (g *GameState) JimmyDoor(direction Direction, player *PlayerCharacter) JimmyDoorState {
	mapType := GetMapTypeByLocation(g.Location)

	newPosition := direction.GetNewPositionInDirection(&g.Position)
	targetTile := g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).GetTileTopMapOnlyTile(newPosition)

	switch targetTile.Index {
	case indexes.LockedDoor, indexes.LockedDoorView:
		if g.isJimmySuccessful(player) {
			var unlockedDoor indexes.SpriteIndex
			if targetTile.Index == indexes.LockedDoor {
				unlockedDoor = indexes.RegularDoor
			} else {
				unlockedDoor = indexes.RegularDoorView
			}
			g.LayeredMaps.GetLayeredMap(mapType, g.Floor).SetTile(MapOverrideLayer, newPosition, unlockedDoor)
			return JimmyUnlocked
		} else {
			//g.Provisions.QtyKeys = helpers.Max(g.Provisions.QtyKeys-1, 0)
			return JimmyBrokenPick
		}
	case indexes.MagicLockDoor, indexes.MagicLockDoorWithView:
		g.Provisions.QtyKeys = helpers.Max(g.Provisions.QtyKeys-1, 0)
		return JimmyLockedMagical
	default:
		return JimmyNotADoor
	}

}

func (g *GameState) isJimmySuccessful(player *PlayerCharacter) bool {
	// TODO: actual jimmy logic required, currently forced to 50%
	n := rand.Int()
	return n%2 == 0
}
