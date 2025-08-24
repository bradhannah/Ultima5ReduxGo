package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionLookSmallMap(direction references.Direction) bool {
	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	topTile := g.GetLayeredMapByCurrentLocation().GetTopTile(newPosition)
	lookDescription := g.GameReferences.LookReferences.GetTileLookDescription(topTile.Index)
	g.SystemCallbacks.Message.AddRowStr("Thou dost see " + lookDescription)

	// Special cases for interactive tiles
	switch topTile.Index {
	case indexes.Clock1, indexes.Clock2:
		g.SystemCallbacks.Message.AppendToCurrentRowStr(g.DateTime.GetTimeAsString())
	}

	g.SystemCallbacks.Flow.AdvanceTime(1) // Looking takes time
	return true
}

func (g *GameState) ActionLookLargeMap(direction references.Direction) bool {
	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	topTile := g.GetLayeredMapByCurrentLocation().GetTopTile(newPosition)
	lookDescription := g.GameReferences.LookReferences.GetTileLookDescription(topTile.Index)
	g.SystemCallbacks.Message.AddRowStr("Thou dost see " + lookDescription)

	// Special cases for interactive tiles
	switch topTile.Index {
	case indexes.Clock1, indexes.Clock2:
		g.SystemCallbacks.Message.AppendToCurrentRowStr(g.DateTime.GetTimeAsString())
	}

	g.SystemCallbacks.Flow.AdvanceTime(1) // Looking takes time
	return true
}

func (g *GameState) ActionLookCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Look command - see Commands.md Look section
	// Combat map variant of look command
	return true
}
