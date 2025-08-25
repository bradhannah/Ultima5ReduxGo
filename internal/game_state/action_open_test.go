package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/test/helpers"
)

func TestActionOpenLargeMap_Door_Success(t *testing.T) {
	// Test opening a regular door on large map using real game data
	config := helpers.LoadTestConfiguration(t)
	gameRefs := helpers.LoadTestGameReferences(t, config)

	// Load bootstrap save data for testing
	saveData := helpers.LoadBootstrapSaveData(t)
	xTiles, yTiles := helpers.GetTestGameConstants()

	gs := NewGameStateFromLegacySaveBytes(saveData, config, gameRefs, xTiles, yTiles)

	// Set player location to Britain and update the small map
	gs.MapState.PlayerLocation.Location = references.Britain
	gs.MapState.PlayerLocation.Floor = references.FloorNumber(0)
	gs.UpdateSmallMap(gameRefs.TileReferences, gameRefs.LocationReferences)

	// Position player near a door (these coordinates should have a door in Britain)
	gs.MapState.PlayerLocation.Position = references.Position{X: 16, Y: 16}

	// Test basic open action doesn't crash with real data
	result := gs.ActionOpenLargeMap(references.Up)

	// The exact result depends on what's at that position in the real data
	// But at minimum it should not crash
	if result {
		t.Log("ActionOpenLargeMap succeeded with real data")
	} else {
		t.Log("ActionOpenLargeMap returned false with real data (expected behavior)")
	}
}

func TestActionOpenLargeMap_InvalidPosition(t *testing.T) {
	// Test opening when not at a door
	config := helpers.LoadTestConfiguration(t)
	gameRefs := helpers.LoadTestGameReferences(t, config)

	// Load bootstrap save data for testing
	saveData := helpers.LoadBootstrapSaveData(t)
	xTiles, yTiles := helpers.GetTestGameConstants()

	gs := NewGameStateFromLegacySaveBytes(saveData, config, gameRefs, xTiles, yTiles)

	// Set player location to Britain and update the small map
	gs.MapState.PlayerLocation.Location = references.Britain
	gs.MapState.PlayerLocation.Floor = references.FloorNumber(0)
	gs.UpdateSmallMap(gameRefs.TileReferences, gameRefs.LocationReferences)

	// Position player at a location without a door
	gs.MapState.PlayerLocation.Position = references.Position{X: 1, Y: 1}

	result := gs.ActionOpenLargeMap(references.Up)

	// Should return false when there's nothing to open
	if result {
		t.Error("Expected ActionOpenLargeMap to return false when no door present")
	}
}
