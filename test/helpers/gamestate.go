package helpers

import (
	"os"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// Helper for loading test configuration (assumes Ultima V data is available)
func LoadTestConfiguration(t *testing.T) *config.UltimaVConfiguration {
	t.Helper()

	// Try common test data locations
	testPaths := []string{
		"/Users/bradhannah/games/Ultima_5/Gold/", // Developer setup
		"./OLD/",                                 // Relative OLD directory
		os.Getenv("ULTIMA5_DATA_PATH"),           // Environment variable
	}

	for _, path := range testPaths {
		if path != "" && PathExists(path) {
			// Set the data file path for the configuration
			os.Setenv("ULTIMA5_DATA_PATH", path) // This might be used by the config
			config := config.NewUltimaVConfiguration()
			// Verify the configuration loaded properly by checking if we can access data files
			if PathExists(config.GetFileWithFullPath("LOOK2.DAT")) {
				return config
			}
		}
	}

	t.Skip("Ultima V game data not found - ensure Ultima V data files are in ./OLD/ or set ULTIMA5_DATA_PATH environment variable")
	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func LoadTestGameReferences(t *testing.T, config *config.UltimaVConfiguration) *references.GameReferences {
	t.Helper()

	gameRefs, err := references.NewGameReferences(config)
	if err != nil {
		t.Fatalf("Failed to load game references: %v", err)
	}

	return gameRefs
}

// GetTestGameConstants returns the constants needed for GameState creation
func GetTestGameConstants() (int, int) {
	const (
		xTilesVisibleOnGameScreen = 19
		yTilesVisibleOnGameScreen = 13
	)
	return xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen
}

// LoadBootstrapSaveData loads the bootstrap save file for testing
func LoadBootstrapSaveData(t *testing.T) []byte {
	t.Helper()

	// Try to load the bootstrap save file
	saveFilePath := "/Users/bradhannah/GitHub/Ultima5ReduxGo/testdata/britain2_SAVED.GAM"
	if PathExists(saveFilePath) {
		data, err := os.ReadFile(saveFilePath)
		if err != nil {
			t.Fatalf("Failed to load bootstrap save file: %v", err)
		}
		return data
	}

	t.Skip("Bootstrap save file not found at testdata/britain2_SAVED.GAM")
	return nil
}

// Note: CreateMockSystemCallbacks moved to individual test files to avoid import cycle
