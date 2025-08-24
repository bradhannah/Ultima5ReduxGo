package game_state

// ActionIgnite handles torch ignition - non-directional action
func (g *GameState) ActionIgnite() bool {
	if !g.PartyState.Inventory.Provisions.Torches.DecrementByOne() {
		g.SystemCallbacks.Message.AddRowStr("None owned!")
		return false
	}
	g.MapState.Lighting.LightTorch()
	g.SystemCallbacks.Message.AddRowStr("Ignite Torch!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return true
}

// IgniteTorch - deprecated, use ActionIgnite instead
// TODO: Remove once all callers are updated
func (g *GameState) IgniteTorch() bool {
	return g.ActionIgnite()
}
