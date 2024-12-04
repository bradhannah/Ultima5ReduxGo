package game_state

import (
	"golang.org/x/exp/rand"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type JimmyDoorState int

const (
	JimmyDoorNoLockPicks JimmyDoorState = iota
	JimmyUnlocked
	JimmyLockedMagical
	JimmyNotADoor
	JimmyBrokenPick
)

func (g *GameState) JimmyDoor(direction references.Direction, player *PlayerCharacter) JimmyDoorState {
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
			g.LayeredMaps.GetLayeredMap(mapType, g.Floor).SetTileByLayer(MapOverrideLayer, newPosition, unlockedDoor)
			return JimmyUnlocked
		} else {
			// g.ProvisionsQuantity.Keys = helpers.Max(g.ProvisionsQuantity.Keys-1, 0)
			return JimmyBrokenPick
		}
	case indexes.MagicLockDoor, indexes.MagicLockDoorWithView:
		g.Inventory.Provisions.Keys = helpers.Max(g.Inventory.Provisions.Keys-1, 0)
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
