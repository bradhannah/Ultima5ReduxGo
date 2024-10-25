package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) KlimbSmallMap(direction references.Direction) bool {
	newPosition := direction.GetNewPositionInDirection(&g.Position)
	targetTile := g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).GetTileTopMapOnlyTile(newPosition)
	if targetTile.Index == indexes.FenceHoriz || targetTile.Index == indexes.FenceVert {
		g.Position = *newPosition
		return true
	}
	return false
}

func (g *GameState) KlimbLargeMap() bool {
	return true
}
