// Test helpers for integration testing with real Ultima V data
// Following patterns established in TESTING.md
package game_state

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// IntegrationTestBuilder provides a fluent API for creating integration tests
// with real game data, comprehensive state setup, and proper cleanup.
type IntegrationTestBuilder struct {
	t               *testing.T
	gameConfig      *config.UltimaVConfiguration
	gameRefs        *references.GameReferences
	gameState       *GameState
	location        references.Location
	floor           references.FloorNumber
	playerX         int
	playerY         int
	systemCallbacks *MockSystemCallbacks
	skipIfNoData    bool
}

// Global cache for game references to avoid reloading in each test
var (
	cachedGameRefs *references.GameReferences
	cachedConfig   *config.UltimaVConfiguration
	cacheOnce      sync.Once
	cacheError     error
)

// NewIntegrationTestBuilder creates a new builder for integration tests
func NewIntegrationTestBuilder(t *testing.T) *IntegrationTestBuilder {
	return &IntegrationTestBuilder{
		t:            t,
		location:     references.Britannia_Underworld,
		floor:        references.FloorNumber(0),
		playerX:      50,
		playerY:      50,
		skipIfNoData: true,
	}
}

// WithLocation sets the initial game location
func (b *IntegrationTestBuilder) WithLocation(location references.Location) *IntegrationTestBuilder {
	b.location = location
	return b
}

// WithFloor sets the floor number within the location
func (b *IntegrationTestBuilder) WithFloor(floor references.FloorNumber) *IntegrationTestBuilder {
	b.floor = floor
	return b
}

// WithPlayerAt sets the initial player position
func (b *IntegrationTestBuilder) WithPlayerAt(x, y int) *IntegrationTestBuilder {
	b.playerX = x
	b.playerY = y
	return b
}

// WithSystemCallbacks sets up mock system callbacks for testing
func (b *IntegrationTestBuilder) WithSystemCallbacks() *IntegrationTestBuilder {
	b.systemCallbacks = NewMockSystemCallbacks(b.t)
	return b
}

// FailIfNoGameData changes the default behavior to fail instead of skip when game data isn't found
func (b *IntegrationTestBuilder) FailIfNoGameData() *IntegrationTestBuilder {
	b.skipIfNoData = false
	return b
}

// Build creates the fully configured GameState for testing
func (b *IntegrationTestBuilder) Build() (*GameState, *MockSystemCallbacks) {
	b.t.Helper()

	// Load game configuration and references (cached)
	if err := b.loadGameData(); err != nil {
		if b.skipIfNoData {
			b.t.Skip("Ultima V game data not found - set ULTIMA5_DATA_PATH or place in standard locations")
		} else {
			b.t.Fatalf("Failed to load Ultima V game data: %v", err)
		}
		return nil, nil
	}

	// Create game state from real save file
	saveFilePath := filepath.Join(b.gameConfig.SavedConfigData.DataFilePath, "SAVED.GAM")
	b.gameState = NewGameStateFromLegacySaveFile(
		saveFilePath,
		b.gameConfig,
		b.gameRefs,
		19, 13) // Standard screen dimensions for testing

	// Set up initial position and location
	b.gameState.MapState.PlayerLocation.Location = b.location
	b.gameState.MapState.PlayerLocation.Floor = b.floor
	b.gameState.MapState.PlayerLocation.Position = references.Position{
		X: references.Coordinate(b.playerX),
		Y: references.Coordinate(b.playerY),
	}

	// Set up mock system callbacks if requested
	if b.systemCallbacks != nil {
		b.gameState.SystemCallbacks = b.systemCallbacks.ToSystemCallbacks()
	}

	// Update to specified map type and load NPCs/objects
	if b.location != references.Britannia_Underworld {
		b.gameState.UpdateSmallMap(b.gameRefs.TileReferences, b.gameRefs.LocationReferences)
	}

	b.t.Logf("Integration test initialized: location=%v, floor=%v, player=(%d,%d)",
		b.location, b.floor, b.playerX, b.playerY)

	return b.gameState, b.systemCallbacks
}

// loadGameData loads game configuration and references with caching
func (b *IntegrationTestBuilder) loadGameData() error {
	cacheOnce.Do(func() {
		// Try multiple common locations for Ultima V data
		testPaths := []string{
			"/Users/bradhannah/games/Ultima_5/Gold/", // Developer setup
			"./testdata/",                            // Relative test data
			os.Getenv("ULTIMA5_DATA_PATH"),           // Environment variable
		}

		for _, path := range testPaths {
			if path == "" {
				continue
			}

			// Check for key Ultima V files to verify path
			if !b.pathHasUltimaVData(path) {
				continue
			}

			// Create configuration
			cachedConfig = config.NewUltimaVConfiguration()
			cachedConfig.SavedConfigData.DataFilePath = path

			// Load game references
			var err error
			cachedGameRefs, err = references.NewGameReferences(cachedConfig)
			if err != nil {
				cacheError = err
				continue
			}

			b.t.Logf("Loaded Ultima V game data from: %s", path)
			return
		}

		cacheError = os.ErrNotExist
	})

	if cacheError != nil {
		return cacheError
	}

	b.gameConfig = cachedConfig
	b.gameRefs = cachedGameRefs
	return nil
}

// pathHasUltimaVData checks if a path contains key Ultima V data files
func (b *IntegrationTestBuilder) pathHasUltimaVData(path string) bool {
	requiredFiles := []string{
		"LOOK2.DAT",
		"SAVED.GAM",
		"CASTLE.TLK",
		"MISCMAPS.DAT",
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(filepath.Join(path, file)); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// QuickIntegrationTest creates a basic integration test setup with minimal configuration
func QuickIntegrationTest(t *testing.T, location references.Location) (*GameState, *MockSystemCallbacks) {
	return NewIntegrationTestBuilder(t).
		WithLocation(location).
		WithSystemCallbacks().
		Build()
}

// QuickSmallMapTest creates an integration test specifically for small map testing
func QuickSmallMapTest(t *testing.T, location references.Location, x, y int) (*GameState, *MockSystemCallbacks) {
	return NewIntegrationTestBuilder(t).
		WithLocation(location).
		WithPlayerAt(x, y).
		WithSystemCallbacks().
		Build()
}
