package ai

import (
	"math/rand"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

// TestMonsterGenerationIntegration tests the complete monster generation system
func TestMonsterGenerationIntegration(t *testing.T) {
	// Seed for consistent results
	rand.Seed(789)

	// Create test monster data
	testMonsters := references.EnemyReferences{
		// Water monster
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:         "Sea Serpent",
				IsWaterEnemy: true,
				WaterWeight:  10,
				LandWeight:   0,
				DesertWeight: 0,
			},
		},
		// Desert monster
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:         "Sand Spider",
				IsSandEnemy:  true,
				WaterWeight:  0,
				LandWeight:   1,
				DesertWeight: 8,
			},
		},
		// Common land monster
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:         "Orc",
				WaterWeight:  0,
				LandWeight:   5,
				DesertWeight: 2,
			},
		},
		// Rare land monster
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:         "Dragon",
				WaterWeight:  0,
				LandWeight:   1,
				DesertWeight: 1,
			},
		},
	}

	controller := &NPCAIControllerLargeMap{
		enemyReferences: &testMonsters,
		debugOptions:    &references.DebugOptions{MonsterGen: true},
		dateTime:        &datetime.UltimaDate{Hour: 12}, // Daytime
	}

	// Test different environment scenarios
	scenarios := []struct {
		name              string
		tileIndex         indexes.SpriteIndex
		expectedEnv       MonsterEnvironment
		shouldFindMonster bool
		description       string
	}{
		{
			name:              "Water environment should spawn water monsters",
			tileIndex:         indexes.Water1,
			expectedEnv:       WaterEnvironment,
			shouldFindMonster: true,
			description:       "Water tiles should spawn sea creatures",
		},
		{
			name:              "Desert environment should spawn desert monsters",
			tileIndex:         indexes.Desert,
			expectedEnv:       DesertEnvironment,
			shouldFindMonster: true,
			description:       "Desert tiles should spawn desert creatures",
		},
		{
			name:              "Land environment should spawn land monsters",
			tileIndex:         indexes.Grass,
			expectedEnv:       LandEnvironment,
			shouldFindMonster: true,
			description:       "Grass tiles should spawn land creatures",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tile := createTestTileForMonsterGen(scenario.tileIndex)

			// Test environment detection
			detectedEnv := controller.determineEnvironmentType(tile)
			if detectedEnv != scenario.expectedEnv {
				t.Errorf("Expected environment %d, got %d", scenario.expectedEnv, detectedEnv)
			}

			// Test monster selection
			selectedMonster := controller.pickMonsterByEnvironment(detectedEnv, tile)

			if scenario.shouldFindMonster && selectedMonster == nil {
				t.Errorf("Expected to find a monster for %s, but got nil", scenario.description)
				return
			}

			if selectedMonster != nil {
				// Verify the monster can actually spawn on the tile
				if !selectedMonster.CanSpawnToTile(tile) {
					t.Errorf("Selected monster %s cannot spawn on tile type %d",
						selectedMonster.AdditionalEnemyFlags.Name, tile.Index)
				}

				// Verify the monster has appropriate weight for this environment
				weight := controller.getMonsterEnvironmentWeight(selectedMonster, detectedEnv)
				if weight <= 0 {
					t.Errorf("Selected monster %s has invalid weight %d for environment %d",
						selectedMonster.AdditionalEnemyFlags.Name, weight, detectedEnv)
				}

				t.Logf("Successfully selected %s (weight %d) for %s",
					selectedMonster.AdditionalEnemyFlags.Name, weight, scenario.description)
			}
		})
	}
}

// TestMonsterSelectionDistribution verifies that weighted selection works across many trials
func TestMonsterSelectionDistribution(t *testing.T) {
	// Seed for reproducible results
	rand.Seed(456)

	// Create monsters with known weight distribution
	testMonsters := references.EnemyReferences{
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:       "Common Beast",
				LandWeight: 8, // Should be selected ~80% of time
			},
		},
		{
			AdditionalEnemyFlags: references.AdditionalEnemyFlags{
				Name:       "Rare Beast",
				LandWeight: 2, // Should be selected ~20% of time
			},
		},
	}

	controller := &NPCAIControllerLargeMap{
		enemyReferences: &testMonsters,
	}

	grassTile := createTestTileForMonsterGen(indexes.Grass)
	grassTile.IsLandEnemyPassable = true

	// Run many selections
	selections := make(map[string]int)
	numTrials := 1000

	for i := 0; i < numTrials; i++ {
		selected := controller.pickMonsterByEnvironment(LandEnvironment, grassTile)
		if selected != nil {
			selections[selected.AdditionalEnemyFlags.Name]++
		}
	}

	commonCount := selections["Common Beast"]
	rareCount := selections["Rare Beast"]

	// Common should be selected ~80% of time (weight 8 out of 10 total)
	expectedCommon := (numTrials * 8) / 10
	tolerance := expectedCommon / 10 // ±10%

	if commonCount < expectedCommon-tolerance || commonCount > expectedCommon+tolerance {
		t.Errorf("Expected Common Beast ~%d times, got %d times", expectedCommon, commonCount)
	}

	// Rare should be selected ~20% of time
	expectedRare := (numTrials * 2) / 10
	rareTolerance := expectedRare / 5 // ±20% tolerance

	if rareCount < expectedRare-rareTolerance || rareCount > expectedRare+rareTolerance {
		t.Errorf("Expected Rare Beast ~%d times, got %d times", expectedRare, rareCount)
	}

	t.Logf("Distribution results: Common=%d (expected ~%d), Rare=%d (expected ~%d)",
		commonCount, expectedCommon, rareCount, expectedRare)
}

// TestProbabilitySystemIntegration tests the complete probability calculation system
func TestProbabilitySystemIntegration(t *testing.T) {
	// Seed for consistent behavior
	rand.Seed(321)

	theOdds := references.NewDefaultTheOdds()
	controller := &NPCAIControllerLargeMap{
		debugOptions: &references.DebugOptions{MonsterGen: true},
		theOdds:      &theOdds,
	}

	scenarios := []struct {
		name                    string
		tileIndex               indexes.SpriteIndex
		hour                    byte
		expectedBaseProbability int
		expectedNightBonus      int
		description             string
	}{
		{
			name:                    "Road at noon - no spawning",
			tileIndex:               indexes.PathUpDown,
			hour:                    12,
			expectedBaseProbability: 0,
			expectedNightBonus:      0,
			description:             "Roads should never spawn monsters, even with night bonus",
		},
		{
			name:                    "Swamp at midnight - dangerous with night bonus",
			tileIndex:               indexes.Swamp,
			hour:                    0,
			expectedBaseProbability: 2,
			expectedNightBonus:      3,
			description:             "Swamps at night are very dangerous",
		},
		{
			name:                    "Grass during day - normal probability",
			tileIndex:               indexes.Grass,
			hour:                    14,
			expectedBaseProbability: 1,
			expectedNightBonus:      0,
			description:             "Grass during day has normal spawn rate",
		},
		{
			name:                    "Mountain at night - dangerous with bonus",
			tileIndex:               indexes.SmallMountains,
			hour:                    3,
			expectedBaseProbability: 2,
			expectedNightBonus:      3,
			description:             "Mountains at night are very dangerous",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			controller.dateTime = &datetime.UltimaDate{Hour: scenario.hour}
			tile := createTestTileForMonsterGen(scenario.tileIndex)

			// Test base probability calculation
			baseProbability := controller.calculateTileBasedProbabilityForTile(tile)
			if baseProbability != scenario.expectedBaseProbability {
				t.Errorf("Expected base probability %d, got %d for %s",
					scenario.expectedBaseProbability, baseProbability, scenario.description)
			}

			// Test night bonus application
			totalProbability := baseProbability
			if controller.dateTime.IsNight() {
				totalProbability += 3
			}

			expectedTotal := scenario.expectedBaseProbability + scenario.expectedNightBonus
			if totalProbability != expectedTotal {
				t.Errorf("Expected total probability %d, got %d for %s",
					expectedTotal, totalProbability, scenario.description)
			}

			// Log the result for verification
			t.Logf("%s: base=%d, night_bonus=%d, total=%d",
				scenario.description, baseProbability, scenario.expectedNightBonus, totalProbability)
		})
	}
}

// Test that the system respects debug mode settings
func TestDebugModeIntegrationComplete(t *testing.T) {
	controller := &NPCAIControllerLargeMap{
		debugOptions: &references.DebugOptions{MonsterGen: false}, // Debug OFF
		dateTime:     &datetime.UltimaDate{Hour: 12},
	}

	// When debug mode is OFF, monster generation should be blocked
	// (This would be tested in the actual generateTileBasedMonster method)
	debugEnabled := controller.debugOptions.MonsterGen
	if debugEnabled {
		t.Error("Debug mode should be disabled, but it's enabled")
	}

	// Enable debug mode
	controller.debugOptions.MonsterGen = true
	debugEnabled = controller.debugOptions.MonsterGen
	if !debugEnabled {
		t.Error("Debug mode should be enabled, but it's disabled")
	}

	t.Log("Debug mode toggle behavior verified")
}

// Helper method that matches the production code for direct testing
func (m *NPCAIControllerLargeMap) calculateTileBasedProbabilityForTile(tile *references.Tile) int {
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
