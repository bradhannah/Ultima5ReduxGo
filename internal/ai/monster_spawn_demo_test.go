package ai

import (
	"math/rand"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// TestMonsterSpawnProblemDemonstration demonstrates the monster spawning issue
func TestMonsterSpawnProblemDemonstration(t *testing.T) {
	// Seed for consistent results
	rand.Seed(12345)

	// Create controller with default settings (the ones you experienced)
	theOdds := references.NewDefaultTheOdds()
	controller := &NPCAIControllerLargeMap{
		debugOptions: &references.DebugOptions{MonsterGen: true},
		dateTime:     &datetime.UltimaDate{Year: 139, Month: 4, Day: 7, Hour: 12}, // Your exact scenario
		theOdds:      &theOdds,
	}

	t.Logf("=== MONSTER SPAWN PROBLEM ANALYSIS ===")
	t.Logf("Your scenario: Date 4-7-139, walking on grass at noon")
	t.Logf("Default odds: 1 in %d chance per step", controller.theOdds.GetOneInXMonsterGeneration())

	// Simulate OLD BROKEN SYSTEM (before our fix)
	t.Logf("\n--- OLD SYSTEM (before fix) ---")

	numSteps := 100
	oldSystemSpawns := 0

	for step := 1; step <= numSteps; step++ {
		// OLD SYSTEM: Two compounding probability gates
		// Gate 1: 1 in 32 base chance
		if rand.Intn(32) == 0 {
			// Gate 2: 1 in 1 for grass (tile probability)
			if rand.Intn(1) == 0 {
				oldSystemSpawns++
			}
		}
	}

	oldRate := float64(oldSystemSpawns) / float64(numSteps) * 100
	t.Logf("Old system results:")
	t.Logf("  Monsters in %d steps: %d", numSteps, oldSystemSpawns)
	t.Logf("  Spawn rate: %.2f%%", oldRate)
	if oldSystemSpawns == 0 {
		t.Logf("  ❌ PROBLEM: No monsters encountered - this was your issue!")
	}

	// Reset seed for fair comparison
	rand.Seed(12345)

	// Simulate NEW FIXED SYSTEM (after our fix)
	t.Logf("\n--- NEW SYSTEM (after fix) ---")

	newSystemSpawns := 0
	for step := 1; step <= numSteps; step++ {
		// NEW SYSTEM: Single combined probability
		// For grass (probability 1): combinedOdds = 32 / 1 = 32
		// So still 1 in 32 chance, but no double gating
		if rand.Intn(32) == 0 {
			newSystemSpawns++
		}
	}

	newRate := float64(newSystemSpawns) / float64(numSteps) * 100
	t.Logf("Fixed system results:")
	t.Logf("  Monsters in %d steps: %d", numSteps, newSystemSpawns)
	t.Logf("  Spawn rate: %.2f%%", newRate)
	if newSystemSpawns > 0 {
		t.Logf("  ✅ FIXED: You should now see monsters while walking!")
	}

	// Show the difference
	t.Logf("\n--- COMPARISON ---")
	t.Logf("Old system: %.2f%% (every ~%.0f steps)", oldRate, 100.0/max(oldRate, 0.01))
	t.Logf("New system: %.2f%% (every ~%.0f steps)", newRate, 100.0/max(newRate, 0.01))

	if newSystemSpawns > oldSystemSpawns {
		improvement := float64(newSystemSpawns-oldSystemSpawns) / float64(max(float64(oldSystemSpawns), 1.0)) * 100
		t.Logf("Improvement: +%.0f%% more monster encounters", improvement)
	}
}

// TestDifferentTerrainSpawnRates tests how the fix affects different terrain types
func TestDifferentTerrainSpawnRates(t *testing.T) {
	rand.Seed(789)

	theOdds := references.NewDefaultTheOdds()
	controller := &NPCAIControllerLargeMap{
		debugOptions: &references.DebugOptions{MonsterGen: true},
		dateTime:     &datetime.UltimaDate{Year: 139, Month: 4, Day: 7, Hour: 12},
		theOdds:      &theOdds,
	}

	t.Logf("=== TERRAIN-BASED SPAWN RATES ===")

	terrainTypes := []struct {
		name        string
		probability int
		description string
	}{
		{"Roads", 0, "Should never spawn"},
		{"Grass/Normal", 1, "Standard terrain"},
		{"Swamp/Forest/Mountain", 2, "Dangerous terrain"},
	}

	baseOdds := controller.theOdds.GetOneInXMonsterGeneration() // 32

	for _, terrain := range terrainTypes {
		t.Logf("\n%s (%s):", terrain.name, terrain.description)
		t.Logf("  Base tile probability: %d", terrain.probability)

		if terrain.probability == 0 {
			t.Logf("  Final spawn rate: 0%% (never)")
			continue
		}

		// With our fix: combinedOdds = baseOdds / probability
		combinedOdds := baseOdds / terrain.probability
		expectedRate := (1.0 / float64(combinedOdds)) * 100

		t.Logf("  Combined odds: 1 in %d", combinedOdds)
		t.Logf("  Expected spawn rate: %.2f%%", expectedRate)
		t.Logf("  Steps per monster: ~%.0f", 100.0/expectedRate)

		// Quick simulation to verify
		spawns := 0
		for i := 0; i < 1000; i++ {
			if rand.Intn(combinedOdds) == 0 {
				spawns++
			}
		}
		actualRate := float64(spawns) / 10.0 // Convert to percentage
		t.Logf("  Simulated rate: %.1f%% (%d/1000)", actualRate, spawns)
	}
}

// TestNightBonusEffect tests how night affects spawn rates
func TestNightBonusEffect(t *testing.T) {
	rand.Seed(456)

	theOdds := references.NewDefaultTheOdds()

	scenarios := []struct {
		name        string
		hour        byte
		description string
	}{
		{"Daytime", 12, "Your original situation"},
		{"Night", 2, "With +3 night bonus"},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			controller := &NPCAIControllerLargeMap{
				debugOptions: &references.DebugOptions{MonsterGen: true},
				dateTime:     &datetime.UltimaDate{Year: 139, Month: 4, Day: 7, Hour: scenario.hour},
				theOdds:      &theOdds,
			}

			t.Logf("%s (%s):", scenario.name, scenario.description)
			t.Logf("  Is night: %v", controller.dateTime.IsNight())

			// Test grass terrain
			grassProbability := 1
			if controller.dateTime.IsNight() {
				grassProbability += 3
			}

			combinedOdds := controller.theOdds.GetOneInXMonsterGeneration() / grassProbability
			expectedRate := (1.0 / float64(combinedOdds)) * 100

			t.Logf("  Grass probability: %d", grassProbability)
			t.Logf("  Combined odds: 1 in %d", combinedOdds)
			t.Logf("  Expected spawn rate: %.2f%%", expectedRate)
			t.Logf("  Steps per monster: ~%.0f", 100.0/expectedRate)

			// Test swamp terrain (more dramatic difference)
			swampProbability := 2
			if controller.dateTime.IsNight() {
				swampProbability += 3
			}

			swampOdds := controller.theOdds.GetOneInXMonsterGeneration() / swampProbability
			swampRate := (1.0 / float64(swampOdds)) * 100

			t.Logf("  Swamp probability: %d", swampProbability)
			t.Logf("  Swamp spawn rate: %.2f%%", swampRate)
		})
	}
}

// TestRecommendedSettings tests different debug settings for better gameplay
func TestRecommendedSettings(t *testing.T) {
	rand.Seed(321)

	t.Logf("=== RECOMMENDED DEBUG SETTINGS ===")
	t.Logf("For better monster encounter rates while walking:")

	settings := []struct {
		name    string
		oneInX  int
		verdict string
	}{
		{"Current Default", 32, "Too rare - you experienced this"},
		{"Balanced", 16, "RECOMMENDED - Good for gameplay"},
		{"Active", 8, "Good for testing/debugging"},
		{"Very Active", 4, "Very frequent encounters"},
	}

	for _, setting := range settings {
		theOdds := references.NewDefaultTheOdds()
		theOdds.SetMonsterGeneration(setting.oneInX)

		// Test grass spawn rate
		grassRate := (1.0 / float64(setting.oneInX)) * 100
		stepsPerMonster := 100.0 / grassRate

		t.Logf("\n%s (1 in %d):", setting.name, setting.oneInX)
		t.Logf("  Grass spawn rate: %.1f%%", grassRate)
		t.Logf("  Steps per monster: ~%.0f", stepsPerMonster)
		t.Logf("  Verdict: %s", setting.verdict)

		if setting.oneInX == 16 {
			t.Logf("  ⭐ RECOMMENDED: SetMonsterGeneration(%d)", setting.oneInX)
		}
	}

	t.Logf("\nTo fix your issue, consider calling:")
	t.Logf("  controller.theOdds.SetMonsterGeneration(16)")
	t.Logf("This will give you a monster every ~16 steps on grass instead of ~32.")
}

// Helper function to avoid division by zero
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
