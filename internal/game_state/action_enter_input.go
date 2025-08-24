package game_state

// ActionEnterInput handles the Enter key input - placeholder for confirm/interact actions
func (g *GameState) ActionEnterInput() bool {
	// TODO: Implement Enter key functionality based on context
	// This could be used for confirming selections, interacting with current position, etc.
	g.SystemCallbacks.Message.AddRowStr("Enter")
	return true
}
