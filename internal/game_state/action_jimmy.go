package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionJimmySmallMap(direction references.Direction) bool {
	// Use existing JimmyDoor logic but with first character by default
	// TODO: Allow character selection as per Commands.md Jimmy section
	jimmyResult := g.JimmyDoor(direction, &g.PartyState.Characters[0])
	return jimmyResult == JimmyUnlocked
}

func (g *GameState) ActionJimmyLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Jimmy command - see Commands.md Jimmy section
	// Large map variant of jimmy command - may need different targeting logic
	return g.ActionJimmySmallMap(direction) // Temporary fallback
}

func (g *GameState) ActionJimmyCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Jimmy command - see Commands.md Jimmy section
	// Combat map variant of jimmy command - likely limited functionality
	return true
}
