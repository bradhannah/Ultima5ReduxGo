package game_state

func (g *GameState) ActionMixSmallMap() bool {
	// TODO: Check if party has any reagents
	// For now, just indicate no reagents until reagent system is implemented
	g.SystemCallbacks.Message.AddRowStr("None of that!")
	return false
}

func (g *GameState) ActionMixLargeMap() bool {
	// TODO: Check if party has any reagents
	// For now, just indicate no reagents until reagent system is implemented
	g.SystemCallbacks.Message.AddRowStr("None of that!")
	return false
}

func (g *GameState) ActionMixCombatMap() bool {
	// Mixing reagents during combat is typically not allowed
	g.SystemCallbacks.Message.AddRowStr("Not here!")
	return false
}

func (g *GameState) ActionMixDungeonMap() bool {
	// TODO: Check if party has any reagents
	// For now, just indicate no reagents until reagent system is implemented
	g.SystemCallbacks.Message.AddRowStr("None of that!")
	return false
}
