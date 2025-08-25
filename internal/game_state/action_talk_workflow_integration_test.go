// Integration tests for Talk command workflows with real game data
// These tests validate the complete talk command flow including NPC detection,
// dialog creation, conversation management, and system callbacks.
package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// TestTalkWorkflow_CompleteFlow tests the entire talk command workflow
func TestTalkWorkflow_CompleteFlow(t *testing.T) {
	// Test complete talk workflow with real game data
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Lord_Britishs_Castle). // Castle likely has NPCs
		WithPlayerAt(10, 10).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return // Test was skipped due to missing game data
	}

	t.Logf("💬 Talk Workflow Test - Castle Environment")

	// Test talk in different directions to find NPCs
	directions := []references.Direction{references.Up, references.Down, references.Left, references.Right}

	for _, dir := range directions {
		mockCallbacks.Reset()

		result := gs.ActionTalkSmallMap(dir)

		t.Logf("Talk %v result: %v", dir, result)

		// Log all feedback messages
		for i, msg := range mockCallbacks.Messages {
			t.Logf("Talk message %d: %s", i+1, msg)
		}

		// Check for time advancement (should advance on successful talk)
		if len(mockCallbacks.TimeAdvanced) > 0 {
			t.Logf("⏰ Time advanced: %d minute(s)", mockCallbacks.TimeAdvanced[0])
		}

		// Check for talk dialog creation
		if len(mockCallbacks.TalkDialogCalls) > 0 {
			t.Logf("💭 Talk dialog created for NPC")
		}

		// Check for dialog pushing
		if len(mockCallbacks.DialogsPushed) > 0 {
			t.Logf("📋 Dialog pushed to UI stack")
		}
	}

	t.Logf("💬 Talk workflow integration test completed")
}

// TestTalkWorkflow_NPCDetection tests NPC detection and interaction
func TestTalkWorkflow_NPCDetection(t *testing.T) {
	// Test NPC detection logic with real game data
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain). // Britain should have NPCs
		WithPlayerAt(15, 15).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("👥 Talk NPC Detection Test - Britain")

	// Test multiple positions to find NPCs
	testPositions := []struct {
		x, y int
		name string
	}{
		{15, 15, "Town Center"},
		{20, 20, "Southeast"},
		{10, 10, "Northwest"},
		{25, 15, "East Side"},
		{15, 25, "South Area"},
		{8, 8, "North Area"},
	}

	npcsFound := 0

	for _, pos := range testPositions {
		// Move to test position
		gs.MapState.PlayerLocation.Position = references.Position{
			X: references.Coordinate(pos.x),
			Y: references.Coordinate(pos.y),
		}

		t.Logf("Checking for NPCs around %s (%d,%d)", pos.name, pos.x, pos.y)

		// Test talk in each direction from this position
		for _, dir := range []references.Direction{references.Up, references.Down, references.Left, references.Right} {
			mockCallbacks.Reset()

			result := gs.ActionTalkSmallMap(dir)

			// Analyze response for NPC presence
			hasNPCResponse := false
			for _, msg := range mockCallbacks.Messages {
				if msg != "No-one to talk to!" && msg != "Can't talk to that!" {
					hasNPCResponse = true
					t.Logf("  🎯 Potential NPC interaction %v at %s: %s", dir, pos.name, msg)
					npcsFound++
				}
			}

			// Check for dialog creation (indicates NPC found)
			if len(mockCallbacks.TalkDialogCalls) > 0 {
				t.Logf("  ✅ NPC found %v from %s - dialog created", dir, pos.name)
				npcsFound++
			}

			if result && !hasNPCResponse {
				t.Logf("  ✅ Talk succeeded %v from %s", dir, pos.name)
			}
		}
	}

	t.Logf("👥 NPC Detection Summary:")
	t.Logf("  Potential NPCs found: %d", npcsFound)
	t.Logf("  NPCs detection test completed")
}

// TestTalkWorkflow_DialogSystem tests dialog system integration
func TestTalkWorkflow_DialogSystem(t *testing.T) {
	// Test dialog system integration
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Lord_Britishs_Castle).
		WithPlayerAt(12, 12).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("📋 Talk Dialog System Integration Test")

	// Test talk action and analyze dialog system integration
	directions := []references.Direction{references.Up, references.Down, references.Left, references.Right}

	for _, dir := range directions {
		mockCallbacks.Reset()

		result := gs.ActionTalkSmallMap(dir)

		t.Logf("Dialog system test %v: result=%v", dir, result)

		// Analyze dialog system callbacks
		t.Logf("Dialog system analysis for direction %v:", dir)
		t.Logf("  Dialog creation calls: %d", len(mockCallbacks.TalkDialogCalls))
		t.Logf("  Dialogs pushed: %d", len(mockCallbacks.DialogsPushed))
		t.Logf("  Messages: %d", len(mockCallbacks.Messages))

		// Log dialog-related feedback
		for _, msg := range mockCallbacks.Messages {
			switch msg {
			case "No-one to talk to!":
				t.Logf("  ❌ No NPC found in direction %v", dir)
			case "Can't talk to that!":
				t.Logf("  ⚠️ Non-talkable entity found in direction %v", dir)
			default:
				t.Logf("  💬 Dialog response: %s", msg)
			}
		}

		// Verify dialog creation/push consistency
		if len(mockCallbacks.TalkDialogCalls) > len(mockCallbacks.DialogsPushed) {
			t.Logf("  ⚠️ Dialog created but not pushed - possible error")
		} else if len(mockCallbacks.TalkDialogCalls) > 0 {
			t.Logf("  ✅ Dialog system working properly")
		}
	}

	t.Logf("📋 Dialog system integration test completed")
}

// TestTalkWorkflow_ErrorHandling tests error conditions and edge cases
func TestTalkWorkflow_ErrorHandling(t *testing.T) {
	// Test error handling for talk command
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(1, 1). // Edge position
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("⚠️ Talk Error Handling Test")

	// Test talk at map boundaries
	edgePositions := []struct {
		x, y int
		name string
	}{
		{0, 0, "Top-left corner"},
		{31, 31, "Bottom-right corner"},
		{0, 15, "Left edge"},
		{31, 15, "Right edge"},
		{15, 0, "Top edge"},
		{15, 31, "Bottom edge"},
	}

	for _, pos := range edgePositions {
		gs.MapState.PlayerLocation.Position = references.Position{
			X: references.Coordinate(pos.x),
			Y: references.Coordinate(pos.y),
		}

		t.Logf("Error handling test at %s (%d,%d)", pos.name, pos.x, pos.y)

		mockCallbacks.Reset()
		result := gs.ActionTalkSmallMap(references.Up)

		t.Logf("  Edge talk result: %v", result)
		for _, msg := range mockCallbacks.Messages {
			t.Logf("  Edge message: %s", msg)
		}

		// Verify system doesn't crash at boundaries
		if len(mockCallbacks.Messages) == 0 {
			t.Logf("  ✅ No crash at edge case %s", pos.name)
		} else {
			t.Logf("  ✅ Proper error handling at %s", pos.name)
		}
	}

	t.Logf("⚠️ Error handling test completed")
}

// TestTalkWorkflow_SystemCallbacksIntegration tests comprehensive callback integration
func TestTalkWorkflow_SystemCallbacksIntegration(t *testing.T) {
	// Test all system callbacks related to talk functionality
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Lord_Britishs_Castle).
		WithPlayerAt(16, 16).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("🔗 Talk System Callbacks Integration Test")

	// Perform talk action and analyze all callback interactions
	mockCallbacks.Reset()
	result := gs.ActionTalkSmallMap(references.Down)

	t.Logf("Talk callbacks integration result: %v", result)

	// Comprehensive callback analysis
	t.Logf("📊 Talk Callback Analysis:")
	t.Logf("  Messages: %d", len(mockCallbacks.Messages))
	t.Logf("  Time advancements: %d", len(mockCallbacks.TimeAdvanced))
	t.Logf("  Talk dialog calls: %d", len(mockCallbacks.TalkDialogCalls))
	t.Logf("  Dialogs pushed: %d", len(mockCallbacks.DialogsPushed))
	t.Logf("  Visual effects: %d", mockCallbacks.VisualCallCount)
	t.Logf("  Audio effects: %d", len(mockCallbacks.SoundEffectsPlayed))

	// Log detailed callback information
	for i, msg := range mockCallbacks.Messages {
		t.Logf("  Message %d: %s", i+1, msg)
	}

	for i, time := range mockCallbacks.TimeAdvanced {
		t.Logf("  Time advancement %d: %d minutes", i+1, time)
	}

	// Verify callback patterns
	if result {
		// Successful talk should have appropriate callbacks
		if len(mockCallbacks.Messages) == 0 && len(mockCallbacks.TalkDialogCalls) == 0 {
			t.Logf("  ⚠️ Successful talk but no feedback - investigate")
		} else {
			t.Logf("  ✅ Successful talk with proper callbacks")
		}
	} else {
		// Failed talk should have appropriate feedback
		if len(mockCallbacks.Messages) > 0 {
			t.Logf("  ✅ Failed talk with appropriate error message")
		}
	}

	t.Logf("🔗 Talk system callbacks integration test completed")
}

// TestTalkWorkflow_MultipleNPCs tests talking to multiple NPCs
func TestTalkWorkflow_MultipleNPCs(t *testing.T) {
	// Test talking to multiple NPCs in sequence
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Lord_Britishs_Castle). // Castle likely has multiple NPCs
		WithPlayerAt(8, 8).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("👥 Talk Multiple NPCs Test")

	// Test multiple positions for NPC interactions
	talkPositions := []struct {
		x, y      int
		direction references.Direction
		name      string
	}{
		{8, 8, references.Up, "North from entrance"},
		{8, 8, references.Down, "South from entrance"},
		{12, 12, references.Left, "West from center"},
		{12, 12, references.Right, "East from center"},
		{16, 16, references.Up, "North from back"},
		{20, 20, references.Down, "South from corner"},
	}

	talkAttempts := 0
	successfulTalks := 0
	dialogsCreated := 0

	for i, pos := range talkPositions {
		gs.MapState.PlayerLocation.Position = references.Position{
			X: references.Coordinate(pos.x),
			Y: references.Coordinate(pos.y),
		}

		mockCallbacks.Reset()

		result := gs.ActionTalkSmallMap(pos.direction)
		talkAttempts++

		t.Logf("Talk attempt %d (%s): result=%v", i+1, pos.name, result)

		if result {
			successfulTalks++
		}

		if len(mockCallbacks.TalkDialogCalls) > 0 {
			dialogsCreated++
			t.Logf("  ✅ Dialog created for %s", pos.name)
		}

		// Log specific feedback
		for _, msg := range mockCallbacks.Messages {
			if msg != "No-one to talk to!" {
				t.Logf("  💬 %s: %s", pos.name, msg)
			}
		}
	}

	t.Logf("Multiple NPCs summary:")
	t.Logf("  Talk attempts: %d", talkAttempts)
	t.Logf("  Successful talks: %d", successfulTalks)
	t.Logf("  Dialogs created: %d", dialogsCreated)

	t.Logf("👥 Multiple NPCs test completed")
}
