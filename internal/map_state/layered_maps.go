package map_state

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type LayeredMaps struct {
	layeredMaps map[references2.GeneralMapType]map[references2.FloorNumber]*LayeredMap
}

func NewLayeredMaps(tileRefs *references2.Tiles,
	overworldRef *references2.LargeMapReference,
	underworldRef *references2.LargeMapReference,
	xTilesInMap int,
	yTilesInMap int,
) *LayeredMaps {
	maps := &LayeredMaps{}
	maps.layeredMaps = make(map[references2.GeneralMapType]map[references2.FloorNumber]*LayeredMap)
	maps.layeredMaps[references2.LargeMapType] = make(map[references2.FloorNumber]*LayeredMap)
	maps.layeredMaps[references2.LargeMapType][0] = newLayeredMap(references2.XLargeMapTiles, references2.YLargeMapTiles, tileRefs, xTilesInMap, yTilesInMap, true)
	maps.layeredMaps[references2.LargeMapType][-1] = newLayeredMap(references2.XLargeMapTiles, references2.YLargeMapTiles, tileRefs, xTilesInMap, yTilesInMap, true)

	maps.SetInitialLargeMap(0, overworldRef)
	maps.SetInitialLargeMap(-1, underworldRef)

	return maps
}

func GetMapTypeByLocation(location references2.Location) references2.GeneralMapType {
	if location == references2.Britannia_Underworld {
		return references2.LargeMapType
	}
	return references2.SmallMapType
}

func (l *LayeredMaps) GetLayeredMap(generalMapType references2.GeneralMapType, floorNumber references2.FloorNumber) *LayeredMap {
	return l.layeredMaps[generalMapType][floorNumber]
}

func (l *LayeredMaps) SetInitialLargeMap(nFloor references2.FloorNumber, lrg *references2.LargeMapReference) {
	for x := references2.Coordinate(0); x < references2.XLargeMapTiles; x++ {
		for y := references2.Coordinate(0); y < references2.YLargeMapTiles; y++ {
			l.layeredMaps[references2.LargeMapType][nFloor].SetTileByLayer(MapLayer, &references2.Position{X: x, Y: y}, lrg.GetSpriteIndex(x, y)) // slr.GetTileNumber(floor, x, y))
		}
	}
}

func (l *LayeredMaps) ResetAndCreateSmallMap(
	slr *references2.SmallLocationReference,
	tileRefs *references2.Tiles,
	xTilesInMap int,
	yTilesInMap int,
) {
	l.layeredMaps[references2.SmallMapType] = make(map[references2.FloorNumber]*LayeredMap)
	for _, floor := range slr.ListOfFloors {
		l.layeredMaps[references2.SmallMapType][floor] = newLayeredMap(references2.XSmallMapTiles, references2.YSmallMapTiles, tileRefs, xTilesInMap, yTilesInMap, false)
		theFloor := l.layeredMaps[references2.SmallMapType][floor]
		// populate the MapLayer immediately
		for x := references2.Coordinate(0); x < references2.XSmallMapTiles; x++ {
			for y := references2.Coordinate(0); y < references2.YSmallMapTiles; y++ {
				pos := references2.Position{X: x, Y: y}
				theFloor.SetTileByLayer(MapLayer, &pos, slr.GetTileNumber(pos, floor))
			}
		}
	}
}

func (l *LayeredMaps) GetTileRefByPosition(mapType references2.GeneralMapType, mapLayer LayerType, pos *references2.Position, nFloor references2.FloorNumber) *references2.Tile {
	index := l.layeredMaps[mapType][nFloor].layers[mapLayer][pos.X][pos.Y]
	return l.layeredMaps[mapType][nFloor].tileRefs.GetTile(index)
}

func (l *LayeredMaps) GetTileTopMapOnlyTileByPosition(mapType references2.GeneralMapType, pos *references2.Position, nFloor references2.FloorNumber) *references2.Tile {
	return l.GetLayeredMap(mapType, nFloor).GetTileTopMapOnlyTile(pos)
}

func (l *LayeredMaps) GetSmallMapFloorKlimbOffset(position references2.Position, nFloor references2.FloorNumber) int {
	tile := l.GetTileTopMapOnlyTileByPosition(references2.SmallMapType, &position, nFloor)
	if !tile.Index.IsStairs() {
		return 0
	}

	if l.HasLowerFloor(nFloor) {
		lowerTile := l.GetTileTopMapOnlyTileByPosition(references2.SmallMapType, &position, nFloor-1)
		if lowerTile.Index.IsStairs() {
			return -1
		}
	}
	return 1
}

func (l *LayeredMaps) HasLowerFloor(currentFloor references2.FloorNumber) bool {
	_, exists := l.layeredMaps[references2.SmallMapType][currentFloor-1]
	return exists
}
