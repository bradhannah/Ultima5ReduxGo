package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func TestIsPushable_Chair_Success(t *testing.T) {
	// Test chair tile
	chairTile := &references.Tile{Index: indexes.ChairFacingUp}

	result := chairTile.IsPushable()

	if !result {
		t.Error("Expected chair to be pushable")
	}
}

func TestIsPushable_Cannon_Success(t *testing.T) {
	// Test cannon tile
	cannonTile := &references.Tile{Index: indexes.CannonFacingUp}

	result := cannonTile.IsPushable()

	if !result {
		t.Error("Expected cannon to be pushable")
	}
}

func TestIsPushable_Barrel_Success(t *testing.T) {
	// Test barrel tile
	barrelTile := &references.Tile{Index: indexes.Barrel}

	result := barrelTile.IsPushable()

	if !result {
		t.Error("Expected barrel to be pushable")
	}
}

func TestIsPushable_NotPushable_Floor(t *testing.T) {
	// Test floor tile - should not be pushable
	floorTile := &references.Tile{Index: indexes.BrickFloor}

	result := floorTile.IsPushable()

	if result {
		t.Error("Expected floor to not be pushable")
	}
}

func TestIsPushable_NotPushable_Wall(t *testing.T) {
	// Test wall tile - should not be pushable
	wallTile := &references.Tile{Index: indexes.StoneBrickWall}

	result := wallTile.IsPushable()

	if result {
		t.Error("Expected wall to not be pushable")
	}
}

func TestIsPushable_TableMiddle_Success(t *testing.T) {
	tableMiddleTile := &references.Tile{Index: indexes.TableMiddle}

	result := tableMiddleTile.IsPushable()

	if !result {
		t.Error("Expected TableMiddle to be pushable")
	}
}

func TestIsPushable_Chest_Success(t *testing.T) {
	chestTile := &references.Tile{Index: indexes.Chest}

	result := chestTile.IsPushable()

	if !result {
		t.Error("Expected Chest to be pushable")
	}
}

func TestIsPushable_EndTable_Success(t *testing.T) {
	endTableTile := &references.Tile{Index: indexes.EndTable}

	result := endTableTile.IsPushable()

	if !result {
		t.Error("Expected EndTable to be pushable")
	}
}

func TestIsPushable_Plant_Success(t *testing.T) {
	plantTile := &references.Tile{Index: indexes.Plant}

	result := plantTile.IsPushable()

	if !result {
		t.Error("Expected Plant to be pushable")
	}
}

func TestActionPushSmallMap_Chair_PushForward_Success(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestActionPushSmallMap_Chair_PushForward_Blocked(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestActionPushSmallMap_Chair_PullSwap_Success(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestActionPushSmallMap_NotPushable_WontBudge(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestActionPushSmallMap_Cannon_RequiresHexFloor(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestActionPushLargeMap_NotImplemented(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestActionPushCombatMap_NotImplemented(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestObjectPresentAt_Stub(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestCheckTrap_Stub(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestActionPushSmallMap_Barrel_PushForward_Success(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}

func TestActionPushSmallMap_Barrel_PushForward_Blocked_Pull(t *testing.T) {
	t.Skip("Converting to use real game data - see TESTING.md")
}
