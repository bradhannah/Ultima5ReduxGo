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
		g.LastLargeMapPosition = g.Position
		g.LastLargeMapFloor = g.Floor
		g.Position = references.Position{
			X: smallMapStartingPositionX,
			Y: smallMapStartingPositionY,
		}
		g.Location = slr.Location
		g.Floor = smallMapStartingPositionFloor
		g.UpdateSmallMap(g.GameReferences.TileReferences, g.GameReferences.LocationReferences)
	}
}
