package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

func (g *GameState) UpdateSmallMap(tileRefs *references.Tiles, locs *references.LocationReferences) {
	slr := locs.GetLocationReference(g.Location)
	g.LayeredMaps.ResetAndCreateSmallMap(
		slr,
		tileRefs,
		// npcRefs,
		g.XTilesInMap,
		g.YTilesInMap)
	g.NPCAIController = *NewNPCAIController(slr, tileRefs, g)
	g.NPCAIController.PopulateMapFirstLoad(
		g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor),
		g.DateTime,
	)
}