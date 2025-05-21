package map_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *MapState) SmallMapProcessTurnDoors() {
	if g.openDoorPos != nil {
		if g.openDoorTurns == 0 {
			tile := g.LayeredMaps.GetTileRefByPosition(references.SmallMapType, MapLayer, g.openDoorPos, g.PlayerLocation.Floor)
			var doorTileIndex indexes.SpriteIndex
			if tile.Index.IsWindowedDoor() {
				doorTileIndex = indexes.RegularDoorView
			} else {
				doorTileIndex = indexes.RegularDoor
			}
			g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.PlayerLocation.Floor).SetTileByLayer(MapOverrideLayer, g.openDoorPos, doorTileIndex)
			g.openDoorPos = nil
		} else {
			g.openDoorTurns--
		}
	}
}
