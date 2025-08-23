package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
	smallMap := g.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.MapState.PlayerLocation.Floor)

	pushableThingPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	pushableThingTile := smallMap.GetTopTile(pushableThingPosition)

	farSideOfPushableThingPosition := direction.GetNewPositionInDirection(pushableThingPosition)

	farSideTile := smallMap.GetTopTile(farSideOfPushableThingPosition)
	bFarSideAccessible := !g.IsOutOfBounds(*farSideOfPushableThingPosition) && farSideTile.Index.IsPushableFloor()

	// chair pushing is different
	if pushableThingTile.IsChair() {
		if bFarSideAccessible {
			smallMap.SwapTiles(pushableThingPosition, farSideOfPushableThingPosition)
			smallMap.SetTileByLayer(map_state.MapLayer, farSideOfPushableThingPosition, g.GameReferences.TileReferences.GetChairByPushDirection(direction).Index)
		} else {
			smallMap.SwapTiles(&g.MapState.PlayerLocation.Position, pushableThingPosition)
			smallMap.SetTileByLayer(map_state.MapLayer, &g.MapState.PlayerLocation.Position, g.GameReferences.TileReferences.GetChairByPushDirection(direction.GetOppositeDirection()).Index)
		}
		// move avatar to the swapped spot
		g.MapState.PlayerLocation.Position = *pushableThingPosition
		return true
	} else if pushableThingTile.IsCannon() {
		if bFarSideAccessible {
			smallMap.SwapTiles(pushableThingPosition, farSideOfPushableThingPosition)
			smallMap.SetTileByLayer(map_state.MapLayer, farSideOfPushableThingPosition, g.GameReferences.TileReferences.GetCannonByPushDirection(direction).Index)
		} else {
			smallMap.SwapTiles(&g.MapState.PlayerLocation.Position, pushableThingPosition)
			smallMap.SetTileByLayer(map_state.MapLayer, &g.MapState.PlayerLocation.Position, g.GameReferences.TileReferences.GetCannonByPushDirection(direction.GetOppositeDirection()).Index)
		}
		g.MapState.PlayerLocation.Position = *pushableThingPosition
		return true
	}

	if !bFarSideAccessible {
		smallMap.SwapTiles(&g.MapState.PlayerLocation.Position, pushableThingPosition)
	} else {
		smallMap.SwapTiles(pushableThingPosition, farSideOfPushableThingPosition)
	}
	// move avatar to the swapped spot
	g.MapState.PlayerLocation.Position = *pushableThingPosition
	return true
}
