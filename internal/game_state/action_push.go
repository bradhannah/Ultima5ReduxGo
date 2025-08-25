package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
	// Auto-shut doors if needed (equivalent to shut_previous_door)
	g.MapState.SmallMapProcessTurnDoors()

	pushableThingPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	smallMap := g.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.MapState.PlayerLocation.Floor)
	pushableThingTile := smallMap.GetTopTile(pushableThingPosition)

	// Detailed validation with game state
	if g.ObjectPresentAt(pushableThingPosition) || !pushableThingTile.IsPushable() {
		g.SystemCallbacks.Message.AddRowStr("Won't budge!")
		return false
	}

	// Determine legal floor based on object type
	legalFloorIndex := indexes.SpriteIndex(indexes.BrickFloor)
	if pushableThingTile.IsCannon() {
		legalFloorIndex = indexes.SpriteIndex(indexes.HexMetalGridFloor) // Cannons require hex metal floor
	}

	farSideOfPushableThingPosition := direction.GetNewPositionInDirection(pushableThingPosition)
	farSideTile := smallMap.GetTopTile(farSideOfPushableThingPosition)
	playerTile := smallMap.GetTopTile(&g.MapState.PlayerLocation.Position)

	// Try to push object forward into farSide
	if !g.IsOutOfBounds(*farSideOfPushableThingPosition) &&
		!g.ObjectPresentAt(farSideOfPushableThingPosition) &&
		farSideTile.Index == legalFloorIndex {

		// Execute push
		g.pushIt(smallMap, pushableThingTile, farSideTile, pushableThingPosition, farSideOfPushableThingPosition, direction)
		g.MapState.PlayerLocation.Position = *pushableThingPosition

		// User feedback via callbacks
		g.SystemCallbacks.Message.AddRowStr("Pushed!")
		g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true

	} else if playerTile.Index == legalFloorIndex {
		// Try swapping (pulling)
		g.swapIt(smallMap, playerTile, pushableThingTile, &g.MapState.PlayerLocation.Position, pushableThingPosition, direction)
		g.MapState.PlayerLocation.Position = *pushableThingPosition

		// User feedback via callbacks
		g.SystemCallbacks.Message.AddRowStr("Pulled!")
		g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true
	}

	// Failed to push or pull
	g.SystemCallbacks.Message.AddRowStr("Won't budge!")
	return false
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

func (g *GameState) ActionPushDungeonMap(direction references.Direction) bool {
	// TODO: Implement dungeon map Push command - see Commands.md Push section
	// Dungeon map variant of push command
	return true
}

func (g *GameState) ActionPushCombatMap(direction references.Direction) bool {
	// Combat map push includes combat-specific logic per pseudocode

	// TODO: Implement when combat system is available:
	// - save_avatar_pos(); map_to_active_combatant()

	pushableThingPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)

	// Check for stepping trap guarding push
	if g.CheckTrap(pushableThingPosition) {
		// Trap triggered, don't continue with push
		return false
	}

	// Use same push logic as small map for now
	// TODO: Replace with combat-specific logic when combat maps are implemented
	result := g.ActionPushSmallMap(direction)

	// TODO: Implement when combat system is available:
	// - sync_combatant_and_object_positions(dx, dy) after push
	// - restore_avatar_pos(); update_screen()

	return result
}

// Helper methods - TODO: Move these to appropriate files when systems are implemented

func (g *GameState) ObjectPresentAt(position *references.Position) bool {
	// Check if any NPC/object is present at this position
	// This implements the original Ultima V looklist() functionality
	mapUnit := g.CurrentNPCAIController.GetNpcs().GetMapUnitAtPositionOrNil(*position)
	return mapUnit != nil
}

func (g *GameState) CheckTrap(position *references.Position) bool {
	// TODO: Implement trap checking when trap system is available
	return false // Temporary stub
}

// Note: IsInCombat() check removed from ActionPushSmallMap since small maps are never combat maps
// Combat-specific logic belongs in ActionPushCombatMap
