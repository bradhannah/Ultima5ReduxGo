package game_state

func (g *GameState) ActionReadySmallMap() bool {
	// TODO: Implement small map Ready command - see Commands.md Ready section
	// Should handle:
	// - Equipment slot management (2 hands, 1 armor, 1 head, 3 finger rings, 1 amulet)
	// - Weight limit enforcement
	// - Ammo requirement checks for ranged weapons
	// - Negative: "No armor change!" in combat
	// - Negative: "Missing ammo" for bows/xbows without arrows/bolts
	// - Character selection and inventory management
	return true
}

func (g *GameState) ActionReadyLargeMap() bool {
	// TODO: Implement large map Ready command - see Commands.md Ready section
	// Large map variant of ready command
	return true
}

func (g *GameState) ActionReadyCombatMap() bool {
	// Armor cannot be changed during combat, but weapons can be readied
	// TODO: Implement weapon readying logic
	g.SystemCallbacks.Message.AddRowStr("Ready weapon...")
	return false
}

func (g *GameState) ActionReadyDungeonMap() bool {
	// TODO: Implement dungeon map Ready command - see Commands.md Ready section
	// Dungeon map variant of ready command
	return true
}
