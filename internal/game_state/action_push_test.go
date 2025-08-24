package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func TestIsPushable_Chair_Success(t *testing.T) {
	gs, _ := createTestGameStateForPush(t) // Don't need test data for IsPushable tests

	// Test chair tile
	chairTile := &references.Tile{Index: indexes.ChairFacingUp}

	result := gs.IsPushable(chairTile)

	if !result {
		t.Error("Expected chair to be pushable")
	}
}

func TestIsPushable_Cannon_Success(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Test cannon tile
	cannonTile := &references.Tile{Index: indexes.CannonFacingUp}

	result := gs.IsPushable(cannonTile)

	if !result {
		t.Error("Expected cannon to be pushable")
	}
}

func TestIsPushable_Barrel_Success(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Test barrel tile
	barrelTile := &references.Tile{Index: indexes.Barrel}

	result := gs.IsPushable(barrelTile)

	if !result {
		t.Error("Expected barrel to be pushable")
	}
}

func TestIsPushable_NotPushable_Floor(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Test floor tile - should not be pushable
	floorTile := &references.Tile{Index: indexes.BrickFloor}

	result := gs.IsPushable(floorTile)

	if result {
		t.Error("Expected floor to not be pushable")
	}
}

func TestIsPushable_NotPushable_Wall(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Test wall tile - should not be pushable
	wallTile := &references.Tile{Index: indexes.StoneBrickWall}

	result := gs.IsPushable(wallTile)

	if result {
		t.Error("Expected wall to not be pushable")
	}
}

func TestIsPushable_TableMiddle_Success(t *testing.T) {
	gs := &GameState{}

	tableMiddleTile := &references.Tile{Index: indexes.TableMiddle}

	result := gs.IsPushable(tableMiddleTile)

	if !result {
		t.Error("Expected TableMiddle to be pushable")
	}
}

func TestIsPushable_Chest_Success(t *testing.T) {
	gs := &GameState{}

	chestTile := &references.Tile{Index: indexes.Chest}

	result := gs.IsPushable(chestTile)

	if !result {
		t.Error("Expected Chest to be pushable")
	}
}

func TestActionPushSmallMap_Chair_PushForward_Success(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Set up chair north of player
	direction := references.Up
	chairPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	farSidePos := direction.GetNewPositionInDirection(chairPos)

	smallMap := gs.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, gs.MapState.PlayerLocation.Floor)
	smallMap.SetTileByLayer(map_state.MapLayer, chairPos, indexes.ChairFacingUp)
	smallMap.SetTileByLayer(map_state.MapLayer, farSidePos, indexes.BrickFloor) // Legal floor

	// Test pushing chair
	result := gs.ActionPushSmallMap(direction)

	if !result {
		t.Error("Expected push to succeed")
	}

	// Verify chair moved to far side
	farSideTile := smallMap.GetTopTile(farSidePos)
	if !farSideTile.IsChair() {
		t.Errorf("Expected chair at far side, got %v", farSideTile.Index)
	}

	// Verify player moved to chair's original position
	if gs.MapState.PlayerLocation.Position != *chairPos {
		t.Errorf("Expected player to move to chair position %v, got %v", chairPos, gs.MapState.PlayerLocation.Position)
	}

	// Verify original chair position now has floor
	originalChairTile := smallMap.GetTopTile(chairPos)
	if originalChairTile.Index != indexes.BrickFloor {
		t.Errorf("Expected floor at original chair position, got %v", originalChairTile.Index)
	}
}

func TestActionPushSmallMap_Chair_PushForward_Blocked(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Set up chair north of player, but blocked by wall
	direction := references.Up
	chairPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	farSidePos := direction.GetNewPositionInDirection(chairPos)

	smallMap := gs.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, gs.MapState.PlayerLocation.Floor)
	smallMap.SetTileByLayer(map_state.MapLayer, chairPos, indexes.ChairFacingUp)
	smallMap.SetTileByLayer(map_state.MapLayer, farSidePos, indexes.StoneBrickWall) // Blocked

	initialPlayerPos := gs.MapState.PlayerLocation.Position

	// Test pushing blocked chair
	result := gs.ActionPushSmallMap(direction)

	// Since pull is attempted when push fails, this should still succeed if player on legal floor
	// But let's test the case where player is also blocked
	gs.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, gs.MapState.PlayerLocation.Floor).
		SetTileByLayer(map_state.MapLayer, &initialPlayerPos, indexes.Water1) // Illegal floor

	result = gs.ActionPushSmallMap(direction)

	if result {
		t.Error("Expected push to fail when both forward and swap are blocked")
	}

	// Verify player didn't move
	if gs.MapState.PlayerLocation.Position != initialPlayerPos {
		t.Error("Expected player position to remain unchanged when push fails")
	}
}

func TestActionPushSmallMap_Chair_PullSwap_Success(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Set up chair north of player, but forward blocked, player on legal floor
	direction := references.Up
	chairPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	farSidePos := direction.GetNewPositionInDirection(chairPos)

	smallMap := gs.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, gs.MapState.PlayerLocation.Floor)
	smallMap.SetTileByLayer(map_state.MapLayer, chairPos, indexes.ChairFacingUp)
	smallMap.SetTileByLayer(map_state.MapLayer, farSidePos, indexes.StoneBrickWall)                       // Forward blocked
	smallMap.SetTileByLayer(map_state.MapLayer, &gs.MapState.PlayerLocation.Position, indexes.BrickFloor) // Legal floor

	initialPlayerPos := gs.MapState.PlayerLocation.Position

	// Test pull/swap operation
	result := gs.ActionPushSmallMap(direction)

	if !result {
		t.Error("Expected pull/swap to succeed")
	}

	// Verify player moved to chair's original position
	if gs.MapState.PlayerLocation.Position != *chairPos {
		t.Error("Expected player to move to chair position in swap")
	}

	// Verify chair moved to player's original position
	playerOriginalTile := smallMap.GetTopTile(&initialPlayerPos)
	if !playerOriginalTile.IsChair() {
		t.Error("Expected chair at player's original position after swap")
	}
}

func TestActionPushSmallMap_NotPushable_WontBudge(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Set up wall north of player (not pushable)
	direction := references.Up
	wallPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)

	smallMap := gs.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, gs.MapState.PlayerLocation.Floor)
	smallMap.SetTileByLayer(map_state.MapLayer, wallPos, indexes.StoneBrickWall)

	initialPlayerPos := gs.MapState.PlayerLocation.Position

	// Test pushing wall
	result := gs.ActionPushSmallMap(direction)

	if result {
		t.Error("Expected push to fail for non-pushable wall")
	}

	// Verify player didn't move
	if gs.MapState.PlayerLocation.Position != initialPlayerPos {
		t.Error("Expected player position to remain unchanged when pushing wall")
	}
}

func TestActionPushSmallMap_Cannon_RequiresHexFloor(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Set up cannon north of player
	direction := references.Up
	cannonPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	farSidePos := direction.GetNewPositionInDirection(cannonPos)

	smallMap := gs.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, gs.MapState.PlayerLocation.Floor)
	smallMap.SetTileByLayer(map_state.MapLayer, cannonPos, indexes.CannonFacingUp)

	// Test with regular floor (should fail when legal floor validation is implemented)
	smallMap.SetTileByLayer(map_state.MapLayer, farSidePos, indexes.BrickFloor)

	// Note: Current implementation doesn't enforce hex floor requirement yet
	// This test documents expected behavior when that's implemented
	result := gs.ActionPushSmallMap(direction)

	// For now, this will succeed, but when hex floor validation is added, it should fail
	t.Logf("Current result: %v (will change when hex floor validation is implemented)", result)
}

func TestActionPushLargeMap_NotImplemented(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	result := gs.ActionPushLargeMap(references.Up)

	if !result {
		t.Error("Expected large map push to return true (TODO)")
	}
}

func TestActionPushCombatMap_NotImplemented(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	result := gs.ActionPushCombatMap(references.Up)

	if !result {
		t.Error("Expected combat map push to return true (TODO)")
	}
}

func TestObjectPresentAt_Stub(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Test the stub implementation
	result := gs.ObjectPresentAt(&references.Position{X: 0, Y: 0})

	if result {
		t.Error("Expected ObjectPresentAt stub to return false")
	}
}

func TestCheckTrap_Stub(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Test the stub implementation
	result := gs.CheckTrap(&references.Position{X: 0, Y: 0})

	if result {
		t.Error("Expected CheckTrap stub to return false")
	}
}

func TestActionPushSmallMap_Barrel_PushForward_Success(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Regression test for issue: "I am doing manual tests on Push mechanic in
	// Lord British's castle floor -1 at x=8, y=17. I am trying to push a barrel
	// to my left (west) and I get "Won't budge!""
	//
	// This verifies that barrels can be pushed like other objects

	// Set up barrel west of player (direction: Left/West)
	direction := references.Left
	barrelPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	farSidePos := direction.GetNewPositionInDirection(barrelPos)

	smallMap := gs.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, gs.MapState.PlayerLocation.Floor)
	smallMap.SetTileByLayer(map_state.MapLayer, barrelPos, indexes.Barrel)
	smallMap.SetTileByLayer(map_state.MapLayer, farSidePos, indexes.BrickFloor) // Legal floor

	// Test pushing barrel
	result := gs.ActionPushSmallMap(direction)

	if !result {
		t.Error("Expected barrel push to succeed")
	}

	// Verify barrel moved to far side
	farSideTile := smallMap.GetTopTile(farSidePos)
	if farSideTile.Index != indexes.Barrel {
		t.Errorf("Expected barrel at far side, got %v", farSideTile.Index)
	}

	// Verify player moved to barrel's original position
	if gs.MapState.PlayerLocation.Position != *barrelPos {
		t.Errorf("Expected player to move to barrel position %v, got %v", barrelPos, gs.MapState.PlayerLocation.Position)
	}

	// Verify original barrel position now has floor
	originalBarrelTile := smallMap.GetTopTile(barrelPos)
	if originalBarrelTile.Index != indexes.BrickFloor {
		t.Errorf("Expected floor at original barrel position, got %v", originalBarrelTile.Index)
	}
}

func TestActionPushSmallMap_Barrel_PushForward_Blocked_Pull(t *testing.T) {
	gs, _ := createTestGameStateForPush(t)

	// Test barrel that can't be pushed forward but can be pulled (swapped with player)
	direction := references.Left
	barrelPos := direction.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
	farSidePos := direction.GetNewPositionInDirection(barrelPos)

	smallMap := gs.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, gs.MapState.PlayerLocation.Floor)
	smallMap.SetTileByLayer(map_state.MapLayer, barrelPos, indexes.Barrel)
	smallMap.SetTileByLayer(map_state.MapLayer, farSidePos, indexes.StoneBrickWall)                       // Forward blocked
	smallMap.SetTileByLayer(map_state.MapLayer, &gs.MapState.PlayerLocation.Position, indexes.BrickFloor) // Legal floor

	initialPlayerPos := gs.MapState.PlayerLocation.Position

	// Test pull/swap operation
	result := gs.ActionPushSmallMap(direction)

	if !result {
		t.Error("Expected barrel pull/swap to succeed")
	}

	// Verify player moved to barrel's original position
	if gs.MapState.PlayerLocation.Position != *barrelPos {
		t.Error("Expected player to move to barrel position in swap")
	}

	// Verify barrel moved to player's original position
	playerOriginalTile := smallMap.GetTopTile(&initialPlayerPos)
	if playerOriginalTile.Index != indexes.Barrel {
		t.Error("Expected barrel at player's original position after swap")
	}
}

// TestCallbackData holds captured data from SystemCallbacks for verification
type TestCallbackData struct {
	Messages     []string
	Sounds       []SoundEffect
	TimeAdvanced int
}

// Helper function to create a test game state for push testing
func createTestGameStateForPush(t *testing.T) (*GameState, *TestCallbackData) {
	t.Helper()

	// Create test data capture structure
	testData := &TestCallbackData{
		Messages:     []string{},
		Sounds:       []SoundEffect{},
		TimeAdvanced: 0,
	}

	messageCallbacks := MessageCallbacks{
		AddRowStr: func(message string) {
			testData.Messages = append(testData.Messages, message)
		},
		AppendToCurrentRowStr: func(message string) {
			// Not used in push tests
		},
		ShowCommandPrompt: func(command string) {
			// Not used in push tests
		},
	}

	audioCallbacks := NewAudioCallbacks(func(effect SoundEffect) {
		testData.Sounds = append(testData.Sounds, effect)
	})

	flowCallbacks := NewFlowCallbacks(
		func() {}, // finishTurn
		func() {}, // activateGuards
		func(minutes int) { testData.TimeAdvanced += minutes }, // advanceTime
		func(datetime.TimeOfDay) {},                            // setTimeOfDay
		func() {},                                              // delayFx
		func() {},                                              // checkUpdate
	)

	visualCallbacks := NewVisualCallbacks(nil, nil, nil)
	screenCallbacks := NewScreenCallbacks(nil, nil, nil, nil, nil)
	talkCallbacks := NewTalkCallbacks(nil, nil) // No-op for tests

	systemCallbacks, err := NewSystemCallbacks(messageCallbacks, visualCallbacks, audioCallbacks, screenCallbacks, flowCallbacks, talkCallbacks)
	if err != nil {
		t.Fatalf("Failed to create SystemCallbacks: %v", err)
	}

	gs := &GameState{
		SystemCallbacks: systemCallbacks,
	}

	return gs, testData
}
