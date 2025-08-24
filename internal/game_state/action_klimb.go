package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionKlimbSmallMap(direction references.Direction) bool {
	// Check current position for ladder/grate first (player stands on these)
	currentTile := g.MapState.LayeredMaps.GetTileRefByPosition(
		references.SmallMapType,
		map_state.MapLayer,
		&g.MapState.PlayerLocation.Position,
		g.MapState.PlayerLocation.Floor)

	switch currentTile.Index {
	case indexes.AvatarOnLadderDown, indexes.LadderDown, indexes.Grate:
		locationRef := g.GameReferences.LocationReferences.GetLocationReference(g.MapState.PlayerLocation.Location)
		if locationRef.CanGoDownOneFloor(g.MapState.PlayerLocation.Floor) {
			g.MapState.PlayerLocation.Floor--
			g.UpdateSmallMap(g.GameReferences.TileReferences, g.GameReferences.LocationReferences)
			g.SystemCallbacks.Message.AddRowStr("Klimb-Down!")
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		} else {
			g.SystemCallbacks.Message.AddRowStr("Can't go lower!")
			return false
		}

	case indexes.AvatarOnLadderUp, indexes.LadderUp:
		locationRef := g.GameReferences.LocationReferences.GetLocationReference(g.MapState.PlayerLocation.Location)
		if locationRef.CanGoUpOneFloor(g.MapState.PlayerLocation.Floor) {
			g.MapState.PlayerLocation.Floor++
			g.UpdateSmallMap(g.GameReferences.TileReferences, g.GameReferences.LocationReferences)
			g.SystemCallbacks.Message.AddRowStr("Klimb-Up!")
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		} else {
			g.SystemCallbacks.Message.AddRowStr("Can't go higher!")
			return false
		}
	}

	// Check target position for fence climbing
	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	targetTile := g.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.MapState.PlayerLocation.Floor).GetTileTopMapOnlyTile(newPosition)
	if targetTile.Index == indexes.FenceHoriz || targetTile.Index == indexes.FenceVert {
		g.MapState.PlayerLocation.Position = *newPosition
		g.SystemCallbacks.Message.AddRowStr("Klimb!")
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true
	}

	g.SystemCallbacks.Message.AddRowStr("What?")
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
