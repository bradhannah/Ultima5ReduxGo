package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// Integration test for exit vehicle functionality with real data
func TestActionExit_IntegrationWithRealData(t *testing.T) {
	// Test exit vehicle functionality using real game data
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(10, 10).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return // Test was skipped due to missing game data
	}

	// Test exit actions in different contexts
	testScenarios := []struct {
		name     string
		location references.Location
	}{
		{"Britain_Exit", references.Britain},
		{"Castle_Exit", references.Lord_Britishs_Castle},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			mockCallbacks.Reset()
			gs.MapState.PlayerLocation.Location = scenario.location
			gs.UpdateSmallMap(gs.GameReferences.TileReferences, gs.GameReferences.LocationReferences)

			// Test exit action
			result := gs.ActionExit()

			t.Logf("Exit result for %s: %v", scenario.name, result)

			// Log any messages or feedback
			for i, msg := range mockCallbacks.Messages {
				t.Logf("Exit message %d: %s", i+1, msg)
			}

			// Log time advancement (if any)
			if len(mockCallbacks.TimeAdvanced) > 0 {
				totalTime := 0
				for _, time := range mockCallbacks.TimeAdvanced {
					totalTime += time
				}
				t.Logf("Time advanced: %d minutes", totalTime)
			}
		})
	}
}

// Integration test for vehicle exit functionality
func TestActionExit_VehicleExitIntegration(t *testing.T) {
	// Test exiting vehicles with real game data
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(15, 15).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return // Test was skipped due to missing game data
	}

	// Test different vehicle states
	// Note: This tests the exit logic regardless of actual vehicle presence
	// as it validates the system's response to exit attempts

	result := gs.ActionExit()
	t.Logf("Vehicle exit test result: %v", result)

	// Validate system response
	if len(mockCallbacks.Messages) > 0 {
		t.Logf("Exit feedback: %s", mockCallbacks.Messages[len(mockCallbacks.Messages)-1])
	}

	// Check for appropriate time advancement
	mockCallbacks.AssertTimeAdvanced(0) // Should be 0 if no exit occurred

	t.Logf("Vehicle exit integration test completed with real data")
}
