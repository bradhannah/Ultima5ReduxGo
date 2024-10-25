package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

func (g *GameState) ActionPushSmallMap(direction Direction) bool {
	pushableThingPosition := direction.GetNewPositionInDirection(&g.Position)
	farSideOfPushableThingPosition := direction.GetNewPositionInDirection(pushableThingPosition)

	// if the pushableThingPosition is out of bounds OR blocked
	if g.IsOutOfBounds(*farSideOfPushableThingPosition) || !g.LayeredMaps.GetTileTopMapOnlyTileByPosition(references.SmallMapType, farSideOfPushableThingPosition, g.Floor).IsWalkingPassable {
		// just swap
		g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).SwapTiles(&g.Position, pushableThingPosition)
		// move avatar to the swapped spot
		g.Position = *pushableThingPosition
		return true
	}
	// it is in bounds, so
	g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).SwapTiles(pushableThingPosition, farSideOfPushableThingPosition)
	// move avatar to the swapped spot
	g.Position = *pushableThingPosition
	return true
}
