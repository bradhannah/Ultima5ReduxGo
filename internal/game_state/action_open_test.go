package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func TestActionOpenLargeMap_Door_Success(t *testing.T) {
	// Test opening a regular door on large map
	gs := createTestGameStateWithMap(t)

	// Set up a regular door to the north
	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	layeredMap := gs.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.RegularDoor)

	// Test opening the door
	result := gs.ActionOpenLargeMap(direction)

	if !result {
		t.Error("Expected door to open successfully")
	}

	// Verify door state changed to floor
	topTile := layeredMap.GetTopTile(targetPos)
	if topTile.Index != indexes.BrickFloor {
		t.Errorf("Expected door to be opened (BrickFloor), got %v", topTile.Index)
	}

	// Verify timed door tracking is set
	if gs.MapState.GetOpenDoorPos() == nil {
		t.Error("Expected open door position to be tracked")
	}
	if gs.MapState.GetOpenDoorTurns() != 2 {
		t.Errorf("Expected 2 turns for door closure, got %d", gs.MapState.GetOpenDoorTurns())
	}
}

func TestActionOpenLargeMap_Door_Locked(t *testing.T) {
	// Test attempting to open a locked door
	gs := createTestGameStateWithMap(t)

	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	layeredMap := gs.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.LockedDoor)

	result := gs.ActionOpenLargeMap(direction)

	if result {
		t.Error("Expected locked door to NOT open")
	}

	// Verify door state unchanged
	topTile := layeredMap.GetTopTile(targetPos)
	if topTile.Index != indexes.LockedDoor {
		t.Errorf("Expected door to remain locked, got %v", topTile.Index)
	}
}

func TestActionOpenLargeMap_Door_MagicallyLocked(t *testing.T) {
	// Test attempting to open a magically locked door
	gs := createTestGameStateWithMap(t)

	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	layeredMap := gs.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.MagicLockDoor)

	result := gs.ActionOpenLargeMap(direction)

	if result {
		t.Error("Expected magically locked door to NOT open")
	}

	// Verify door state unchanged
	topTile := layeredMap.GetTopTile(targetPos)
	if topTile.Index != indexes.MagicLockDoor {
		t.Errorf("Expected door to remain magically locked, got %v", topTile.Index)
	}
}

func TestActionOpenLargeMap_ItemStack_Present(t *testing.T) {
	// Test when there's an item stack at the target position
	gs := createTestGameStateWithMap(t)

	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)

	// Create an item stack at target position
	itemStack := references.CreateNewItemStack(references.GoldPieces)
	gs.ItemStacksMap.Push(targetPos, &itemStack)

	result := gs.ActionOpenLargeMap(direction)

	if result {
		t.Error("Expected item stack to prevent opening (not implemented yet)")
	}
}

func TestActionOpenLargeMap_NotOpenable(t *testing.T) {
	// Test attempting to open something that's not openable
	gs := createTestGameStateWithMap(t)

	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	layeredMap := gs.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.BrickFloor) // Just floor

	result := gs.ActionOpenLargeMap(direction)

	if result {
		t.Error("Expected floor tile to not be openable")
	}
}

func TestActionOpenSmallMap_NotImplemented(t *testing.T) {
	// Test that small map open returns true (TODO implementation)
	gs := createTestGameStateWithMap(t)

	result := gs.ActionOpenSmallMap(references.North)

	if !result {
		t.Error("Expected small map open to return true (TODO)")
	}
}

func TestActionOpenCombatMap_NotImplemented(t *testing.T) {
	// Test that combat map open returns true (TODO implementation)
	gs := createTestGameStateWithMap(t)

	result := gs.ActionOpenCombatMap(references.North)

	if !result {
		t.Error("Expected combat map open to return true (TODO)")
	}
}

// Helper function to create a test game state with minimal map setup
func createTestGameStateWithMap(t *testing.T) *GameState {
	t.Helper()

	// This would need to be implemented based on your existing test setup patterns
	// For now, return a minimal GameState that won't crash
	gs := &GameState{}

	// You'll need to initialize the required components based on your architecture
	// This is a placeholder that shows the testing approach

	return gs
}
