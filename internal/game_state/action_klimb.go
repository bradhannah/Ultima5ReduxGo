package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionKlimbSmallMap(direction references.Direction) bool {
	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	targetTile := g.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.MapState.PlayerLocation.Floor).GetTileTopMapOnlyTile(newPosition)
	if targetTile.Index == indexes.FenceHoriz || targetTile.Index == indexes.FenceVert {
		g.MapState.PlayerLocation.Position = *newPosition
		return true
	}
	return false
}

func (g *GameState) ActionKlimbLargeMap(direction references.Direction) bool {
	return true
}

func (g *GameState) ActionKlimbCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Klimb command - see Commands.md Klimb section
	// Combat map variant of klimb command - likely limited functionality
	return true
}
