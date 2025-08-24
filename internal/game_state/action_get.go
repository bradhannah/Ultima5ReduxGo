package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionGetSmallMap(direction references.Direction) bool {
	return true
}

func (g *GameState) ActionGetLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Get command - see Commands.md Get section
	// Large map variant of get command
	return true
}

func (g *GameState) ActionGetCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Get command - see Commands.md Get section
	// Combat map variant of get command
	return true
}
