package game_state

func (g *GameState) ActionCastSmallMap() bool {
	// TODO: Check if party has any spell mixtures
	// For now, just indicate no spells are mixed until spell system is implemented
	g.SystemCallbacks.Message.AddRowStr("None mixed!")
	return false
}

func (g *GameState) ActionCastLargeMap() bool {
	// TODO: Check if party has any spell mixtures
	// For now, just indicate no spells are mixed until spell system is implemented
	g.SystemCallbacks.Message.AddRowStr("None mixed!")
	return false
}

func (g *GameState) ActionCastCombatMap() bool {
	// TODO: Check if party has any spell mixtures
	// For now, just indicate no spells are mixed until spell system is implemented
	g.SystemCallbacks.Message.AddRowStr("None mixed!")
	return false
}

func (g *GameState) ActionCastDungeonMap() bool {
	// TODO: Check if party has any spell mixtures
	// For now, just indicate no spells are mixed until spell system is implemented
	g.SystemCallbacks.Message.AddRowStr("None mixed!")
	return false
}
