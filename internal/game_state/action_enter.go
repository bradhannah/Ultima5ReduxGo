package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const smallMapStartingPositionX = 15
const smallMapStartingPositionY = 30
const smallMapStartingPositionFloor = 0

func (g *GameState) EnterBuilding(
	slr *references.SmallLocationReference,
	tileRefs *references.Tiles,
) {
	if slr.Location != references.EmptyLocation {
		g.LastLargeMapPosition = g.MapState.PlayerLocation.Position
		g.LastLargeMapFloor = g.MapState.PlayerLocation.Floor
		g.MapState.PlayerLocation.Position = references.Position{
			X: smallMapStartingPositionX,
			Y: smallMapStartingPositionY,
		}
		g.MapState.PlayerLocation.Location = slr.Location
		g.MapState.PlayerLocation.Floor = smallMapStartingPositionFloor
		g.UpdateSmallMap(g.GameReferences.TileReferences, g.GameReferences.LocationReferences)
	}
}
