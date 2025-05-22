package map_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
)

func (m *MapState) SmallMapProcessTurnDoors() {
	if m.openDoorPos != nil {
		if m.openDoorTurns == 0 {
			tile := m.LayeredMaps.GetTileRefByPosition(references.SmallMapType, MapLayer, m.openDoorPos, m.PlayerLocation.Floor)
			var doorTileIndex indexes.SpriteIndex
			if tile.Index.IsWindowedDoor() {
				doorTileIndex = indexes.RegularDoorView
			} else {
				doorTileIndex = indexes.RegularDoor
			}
			m.LayeredMaps.GetLayeredMap(references.SmallMapType, m.PlayerLocation.Floor).SetTileByLayer(MapOverrideLayer, m.openDoorPos, doorTileIndex)
			m.openDoorPos = nil
		} else {
			m.openDoorTurns--
		}
	}
}
