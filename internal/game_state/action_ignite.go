package game_state

func (g *GameState) IgniteTorch() bool {
	if g.PartyState.Inventory.Provisions.Torches == 0 {
		return false
	}
	g.PartyState.Inventory.Provisions.Torches--
	g.MapState.Lighting.LightTorch()
	return true
}
