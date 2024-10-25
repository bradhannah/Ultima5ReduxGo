package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
	smallMap := g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor)

	pushableThingPosition := direction.GetNewPositionInDirection(&g.Position)
	pushableThingTile := smallMap.GetTileTopMapOnlyTile(pushableThingPosition)

	farSideOfPushableThingPosition := direction.GetNewPositionInDirection(pushableThingPosition)

	farSideTile := smallMap.GetTileTopMapOnlyTile(farSideOfPushableThingPosition)
	bFarSideAccessible := !(g.IsOutOfBounds(*farSideOfPushableThingPosition) || !farSideTile.IsWalkingPassable)

	// chair pushing is different
	if pushableThingTile.IsChair() {
		if bFarSideAccessible {
			// it is in bounds, so
			smallMap.SwapTiles(pushableThingPosition, farSideOfPushableThingPosition)
			smallMap.SetTile(MapOverrideLayer, farSideOfPushableThingPosition, g.GameReferences.TileReferences.GetChairByPushDirection(direction).Index)
			// move avatar to the swapped spot
		} else {
			// turn it around
			smallMap.SwapTiles(&g.Position, pushableThingPosition)
			smallMap.SetTile(MapOverrideLayer, &g.Position, g.GameReferences.TileReferences.GetChairByPushDirection(direction.GetOppositeDirection()).Index)

		}
		// move avatar to the swapped spot
		g.Position = *pushableThingPosition
		return true
	}

	// if the pushableThingPosition is out of bounds OR blocked
	if !bFarSideAccessible {
		// just swap
		smallMap.SwapTiles(&g.Position, pushableThingPosition)
		// move avatar to the swapped spot
		g.Position = *pushableThingPosition
		return true
	}
	// it is in bounds, so
	smallMap.SwapTiles(pushableThingPosition, farSideOfPushableThingPosition)
	// move avatar to the swapped spot
	g.Position = *pushableThingPosition
	return true
}
