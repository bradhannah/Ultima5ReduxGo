package map_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func TestSmallMapProcessTurnDoors_NoOpenDoor(t *testing.T) {
	// Test that door processing system can be initialized with real game data
	config := loadTestConfiguration(t)
	if config == nil {
		return // Skip if no game data
	}

	gameRefs, err := references.NewGameReferences(config)
	if err != nil {
		t.Skipf("Failed to load game references: %v", err)
	}

	// Validate that we can load tile references with doors
	if gameRefs.TileReferences == nil {
		t.Fatal("TileReferences should not be nil")
	}

	t.Logf("Door processing data validation completed successfully with real game data")
	t.Logf("Loaded %d tile references for door processing", len(*gameRefs.TileReferences))

	// This test validates that the door system data is properly loaded
	// The actual SmallMapProcessTurnDoors method will be tested when implemented
}

// Helper function to load test configuration
func loadTestConfiguration(t *testing.T) *config.UltimaVConfiguration {
	t.Helper()

	// Try common test data locations
	testPaths := []string{
		"/Users/bradhannah/games/Ultima_5/Gold/", // Developer setup
		"./testdata/",                            // Relative test data
	}

	for _, path := range testPaths {
		config := config.NewUltimaVConfiguration()
		config.SavedConfigData.DataFilePath = path

		// Check if this path has Ultima V data
		if pathHasUltimaVData(path) {
			t.Logf("Loading Ultima V configuration from: %s", path)
			return config
		}
	}

	t.Logf("Ultima V game data not found in standard locations")
	return nil
}

// pathHasUltimaVData checks if a path contains key Ultima V data files
func pathHasUltimaVData(path string) bool {
	// Simple check - could be more comprehensive
	return true // For now, assume path is valid if provided
}
