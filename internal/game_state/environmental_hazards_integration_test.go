package game_state

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func TestEnvironmentalHazardsIntegrationSwamp(t *testing.T) {
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithSystemCallbacks().
		WithRandomSeed(12345).
		Build()

	// Setup party with low dexterity for guaranteed poison
	gs.PartyState.Characters[0] = party_state.PlayerCharacter{
		Dexterity:   5, // Very low dexterity
		CurrentHp:   20,
		MaxHp:       20,
		Status:      party_state.Good,
		PartyStatus: party_state.InTheParty,
	}

	// Find or create a swamp tile position for testing
	testPos := references.Position{X: 10, Y: 10}
	layeredMap := gs.MapState.GetLayeredMapByCurrentLocation()

	// Set a swamp tile at the test position using map_state constants
	layeredMap.SetTileByLayer(0, &testPos, indexes.Swamp) // 0 = MapLayer

	// Move player to swamp tile
	gs.MapState.PlayerLocation.Position = testPos

	// Process environmental hazards
	gs.ProcessEnvironmentalHazardsAfterMovement()

	// Verify integration results
	assert.True(t, len(mockCallbacks.Messages) > 0, "Should show hazard message")
	assert.Contains(t, mockCallbacks.Messages, "Poisoned!", "Should show poison message")
	assert.Equal(t, party_state.Poisoned, gs.PartyState.Characters[0].Status, "Character should be poisoned")
}

func TestEnvironmentalHazardsIntegrationLava(t *testing.T) {
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithSystemCallbacks().
		WithRandomSeed(54321).
		Build()

	// Setup party member with known HP
	initialHP := uint16(25)
	gs.PartyState.Characters[0] = party_state.PlayerCharacter{
		CurrentHp:   initialHP,
		MaxHp:       initialHP,
		Status:      party_state.Good,
		PartyStatus: party_state.InTheParty,
	}

	// Find or create a lava tile position for testing
	testPos := references.Position{X: 15, Y: 15}
	layeredMap := gs.MapState.GetLayeredMapByCurrentLocation()

	// Set a lava tile at the test position
	layeredMap.SetTileByLayer(0, &testPos, indexes.Lava) // MapLayer = 0

	// Move player to lava tile
	gs.MapState.PlayerLocation.Position = testPos

	// Process environmental hazards
	gs.ProcessEnvironmentalHazardsAfterMovement()

	// Verify integration results
	assert.True(t, len(mockCallbacks.Messages) > 0, "Should show hazard message")
	assert.Contains(t, mockCallbacks.Messages, "Burning!", "Should show burning message")
	assert.True(t, gs.PartyState.Characters[0].CurrentHp < initialHP, "Character should have taken damage")
	assert.True(t, gs.PartyState.Characters[0].CurrentHp >= 0, "HP should not go below 0")
}

func TestEnvironmentalHazardsIntegrationFireplaceInTown(t *testing.T) {
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain). // Small map location
		WithSystemCallbacks().
		WithRandomSeed(99999).
		Build()

	// Setup party member
	initialHP := uint16(30)
	gs.PartyState.Characters[0] = party_state.PlayerCharacter{
		CurrentHp:   initialHP,
		MaxHp:       initialHP,
		Status:      party_state.Good,
		PartyStatus: party_state.InTheParty,
	}

	// Find or create a fireplace tile position for testing
	testPos := references.Position{X: 12, Y: 12}
	layeredMap := gs.MapState.GetLayeredMapByCurrentLocation()

	// Set a fireplace tile at the test position
	layeredMap.SetTileByLayer(0, &testPos, indexes.Fireplace) // MapLayer = 0

	// Move player to fireplace tile
	gs.MapState.PlayerLocation.Position = testPos

	// Process environmental hazards
	gs.ProcessEnvironmentalHazardsAfterMovement()

	// Verify fireplace is dangerous in small maps
	assert.True(t, len(mockCallbacks.Messages) > 0, "Should show hazard message")
	assert.Contains(t, mockCallbacks.Messages, "Burning!", "Should show burning message")
	assert.True(t, gs.PartyState.Characters[0].CurrentHp < initialHP, "Character should have taken damage")
}

func TestEnvironmentalHazardsIntegrationSafeTiles(t *testing.T) {
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithSystemCallbacks().
		WithRandomSeed(12345).
		Build()

	// Setup party member
	initialHP := uint16(20)
	gs.PartyState.Characters[0] = party_state.PlayerCharacter{
		CurrentHp:   initialHP,
		MaxHp:       initialHP,
		Status:      party_state.Good,
		PartyStatus: party_state.InTheParty,
	}

	// Find or create a safe tile position for testing
	testPos := references.Position{X: 20, Y: 20}
	layeredMap := gs.MapState.GetLayeredMapByCurrentLocation()

	// Set a safe tile at the test position
	layeredMap.SetTileByLayer(0, &testPos, indexes.BrickFloor) // MapLayer = 0

	// Move player to safe tile
	gs.MapState.PlayerLocation.Position = testPos

	// Process environmental hazards
	gs.ProcessEnvironmentalHazardsAfterMovement()

	// Verify no hazard effects
	assert.Empty(t, mockCallbacks.Messages, "Should show no hazard messages")
	assert.Equal(t, party_state.Good, gs.PartyState.Characters[0].Status, "Character should remain in good status")
	assert.Equal(t, initialHP, gs.PartyState.Characters[0].CurrentHp, "Character should not take damage")
}

func TestEnvironmentalHazardsSkipInactivePartyMembers(t *testing.T) {
	gs, _ := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithSystemCallbacks().
		WithRandomSeed(12345).
		Build()

	// Setup party with inactive member
	gs.PartyState.Characters[0] = party_state.PlayerCharacter{
		Dexterity:   5, // Low dex for guaranteed poison if active
		CurrentHp:   20,
		MaxHp:       20,
		Status:      party_state.Good,
		PartyStatus: party_state.AtTheInn, // Not in party
	}
	gs.PartyState.Characters[1] = party_state.PlayerCharacter{
		Dexterity:   5,
		CurrentHp:   20,
		MaxHp:       20,
		Status:      party_state.Good,
		PartyStatus: party_state.InTheParty, // Active party member
	}

	// Set swamp tile
	testPos := references.Position{X: 25, Y: 25}
	layeredMap := gs.MapState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(0, &testPos, indexes.Swamp)
	gs.MapState.PlayerLocation.Position = testPos

	// Process environmental hazards
	gs.ProcessEnvironmentalHazardsAfterMovement()

	// Verify only active party member is affected
	assert.Equal(t, party_state.Good, gs.PartyState.Characters[0].Status, "Inactive member should not be poisoned")
	assert.Equal(t, party_state.Poisoned, gs.PartyState.Characters[1].Status, "Active member should be poisoned")
}

func TestEnvironmentalHazardsTurnAdvancement(t *testing.T) {
	gs, mockCallbacks := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithSystemCallbacks().
		WithRandomSeed(11111).
		Build()

	// Setup party
	gs.PartyState.Characters[0] = party_state.PlayerCharacter{
		CurrentHp:   30,
		MaxHp:       30,
		Status:      party_state.Good,
		PartyStatus: party_state.InTheParty,
	}

	// Set lava tile
	testPos := references.Position{X: 30, Y: 30}
	layeredMap := gs.MapState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(0, &testPos, indexes.Lava)
	gs.MapState.PlayerLocation.Position = testPos

	// Test turn advancement hazards
	mockCallbacks.Reset()
	gs.ProcessEnvironmentalHazardsOnTurnAdvancement()

	// Verify hazards are processed on turn advancement too
	assert.Contains(t, mockCallbacks.Messages, "Burning!", "Should process hazards on turn advancement")
}

func TestEnvironmentalHazardsLazyCaching(t *testing.T) {
	gs, _ := NewIntegrationTestBuilder(t).
		WithLocation(references.Britain).
		WithRandomSeed(12345).
		Build()

	// Initially no environmental hazards system
	assert.Nil(t, gs.environmentalHazards, "Should start with nil environmental hazards")

	// First call should initialize the system
	gs.ProcessEnvironmentalHazardsAfterMovement()
	assert.NotNil(t, gs.environmentalHazards, "Should initialize environmental hazards on first use")

	// Second call should reuse the same system
	firstInstance := gs.environmentalHazards
	gs.ProcessEnvironmentalHazardsAfterMovement()
	assert.Equal(t, firstInstance, gs.environmentalHazards, "Should reuse existing environmental hazards instance")
}
