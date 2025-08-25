package game_state

func (g *GameState) ActionViewSmallMap() bool {
	if !g.PartyState.Inventory.Provisions.Gems.HasSome() {
		g.SystemCallbacks.Message.AddRowStr("You have none!")
		return false
	}

	g.PartyState.Inventory.Provisions.Gems.DecrementByOne()
	// TODO: Implement small map overhead view display
	g.SystemCallbacks.Message.AddRowStr("View area!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return true
}

func (g *GameState) ActionViewLargeMap() bool {
	if !g.PartyState.Inventory.Provisions.Gems.HasSome() {
		g.SystemCallbacks.Message.AddRowStr("You have none!")
		return false
	}

	g.PartyState.Inventory.Provisions.Gems.DecrementByOne()
	// TODO: Implement overworld area view display
	g.SystemCallbacks.Message.AddRowStr("View area!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return true
}

func (g *GameState) ActionViewCombatMap() bool {
	if !g.PartyState.Inventory.Provisions.Gems.HasSome() {
		g.SystemCallbacks.Message.AddRowStr("You have none!")
		return false
	}

	g.PartyState.Inventory.Provisions.Gems.DecrementByOne()
	// TODO: Implement combat map tactical overview display
	g.SystemCallbacks.Message.AddRowStr("View area!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return true
}

func (g *GameState) ActionViewDungeonMap() bool {
	if !g.PartyState.Inventory.Provisions.Gems.HasSome() {
		g.SystemCallbacks.Message.AddRowStr("You have none!")
		return false
	}

	g.PartyState.Inventory.Provisions.Gems.DecrementByOne()
	// TODO: Implement dungeon level cell layout rendering
	g.SystemCallbacks.Message.AddRowStr("View dungeon!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return true
}
