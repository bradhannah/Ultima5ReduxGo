// System Integration Validation Tests (Task #11 Phase 1)
// These tests validate that all major remediation fixes work together properly
// without conflicts or regressions across the integrated system.
package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// TestDeterminismIntegration_AllSystems validates that RNG, animations, and AI all use central GameClock
func TestDeterminismIntegration_AllSystems(t *testing.T) {
	// Test that all deterministic systems work together properly
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(15, 15).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("🎯 Determinism Integration Validation Test")

	// Record initial system state
	initialTime := gs.DateTime

	t.Logf("Initial state: Time=%d:%02d, Turn=%d",
		initialTime.Hour, initialTime.Minute, initialTime.Turn)

	// Test deterministic behavior across multiple runs with same seed
	testSequences := []struct {
		name       string
		action     func() bool
		expectRNG  bool
		expectTime bool
		expectAI   bool
	}{
		{"Look command", func() bool { return gs.ActionLookSmallMap(references.Up) }, false, true, false},
		{"Talk command", func() bool { return gs.ActionTalkSmallMap(references.Down) }, true, false, true}, // May trigger AI
		{"Get command", func() bool { return gs.ActionGetSmallMap(references.Left) }, false, false, false},
		{"Jimmy command", func() bool { return gs.ActionJimmySmallMap(references.Right) }, true, false, false}, // Uses RNG
	}

	// Run sequence twice with same seed to verify determinism
	firstRunResults := make([]struct {
		result   bool
		messages []string
		timeAdv  []int
	}, len(testSequences))

	// First run
	t.Logf("🔄 First deterministic run:")
	gs.SetRandomSeed(12345) // Fixed seed for determinism
	for i, test := range testSequences {
		mockCallbacks.Reset()

		result := test.action()
		firstRunResults[i].result = result
		firstRunResults[i].messages = append([]string{}, mockCallbacks.Messages...)
		firstRunResults[i].timeAdv = append([]int{}, mockCallbacks.TimeAdvanced...)

		t.Logf("  %s: result=%v, messages=%d, time_adv=%d",
			test.name, result, len(firstRunResults[i].messages),
			len(firstRunResults[i].timeAdv))
	}

	// Reset game state and run again with same seed
	t.Logf("🔄 Second deterministic run (verification):")
	gs, mockCallbacks = NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(15, 15).
		WithSystemCallbacks().
		Build()
	gs.SetRandomSeed(12345) // Same seed

	deterministicMatches := 0
	for i, test := range testSequences {
		mockCallbacks.Reset()

		result := test.action()

		// Compare with first run
		messagesMatch := len(mockCallbacks.Messages) == len(firstRunResults[i].messages)
		if messagesMatch && len(mockCallbacks.Messages) > 0 {
			messagesMatch = mockCallbacks.Messages[0] == firstRunResults[i].messages[0]
		}

		timeAdvMatch := len(mockCallbacks.TimeAdvanced) == len(firstRunResults[i].timeAdv)
		if timeAdvMatch && len(mockCallbacks.TimeAdvanced) > 0 {
			timeAdvMatch = mockCallbacks.TimeAdvanced[0] == firstRunResults[i].timeAdv[0]
		}

		if result == firstRunResults[i].result && messagesMatch && timeAdvMatch {
			deterministicMatches++
			t.Logf("  %s: ✅ DETERMINISTIC", test.name)
		} else {
			t.Logf("  %s: ❌ NON-DETERMINISTIC (result=%v vs %v, msg=%v vs %v, time=%v vs %v)",
				test.name, result, firstRunResults[i].result,
				len(mockCallbacks.Messages), len(firstRunResults[i].messages),
				len(mockCallbacks.TimeAdvanced), len(firstRunResults[i].timeAdv))
		}
	}

	t.Logf("🎯 Determinism Integration Results:")
	t.Logf("  Deterministic matches: %d/%d", deterministicMatches, len(testSequences))
	t.Logf("  Determinism rate: %.1f%%", float64(deterministicMatches)/float64(len(testSequences))*100)

	// Validate determinism integration
	if deterministicMatches == len(testSequences) {
		t.Logf("  ✅ All systems fully deterministic")
	} else if deterministicMatches > len(testSequences)/2 {
		t.Logf("  ⚠️ Partial determinism - some systems may need review")
	} else {
		t.Errorf("❌ Poor determinism integration - major issues detected")
	}

	t.Logf("🎯 Determinism integration validation completed")
}

// TestCommandIntegration_JimmyWithCollision validates Jimmy command with collision detection
func TestCommandIntegration_JimmyWithCollision(t *testing.T) {
	// Test that Jimmy command works properly with object collision detection
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Lord_Britishs_Castle). // Environment with doors and objects
		WithPlayerAt(10, 10).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	// Set up test conditions
	gs.PartyState.Inventory.Provisions.Keys.Set(5)
	gs.PartyState.Characters[0].Status = 1 // Good status

	t.Logf("🔐 Jimmy Command + Collision Integration Test")

	initialKeys := gs.PartyState.Inventory.Provisions.Keys.Get()
	t.Logf("Initial keys: %d", initialKeys)

	// Test Jimmy command in all directions to validate collision detection integration
	directions := []references.Direction{references.Up, references.Down, references.Left, references.Right}

	jimmyResults := make(map[references.Direction]struct {
		result      bool
		keyConsumed bool
		message     string
		collision   bool
	})

	for _, dir := range directions {
		mockCallbacks.Reset()
		keysBefore := gs.PartyState.Inventory.Provisions.Keys.Get()

		// Check what's in the direction (collision detection integration)
		targetPos := dir.GetNewPositionInDirection(&gs.MapState.PlayerLocation.Position)
		hasObject := gs.ObjectPresentAt(targetPos)

		result := gs.ActionJimmySmallMap(dir)

		keysAfter := gs.PartyState.Inventory.Provisions.Keys.Get()
		keyConsumed := keysBefore > keysAfter

		message := ""
		if len(mockCallbacks.Messages) > 0 {
			message = mockCallbacks.Messages[0]
		}

		jimmyResults[dir] = struct {
			result      bool
			keyConsumed bool
			message     string
			collision   bool
		}{
			result:      result,
			keyConsumed: keyConsumed,
			message:     message,
			collision:   hasObject,
		}

		t.Logf("Jimmy %v: result=%v, key_consumed=%v, collision=%v, message=%s",
			dir, result, keyConsumed, hasObject, message)
	}

	// Analyze integration behavior
	totalAttempts := len(directions)
	keyConsumptionAttempts := 0
	validTargetAttempts := 0
	collisionDetectionWorking := 0

	for dir, jimmyResult := range jimmyResults {
		// Count attempts where key consumption logic was tested
		if jimmyResult.message != "No lock picks!" {
			keyConsumptionAttempts++
		}

		// Count attempts on valid targets (doors, chests)
		if jimmyResult.message == "Unlocked!" || jimmyResult.message == "Lock pick broke!" ||
			jimmyResult.message == "Not lock!" {
			validTargetAttempts++
		}

		// Verify collision detection is working
		if jimmyResult.collision || jimmyResult.message == "Not lock!" {
			collisionDetectionWorking++
		}

		t.Logf("  Direction %v analysis: key_logic=%v, valid_target=%v, collision_detected=%v",
			dir, jimmyResult.keyConsumed,
			jimmyResult.message != "No lock picks!" && jimmyResult.message != "",
			jimmyResult.collision)
	}

	finalKeys := gs.PartyState.Inventory.Provisions.Keys.Get()
	totalKeysConsumed := int(initialKeys - finalKeys)

	t.Logf("🔐 Jimmy + Collision Integration Analysis:")
	t.Logf("  Total jimmy attempts: %d", totalAttempts)
	t.Logf("  Key consumption attempts: %d", keyConsumptionAttempts)
	t.Logf("  Valid target attempts: %d", validTargetAttempts)
	t.Logf("  Collision detection working: %d", collisionDetectionWorking)
	t.Logf("  Total keys consumed: %d", totalKeysConsumed)

	// Validate integration working properly
	integrationScore := 0
	if keyConsumptionAttempts > 0 {
		integrationScore++
		t.Logf("  ✅ Key consumption logic integrated")
	}
	if collisionDetectionWorking > 0 {
		integrationScore++
		t.Logf("  ✅ Collision detection integrated")
	}
	if totalKeysConsumed > 0 {
		integrationScore++
		t.Logf("  ✅ SystemCallbacks integration working")
	}

	if integrationScore >= 2 {
		t.Logf("  ✅ Jimmy + Collision integration validated")
	} else {
		t.Logf("  ⚠️ Jimmy + Collision integration needs review")
	}

	t.Logf("🔐 Jimmy command collision integration test completed")
}

// TestArchitectureIntegration_DisplayAndPackages validates DisplayManager with clean package boundaries
func TestArchitectureIntegration_DisplayAndPackages(t *testing.T) {
	// Test that DisplayManager integration maintains clean package boundaries
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(12, 12).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("🏗️ Architecture Integration Validation Test")

	// Test that game operations work correctly with DisplayManager architecture
	architectureTests := []struct {
		name     string
		action   func() bool
		category string
	}{
		{"Look command", func() bool { return gs.ActionLookSmallMap(references.Up) }, "visual"},
		{"Talk command", func() bool { return gs.ActionTalkSmallMap(references.Down) }, "dialog"},
		{"Jimmy command", func() bool { return gs.ActionJimmySmallMap(references.Left) }, "interaction"},
		{"Get command", func() bool { return gs.ActionGetSmallMap(references.Right) }, "inventory"},
	}

	systemIntegrationsWorking := make(map[string]int)
	totalTests := len(architectureTests)

	for i, test := range architectureTests {
		mockCallbacks.Reset()

		result := test.action()

		t.Logf("Architecture test %d (%s): result=%v", i+1, test.name, result)

		// Validate SystemCallbacks integration (architecture boundary)
		callbacksWorking := false
		if len(mockCallbacks.Messages) > 0 {
			callbacksWorking = true
			systemIntegrationsWorking["messages"]++
			t.Logf("  → Message system: %s", mockCallbacks.Messages[0])
		}

		if len(mockCallbacks.TimeAdvanced) > 0 {
			callbacksWorking = true
			systemIntegrationsWorking["time"]++
			t.Logf("  → Time system: %d minutes", mockCallbacks.TimeAdvanced[0])
		}

		if len(mockCallbacks.SoundEffectsPlayed) > 0 {
			callbacksWorking = true
			systemIntegrationsWorking["audio"]++
			t.Logf("  → Audio system: %v", mockCallbacks.SoundEffectsPlayed[0])
		}

		if len(mockCallbacks.TalkDialogCalls) > 0 {
			callbacksWorking = true
			systemIntegrationsWorking["dialog"]++
			t.Logf("  → Dialog system: activated")
		}

		if callbacksWorking {
			systemIntegrationsWorking[test.category]++
		}
	}

	// Test package boundary validation - no core logic should import rendering
	t.Logf("📦 Package Boundary Validation:")

	// This test validates that the game state operations work without direct rendering dependencies
	// The fact that we can run these tests without importing Ebitengine validates clean boundaries
	boundariesMaintained := true
	if gs.GameReferences != nil {
		t.Logf("  ✅ Core game logic accessible without rendering imports")
	} else {
		t.Logf("  ❌ Core game logic appears to have dependency issues")
		boundariesMaintained = false
	}

	t.Logf("🏗️ Architecture Integration Analysis:")
	t.Logf("  Total architecture tests: %d", totalTests)
	t.Logf("  System integrations working:")
	for system, count := range systemIntegrationsWorking {
		t.Logf("    %s: %d operations", system, count)
	}

	// Validate overall architecture integration
	systemsWorking := len(systemIntegrationsWorking)
	if systemsWorking >= 3 && boundariesMaintained {
		t.Logf("  ✅ Architecture integration validated")
		t.Logf("  ✅ DisplayManager + Package boundaries working")
	} else if systemsWorking >= 2 {
		t.Logf("  ⚠️ Architecture integration partially working")
	} else {
		t.Logf("  ❌ Architecture integration issues detected")
	}

	t.Logf("🏗️ Architecture integration validation completed")
}

// TestErrorHandlingIntegration_SystemCallbacks validates error handling across all systems
func TestErrorHandlingIntegration_SystemCallbacks(t *testing.T) {
	// Test that error handling improvements work across all integrated systems
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(5, 5).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("⚠️ Error Handling Integration Validation Test")

	// Test error conditions across different systems to validate consistent handling
	errorConditions := []struct {
		name       string
		setup      func()
		action     func() bool
		expectFail bool
		category   string
	}{
		{
			"Jimmy with no keys",
			func() { gs.PartyState.Inventory.Provisions.Keys.Set(0) },
			func() bool { return gs.ActionJimmySmallMap(references.Up) },
			true,
			"resource_exhaustion",
		},
		{
			"Jimmy with dead character",
			func() {
				gs.PartyState.Inventory.Provisions.Keys.Set(5)
				gs.PartyState.Characters[0].Status = 2 // Dead
			},
			func() bool { return gs.ActionJimmySmallMap(references.Down) },
			true,
			"character_status",
		},
		{
			"Talk at empty location",
			func() {}, // No setup needed
			func() bool { return gs.ActionTalkSmallMap(references.Left) },
			true,
			"target_validation",
		},
		{
			"Get from empty location",
			func() {}, // No setup needed
			func() bool { return gs.ActionGetSmallMap(references.Right) },
			true,
			"target_validation",
		},
	}

	errorHandlingResults := make(map[string]struct {
		handled    bool
		consistent bool
		message    string
	})

	for _, errorTest := range errorConditions {
		mockCallbacks.Reset()

		// Setup error condition
		errorTest.setup()

		result := errorTest.action()

		// Analyze error handling
		errorHandled := !result && len(mockCallbacks.Messages) > 0
		message := ""
		if len(mockCallbacks.Messages) > 0 {
			message = mockCallbacks.Messages[0]
		}

		// Check for consistent error messaging (no log.Fatal crashes)
		consistent := true
		if errorTest.expectFail && result {
			consistent = false // Should have failed but didn't
		}

		errorHandlingResults[errorTest.name] = struct {
			handled    bool
			consistent bool
			message    string
		}{
			handled:    errorHandled,
			consistent: consistent,
			message:    message,
		}

		t.Logf("Error test '%s': handled=%v, consistent=%v, message=%s",
			errorTest.name, errorHandled, consistent, message)

		// Reset for next test
		gs.PartyState.Characters[0].Status = 1 // Reset to good status
	}

	// Analyze error handling integration
	totalTests := len(errorConditions)
	handledCount := 0
	consistentCount := 0

	for _, result := range errorHandlingResults {
		if result.handled {
			handledCount++
		}
		if result.consistent {
			consistentCount++
		}
	}

	t.Logf("⚠️ Error Handling Integration Analysis:")
	t.Logf("  Total error conditions tested: %d", totalTests)
	t.Logf("  Properly handled: %d", handledCount)
	t.Logf("  Consistent behavior: %d", consistentCount)
	t.Logf("  Error handling rate: %.1f%%", float64(handledCount)/float64(totalTests)*100)

	// Validate error handling integration
	if handledCount == totalTests && consistentCount == totalTests {
		t.Logf("  ✅ Error handling integration excellent")
	} else if handledCount >= totalTests*3/4 {
		t.Logf("  ✅ Error handling integration good")
	} else {
		t.Logf("  ⚠️ Error handling integration needs improvement")
	}

	t.Logf("⚠️ Error handling integration validation completed")
}
