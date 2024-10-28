package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) EnterBuilding(slr *references.SmallLocationReference, tileRefs *references.Tiles) {
	if slr.Location != references.EmptyLocation {
		g.LastLargeMapPosition = g.Position
		g.Position = references.Position{
			X: 15,
			Y: 30,
		}
		g.Location = slr.Location
		g.Floor = 0
		g.LayeredMaps.ResetAndCreateSmallMap(slr, tileRefs, g.XTilesInMap, g.YTilesInMap)
	}
}
