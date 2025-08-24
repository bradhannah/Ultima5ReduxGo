package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionTalkSmallMap(direction references.Direction) bool {
	// TODO: Implement small map Talk command - see Commands.md Talk section
	// Should handle:
	// - Freed NPC acknowledgement (stocks/manacles with karma +2)
	// - Standard conversation initiation with NPCs
	// - Integration with conversation system
	// - Target validation (NPC present at target tile)
	return true
}

func (g *GameState) ActionTalkLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Talk command - see Commands.md Talk section
	// Large map variant of talk command
	return true
}

func (g *GameState) ActionTalkCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Talk command - see Commands.md Talk section
	// Combat map variant of talk command - likely disabled during combat
	return true
}
