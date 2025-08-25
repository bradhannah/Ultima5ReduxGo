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

	// Special items not implemented yet
	g.SystemCallbacks.Message.AddRowStr("Nothing happens.")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return false
}

func (g *GameState) ActionUseLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Use command - see Commands.md Use section
	// Large map variant of use command

	// Special items not implemented yet
	g.SystemCallbacks.Message.AddRowStr("Nothing happens.")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return false
}

func (g *GameState) ActionUseCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Use command - see Commands.md Use section
	// Combat map variant of use command - items/spells during combat

	// Not allowed during combat
	g.SystemCallbacks.Message.AddRowStr("Not now!")
	return false
}

func (g *GameState) ActionUseDungeonMap(direction references.Direction) bool {
	// TODO: Implement dungeon map Use command - see Commands.md Use section
	// Dungeon map variant of use command
	// Should handle special dungeon items and fixtures

	// Special items not implemented yet
	g.SystemCallbacks.Message.AddRowStr("Nothing happens.")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return false
}
