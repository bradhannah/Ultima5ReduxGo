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
	topTile := g.GetCurrentLayeredMapAvatarTopTile()

	// we care about speed factor only for large maps
	g.DateTime.Advance(topTile.SpeedFactor)
}

func (g *GameState) smallMapProcessEndOfTurn() {
	g.smallMapProcessTurnDoors()
	// small maps are always 1 minute per turn
	g.DateTime.Advance(DefaultSmallMapMinutesPerTurn)

	g.smallMapProcessNPCs()
}

func (g *GameState) smallMapProcessNPCs() {
	g.CurrentNPCAIController.CalculateNextNPCPositions()
}

func (g *GameState) smallMapProcessTurnDoors() {
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
}
