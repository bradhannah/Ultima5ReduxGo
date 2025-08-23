package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionUseSmallMap(direction references.Direction) bool {
	// TODO: Implement small map Use command - see Commands.md Use section
	// Should handle:
	// - Magic Carpet usage (boarding conditions, consumption)
	// - Skull Key usage (directional In Ex Por semantics)
	// - Amulet/Crown equipping (persistent effects)
	// - Sceptre usage (clear chests/dispel fields)
	// - Spyglass usage (astronomy at night)
	// - Item validation and context gating
	return true
}

func (g *GameState) ActionUseLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Use command - see Commands.md Use section
	// Large map variant of use command
	return true
}
