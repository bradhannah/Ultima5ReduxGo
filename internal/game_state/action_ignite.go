package game_state

func (g *GameState) IgniteTorch() bool {
	//if g.PartyState.Inventory.Provisions.Torches.HasSome() {
	//	return false
	//}
	if !g.PartyState.Inventory.Provisions.Torches.DecrementByOne() {
		return false
	}
	g.MapState.Lighting.LightTorch()
	return true
}
