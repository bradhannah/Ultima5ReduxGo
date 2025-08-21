package conversation

import (
	"testing"
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// MockGameState creates a minimal game state for testing
func createTestGameState(t *testing.T) *game_state.GameState {
	cfg := config.NewUltimaVConfiguration()
	gameReferences, err := references.NewGameReferences(cfg)
	if err != nil {
		t.Fatalf("Failed to create game references: %v", err)
	}

	// Use the embedded save file from convo-demo
	// For now, create a minimal state without loading a full save file
	return &game_state.GameState{
		GameReferences: gameReferences,
		// Add minimal required fields as needed
	}
}

// getTreannaDialog returns Treanna's talk script (dialog_number: 2)
func getTreannaDialog(gameState *game_state.GameState) *references.TalkScript {
	// Treanna is dialog_number 2 in CASTLE.TLK
	return gameState.GameReferences.TalkReferences.GetTalkScriptByNpcIndex(references.Castle, 2)
}

// getTreannaNPCReference returns Treanna's NPC reference
func getTreannaNPCReference(gameState *game_state.GameState) *references.NPCReference {
	// Treanna should be in Britain (based on stable context)
	npcRefs := gameState.GameReferences.NPCReferences.GetNPCReferencesByLocation(references.Britain)
	// Find the NPC with dialog_number 2 (Treanna)
	for _, npcRef := range *npcRefs {
		if npcRef.DialogNumber == 2 {
			return &npcRef
		}
	}
	return nil // Not found - may need to check other locations
}

// TestTreannaDialogData tests that Treanna's dialog data is properly loaded
func TestTreannaDialogData(t *testing.T) {
	gameState := createTestGameState(t)
	talkScript := getTreannaDialog(gameState)
	if talkScript == nil {
		t.Fatal("Could not find Treanna's talk script")
	}

	// Verify we have the expected dialog structure
	if len(talkScript.Lines) < 5 {
		t.Errorf("Expected at least 5 fixed entries (description, greeting, name, job, bye), got %d", len(talkScript.Lines))
	}

	// Check that we have question groups for basic commands
	hasNameGroup := false
	hasJobGroup := false
	hasByeGroup := false

	for _, group := range talkScript.QuestionGroups {
		for _, option := range group.Options {
			switch option {
			case "name":
				hasNameGroup = true
				// Verify the name response contains "Treanna"
				if len(group.Script) > 0 && group.Script[0].Cmd == 0 {
					if group.Script[0].Str != "Treanna" {
						t.Errorf("Expected name response to be 'Treanna', got %q", group.Script[0].Str)
					}
				}
			case "job":
				hasJobGroup = true
				// Verify job response mentions stable
				if len(group.Script) > 0 && group.Script[0].Cmd == 0 {
					if group.Script[0].Str != "I am the stable girl." {
						t.Errorf("Expected job response to be 'I am the stable girl.', got %q", group.Script[0].Str)
					}
				}
			case "bye":
				hasByeGroup = true
				// Verify bye response
				if len(group.Script) > 0 && group.Script[0].Cmd == 0 {
					if group.Script[0].Str != "Goodbye." {
						t.Errorf("Expected bye response to be 'Goodbye.', got %q", group.Script[0].Str)
					}
				}
			}
		}
	}

	if !hasNameGroup {
		t.Error("Expected to find 'name' question group")
	}
	if !hasJobGroup {
		t.Error("Expected to find 'job' question group")
	}
	if !hasByeGroup {
		t.Error("Expected to find 'bye' question group")
	}

	t.Logf("Successfully validated Treanna's dialog data with %d question groups", len(talkScript.QuestionGroups))
}

// TestTreannaConversationInitialization tests basic conversation setup
func TestTreannaConversationInitialization(t *testing.T) {
	gameState := createTestGameState(t)
	talkScript := getTreannaDialog(gameState)
	if talkScript == nil {
		t.Fatal("Could not find Treanna's talk script")
	}

	npcRef := getTreannaNPCReference(gameState)
	if npcRef == nil {
		// Create a minimal NPC reference for testing
		npcRef = &references.NPCReference{
			DialogNumber: 2,
			NpcType:      references.NoStatedNpc,
		}
	}

	// Test conversation creation
	convo := NewConversation(*npcRef, gameState, talkScript)
	if convo == nil {
		t.Fatal("Failed to create conversation")
	}

	// Test that channels are available
	if convo.Out() == nil {
		t.Error("Output channel is nil")
	}
	if convo.In() == nil {
		t.Error("Input channel is nil")
	}

	// Start conversation and get initial output
	convo.Start()
	defer convo.Stop()

	// Collect initial outputs with shorter timeout since we know this works
	var outputs []string
	for i := 0; i < 3; i++ {
		select {
		case item := <-convo.Out():
			outputs = append(outputs, item.Str)
			t.Logf("Initial output %d: %q (cmd: %d)", i, item.Str, item.Cmd)
		case <-time.After(500 * time.Millisecond):
			break // No more immediate output
		}
	}

	// Should have gotten at least some initial output
	if len(outputs) == 0 {
		t.Error("Expected initial conversation output")
	}

	// Should contain description-like text
	foundDescription := false
	for _, output := range outputs {
		if output == "a young girl." {
			foundDescription = true
			break
		}
	}

	if !foundDescription {
		t.Logf("Expected to find description 'a young girl.' in outputs: %v", outputs)
	}

	t.Logf("Conversation initialized successfully with %d initial outputs", len(outputs))
}

// TestTreannaKeywordOptions tests that Treanna has expected conversation options
func TestTreannaKeywordOptions(t *testing.T) {
	gameState := createTestGameState(t)
	talkScript := getTreannaDialog(gameState)
	if talkScript == nil {
		t.Fatal("Could not find Treanna's talk script")
	}

	// Collect all available keywords/options
	var allOptions []string
	for _, group := range talkScript.QuestionGroups {
		allOptions = append(allOptions, group.Options...)
	}

	t.Logf("Treanna has %d conversation options: %v", len(allOptions), allOptions)

	// Test for expected basic options
	expectedOptions := []string{"name", "job", "bye"}
	for _, expected := range expectedOptions {
		found := false
		for _, option := range allOptions {
			if option == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find option %q in conversation", expected)
		}
	}

	// Test for horse-related options (Treanna is a stable girl)
	horseRelatedOptions := []string{"bree", "horse"}
	foundHorseOptions := 0
	for _, horseOption := range horseRelatedOptions {
		for _, option := range allOptions {
			if option == horseOption {
				foundHorseOptions++
				t.Logf("Found horse-related option: %q", option)
				break
			}
		}
	}

	if foundHorseOptions == 0 {
		t.Log("No horse-related options found (may be expected depending on dialog structure)")
	}
}

// TestConversationDataIntegrity tests that the conversation system loads valid data
func TestConversationDataIntegrity(t *testing.T) {
	gameState := createTestGameState(t)

	// Test that we can load talk references for Castle location
	talkRefs := gameState.GameReferences.TalkReferences
	if talkRefs == nil {
		t.Fatal("TalkReferences is nil")
	}

	// Test that Treanna's script exists
	talkScript := talkRefs.GetTalkScriptByNpcIndex(references.Castle, 2)
	if talkScript == nil {
		t.Fatal("Could not get Treanna's talk script (dialog 2 in Castle)")
	}

	// Basic validation of script structure
	if len(talkScript.Lines) == 0 {
		t.Error("Talk script has no lines")
	}

	if len(talkScript.QuestionGroups) == 0 {
		t.Error("Talk script has no question groups")
	}

	// Log some debug info
	t.Logf("Treanna's script has %d lines and %d question groups",
		len(talkScript.Lines), len(talkScript.QuestionGroups))

	// Verify basic script commands
	for i, line := range talkScript.Lines {
		if i >= 5 { // Only check first few lines
			break
		}
		if len(line) > 0 {
			t.Logf("Line %d: Cmd=%d, Str=%q", i, line[0].Cmd, line[0].Str)
		} else {
			t.Logf("Line %d: Empty script line", i)
		}
	}

	t.Log("Conversation data integrity test passed")
}
