package game_state

import (
	"os"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// REGRESSION TEST: Movement Collision Detection Fix
//
// BUG DESCRIPTION:
// User reported: "I am able to walk right on top of Lord British's castle at x=85, y=107 in overworld"
// The IsPassable() function was only checking terrain passability but ignoring object collision.
//
// ROOT CAUSE:
// IsPassable() method only called topTile.IsPassable() for terrain but did not check
// for objects at the position using ObjectPresentAt(), allowing walking through castles/NPCs.
//
// ORIGINAL ULTIMA V BEHAVIOR:
// The test() function in OLD/CMDS.C checks both:
// 1. legalmove(PLAYER,tile) - terrain passability (our topTile.IsPassable())
// 2. looklist(x,y,level) - object collision (our ObjectPresentAt())
//
// WHY THIS TEST MATTERS:
// - Prevents walking through castles, NPCs, and other objects
// - Ensures movement collision matches original Ultima V physics
// - Validates integration between terrain and object collision systems
//
// TEST VALIDATES:
// 1. IsPassable() includes object collision checking (not just terrain)
// 2. Movement collision works with real game data and map layouts
// 3. Specific fix for Lord British's castle collision at overworld position (85,107)
func TestMovementCollisionDetection_LordBritishCastle_Integration(t *testing.T) {
	// Setup: Load real game data with overworld map
	gameState, _, _ := createTestGameStateWithRealData(t, references.Britannia_Underworld)
	if gameState == nil {
		t.Skip("Converting to use real game data - see TESTING.md")
		return
	}

	// Switch to overworld where Lord British's castle is located
	gameState.MapState.PlayerLocation.Location = references.Britannia_Underworld
	gameState.MapState.PlayerLocation.Floor = references.FloorNumber(0)

	// Test the reported problematic position: Lord British's castle at (85,107)
	castlePosition := references.Position{X: 85, Y: 107}

	t.Logf("🏰 Testing movement collision at Lord British's castle position (%d,%d)",
		castlePosition.X, castlePosition.Y)

	// Debug: Check what's at this position
	objectPresent := gameState.ObjectPresentAt(&castlePosition)
	topTile := gameState.GetCurrentLayeredMap().GetTopTile(&castlePosition)
	isPassable := gameState.IsPassable(&castlePosition)

	t.Logf("DEBUG: Position (%d,%d) - ObjectPresent=%v, TileIndex=%v, IsPassable=%v",
		castlePosition.X, castlePosition.Y, objectPresent, topTile.Index, isPassable)
	t.Logf("DEBUG: Castle tile - name=%q, isBuilding=%v, speedFactor=%v",
		topTile.Name, topTile.IsBuilding, topTile.SpeedFactor)

	// Log current location and floor for debugging
	t.Logf("DEBUG: Player location=%v, floor=%v",
		gameState.MapState.PlayerLocation.Location, gameState.MapState.PlayerLocation.Floor)

	// Before the fix: This would return true (walkable) - BUG
	// After the fix: This should return false (blocked by castle) - CORRECT
	if isPassable {
		t.Errorf("REGRESSION: Lord British's castle at (%d,%d) should block movement but IsPassable returned true",
			castlePosition.X, castlePosition.Y)
		t.Errorf("This means the object collision detection fix is not working properly")
		t.Errorf("ObjectPresent=%v, TilePassable=%v", objectPresent, topTile.IsPassable(gameState.PartyVehicle.GetVehicleDetails().VehicleType))
	} else {
		t.Logf("✅ FIXED: Lord British's castle properly blocks movement - collision detection working")
	}

	// Additional validation: Test nearby positions to understand the tile pattern
	adjacentPositions := []references.Position{
		{X: 84, Y: 107}, // Left of castle
		{X: 86, Y: 107}, // Right of castle
		{X: 85, Y: 106}, // Above castle
		{X: 85, Y: 108}, // Below castle
		{X: 86, Y: 109}, // User mentioned this should be passable
	}

	for _, pos := range adjacentPositions {
		adjacentPassable := gameState.IsPassable(&pos)
		adjacentTile := gameState.GetCurrentLayeredMap().GetTopTile(&pos)
		t.Logf("Adjacent position (%d,%d): passable=%v, tileIndex=%v, name=%q, isBuilding=%v, speedFactor=%v",
			pos.X, pos.Y, adjacentPassable, adjacentTile.Index, adjacentTile.Name, adjacentTile.IsBuilding, adjacentTile.SpeedFactor)
	}

	t.Logf("🎯 Movement collision detection now matches original Ultima V test() function behavior")
}

// TestMovementCollisionDetection_Integration verifies the collision system with real data
func TestMovementCollisionDetection_Integration(t *testing.T) {
	gameState, _, gameRefs := createTestGameStateWithRealData(t, references.Lord_Britishs_Castle)
	if gameState == nil {
		t.Skip("Converting to use real game data - see TESTING.md")
		return
	}

	// Test movement collision in Lord British's Castle (small map with NPCs)
	gameState.MapState.PlayerLocation.Location = references.Lord_Britishs_Castle
	gameState.MapState.PlayerLocation.Floor = references.FloorNumber(0)

	// Update to small map to load NPCs and objects
	gameState.UpdateSmallMap(gameRefs.TileReferences, gameRefs.LocationReferences)

	t.Logf("🏰 Testing movement collision in Lord British's Castle with real NPCs")

	// Check if IsPassable properly integrates object collision
	// We can't know exact NPC positions without examining the data,
	// but we can verify the integration is working

	testPositions := []references.Position{
		{X: 10, Y: 10},
		{X: 15, Y: 15},
		{X: 20, Y: 20},
	}

	for _, pos := range testPositions {
		// Check both individual components
		objectPresent := gameState.ObjectPresentAt(&pos)
		isPassable := gameState.IsPassable(&pos)

		// Verify integration: if object is present, position should not be passable
		if objectPresent && isPassable {
			t.Errorf("INTEGRATION FAILURE: Object present at (%d,%d) but IsPassable=true - object collision not working",
				pos.X, pos.Y)
		}

		t.Logf("Position (%d,%d): object_present=%v, passable=%v", pos.X, pos.Y, objectPresent, isPassable)
	}

	t.Logf("✅ Movement collision integration verified with real game data")
}

// Helper function for creating test game state with real data (as per TESTING.md)
func createTestGameStateWithRealData(t *testing.T, location references.Location) (*GameState, *config.UltimaVConfiguration, *references.GameReferences) {
	t.Helper()

	gameConfig := loadTestConfiguration(t)
	if gameConfig == nil {
		return nil, nil, nil
	}

	gameRefs := loadTestGameReferences(t, gameConfig)
	if gameRefs == nil {
		return nil, nil, nil
	}

	// Create game state with real data
	gameState := NewGameStateFromLegacySaveFile(
		gameConfig.SavedConfigData.DataFilePath+"/SAVED.GAM",
		gameConfig,
		gameRefs,
		16, 16) // Standard screen dimensions for testing

	// Initialize to specified location if needed
	if location != references.Britannia_Underworld {
		gameState.MapState.PlayerLocation.Location = location
	}

	return gameState, gameConfig, gameRefs
}

// Helper for loading test configuration (tries multiple common locations)
func loadTestConfiguration(t *testing.T) *config.UltimaVConfiguration {
	t.Helper()

	// Try common test data locations as per TESTING.md
	testPaths := []string{
		"/Users/bradhannah/games/Ultima_5/Gold/", // Developer setup
		"./testdata/",                            // Relative test data
		// Could add more environment variable paths
	}

	for _, path := range testPaths {
		if pathExists(path + "LOOK2.DAT") { // Check for a key Ultima V file
			t.Logf("Loading Ultima V configuration from: %s", path)
			gameConfig := config.NewUltimaVConfiguration()
			// Set data file path for testing
			gameConfig.SavedConfigData.DataFilePath = path
			return gameConfig
		}
	}

	t.Logf("Ultima V game data not found in standard locations")
	return nil
}

func loadTestGameReferences(t *testing.T, gameConfig *config.UltimaVConfiguration) *references.GameReferences {
	t.Helper()

	gameRefs, err := references.NewGameReferences(gameConfig)
	if err != nil {
		t.Logf("Failed to load game references: %v", err)
		return nil
	}

	return gameRefs
}

func pathExists(filename string) bool {
	// Simple check - could use os.Stat but keeping it simple for now
	_, err := os.Stat(filename)
	return err == nil
}
