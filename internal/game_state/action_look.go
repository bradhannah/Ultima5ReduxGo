package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionLookSmallMap(direction references.Direction) bool {
	return true
}

func (g *GameState) ActionLookLargeMap(direction references.Direction) bool {
	return true
}

func (g *GameState) ActionLookCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Look command - see Commands.md Look section
	// Combat map variant of look command
	return true
}
