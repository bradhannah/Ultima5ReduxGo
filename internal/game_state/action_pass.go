package game_state

// ActionPass handles the Pass command - player chooses to do nothing this turn
func (g *GameState) ActionPass() bool {
	g.SystemCallbacks.Message.AddRowStr("Pass")
	// Pass advances time as the player is choosing to do nothing this turn
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return true
}
