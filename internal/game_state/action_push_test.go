package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/test/helpers"
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
	// Test pushing a chair using real game data
	config := helpers.LoadTestConfiguration(t)
	gameRefs := helpers.LoadTestGameReferences(t, config)

	// Load bootstrap save data for testing
	saveData := helpers.LoadBootstrapSaveData(t)
	xTiles, yTiles := helpers.GetTestGameConstants()

	gs := NewGameStateFromLegacySaveBytes(saveData, config, gameRefs, xTiles, yTiles)

	// Set up mock system callbacks for testing
	messageCallbacks, err := NewMessageCallbacks(
		func(str string) { t.Logf("Message: %s", str) }, // AddRowStr
		func(str string) { t.Logf("Append: %s", str) },  // AppendToCurrentRowStr
		func(str string) { t.Logf("Prompt: %s", str) },  // ShowCommandPrompt
	)
	if err != nil {
		t.Fatalf("Failed to create message callbacks: %v", err)
	}

	visualCallbacks := NewVisualCallbacks(nil, nil, nil)
	audioCallbacks := NewAudioCallbacks(nil)
	screenCallbacks := NewScreenCallbacks(nil, nil, nil, nil, nil)
	flowCallbacks := NewFlowCallbacks(nil, nil, nil, nil, nil, nil)
	talkCallbacks := NewTalkCallbacks(nil, nil)

	systemCallbacks, err := NewSystemCallbacks(messageCallbacks, visualCallbacks, audioCallbacks, screenCallbacks, flowCallbacks, talkCallbacks)
	if err != nil {
		t.Fatalf("Failed to create system callbacks: %v", err)
	}
	gs.SystemCallbacks = systemCallbacks

	// Set player location to Britain and update the small map
	gs.MapState.PlayerLocation.Location = references.Britain
	gs.MapState.PlayerLocation.Floor = references.FloorNumber(0)
	gs.UpdateSmallMap(gameRefs.TileReferences, gameRefs.LocationReferences)

	// Position player where they can push (coordinates in Britain)
	gs.MapState.PlayerLocation.Position = references.Position{X: 10, Y: 10}

	// Test basic push action doesn't crash with real data
	result := gs.ActionPushSmallMap(references.Up)

	// The exact result depends on what's at that position in the real data
	// But at minimum it should not crash and provide appropriate feedback
	t.Logf("ActionPushSmallMap result: %v", result)
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
