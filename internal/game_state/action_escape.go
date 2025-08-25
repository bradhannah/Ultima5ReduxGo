package game_state

func (g *GameState) ActionEscapeSmallMap() bool {
	// TODO: Implement small map Escape command - see Commands.md Escape section
	// Should handle:
	// - Exit combat screen after victory
	// - Negative: "No effect!" if not in post-combat state
	// - Return to previous map location
	// - Clean up combat state
	return true
}

func (g *GameState) ActionEscapeLargeMap() bool {
	// TODO: Implement large map Escape command - see Commands.md Escape section
	// Large map variant - typically not used
	return true
}

func (g *GameState) ActionEscapeCombatMap() bool {
	// TODO: Implement combat map Escape command - see Commands.md Escape section
	// Should handle:
	// - Check if battle is won (all enemies defeated)
	// - Exit combat screen if victory achieved
	// - Negative: "Not yet!" if enemies remain
	// - Loot collection before exit
	return true
}

func (g *GameState) ActionEscapeDungeonMap() bool {
	// TODO: Implement dungeon map Escape command - see Commands.md Escape section
	// Dungeon map variant - exit from dungeon combat
	return true
}
