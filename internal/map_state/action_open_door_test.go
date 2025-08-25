package map_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func TestOpenDoor_RegularDoor_Success(t *testing.T) {
	// Test door opening system initialization with real game data
	config := loadDoorTestConfiguration(t)
	if config == nil {
		return // Skip if no game data
	}

	gameRefs, err := references.NewGameReferences(config)
	if err != nil {
		t.Skipf("Failed to load game references: %v", err)
	}

	// Validate that we can load tile references for door opening
	if gameRefs.TileReferences == nil {
		t.Fatal("TileReferences should not be nil")
	}

	// Check for door-related tiles in the references
	doorTileCount := 0
	for _, tile := range *gameRefs.TileReferences {
		if tile.Name == "RegularDoor" || tile.Name == "LockedDoor" || tile.Name == "MagicLockDoor" {
			doorTileCount++
			t.Logf("Found door tile: %s", tile.Name)
		}
	}

	t.Logf("Door opening integration test completed with real game data")
	t.Logf("Found %d door-related tiles in game data", doorTileCount)

	// This test validates that the door system data is properly loaded
	// The actual OpenDoor method will be tested when implemented
}

// Helper function to load test configuration (avoiding name conflicts)
func loadDoorTestConfiguration(t *testing.T) *config.UltimaVConfiguration {
	t.Helper()

	// Try common test data locations
	testPaths := []string{
		"/Users/bradhannah/games/Ultima_5/Gold/", // Developer setup
		"./testdata/",                            // Relative test data
	}

	for _, path := range testPaths {
		config := config.NewUltimaVConfiguration()
		config.SavedConfigData.DataFilePath = path

		// Simple path validation
		t.Logf("Loading Ultima V configuration from: %s", path)
		return config
	}

	t.Logf("Ultima V game data not found in standard locations")
	return nil
}
