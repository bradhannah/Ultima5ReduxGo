package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
	smallMap := g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor)

	pushableThingPosition := direction.GetNewPositionInDirection(&g.Position)
	pushableThingTile := smallMap.GetTileTopMapOnlyTile(pushableThingPosition)

	farSideOfPushableThingPosition := direction.GetNewPositionInDirection(pushableThingPosition)

	farSideTile := smallMap.GetTileTopMapOnlyTile(farSideOfPushableThingPosition)
	bFarSideAccessible := !g.IsOutOfBounds(*farSideOfPushableThingPosition) && farSideTile.Index == indexes.BrickFloor

	// chair pushing is different
	if pushableThingTile.IsChair() {
		if bFarSideAccessible {
			smallMap.SwapTiles(pushableThingPosition, farSideOfPushableThingPosition)
			smallMap.SetTileByLayer(MapLayer, farSideOfPushableThingPosition, g.GameReferences.TileReferences.GetChairByPushDirection(direction).Index)
		} else {
			smallMap.SwapTiles(&g.Position, pushableThingPosition)
			smallMap.SetTileByLayer(MapLayer, &g.Position, g.GameReferences.TileReferences.GetChairByPushDirection(direction.GetOppositeDirection()).Index)

		}
		// move avatar to the swapped spot
		g.Position = *pushableThingPosition
		return true
	}

	if !bFarSideAccessible {
		smallMap.SwapTiles(&g.Position, pushableThingPosition)
	} else {
		smallMap.SwapTiles(pushableThingPosition, farSideOfPushableThingPosition)
	}
	// move avatar to the swapped spot
	g.Position = *pushableThingPosition
	return true
}