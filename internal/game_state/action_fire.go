package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionFireSmallMap(direction references.Direction) bool {
	// TODO: Implement small map Fire command - see Commands.md Fire (Cannons) section
	// Should handle:
	// - Town cannon firing with directional targeting
	// - Cannonball consumption
	// - Damage calculation and blast radius
	// - Negative: "No cannons!" if not near cannon
	// - Negative: "No cannonballs!" if none in inventory
	return true
}

func (g *GameState) ActionFireLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Fire command - see Commands.md Fire (Cannons) section
	// Should handle:
	// - Ship broadside firing (port/starboard)
	// - Negative: "Fire broadsides only!" if not on ship
	// - Negative: "What?" or "Cannot!" based on context
	// - Frigate special: all 4 cannons fire in broadside
	return true
}

func (g *GameState) ActionFireCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Fire command - see Commands.md Fire (Cannons) section
	// Combat map variant - tactical cannon use
	return true
}

func (g *GameState) ActionFireDungeonMap(direction references.Direction) bool {
	// TODO: Implement dungeon map Fire command - see Commands.md Fire (Cannons) section
	// Dungeon map variant - likely not applicable
	return true
}
