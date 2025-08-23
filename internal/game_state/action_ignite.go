package game_state

// ActionIgnite handles torch ignition - non-directional action
func (g *GameState) ActionIgnite() bool {
	if !g.PartyState.Inventory.Provisions.Torches.DecrementByOne() {
		return false
	}
	g.MapState.Lighting.LightTorch()
	return true
}

// IgniteTorch - deprecated, use ActionIgnite instead
// TODO: Remove once all callers are updated
func (g *GameState) IgniteTorch() bool {
	return g.ActionIgnite()
}
