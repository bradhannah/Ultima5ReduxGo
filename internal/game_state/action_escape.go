package game_state

import "github.com/bradhannah/Ultima5ReduxGo/internal/references"

func (g *GameState) ActionEscapeSmallMap() bool {
	// Escape is only valid from combat screens
	g.SystemCallbacks.Message.AddRowStr("Not here!")
	return false
}

func (g *GameState) ActionEscapeLargeMap() bool {
	// Escape is only valid from combat screens
	g.SystemCallbacks.Message.AddRowStr("Not here!")
	return false
}

func (g *GameState) ActionEscapeCombatMap() bool {
	// Check if in dungeon room combat (can't escape dungeon rooms)
	if g.MapState.PlayerLocation.Location.GetMapType() == references.DungeonMapType {
		g.SystemCallbacks.Message.AddRowStr("Not here!")
		return false
	}

	// TODO: Check if battle is won (victory flag)
	// For now, since combat isn't implemented, always deny escape
	g.SystemCallbacks.Message.AddRowStr("Not yet!")
	return false
}

func (g *GameState) ActionEscapeDungeonMap() bool {
	// Escape is only valid from combat screens
	g.SystemCallbacks.Message.AddRowStr("Not here!")
	return false
}
