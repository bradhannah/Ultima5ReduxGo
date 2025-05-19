package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

func (g *GameState) UpdateSmallMap(tileRefs *references.Tiles, locs *references.LocationReferences) {
	slr := locs.GetLocationReference(g.Location)
	g.LayeredMaps.ResetAndCreateSmallMap(
		slr,
		tileRefs,
		g.XTilesVisibleOnGameScreen,
		g.YTilesVisibleOnGameScreen)
	g.CurrentNPCAIController = NewNPCAIControllerSmallMap(slr, tileRefs, g)
	g.CurrentNPCAIController.PopulateMapFirstLoad()
}
