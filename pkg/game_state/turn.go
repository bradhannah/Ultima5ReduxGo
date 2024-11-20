package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) FinishTurn() {
	switch g.Location.GetMapType() {
	case references.SmallMapType:
		g.smallMapProcessEndOfTurn()
	case references.LargeMapType:
		g.largeMapProcessEndOfTurn()
	default:
		panic("unhandled default case")
	}
}

func (g *GameState) largeMapProcessEndOfTurn() {
	g.DateTime.Advance(DefaultLargeMapMinutesPerTurn)
}

func (g *GameState) smallMapProcessEndOfTurn() {
	if g.openDoorPos != nil {
		if g.openDoorTurns == 0 {
			tile := g.LayeredMaps.GetTileRefByPosition(references.SmallMapType, MapLayer, g.openDoorPos, g.Floor)
			if tile.Index.IsWindowedDoor() {
				g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).SetTileByLayer(MapOverrideLayer, g.openDoorPos, indexes.RegularDoorView)
			} else {
				g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).SetTileByLayer(MapOverrideLayer, g.openDoorPos, indexes.RegularDoor)
			}
			g.openDoorPos = nil
		} else {
			g.openDoorTurns--
		}
	}
	g.DateTime.Advance(DefaultSmallMapMinutesPerTurn)
}
