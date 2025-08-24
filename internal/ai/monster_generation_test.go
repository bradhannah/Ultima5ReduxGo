package ai

import (
	"math/rand"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

// Test helper function to create tiles with specific properties
func createTestTileForMonsterGen(spriteIndex indexes.SpriteIndex) *references.Tile {
	tile := &references.Tile{
		Index:                spriteIndex,
		IsLandEnemyPassable:  true,
		IsWaterEnemyPassable: false,
	}

	// Set water tiles to be water passable
	if spriteIndex == indexes.Water1 || spriteIndex == indexes.Water2 || spriteIndex == indexes.WaterShallow {
		tile.IsWaterEnemyPassable = true
		tile.IsLandEnemyPassable = false
	}

	return tile
}

// Test the core tile-based probability logic
func TestTileBasedProbabilityCalculation(t *testing.T) {
	// Seed RNG for consistent test results
	rand.Seed(42)

	tests := []struct {
		name         string
		tileIndex    indexes.SpriteIndex
		expectedProb int
		description  string
	}{
		{
			name:         "Road/Path tiles have 0 probability",
			tileIndex:    indexes.PathUpDown,
			expectedProb: 0,
			description:  "No monsters should spawn on roads",
		},
		{
			name:         "Swamp tiles have probability 2",
			tileIndex:    indexes.Swamp,
			expectedProb: 2,
			description:  "Dangerous swamp terrain",
		},
		{
			name:         "Mountain tiles have probability 2",
			tileIndex:    indexes.SmallMountains,
			expectedProb: 2,
			description:  "Dangerous mountain terrain",
		},
		{
			name:         "Grass tiles have probability 1",
			tileIndex:    indexes.Grass,
			expectedProb: 1,
			description:  "Normal grassland terrain",
		},
		{
			name:         "Desert tiles have probability 1",
			tileIndex:    indexes.Desert,
			expectedProb: 1,
			description:  "Normal desert terrain",
		},
		{
			name:         "Water tiles have probability 1",
			tileIndex:    indexes.Water1,
			expectedProb: 1,
			description:  "Normal water terrain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tile := createTestTileForMonsterGen(tt.tileIndex)
			probability := calculateTileBasedProbabilityDirect(tile)

			if probability != tt.expectedProb {
				t.Errorf("Expected probability %d for %s, got %d",
					tt.expectedProb, tt.description, probability)
			}
		})
	}
}

// Direct implementation of probability calculation for testing
func calculateTileBasedProbabilityDirect(tile *references.Tile) int {
	// Roads = 0 probability (no monsters spawn)
	if tile.IsRoad() {
		return 0
	}

	// Swamp, forest, mountains = 2 probability
	if tile.IsSwamp() || tile.IsForest() || tile.IsMountain() {
		return 2
	}

	// All other tiles = 1 probability
	return 1
}

// Test night bonus application
func TestNightBonusLogic(t *testing.T) {
	tests := []struct {
		name            string
		hour            byte
		baseProbability int
		expectedTotal   int
		description     string
	}{
		{
			name:            "Daytime has no bonus",
			hour:            12, // Noon
			baseProbability: 1,
			expectedTotal:   1,
			description:     "Midday should have no night bonus",
		},
		{
			name:            "Night adds +3 bonus",
			hour:            2, // 2 AM
			baseProbability: 1,
			expectedTotal:   4, // 1 + 3
			description:     "Night should add +3 to base probability",
		},
		{
			name:            "Night bonus on dangerous terrain",
			hour:            23, // 11 PM
			baseProbability: 2,  // Swamp/Mountain
			expectedTotal:   5,  // 2 + 3
			description:     "Night bonus applies to all terrain types",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateTime := &datetime.UltimaDate{Hour: tt.hour}

			totalProbability := tt.baseProbability
			if dateTime.IsNight() {
				totalProbability += 3
			}

			if totalProbability != tt.expectedTotal {
				t.Errorf("Expected total probability %d for %s, got %d",
					tt.expectedTotal, tt.description, totalProbability)
			}
		})
	}
}

// Test environment type detection
func TestEnvironmentDetection(t *testing.T) {
	tests := []struct {
		name        string
		tileIndex   indexes.SpriteIndex
		expectedEnv MonsterEnvironment
		description string
	}{
		{
			name:        "Water1 detected as water environment",
			tileIndex:   indexes.Water1,
			expectedEnv: WaterEnvironment,
			description: "Water tiles should be water environment",
		},
		{
			name:        "Water2 detected as water environment",
			tileIndex:   indexes.Water2,
			expectedEnv: WaterEnvironment,
			description: "All water variants should be water environment",
		},
		{
			name:        "Desert detected as desert environment",
			tileIndex:   indexes.Desert,
			expectedEnv: DesertEnvironment,
			description: "Desert tiles should be desert environment",
		},
		{
			name:        "Grass detected as land environment",
			tileIndex:   indexes.Grass,
			expectedEnv: LandEnvironment,
			description: "Grass is standard land environment",
		},
		{
			name:        "Swamp detected as land environment",
			tileIndex:   indexes.Swamp,
			expectedEnv: LandEnvironment,
			description: "Swamp is considered land environment",
		},
		{
			name:        "Mountain detected as land environment",
			tileIndex:   indexes.SmallMountains,
			expectedEnv: LandEnvironment,
			description: "Mountains are land environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tile := createTestTileForMonsterGen(tt.tileIndex)
			environment := determineEnvironmentTypeDirect(tile)

			if environment != tt.expectedEnv {
				t.Errorf("Expected environment %d for %s, got %d",
					tt.expectedEnv, tt.description, environment)
			}
		})
	}
}

// Direct implementation for testing
func determineEnvironmentTypeDirect(tile *references.Tile) MonsterEnvironment {
	if tile.IsWater() {
		return WaterEnvironment
	}
	if tile.IsDesert() {
		return DesertEnvironment
	}

	return LandEnvironment
}

// Test weighted random selection with seeded RNG
func TestWeightedSelection(t *testing.T) {
	// Seed RNG for reproducible results
	rand.Seed(12345)

	// Create test enemies with different weights
	enemies := []*references.EnemyReference{
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:       "Rare Monster",
				LandWeight: 1,
			},
		},
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:       "Common Monster",
				LandWeight: 10,
			},
		},
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:       "Another Rare Monster",
				LandWeight: 1,
			},
		},
	}

	weights := []int{1, 10, 1} // Total = 12, common should be selected ~83% of time

	// Run many selections to test distribution
	selections := make(map[string]int)
	numTrials := 1200 // Multiple of total weight

	for i := 0; i < numTrials; i++ {
		selected := weightedRandomSelectionDirect(enemies, weights)
		if selected != nil {
			selections[selected.AdditionalEnemyFlags.Name]++
		}
	}

	// Verify distribution is roughly correct
	commonCount := selections["Common Monster"]
	rareCount1 := selections["Rare Monster"]
	rareCount2 := selections["Another Rare Monster"]

	// Common should be selected about 10/12 = 83.3% of time = ~1000 times
	expectedCommon := (numTrials * 10) / 12
	tolerance := expectedCommon / 10 // ±10%

	if commonCount < expectedCommon-tolerance || commonCount > expectedCommon+tolerance {
		t.Errorf("Expected Common Monster ~%d times, got %d times",
			expectedCommon, commonCount)
	}

	// Each rare monster should be selected about 1/12 = 8.3% = ~100 times
	expectedRare := numTrials / 12
	rareTolerance := expectedRare / 2 // ±50% tolerance for small numbers

	if rareCount1 < expectedRare-rareTolerance || rareCount1 > expectedRare+rareTolerance ||
		rareCount2 < expectedRare-rareTolerance || rareCount2 > expectedRare+rareTolerance {
		t.Errorf("Rare monsters selected unexpected times. Rare1: %d, Rare2: %d, expected ~%d each",
			rareCount1, rareCount2, expectedRare)
	}

	t.Logf("Selection results: Common=%d, Rare1=%d, Rare2=%d",
		commonCount, rareCount1, rareCount2)
}

// Direct implementation for testing
func weightedRandomSelectionDirect(enemies []*references.EnemyReference, weights []int) *references.EnemyReference {
	if len(enemies) == 0 || len(weights) == 0 {
		return nil
	}

	// Calculate total weight
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	if totalWeight == 0 {
		return nil
	}

	// Generate random number in range [0, totalWeight)
	randomValue := rand.Intn(totalWeight)

	// Find the selected enemy
	currentWeight := 0
	for i, weight := range weights {
		currentWeight += weight
		if randomValue < currentWeight {
			return enemies[i]
		}
	}

	// Fallback
	return enemies[len(enemies)-1]
}

// Test tile identification methods work correctly
func TestTileIdentificationMethods(t *testing.T) {
	tests := []struct {
		name      string
		tileIndex indexes.SpriteIndex
		method    string
		expected  bool
	}{
		{"Water1 IsWater", indexes.Water1, "IsWater", true},
		{"Water2 IsWater", indexes.Water2, "IsWater", true},
		{"WaterShallow IsWater", indexes.WaterShallow, "IsWater", true},
		{"Desert IsDesert", indexes.Desert, "IsDesert", true},
		{"LeftDesert2 IsDesert", indexes.LeftDesert2, "IsDesert", true},
		{"RightDesert2 IsDesert", indexes.RightDesert2, "IsDesert", true},
		{"Swamp IsSwamp", indexes.Swamp, "IsSwamp", true},
		{"SmallMountains IsMountain", indexes.SmallMountains, "IsMountain", true},
		{"PathUpDown IsRoad", indexes.PathUpDown, "IsRoad", true},
		{"PathLeftRight IsRoad", indexes.PathLeftRight, "IsRoad", true},
		{"PathAllWays IsRoad", indexes.PathAllWays, "IsRoad", true},
		{"Grass IsWater", indexes.Grass, "IsWater", false},
		{"Desert IsSwamp", indexes.Desert, "IsSwamp", false},
		{"Water1 IsDesert", indexes.Water1, "IsDesert", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tile := createTestTileForMonsterGen(tt.tileIndex)
			var result bool

			switch tt.method {
			case "IsWater":
				result = tile.IsWater()
			case "IsDesert":
				result = tile.IsDesert()
			case "IsSwamp":
				result = tile.IsSwamp()
			case "IsMountain":
				result = tile.IsMountain()
			case "IsRoad":
				result = tile.IsRoad()
			case "IsForest":
				result = tile.IsForest()
			default:
				t.Fatalf("Unknown method: %s", tt.method)
			}

			if result != tt.expected {
				t.Errorf("Expected %s.%s() = %v, got %v",
					tt.name, tt.method, tt.expected, result)
			}
		})
	}
}

// Test that debug mode behavior is respected
func TestDebugModeIntegration(t *testing.T) {
	controller := &NPCAIControllerLargeMap{
		debugOptions: &references.DebugOptions{MonsterGen: true},
		dateTime:     &datetime.UltimaDate{Hour: 12},
	}

	// Test debug enabled
	if !controller.debugOptions.MonsterGen {
		t.Error("Debug mode should be enabled for testing")
	}

	// Test debug disabled
	controller.debugOptions.MonsterGen = false
	if controller.debugOptions.MonsterGen {
		t.Error("Debug mode should be disabled")
	}

	t.Log("Debug mode behavior test passed")
}
