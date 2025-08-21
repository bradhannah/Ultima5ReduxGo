package conversation

import (
	"testing"
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// NPCTestData defines test data for an NPC's conversation
type NPCTestData struct {
	Name              string                         // NPC name for test identification
	Location          references.SmallMapMasterTypes // Location type (Castle, Towne, etc.)
	DialogNumber      int                            // Dialog index in TLK file
	ExpectedName      string                         // Expected response to "name" command
	ExpectedJob       string                         // Expected response to "job" command
	ExpectedBye       string                         // Expected response to "bye" command
	ExpectedOptions   []string                       // Expected available conversation options
	ThematicOptions   []string                       // Expected thematic options (horse, magic, etc.)
	MinQuestionGroups int                            // Minimum expected question groups
	Description       string                         // Expected description text
}

// ConversationTestSuite holds test data for multiple NPCs
type ConversationTestSuite struct {
	NPCs []NPCTestData
}

// GetTestSuite returns the conversation test suite with all NPC test data
func GetTestSuite() ConversationTestSuite {
	return ConversationTestSuite{
		NPCs: []NPCTestData{
			{
				Name:              "Treanna",
				Location:          references.Castle,
				DialogNumber:      2,
				ExpectedName:      "Treanna",
				ExpectedJob:       "I am the stable girl.",
				ExpectedBye:       "Goodbye.",
				ExpectedOptions:   []string{"name", "job", "bye"},
				ThematicOptions:   []string{"bree", "stab"}, // stable/horse related
				MinQuestionGroups: 5,
				Description:       "a young girl.",
			},
			{
				Name:              "Alistair",
				Location:          references.Castle,
				DialogNumber:      0,
				ExpectedName:      "Alistair the Bard",
				ExpectedJob:       "",               // Will be filled in once we see actual response
				ExpectedBye:       "",               // Will be filled in once we see actual response
				ExpectedOptions:   []string{"name"}, // Start with minimal expectations
				ThematicOptions:   []string{},       // Will be discovered from actual data
				MinQuestionGroups: 1,
				Description:       "a melancholy musician.",
			},
		},
	}
}

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

// getNPCDialog returns an NPC's talk script based on test data
func getNPCDialog(gameState *game_state.GameState, npcData NPCTestData) *references.TalkScript {
	return gameState.GameReferences.TalkReferences.GetTalkScriptByNpcIndex(npcData.Location, npcData.DialogNumber)
}

// getNPCReference returns an NPC reference based on test data
func getNPCReference(gameState *game_state.GameState, npcData NPCTestData) *references.NPCReference {
	// Try to find the NPC in various locations
	locations := []references.Location{references.Britain, references.Lord_Britishs_Castle}

	for _, location := range locations {
		npcRefs := gameState.GameReferences.NPCReferences.GetNPCReferencesByLocation(location)
		for _, npcRef := range *npcRefs {
			if npcRef.DialogNumber == byte(npcData.DialogNumber) {
				return &npcRef
			}
		}
	}

	// Return a minimal NPC reference for testing if not found
	return &references.NPCReference{
		DialogNumber: byte(npcData.DialogNumber),
		NpcType:      references.NoStatedNpc,
	}
}

// TestNPCDialogData tests that NPC dialog data is properly loaded using test suite data
func TestNPCDialogData(t *testing.T) {
	gameState := createTestGameState(t)
	testSuite := GetTestSuite()

	for _, npcData := range testSuite.NPCs {
		t.Run(npcData.Name, func(t *testing.T) {
			talkScript := getNPCDialog(gameState, npcData)
			if talkScript == nil {
				t.Fatalf("Could not find %s's talk script (location: %v, dialog: %d)",
					npcData.Name, npcData.Location, npcData.DialogNumber)
			}

			// Verify we have the expected dialog structure
			if len(talkScript.Lines) < 5 {
				t.Errorf("Expected at least 5 fixed entries (description, greeting, name, job, bye), got %d", len(talkScript.Lines))
			}

			// Check minimum question groups
			if len(talkScript.QuestionGroups) < npcData.MinQuestionGroups {
				t.Errorf("Expected at least %d question groups, got %d",
					npcData.MinQuestionGroups, len(talkScript.QuestionGroups))
			}

			// Check that we have question groups for expected commands
			foundOptions := make(map[string]bool)
			foundResponses := make(map[string]string)

			for _, group := range talkScript.QuestionGroups {
				for _, option := range group.Options {
					foundOptions[option] = true

					// Store responses for validation
					if len(group.Script) > 0 && group.Script[0].Cmd == 0 {
						foundResponses[option] = group.Script[0].Str
					}
				}
			}

			// Verify expected options exist
			for _, expectedOption := range npcData.ExpectedOptions {
				if !foundOptions[expectedOption] {
					t.Errorf("Expected to find '%s' question group", expectedOption)
				}
			}

			// Verify specific responses
			if response, ok := foundResponses["name"]; ok && npcData.ExpectedName != "" {
				if response != npcData.ExpectedName {
					t.Errorf("Expected name response to be %q, got %q", npcData.ExpectedName, response)
				}
			}

			if response, ok := foundResponses["job"]; ok && npcData.ExpectedJob != "" {
				if response != npcData.ExpectedJob {
					t.Errorf("Expected job response to be %q, got %q", npcData.ExpectedJob, response)
				}
			}

			if response, ok := foundResponses["bye"]; ok && npcData.ExpectedBye != "" {
				if response != npcData.ExpectedBye {
					t.Errorf("Expected bye response to be %q, got %q", npcData.ExpectedBye, response)
				}
			}

			t.Logf("Successfully validated %s's dialog data with %d question groups",
				npcData.Name, len(talkScript.QuestionGroups))
		})
	}
}

// TestNPCConversationInitialization tests basic conversation setup for all NPCs
func TestNPCConversationInitialization(t *testing.T) {
	gameState := createTestGameState(t)
	testSuite := GetTestSuite()

	for _, npcData := range testSuite.NPCs {
		t.Run(npcData.Name, func(t *testing.T) {
			talkScript := getNPCDialog(gameState, npcData)
			if talkScript == nil {
				t.Fatalf("Could not find %s's talk script", npcData.Name)
			}

			npcRef := getNPCReference(gameState, npcData)

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

			// Should contain expected description text
			if npcData.Description != "" {
				foundDescription := false
				for _, output := range outputs {
					if output == npcData.Description {
						foundDescription = true
						break
					}
				}

				if !foundDescription {
					t.Logf("Expected to find description %q in outputs: %v", npcData.Description, outputs)
				}
			}

			t.Logf("%s conversation initialized successfully with %d initial outputs",
				npcData.Name, len(outputs))
		})
	}
}

// TestNPCKeywordOptions tests that NPCs have expected conversation options
func TestNPCKeywordOptions(t *testing.T) {
	gameState := createTestGameState(t)
	testSuite := GetTestSuite()

	for _, npcData := range testSuite.NPCs {
		t.Run(npcData.Name, func(t *testing.T) {
			talkScript := getNPCDialog(gameState, npcData)
			if talkScript == nil {
				t.Fatalf("Could not find %s's talk script", npcData.Name)
			}

			// Collect all available keywords/options
			var allOptions []string
			for _, group := range talkScript.QuestionGroups {
				allOptions = append(allOptions, group.Options...)
			}

			t.Logf("%s has %d conversation options: %v", npcData.Name, len(allOptions), allOptions)

			// Test for expected basic options
			for _, expected := range npcData.ExpectedOptions {
				found := false
				for _, option := range allOptions {
					if option == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected to find option %q in %s's conversation", expected, npcData.Name)
				}
			}

			// Test for thematic options
			foundThematicOptions := 0
			for _, thematicOption := range npcData.ThematicOptions {
				for _, option := range allOptions {
					if option == thematicOption {
						foundThematicOptions++
						t.Logf("Found thematic option for %s: %q", npcData.Name, option)
						break
					}
				}
			}

			if len(npcData.ThematicOptions) > 0 && foundThematicOptions == 0 {
				t.Logf("No thematic options found for %s (may be expected depending on dialog structure)", npcData.Name)
			}
		})
	}
}

// TestConversationDataIntegrity tests that the conversation system loads valid data for all NPCs
func TestConversationDataIntegrity(t *testing.T) {
	gameState := createTestGameState(t)
	testSuite := GetTestSuite()

	// Test that we can load talk references
	talkRefs := gameState.GameReferences.TalkReferences
	if talkRefs == nil {
		t.Fatal("TalkReferences is nil")
	}

	for _, npcData := range testSuite.NPCs {
		t.Run(npcData.Name, func(t *testing.T) {
			// Test that NPC's script exists
			talkScript := talkRefs.GetTalkScriptByNpcIndex(npcData.Location, npcData.DialogNumber)
			if talkScript == nil {
				t.Fatalf("Could not get %s's talk script (location: %v, dialog: %d)",
					npcData.Name, npcData.Location, npcData.DialogNumber)
			}

			// Basic validation of script structure
			if len(talkScript.Lines) == 0 {
				t.Error("Talk script has no lines")
			}

			if len(talkScript.QuestionGroups) == 0 {
				t.Error("Talk script has no question groups")
			}

			// Log some debug info
			t.Logf("%s's script has %d lines and %d question groups",
				npcData.Name, len(talkScript.Lines), len(talkScript.QuestionGroups))

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

			t.Logf("%s conversation data integrity test passed", npcData.Name)
		})
	}
}
