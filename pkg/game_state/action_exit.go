package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

func (g *GameState) ExitSmallMap() {
	g.Location = references.Britannia_Underworld
	g.Floor = g.LastLargeMapFloor
	g.Position = g.LastLargeMapPosition
	g.CurrentNPCAIController = g.GetCurrentLargeMapNPCAIController()
	g.UpdateLargeMap()
}
