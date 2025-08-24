package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
	// Auto-shut doors if needed (equivalent to shut_previous_door)
	// TODO: Implement door timer shutdown logic when door system is implemented

	pushableThingPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)

	smallMap := g.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.MapState.PlayerLocation.Floor)
	pushableThingTile := smallMap.GetTopTile(pushableThingPosition)

	// CRITICAL FIX: Validate object presence and pushability
	if g.ObjectPresentAt(pushableThingPosition) || !g.IsPushable(pushableThingTile) {
		g.SystemCallbacks.Message.AddRowStr("Won't budge!")
		return false
	}

	// CRITICAL FIX: Legal floor depends on object type
	legalFloorIndex := indexes.BrickFloor
	if pushableThingTile.IsCannon() {
		legalFloorIndex = indexes.HexMetalGridFloor // Cannons require hex metal floor
	}
	_ = legalFloorIndex // TODO: Use this when floor validation is implemented

	farSideOfPushableThingPosition := direction.GetNewPositionInDirection(pushableThingPosition)
	farSideTile := smallMap.GetTopTile(farSideOfPushableThingPosition)
	playerTile := smallMap.GetTopTile(&g.MapState.PlayerLocation.Position)

	// Try to push object forward into farSide
	if !g.IsOutOfBounds(*farSideOfPushableThingPosition) &&
		!g.ObjectPresentAt(farSideOfPushableThingPosition) &&
		farSideTile.Index.IsPushableFloor() {
		// Push object forward
		g.pushIt(smallMap, pushableThingTile, farSideTile, pushableThingPosition, farSideOfPushableThingPosition, direction)
		g.SystemCallbacks.Message.AddRowStr("Pushed!")
		g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
	} else {
		// Try swapping player with object if player stands on legal floor
		if playerTile.Index.IsPushableFloor() {
			g.swapIt(smallMap, playerTile, pushableThingTile, &g.MapState.PlayerLocation.Position, pushableThingPosition, direction)
			g.SystemCallbacks.Message.AddRowStr("Pulled!")
			g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
		} else {
			g.SystemCallbacks.Message.AddRowStr("Won't budge!")
			return false
		}
	}

	// CRITICAL FIX: Only move player forward after successful push/pull validation
	g.MapState.PlayerLocation.Position = *pushableThingPosition

	// Advance time - pushing objects takes 1 minute
	g.SystemCallbacks.Flow.AdvanceTime(1)

	// Screen redraw is handled automatically by the rendering system

	return true
}

// Helper function to handle pushing object forward
func (g *GameState) pushIt(smallMap *map_state.LayeredMap, objTile, floorTile *references.Tile, fromPos, toPos *references.Position, direction references.Direction) {
	smallMap.SwapTiles(fromPos, toPos)
	if objTile.IsChair() {
		smallMap.SetTileByLayer(map_state.MapLayer, toPos, g.GameReferences.TileReferences.GetChairByPushDirection(direction).Index)
	} else if objTile.IsCannon() {
		smallMap.SetTileByLayer(map_state.MapLayer, toPos, g.GameReferences.TileReferences.GetCannonByPushDirection(direction).Index)
	}
}

// Helper function to handle swapping player with object (pull)
func (g *GameState) swapIt(smallMap *map_state.LayeredMap, playerTile, objTile *references.Tile, playerPos, objPos *references.Position, direction references.Direction) {
	smallMap.SwapTiles(playerPos, objPos)
	if objTile.IsChair() {
		smallMap.SetTileByLayer(map_state.MapLayer, playerPos, g.GameReferences.TileReferences.GetChairByPushDirection(direction.GetOppositeDirection()).Index)
	} else if objTile.IsCannon() {
		smallMap.SetTileByLayer(map_state.MapLayer, playerPos, g.GameReferences.TileReferences.GetCannonByPushDirection(direction.GetOppositeDirection()).Index)
	}
}

func (g *GameState) ActionPushLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Push command - see Commands.md Push section
	// Large map variant of push command
	return true
}

func (g *GameState) ActionPushCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Push command - see Commands.md Push section
	// Combat map push includes combat-specific logic:
	// - save_avatar_pos(); map_to_active_combatant()
	// - Check for stepping trap guarding push: CheckTrap(pushableThingPosition)
	// - sync_combatant_and_object_positions(dx, dy) after push
	// - restore_avatar_pos(); update_screen()
	return true // Temporary stub
}

// Helper methods - TODO: Move these to appropriate files when systems are implemented

func (g *GameState) IsPushable(tile *references.Tile) bool {
	// TODO: Implement full pushability check per Commands.md Push section
	// Should check for: PLANT, CHAIR variants, DESK, BARREL, VANITY, PITCHER,
	// DRAWERS, END_TABLE, FOOTLOCKER, CANNON variants
	return tile.IsChair() || tile.IsCannon() // Temporary implementation
}

func (g *GameState) ObjectPresentAt(position *references.Position) bool {
	// TODO: Implement object presence check when object system is available
	return false // Temporary stub
}

func (g *GameState) CheckTrap(position *references.Position) bool {
	// TODO: Implement trap checking when trap system is available
	return false // Temporary stub
}

// Note: IsInCombat() check removed from ActionPushSmallMap since small maps are never combat maps
// Combat-specific logic belongs in ActionPushCombatMap
