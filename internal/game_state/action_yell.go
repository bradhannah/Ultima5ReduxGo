package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionYellSmallMap(direction references.Direction) bool {
	// TODO: Implement small map Yell command - see Commands.md Yell section
	// Should handle:
	// - Ship sail control (hoist/furl on small maps)
	// - Town shadowlord summoning (castle-specific words)
	// - Overworld Words of Power (dungeon unsealing/shrine restoration)
	// - Context validation (town vs overworld vs ship)
	// - Word recognition and effect triggering
	return true
}

func (g *GameState) ActionYellLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Yell command - see Commands.md Yell section
	// Large map variant of yell command - primarily overworld Words of Power
	return true
}

func (g *GameState) ActionYellCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Yell command - see Commands.md Yell section
	// Combat map variant of yell command - likely limited functionality
	return true
}
