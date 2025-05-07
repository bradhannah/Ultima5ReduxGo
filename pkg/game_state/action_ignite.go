package game_state

func (g *GameState) IgniteTorch() bool {
	if g.Inventory.Provisions.Torches == 0 {
		return false
	}
	g.Inventory.Provisions.Torches--
	g.Lighting.LightTorch()
	return true
}
