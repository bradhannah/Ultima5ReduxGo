package game_state

import (
	"golang.org/x/exp/rand"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

type JimmyDoorState int

const (
	JimmyDoorNoLockPicks JimmyDoorState = iota
	JimmyUnlocked
	JimmyLockedMagical
	JimmyNotADoor
	JimmyBrokenPick
)

func (g *GameState) JimmyDoor(direction references2.Direction, player *party_state.PlayerCharacter) JimmyDoorState {
	mapType := map_state.GetMapTypeByLocation(g.MapState.PlayerLocation.Location)

	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	targetTile := g.MapState.LayeredMaps.GetLayeredMap(references2.SmallMapType, g.MapState.PlayerLocation.Floor).GetTileTopMapOnlyTile(newPosition)

	switch targetTile.Index {
	case indexes.LockedDoor, indexes.LockedDoorView:
		if g.isJimmySuccessful(player) {
			var unlockedDoor indexes.SpriteIndex
			if targetTile.Index == indexes.LockedDoor {
				unlockedDoor = indexes.RegularDoor
			} else {
				unlockedDoor = indexes.RegularDoorView
			}
			g.MapState.LayeredMaps.GetLayeredMap(mapType, g.MapState.PlayerLocation.Floor).SetTileByLayer(map_state.MapOverrideLayer, newPosition, unlockedDoor)
			return JimmyUnlocked
		} else {
			// g.ProvisionsQuantity.Keys = helpers.Max(g.ProvisionsQuantity.Keys-1, 0)
			return JimmyBrokenPick
		}
	case indexes.MagicLockDoor, indexes.MagicLockDoorWithView:
		g.PartyState.Inventory.Provisions.Keys = helpers.Max(g.PartyState.Inventory.Provisions.Keys-1, 0)
		return JimmyLockedMagical
	default:
		return JimmyNotADoor
	}

}

func (g *GameState) isJimmySuccessful(player *party_state.PlayerCharacter) bool {
	_ = player
	// TODO: actual jimmy logic required, currently forced to 50%
	n := rand.Int()
	return n%2 == 0
}
