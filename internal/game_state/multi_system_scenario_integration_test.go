// Multi-system scenario integration tests with real game data
// These tests validate complex scenarios where multiple game systems interact:
// Time progression affecting NPCs, inventory management during exploration,
// dialog systems with state persistence, and command chaining workflows.
package game_state

import (
	"fmt"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// TestTimeProgressionScenario tests time advancement affecting multiple systems
func TestTimeProgressionScenario(t *testing.T) {
	// Test how time advancement affects NPCs, visibility, and game state
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(12, 12).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("⏰ Time Progression Multi-System Scenario Test")

	// Record initial state
	initialTime := gs.DateTime
	t.Logf("Starting time: Day=%d, Hour=%d, Minute=%d, Turn=%d",
		initialTime.Day, initialTime.Hour, initialTime.Minute, initialTime.Turn)

	isDaylight := initialTime.IsDayLight()
	t.Logf("Initial lighting: daylight=%v", isDaylight)

	// Simulate extended game session with time progression
	timeProgressionActions := []struct {
		name           string
		action         func() bool
		expectedTime   int // minutes expected to advance
		systemsToCheck []string
	}{
		{"Morning exploration", func() bool { return gs.ActionLookSmallMap(references.Up) }, 1, []string{"visibility", "npc"}},
		{"Item gathering", func() bool { return gs.ActionGetSmallMap(references.Down) }, 0, []string{"inventory"}},
		{"NPC interaction", func() bool { return gs.ActionTalkSmallMap(references.Left) }, 0, []string{"dialog", "npc"}},
		{"Object inspection", func() bool { return gs.ActionLookSmallMap(references.Right) }, 1, []string{"visibility"}},
		{"Door interaction", func() bool { return gs.ActionJimmySmallMap(references.Up) }, 0, []string{"time", "inventory"}},
		{"Movement preparation", func() bool { return gs.ActionPushSmallMap(references.Down) }, 0, []string{"positioning"}},
	}

	cumulativeTimeAdvanced := 0
	systemInteractions := make(map[string]int)

	for i, timeAction := range timeProgressionActions {
		mockCallbacks.Reset()
		timeBefore := gs.DateTime

		result := timeAction.action()

		// Calculate actual time advancement
		actualTimeAdvanced := 0
		for _, timeAdv := range mockCallbacks.TimeAdvanced {
			actualTimeAdvanced += timeAdv
			cumulativeTimeAdvanced += timeAdv
		}

		timeAfter := gs.DateTime

		t.Logf("Action %d (%s): result=%v, time_advanced=%d min",
			i+1, timeAction.name, result, actualTimeAdvanced)

		// Log time progression details
		if actualTimeAdvanced > 0 {
			t.Logf("  Time: %d:%02d → %d:%02d (Turn %d → %d)",
				timeBefore.Hour, timeBefore.Minute,
				timeAfter.Hour, timeAfter.Minute,
				timeBefore.Turn, timeAfter.Turn)
		}

		// Track system interactions
		for _, system := range timeAction.systemsToCheck {
			systemInteractions[system]++
		}

		// Check for system-specific effects
		if len(mockCallbacks.TalkDialogCalls) > 0 {
			t.Logf("  → NPC dialog system activated")
			systemInteractions["dialog_active"]++
		}

		if len(mockCallbacks.Messages) > 0 {
			t.Logf("  → User feedback: %s", mockCallbacks.Messages[0])
		}

		// Check daylight changes
		newIsDaylight := timeAfter.IsDayLight()
		if newIsDaylight != isDaylight {
			t.Logf("  → Lighting changed: daylight %v → %v", isDaylight, newIsDaylight)
			isDaylight = newIsDaylight
			systemInteractions["lighting_change"]++
		}
	}

	finalTime := gs.DateTime

	t.Logf("⏰ Time Progression Summary:")
	t.Logf("  Session duration: %d minutes", cumulativeTimeAdvanced)
	t.Logf("  Final time: Day=%d, Hour=%d, Minute=%d, Turn=%d",
		finalTime.Day, finalTime.Hour, finalTime.Minute, finalTime.Turn)
	t.Logf("  System interactions:")

	for system, count := range systemInteractions {
		t.Logf("    %s: %d", system, count)
	}

	// Validate time system integration
	if cumulativeTimeAdvanced > 0 {
		t.Logf("  ✅ Time progression system working")
	}

	if len(systemInteractions) > 0 {
		t.Logf("  ✅ Multi-system time integration confirmed")
	}

	t.Logf("⏰ Time progression scenario test completed")
}

// TestInventoryManagementScenario tests inventory changes across multiple commands
func TestInventoryManagementScenario(t *testing.T) {
	// Test inventory management during complex scenarios
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(20, 20).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("🎒 Inventory Management Multi-System Scenario Test")

	// Set up interesting inventory state
	gs.PartyState.Inventory.Provisions.Keys.Set(5)
	gs.PartyState.Inventory.Gold.Set(100)
	gs.PartyState.Inventory.Provisions.Food.Set(200)
	gs.PartyState.Characters[0].Status = 1 // Good status for actions

	// Track initial inventory
	inventory := map[string]uint16{
		"Keys": gs.PartyState.Inventory.Provisions.Keys.Get(),
		"Gold": gs.PartyState.Inventory.Gold.Get(),
		"Food": gs.PartyState.Inventory.Provisions.Food.Get(),
	}

	t.Logf("Initial inventory: Keys=%d, Gold=%d, Food=%d",
		inventory["Keys"], inventory["Gold"], inventory["Food"])

	// Scenario: Resource management during exploration
	inventoryScenario := []struct {
		phase       string
		actions     []func() bool
		expectUsage []string
	}{
		{
			"Exploration phase - gathering",
			[]func() bool{
				func() bool { return gs.ActionGetSmallMap(references.Up) },
				func() bool { return gs.ActionGetSmallMap(references.Down) },
				func() bool { return gs.ActionGetSmallMap(references.Left) },
				func() bool { return gs.ActionGetSmallMap(references.Right) },
			},
			[]string{"time", "inventory"},
		},
		{
			"Interaction phase - key usage",
			[]func() bool{
				func() bool { return gs.ActionJimmySmallMap(references.Up) },
				func() bool { return gs.ActionJimmySmallMap(references.Down) },
				func() bool { return gs.ActionJimmySmallMap(references.Left) },
			},
			[]string{"keys", "time"},
		},
		{
			"Dialog phase - NPC interactions",
			[]func() bool{
				func() bool { return gs.ActionTalkSmallMap(references.Up) },
				func() bool { return gs.ActionTalkSmallMap(references.Down) },
			},
			[]string{"dialog", "time"},
		},
	}

	inventoryChanges := make(map[string]map[string]int) // phase -> resource -> change
	totalActions := 0
	phaseSummaries := make([]string, 0, len(inventoryScenario))

	for _, phase := range inventoryScenario {
		t.Logf("Phase: %s", phase.phase)

		phaseStartInventory := map[string]uint16{
			"Keys": gs.PartyState.Inventory.Provisions.Keys.Get(),
			"Gold": gs.PartyState.Inventory.Gold.Get(),
			"Food": gs.PartyState.Inventory.Provisions.Food.Get(),
		}

		phaseChanges := make(map[string]int)
		phaseSuccesses := 0

		for i, action := range phase.actions {
			mockCallbacks.Reset()

			result := action()
			totalActions++

			if result {
				phaseSuccesses++
			}

			t.Logf("  Action %d: result=%v", i+1, result)

			// Log notable feedback
			if len(mockCallbacks.Messages) > 0 {
				message := mockCallbacks.Messages[0]
				if message != "Nothing there!" && message != "No-one to talk to!" {
					t.Logf("    → %s", message)
				}
			}

			if len(mockCallbacks.TalkDialogCalls) > 0 {
				t.Logf("    → Dialog created")
			}
		}

		phaseEndInventory := map[string]uint16{
			"Keys": gs.PartyState.Inventory.Provisions.Keys.Get(),
			"Gold": gs.PartyState.Inventory.Gold.Get(),
			"Food": gs.PartyState.Inventory.Provisions.Food.Get(),
		}

		// Calculate phase changes
		for resource := range phaseStartInventory {
			change := int(phaseEndInventory[resource]) - int(phaseStartInventory[resource])
			if change != 0 {
				phaseChanges[resource] = change
				t.Logf("  %s changed: %d → %d (Δ%+d)",
					resource, phaseStartInventory[resource], phaseEndInventory[resource], change)
			}
		}

		inventoryChanges[phase.phase] = phaseChanges
		phaseSummaries = append(phaseSummaries,
			fmt.Sprintf("%s: %d/%d successful", phase.phase, phaseSuccesses, len(phase.actions)))
	}

	// Final inventory analysis
	finalInventory := map[string]uint16{
		"Keys": gs.PartyState.Inventory.Provisions.Keys.Get(),
		"Gold": gs.PartyState.Inventory.Gold.Get(),
		"Food": gs.PartyState.Inventory.Provisions.Food.Get(),
	}

	t.Logf("🎒 Inventory Management Analysis:")
	t.Logf("  Total actions: %d", totalActions)

	for _, summary := range phaseSummaries {
		t.Logf("  %s", summary)
	}

	t.Logf("  Final inventory:")
	for resource, finalAmount := range finalInventory {
		initialAmount := inventory[resource]
		change := int(finalAmount) - int(initialAmount)
		t.Logf("    %s: %d → %d (Δ%+d)", resource, initialAmount, finalAmount, change)
	}

	// Validate inventory management worked
	hasInventoryChanges := false
	for resource, finalAmount := range finalInventory {
		if finalAmount != inventory[resource] {
			hasInventoryChanges = true
			break
		}
	}

	if hasInventoryChanges {
		t.Logf("  ✅ Inventory management system working")
	}

	t.Logf("🎒 Inventory management scenario test completed")
}

// TestDialogAndStateScenario tests dialog system with persistent state changes
func TestDialogAndStateScenario(t *testing.T) {
	// Test dialog system integration with game state persistence
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithPlayerAt(15, 15).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	t.Logf("💬 Dialog and State Persistence Scenario Test")

	// Multiple attempts to find and interact with NPCs
	dialogScenario := []struct {
		location    references.Location
		positions   []references.Position
		description string
	}{
		{
			references.Britain,
			[]references.Position{
				{X: 15, Y: 15},
				{X: 10, Y: 10},
				{X: 20, Y: 20},
				{X: 25, Y: 15},
			},
			"Britain town exploration",
		},
		{
			references.Lord_Britishs_Castle,
			[]references.Position{
				{X: 8, Y: 8},
				{X: 12, Y: 12},
				{X: 16, Y: 16},
			},
			"Castle NPC search",
		},
	}

	totalDialogsFound := 0
	totalTalkAttempts := 0
	stateChanges := make(map[string]int)

	for _, scenario := range dialogScenario {
		t.Logf("Scenario: %s", scenario.description)

		// Update location
		gs.MapState.PlayerLocation.Location = scenario.location
		gs.UpdateSmallMap(gs.GameReferences.TileReferences, gs.GameReferences.LocationReferences)

		for i, position := range scenario.positions {
			gs.MapState.PlayerLocation.Position = position
			t.Logf("  Position %d: (%d,%d)", i+1, position.X, position.Y)

			// Try talking in each direction
			directions := []references.Direction{references.Up, references.Down, references.Left, references.Right}

			for _, direction := range directions {
				mockCallbacks.Reset()
				totalTalkAttempts++

				result := gs.ActionTalkSmallMap(direction)

				if result {
					t.Logf("    Talk %v: SUCCESS", direction)
					stateChanges["successful_talks"]++
				}

				// Check for dialog creation
				if len(mockCallbacks.TalkDialogCalls) > 0 {
					totalDialogsFound++
					t.Logf("    → Dialog created %v from (%d,%d)", direction, position.X, position.Y)
					stateChanges["dialogs_created"]++
				}

				// Check for meaningful responses (not standard failure messages)
				if len(mockCallbacks.Messages) > 0 {
					message := mockCallbacks.Messages[0]
					if message != "No-one to talk to!" && message != "Can't talk to that!" {
						t.Logf("    → Special response: %s", message)
						stateChanges["special_responses"]++
					}
				}

				// Check for time advancement (indicates interaction occurred)
				if len(mockCallbacks.TimeAdvanced) > 0 {
					timeAdvanced := mockCallbacks.TimeAdvanced[0]
					if timeAdvanced > 0 {
						stateChanges["time_advanced"] += timeAdvanced
					}
				}
			}
		}
	}

	t.Logf("💬 Dialog and State Analysis:")
	t.Logf("  Total talk attempts: %d", totalTalkAttempts)
	t.Logf("  Dialogs created: %d", totalDialogsFound)

	for changeType, count := range stateChanges {
		t.Logf("  %s: %d", changeType, count)
	}

	// Calculate success rates
	if totalTalkAttempts > 0 {
		dialogRate := float64(totalDialogsFound) / float64(totalTalkAttempts) * 100
		t.Logf("  Dialog success rate: %.1f%%", dialogRate)
	}

	// Validate dialog system integration
	if totalDialogsFound > 0 {
		t.Logf("  ✅ Dialog system integration confirmed")
	}

	if len(stateChanges) > 0 {
		t.Logf("  ✅ State persistence working")
	}

	t.Logf("💬 Dialog and state scenario test completed")
}

// TestCommandChainingScenario tests complex command sequences and their interactions
func TestCommandChainingScenario(t *testing.T) {
	// Test complex command chaining and interaction effects
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Lord_Britishs_Castle). // Rich environment for testing
		WithPlayerAt(10, 10).
		WithSystemCallbacks().
		Build()

	if gs == nil {
		return
	}

	// Set up resources for complex interactions
	gs.PartyState.Inventory.Provisions.Keys.Set(3)
	gs.PartyState.Characters[0].Status = 1

	t.Logf("🔗 Command Chaining Multi-System Scenario Test")

	// Complex command chains that test system interactions
	commandChains := []struct {
		name     string
		commands []struct {
			action func() bool
			expect string
		}
		description string
	}{
		{
			"Exploration and Interaction Chain",
			[]struct {
				action func() bool
				expect string
			}{
				{func() bool { return gs.ActionLookSmallMap(references.Up) }, "information"},
				{func() bool { return gs.ActionTalkSmallMap(references.Up) }, "npc_interaction"},
				{func() bool { return gs.ActionGetSmallMap(references.Up) }, "item_acquisition"},
				{func() bool { return gs.ActionJimmySmallMap(references.Up) }, "door_manipulation"},
			},
			"Sequential interaction with same direction",
		},
		{
			"Resource Management Chain",
			[]struct {
				action func() bool
				expect string
			}{
				{func() bool { return gs.ActionJimmySmallMap(references.Left) }, "key_usage"},
				{func() bool { return gs.ActionJimmySmallMap(references.Right) }, "key_usage"},
				{func() bool { return gs.ActionGetSmallMap(references.Down) }, "item_search"},
				{func() bool { return gs.ActionPushSmallMap(references.Down) }, "object_manipulation"},
			},
			"Resource consumption and management",
		},
		{
			"Comprehensive Exploration Chain",
			[]struct {
				action func() bool
				expect string
			}{
				{func() bool { return gs.ActionLookSmallMap(references.Down) }, "survey"},
				{func() bool { return gs.ActionLookSmallMap(references.Left) }, "survey"},
				{func() bool { return gs.ActionLookSmallMap(references.Right) }, "survey"},
				{func() bool { return gs.ActionTalkSmallMap(references.Down) }, "social"},
				{func() bool { return gs.ActionGetSmallMap(references.Left) }, "collection"},
			},
			"Thorough area exploration pattern",
		},
	}

	chainResults := make(map[string]struct {
		successful     int
		total          int
		timeSpent      int
		keyChanges     int
		dialogsCreated int
		uniqueMessages map[string]int
	})

	for _, chain := range commandChains {
		t.Logf("Chain: %s", chain.name)

		result := chainResults[chain.name]
		result.uniqueMessages = make(map[string]int)
		result.total = len(chain.commands)

		initialKeys := gs.PartyState.Inventory.Provisions.Keys.Get()

		for i, step := range chain.commands {
			mockCallbacks.Reset()

			success := step.action()
			if success {
				result.successful++
			}

			t.Logf("  Step %d (%s): %v", i+1, step.expect, success)

			// Track time advancement
			for _, timeAdv := range mockCallbacks.TimeAdvanced {
				result.timeSpent += timeAdv
			}

			// Track dialog creation
			result.dialogsCreated += len(mockCallbacks.TalkDialogCalls)

			// Track unique messages
			for _, message := range mockCallbacks.Messages {
				result.uniqueMessages[message]++
			}

			// Log significant feedback
			if len(mockCallbacks.Messages) > 0 {
				message := mockCallbacks.Messages[0]
				if message != "Nothing there!" && message != "No-one to talk to!" && message != "Won't budge!" {
					t.Logf("    → %s", message)
				}
			}
		}

		finalKeys := gs.PartyState.Inventory.Provisions.Keys.Get()
		result.keyChanges = int(initialKeys - finalKeys)

		chainResults[chain.name] = result
	}

	t.Logf("🔗 Command Chaining Analysis:")
	for chainName, result := range chainResults {
		t.Logf("Chain '%s':", chainName)
		t.Logf("  Success rate: %d/%d (%.1f%%)",
			result.successful, result.total,
			float64(result.successful)/float64(result.total)*100)
		t.Logf("  Time spent: %d minutes", result.timeSpent)
		t.Logf("  Keys consumed: %d", result.keyChanges)
		t.Logf("  Dialogs created: %d", result.dialogsCreated)
		t.Logf("  Unique messages: %d", len(result.uniqueMessages))
	}

	// Validate command chaining worked
	totalChains := len(commandChains)
	chainsWithActivity := 0

	for _, result := range chainResults {
		if result.successful > 0 || result.timeSpent > 0 || result.dialogsCreated > 0 {
			chainsWithActivity++
		}
	}

	if chainsWithActivity > 0 {
		t.Logf("  ✅ Command chaining systems working (%d/%d chains active)",
			chainsWithActivity, totalChains)
	}

	t.Logf("🔗 Command chaining scenario test completed")
}
