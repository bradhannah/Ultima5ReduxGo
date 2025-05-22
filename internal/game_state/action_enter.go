package game_state

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

const smallMapStartingPositionX = 15
const smallMapStartingPositionY = 30
const smallMapStartingPositionFloor = 0

func (g *GameState) EnterBuilding(slr *references2.SmallLocationReference) {
	if slr.Location != references2.EmptyLocation {
		g.LastLargeMapPosition = g.MapState.PlayerLocation.Position
		g.LastLargeMapFloor = g.MapState.PlayerLocation.Floor
		g.MapState.PlayerLocation.Position = references2.Position{
			X: smallMapStartingPositionX,
			Y: smallMapStartingPositionY,
		}
		g.MapState.PlayerLocation.Location = slr.Location
		g.MapState.PlayerLocation.Floor = smallMapStartingPositionFloor
		g.UpdateSmallMap(g.GameReferences.TileReferences, g.GameReferences.LocationReferences)
	}
}
