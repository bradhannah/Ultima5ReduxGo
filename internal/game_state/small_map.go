package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/ai"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) UpdateSmallMap(tileRefs *references.Tiles, locs *references.LocationReferences) {
	slr := locs.GetLocationReference(g.MapState.PlayerLocation.Location)
	g.MapState.LayeredMaps.ResetAndCreateSmallMap(
		slr,
		tileRefs,
		g.MapState.XTilesVisibleOnGameScreen,
		g.MapState.YTilesVisibleOnGameScreen)
	g.CurrentNPCAIController = ai.NewNPCAIControllerSmallMap(slr, tileRefs, &g.MapState, &g.DateTime)
	g.CurrentNPCAIController.PopulateMapFirstLoad()
}
