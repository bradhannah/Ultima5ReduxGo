package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
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
	targetPos := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)

	// Get the tile at the target position
	layeredMap := g.GetLayeredMapByCurrentLocation()
	topTile := layeredMap.GetTopTile(targetPos)

	// Check if there's an item stack at this position (like chests)
	if g.ItemStacksMap.HasItemStackAtPosition(targetPos) {
		// For now, just indicate there's something here but can't be opened on large map
		// This could be expanded later to handle specific cases
		return false
	}

	// Check if the tile is openable (like doors)
	if topTile.Index.IsDoor() {
		// Use the door opening logic from small map
		switch g.MapState.OpenDoor(direction) {
		case map_state.OpenDoorOpened:
			g.SystemCallbacks.Message.AddRowStr("Opened!")
			return true
		case map_state.OpenDoorLocked:
			g.SystemCallbacks.Message.AddRowStr("Locked!")
			return false
		case map_state.OpenDoorLockedMagical:
			g.SystemCallbacks.Message.AddRowStr("Magically Locked!")
			return false
		case map_state.OpenDoorNotADoor:
			// Not a door, continue to check other openable things
			break
		}
	}

	// On large maps, most openable things are handled by entering buildings
	// or specific interactions, so most open attempts will fail
	return false
}

func (g *GameState) ActionOpenCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Open command - see Commands.md Open section
	// Combat map variant of open command - likely limited to chests
	return true
}
