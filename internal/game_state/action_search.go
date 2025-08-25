package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionSearchSmallMap(direction references.Direction) bool {
	// TODO: Implement small map Search command - see Commands.md Search section
	// Should handle:
	// - Moongate tile exclusion
	// - Stone cache searching
	// - Reagent cache searching
	// - Hidden objects (daily skull keys, castle keys, glass swords)
	// - Integration with Secrets.md for location-specific spawns

	// For now, just return not found since search systems aren't implemented
	g.SystemCallbacks.Message.AddRowStr("Not found!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return false
}

func (g *GameState) ActionSearchLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Search command - see Commands.md Search section
	// Large map variant of search command

	// For now, just return not found since search systems aren't implemented
	g.SystemCallbacks.Message.AddRowStr("Not found!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return false
}

func (g *GameState) ActionSearchCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Search command - see Commands.md Search section
	// Combat map variant of search command

	// Not allowed during combat
	g.SystemCallbacks.Message.AddRowStr("Not now!")
	return false
}

func (g *GameState) ActionSearchDungeonMap(direction references.Direction) bool {
	// TODO: Implement dungeon map Search command - see Commands.md Search — Dungeon section
	// Should handle:
	// - Ahead-of-avatar search for secret doors/passages
	// - Light requirement for searching
	// - Secret door discovery mechanics
	// - Hidden passage detection
	// - Time cost for searching

	// Check if there's light (torch or magic light)
	if !g.MapState.Lighting.HasTorchLit() {
		g.SystemCallbacks.Message.AddRowStr("It's too dark!")
		return false
	}

	// For now, just return not found since search systems aren't implemented
	g.SystemCallbacks.Message.AddRowStr("Not found!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return false
}
