package game_state

func (g *GameState) ActionZtatsSmallMap() bool {
	// TODO: Implement small map Ztats command - see Commands.md Ztats (Party Member Stats) section
	// Should handle:
	// - Character selection (1-6 or by name)
	// - Display full character statistics
	// - HP/MP current and max values
	// - Equipment readied in each slot
	// - Experience and level information
	// - Status effects (poisoned, sleeping, dead, etc.)
	return true
}

func (g *GameState) ActionZtatsLargeMap() bool {
	// TODO: Implement large map Ztats command - see Commands.md Ztats (Party Member Stats) section
	// Large map variant of ztats command
	return true
}

func (g *GameState) ActionZtatsCombatMap() bool {
	// TODO: Implement combat map Ztats command - see Commands.md Ztats (Party Member Stats) section
	// Combat map variant - quick stats check during battle
	return true
}

func (g *GameState) ActionZtatsDungeonMap() bool {
	// TODO: Implement dungeon map Ztats command - see Commands.md Ztats (Party Member Stats) section
	// Dungeon map variant of ztats command
	return true
}
