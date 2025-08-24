package game_state

import (
	"golang.org/x/exp/rand"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
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
	// Check if we're in a dungeon to use the appropriate logic
	if g.MapState.PlayerLocation.Location.GetMapType() == references2.DungeonMapType {
		return g.jimmyDungeon(direction, player)
	} else {
		return g.jimmySurface(direction, player)
	}
}

func (g *GameState) jimmySurface(direction references2.Direction, player *party_state.PlayerCharacter) JimmyDoorState {
	mapType := map_state.GetMapTypeByLocation(g.MapState.PlayerLocation.Location)

	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	targetTile := g.MapState.LayeredMaps.GetLayeredMap(references2.SmallMapType, g.MapState.PlayerLocation.Floor).GetTileTopMapOnlyTile(newPosition)

	switch targetTile.Index {
	case indexes.LockedDoor, indexes.LockedDoorView:
		// Check if we have keys before attempting jimmy
		if g.PartyState.Inventory.Provisions.Keys.Get() <= 0 {
			return JimmyDoorNoLockPicks
		}

		// Try to jimmy the door
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
			// Key broke during failed jimmy attempt
			g.PartyState.Inventory.Provisions.Keys.DecrementByOne()
			return JimmyBrokenPick
		}
	case indexes.MagicLockDoor, indexes.MagicLockDoorWithView:
		// Magic doors cannot be jimmied, don't consume keys
		return JimmyLockedMagical
	case indexes.Chest:
		// TODO: Implement chest jimmy logic with trap handling
		// Check if we have keys before attempting jimmy
		if g.PartyState.Inventory.Provisions.Keys.Get() <= 0 {
			return JimmyDoorNoLockPicks
		}

		// Try to jimmy the chest
		if g.isJimmySuccessful(player) {
			// TODO: Open chest and handle contents/traps
			// For now, just report success - full chest opening logic needed
			return JimmyUnlocked
		} else {
			// Key broke during failed jimmy attempt
			g.PartyState.Inventory.Provisions.Keys.DecrementByOne()
			return JimmyBrokenPick
		}
	default:
		return JimmyNotADoor
	}
}

func (g *GameState) jimmyDungeon(direction references2.Direction, player *party_state.PlayerCharacter) JimmyDoorState {
	// TODO: Implement dungeon-specific jimmy logic with ahead-of-avatar targeting
	// Per Commands.md, dungeon doors and chests use ahead-of-avatar targeting
	// and have different tile handling (0xF0 masks for chest detection, 0xA0/0xE0/0xF0 for doors)

	// For now, use surface logic as fallback
	// This needs proper dungeon tile handling and chest support
	return g.jimmySurface(direction, player)
}

func (g *GameState) isJimmySuccessful(player *party_state.PlayerCharacter) bool {
	// Testing override
	if g.jimmySuccessForTesting != nil {
		return g.jimmySuccessForTesting(player)
	}

	_ = player
	// TODO: actual jimmy logic required, currently forced to 50%
	// Should be based on character dexterity and other stats
	n := rand.Int()
	return n%2 == 0
}
