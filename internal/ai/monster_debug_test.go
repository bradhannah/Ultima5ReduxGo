package ai

import (
	"math/rand"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

// TestActualMonsterSpawnRate tests the real-world spawn rates you'd experience
func TestActualMonsterSpawnRate(t *testing.T) {
	// Seed for consistent results
	rand.Seed(12345)

	controller := &NPCAIControllerLargeMap{
		debugOptions: &references.DebugOptions{MonsterGen: true},
		dateTime:     &datetime.UltimaDate{Year: 139, Month: 4, Day: 7, Hour: 12}, // Your test date
	}

	// Set up default odds (1 in 32 chance)
	theOdds := references.NewDefaultTheOdds()
	controller.theOdds = &theOdds

	scenarios := []struct {
		name        string
		tileIndex   indexes.SpriteIndex
		hour        byte
		description string
	}{
		{
			name:        "Grass during day (your situation)",
			tileIndex:   indexes.Grass,
			hour:        12,
			description: "Normal grass during daytime",
		},
		{
			name:        "Grass at night",
			tileIndex:   indexes.Grass,
			hour:        2,
			description: "Grass with night bonus",
		},
		{
			name:        "Swamp during day",
			tileIndex:   indexes.Swamp,
			hour:        12,
			description: "Dangerous swamp terrain",
		},
		{
			name:        "Swamp at night",
			tileIndex:   indexes.Swamp,
			hour:        2,
			description: "Very dangerous swamp at night",
		},
		{
			name:        "Road (should never spawn)",
			tileIndex:   indexes.PathUpDown,
			hour:        12,
			description: "Road - should be 0%",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			controller.dateTime.Hour = scenario.hour

			// Test the probability gates
			numTrials := 10000
			firstGatePass := 0
			secondGatePass := 0

			for i := 0; i < numTrials; i++ {
				// First gate: shouldGenerateTileBasedMonster()
				if controller.shouldGenerateTileBasedMonster() {
					firstGatePass++

					// Second gate: tile-based probability
					tile := createTestTileForMonsterGen(scenario.tileIndex)
					probability := controller.calculateTileBasedProbabilityForTile(tile)

					// Apply night bonus
					if controller.dateTime.IsNight() {
						probability += 3
					}

					// Skip roads (probability 0)
					if probability > 0 {
						// Third gate: roll against tile probability
						if rand.Intn(probability) == 0 { // OneInXOdds equivalent
							secondGatePass++
						}
					}
				}
			}

			firstGatePercent := float64(firstGatePass) / float64(numTrials) * 100
			totalSpawnPercent := float64(secondGatePass) / float64(numTrials) * 100

			t.Logf("%s:", scenario.description)
			t.Logf("  First gate (1 in 32): %d/%d = %.2f%%", firstGatePass, numTrials, firstGatePercent)
			t.Logf("  Total spawn rate: %d/%d = %.3f%%", secondGatePass, numTrials, totalSpawnPercent)
			t.Logf("  Expected steps to see monster: ~%.0f", 100.0/totalSpawnPercent)
		})
	}
}

// TestDebugModeSpawnRate tests with debug mode adjustments
func TestDebugModeSpawnRate(t *testing.T) {
	// Seed for consistent results
	rand.Seed(789)

	controller := &NPCAIControllerLargeMap{
		debugOptions: &references.DebugOptions{MonsterGen: true},
		dateTime:     &datetime.UltimaDate{Year: 139, Month: 4, Day: 7, Hour: 12},
	}

	// Test with different debug settings
	debugSettings := []struct {
		name     string
		oneInX   int
		expected string
	}{
		{
			name:     "Default (1 in 32)",
			oneInX:   32,
			expected: "Very low spawn rate",
		},
		{
			name:     "Debug (1 in 4)",
			oneInX:   4,
			expected: "Much higher spawn rate",
		},
		{
			name:     "High Debug (1 in 2)",
			oneInX:   2,
			expected: "Very high spawn rate",
		},
	}

	for _, setting := range debugSettings {
		t.Run(setting.name, func(t *testing.T) {
			theOdds := references.NewDefaultTheOdds()
			theOdds.SetGenerateLargeMapMonster(setting.oneInX)
			controller.theOdds = &theOdds

			numTrials := 1000
			spawns := 0

			for i := 0; i < numTrials; i++ {
				if controller.shouldGenerateTileBasedMonster() {
					// Simulate grass tile (probability 1)
					if rand.Intn(1) == 0 { // Always pass since probability is 1
						spawns++
					}
				}
			}

			spawnPercent := float64(spawns) / float64(numTrials) * 100
			expectedSteps := 100.0 / spawnPercent

			t.Logf("%s: %d/%d = %.1f%% (every ~%.0f steps)",
				setting.expected, spawns, numTrials, spawnPercent, expectedSteps)
		})
	}
}

// TestRecommendedSpawnRateFix tests a proposed fix
func TestRecommendedSpawnRateFix(t *testing.T) {
	// Seed for consistent results
	rand.Seed(456)

	t.Log("RECOMMENDED FIX:")
	t.Log("The current system has two probability gates which compound to create very low spawn rates.")
	t.Log("Consider one of these approaches:")
	t.Log("")

	// Approach 1: Remove the first gate entirely
	t.Log("Approach 1: Remove shouldGenerateTileBasedMonster() gate entirely")
	t.Log("- Let tile-based probability be the only gate")
	t.Log("- Grass: 100% chance to try spawning")
	t.Log("- Swamp: 50% chance (1 in 2)")
	t.Log("- Swamp at night: 20% chance (1 in 5)")
	t.Log("")

	// Approach 2: Use first gate for frequency control, second for terrain
	t.Log("Approach 2: Use first gate as 'turn frequency', second as terrain modifier")
	t.Log("- First gate: Every 8-16 turns try to spawn")
	t.Log("- Second gate: Terrain determines success rate")
	t.Log("- Result: More predictable spawn timing")
	t.Log("")

	// Approach 3: Combine probabilities differently
	t.Log("Approach 3: Multiply probabilities instead of compounding")
	t.Log("- Base rate: 10% per turn")
	t.Log("- Terrain multiplier: Roads×0, Grass×1, Swamp×2, Night×1.5")
	t.Log("- Result: Grass=10%, Swamp day=20%, Swamp night=30%")

	// Test current vs fixed spawn rates
	currentRate := (1.0 / 32.0) * (1.0 / 1.0) * 100 // 3.125%
	fixedRate := (1.0 / 1.0) * 100                  // 100% if we remove first gate

	t.Logf("")
	t.Logf("Current grass spawn rate: %.3f%% (every ~%.0f steps)", currentRate, 100.0/currentRate)
	t.Logf("Fixed grass spawn rate: %.1f%% (every step if no other monsters)", fixedRate)
}

// Helper method is defined in monster_integration_test.go
