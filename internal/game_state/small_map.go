package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/ai"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) UpdateSmallMap(tileRefs *references2.Tiles, locs *references2.LocationReferences) {
	slr := locs.GetLocationReference(g.MapState.PlayerLocation.Location)
	g.MapState.LayeredMaps.ResetAndCreateSmallMap(
		slr,
		tileRefs,
		g.MapState.XTilesVisibleOnGameScreen,
		g.MapState.YTilesVisibleOnGameScreen)
	g.CurrentNPCAIController = ai.NewNPCAIControllerSmallMap(slr, tileRefs, &g.MapState, &g.DateTime, g)
	g.CurrentNPCAIController.PopulateMapFirstLoad()
}
