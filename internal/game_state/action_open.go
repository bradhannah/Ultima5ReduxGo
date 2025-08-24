package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionOpenSmallMap(direction references.Direction) bool {
	// TODO: Implement small map Open command - see Commands.md Open section
	// Should handle:
	// - Closing previously opened timed doors
	// - Getting direction and target tile
	// - Handling doors (locked/unlocked), windows, footlockers, portcullis
	// - Opening chests via open_chest_surface helper
	// - Timed door auto-close mechanism
	return true
}

func (g *GameState) ActionOpenLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Open command - see Commands.md Open section
	// Large map variant of open command
	return true
}

func (g *GameState) ActionOpenCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Open command - see Commands.md Open section
	// Combat map variant of open command - likely limited to chests
	return true
}
