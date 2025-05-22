package game_state

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionKlimbSmallMap(direction references2.Direction) bool {
	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	targetTile := g.MapState.LayeredMaps.GetLayeredMap(references2.SmallMapType, g.MapState.PlayerLocation.Floor).GetTileTopMapOnlyTile(newPosition)
	if targetTile.Index == indexes.FenceHoriz || targetTile.Index == indexes.FenceVert {
		g.MapState.PlayerLocation.Position = *newPosition
		return true
	}
	return false
}

func (g *GameState) KlimbLargeMap() bool {
	return true
}
