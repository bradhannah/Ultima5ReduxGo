package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"

	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

// MockSystemCallbacks for testing
type MockSystemCallbacks struct {
	Messages []string
}

func (m *MockSystemCallbacks) AddRowStr(message string) {
	m.Messages = append(m.Messages, message)
}

func (m *MockSystemCallbacks) Reset() {
	m.Messages = m.Messages[:0]
}

func createTestPartyWithStats(stats []TestCharacterStats) *party_state.PartyState {
	party := &party_state.PartyState{}
	for i, stat := range stats {
		if i >= party_state.NPlayers {
			break // Don't exceed array size
		}
		party.Characters[i] = party_state.PlayerCharacter{
			Dexterity:   byte(stat.Dexterity),
			CurrentHp:   uint16(stat.HitPoints),
			MaxHp:       uint16(stat.HitPoints),
			Status:      stat.Status,
			PartyStatus: party_state.InTheParty,
		}
	}
	return party
}

type TestCharacterStats struct {
	Dexterity int
	HitPoints int
	Status    party_state.CharacterStatus
}

func TestHazardDetection(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	mockCallbacks := &MockSystemCallbacks{}
	hazards := NewEnvironmentalHazards(rng, mockCallbacks)

	testCases := []struct {
		name         string
		tileIndex    indexes.SpriteIndex
		location     references.GeneralMapType
		expectedType HazardType
	}{
		{"Swamp tile", indexes.Swamp, references.LargeMapType, PoisonSwamp},
		{"Swamp in small map", indexes.Swamp, references.SmallMapType, PoisonSwamp},
		{"Lava tile", indexes.Lava, references.LargeMapType, LavaBurn},
		{"Fireplace in small map", indexes.Fireplace, references.SmallMapType, FireplaceBurn},
		{"Fireplace on large map", indexes.Fireplace, references.LargeMapType, NoHazard},
		{"Regular floor", indexes.BrickFloor, references.SmallMapType, NoHazard},
		{"Nil tile", 0, references.LargeMapType, NoHazard},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var tile *references.Tile
			if tc.tileIndex != 0 {
				tile = &references.Tile{Index: tc.tileIndex}
			}

			hazardType := hazards.detectHazardType(tile, tc.location)
			assert.Equal(t, tc.expectedType, hazardType)
		})
	}
}

func TestSwampPoisonWithLowDexterity(t *testing.T) {
	// Setup deterministic RNG
	rng := rand.New(rand.NewSource(12345))
	mockCallbacks := &MockSystemCallbacks{}
	hazards := NewEnvironmentalHazards(rng, mockCallbacks)

	// Create party with low dexterity character
	party := createTestPartyWithStats([]TestCharacterStats{
		{Dexterity: 5, HitPoints: 20, Status: party_state.Good},
		{Dexterity: 25, HitPoints: 20, Status: party_state.Good},
	})

	swampTile := &references.Tile{Index: indexes.Swamp}
	result := hazards.CheckTileHazards(swampTile, party, references.LargeMapType)

	// Verify poison mechanics
	assert.Equal(t, PoisonSwamp, result.Type)
	assert.True(t, len(result.Affected) > 0, "Low dex character should be poisoned")
	assert.True(t, result.MessageShown)
	assert.Contains(t, mockCallbacks.Messages, "Poisoned!")

	// Verify poisoned character status
	poisonedFound := false
	for _, member := range party.Characters {
		if member.Status == party_state.Poisoned {
			poisonedFound = true
			break
		}
	}
	assert.True(t, poisonedFound, "At least one character should be poisoned")
}

func TestSwampPoisonWithHighDexterity(t *testing.T) {
	// Setup deterministic RNG
	rng := rand.New(rand.NewSource(99999))
	mockCallbacks := &MockSystemCallbacks{}
	hazards := NewEnvironmentalHazards(rng, mockCallbacks)

	// Create party with high dexterity characters
	party := createTestPartyWithStats([]TestCharacterStats{
		{Dexterity: 30, HitPoints: 20, Status: party_state.Good},
		{Dexterity: 30, HitPoints: 20, Status: party_state.Good},
	})

	swampTile := &references.Tile{Index: indexes.Swamp}
	result := hazards.CheckTileHazards(swampTile, party, references.LargeMapType)

	// High dex should have lower chance of poison (though not guaranteed)
	assert.Equal(t, PoisonSwamp, result.Type)
}

func TestSwampSkipsDeadAndPoisonedCharacters(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	mockCallbacks := &MockSystemCallbacks{}
	hazards := NewEnvironmentalHazards(rng, mockCallbacks)

	// Create party with dead and already poisoned characters
	party := createTestPartyWithStats([]TestCharacterStats{
		{Dexterity: 1, HitPoints: 0, Status: party_state.Dead},
		{Dexterity: 1, HitPoints: 20, Status: party_state.Poisoned},
		{Dexterity: 5, HitPoints: 20, Status: party_state.Good},
	})

	swampTile := &references.Tile{Index: indexes.Swamp}
	result := hazards.CheckTileHazards(swampTile, party, references.LargeMapType)

	// Only living, non-poisoned characters should be affected
	assert.Equal(t, PoisonSwamp, result.Type)

	// Verify dead and poisoned characters were skipped
	assert.Equal(t, party_state.Dead, party.Characters[0].Status, "Dead character should remain dead")
	assert.Equal(t, party_state.Poisoned, party.Characters[1].Status, "Poisoned character should remain poisoned")
}

func TestLavaDamageApplication(t *testing.T) {
	// Test lava damage with known RNG seed
	rng := rand.New(rand.NewSource(54321))
	mockCallbacks := &MockSystemCallbacks{}
	hazards := NewEnvironmentalHazards(rng, mockCallbacks)

	party := createTestPartyWithStats([]TestCharacterStats{
		{HitPoints: 20, Status: party_state.Good},
		{HitPoints: 15, Status: party_state.Good},
	})
	initialHP := []int{int(party.Characters[0].CurrentHp), int(party.Characters[1].CurrentHp)}

	lavaTile := &references.Tile{Index: indexes.Lava}
	result := hazards.CheckTileHazards(lavaTile, party, references.LargeMapType)

	// Verify damage application
	assert.Equal(t, LavaBurn, result.Type)
	assert.True(t, result.DamageDealt > 0)
	assert.Contains(t, mockCallbacks.Messages, "Burning!")

	// Verify actual HP reduction
	for i, character := range party.Characters[:2] { // Only check first 2 characters we set up
		if character.PartyStatus == party_state.InTheParty {
			assert.True(t, int(character.CurrentHp) < initialHP[i], "Member %d should have taken damage", i)
			assert.True(t, int(character.CurrentHp) >= 0, "HP should not go below 0")
		}
	}
}

func TestFireplaceDamageOnlyInSmallMaps(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	mockCallbacks := &MockSystemCallbacks{}
	hazards := NewEnvironmentalHazards(rng, mockCallbacks)

	party := createTestPartyWithStats([]TestCharacterStats{
		{HitPoints: 20, Status: party_state.Good},
	})

	fireplaceTile := &references.Tile{Index: indexes.Fireplace}

	// Test in small map (should cause damage)
	mockCallbacks.Reset()
	result := hazards.CheckTileHazards(fireplaceTile, party, references.SmallMapType)
	assert.Equal(t, FireplaceBurn, result.Type)
	assert.Contains(t, mockCallbacks.Messages, "Burning!")

	// Test on large map (should be safe)
	mockCallbacks.Reset()
	result = hazards.CheckTileHazards(fireplaceTile, party, references.LargeMapType)
	assert.Equal(t, NoHazard, result.Type)
	assert.Empty(t, mockCallbacks.Messages)
}

func TestDeadCharacterFromDamage(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	mockCallbacks := &MockSystemCallbacks{}
	hazards := NewEnvironmentalHazards(rng, mockCallbacks)

	// Create character with very low HP
	party := createTestPartyWithStats([]TestCharacterStats{
		{HitPoints: 2, Status: party_state.Good}, // Very low HP, likely to die
	})

	lavaTile := &references.Tile{Index: indexes.Lava}
	result := hazards.CheckTileHazards(lavaTile, party, references.LargeMapType)

	// Verify damage was applied
	assert.Equal(t, LavaBurn, result.Type)
	assert.True(t, result.DamageDealt > 0)

	// Character should be dead if HP dropped to 0
	if party.Characters[0].CurrentHp == 0 {
		assert.Equal(t, party_state.Dead, party.Characters[0].Status)
	}
}

func TestNoHazardForSafeTiles(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	mockCallbacks := &MockSystemCallbacks{}
	hazards := NewEnvironmentalHazards(rng, mockCallbacks)

	party := createTestPartyWithStats([]TestCharacterStats{
		{HitPoints: 20, Status: party_state.Good},
	})

	safeTile := &references.Tile{Index: indexes.BrickFloor}
	result := hazards.CheckTileHazards(safeTile, party, references.SmallMapType)

	assert.Equal(t, NoHazard, result.Type)
	assert.Empty(t, result.Affected)
	assert.False(t, result.MessageShown)
	assert.Equal(t, 0, result.DamageDealt)
	assert.Empty(t, mockCallbacks.Messages)
}
