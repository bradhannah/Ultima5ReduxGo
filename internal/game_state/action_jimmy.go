package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// These Action methods are kept for API consistency but are no longer used by UI
// The UI now calls JimmyDoor directly to avoid the double-call bug

func (g *GameState) ActionJimmySmallMap(direction references.Direction) bool {
	selectedCharacter := g.SelectCharacterForJimmy()
	if selectedCharacter == nil {
		return false
	}
	jimmyResult := g.JimmyDoor(direction, selectedCharacter)
	return jimmyResult == JimmyUnlocked
}

func (g *GameState) ActionJimmyLargeMap(direction references.Direction) bool {
	selectedCharacter := g.SelectCharacterForJimmy()
	if selectedCharacter == nil {
		return false
	}
	jimmyResult := g.JimmyDoor(direction, selectedCharacter)
	return jimmyResult == JimmyUnlocked
}

func (g *GameState) ActionJimmyCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Jimmy command - see Commands.md Jimmy section
	// Combat map variant of jimmy command - likely limited functionality
	return true
}

func (g *GameState) SelectCharacterForJimmy() *party_state.PlayerCharacter {
	// TODO: Implement character selection UI as per Commands.md Jimmy section
	// For now, default to first character if alive and in party
	if len(g.PartyState.Characters) > 0 && g.PartyState.Characters[0].Status == party_state.Good {
		return &g.PartyState.Characters[0]
	}

	// Find first good character if first one isn't available
	for i := range g.PartyState.Characters {
		if g.PartyState.Characters[i].Status == party_state.Good {
			return &g.PartyState.Characters[i]
		}
	}

	return nil // No available characters
}
