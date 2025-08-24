package game_state

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

const (
	smallMapStartingPositionX     = 15
	smallMapStartingPositionY     = 30
	smallMapStartingPositionFloor = 0
)

// ActionEnterLargeMap handles entering buildings from large map - includes location finding
func (g *GameState) ActionEnterLargeMap() bool {
	newLocation := g.GameReferences.LocationReferences.WorldLocations.GetLocationByPosition(
		g.MapState.PlayerLocation.Position)

	if newLocation != references2.EmptyLocation {
		slr := g.GameReferences.LocationReferences.GetLocationReference(newLocation)
		if g.ActionEnter(slr) {
			g.SystemCallbacks.Message.AddRowStr(slr.EnteringText)
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		}
	} else {
		g.SystemCallbacks.Message.AddRowStr("Enter what?")
	}
	return false
}

// ActionEnter handles entering buildings/locations - non-directional action
func (g *GameState) ActionEnter(slr *references2.SmallLocationReference) bool {
	if slr.Location == references2.EmptyLocation {
		return false
	}

	g.LastLargeMapPosition = g.MapState.PlayerLocation.Position
	g.LastLargeMapFloor = g.MapState.PlayerLocation.Floor
	g.MapState.PlayerLocation.Position = references2.Position{
		X: smallMapStartingPositionX,
		Y: smallMapStartingPositionY,
	}
	g.MapState.PlayerLocation.Location = slr.Location
	g.MapState.PlayerLocation.Floor = smallMapStartingPositionFloor
	g.UpdateSmallMap(g.GameReferences.TileReferences, g.GameReferences.LocationReferences)
	return true
}

// EnterBuilding - deprecated, use ActionEnter instead
// TODO: Remove once all callers are updated
func (g *GameState) EnterBuilding(slr *references2.SmallLocationReference) {
	g.ActionEnter(slr)
}
