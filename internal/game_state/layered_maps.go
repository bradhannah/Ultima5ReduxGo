package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type LayeredMaps struct {
	layeredMaps map[references.GeneralMapType]map[references.FloorNumber]*LayeredMap
}

func NewLayeredMaps(tileRefs *references.Tiles,
	overworldRef *references.LargeMapReference,
	underworldRef *references.LargeMapReference,
	xTilesInMap int,
	yTilesInMap int,
) *LayeredMaps {
	maps := &LayeredMaps{}
	maps.layeredMaps = make(map[references.GeneralMapType]map[references.FloorNumber]*LayeredMap)
	maps.layeredMaps[references.LargeMapType] = make(map[references.FloorNumber]*LayeredMap)
	maps.layeredMaps[references.LargeMapType][0] = newLayeredMap(references.XLargeMapTiles, references.YLargeMapTiles, tileRefs, xTilesInMap, yTilesInMap, true)
	maps.layeredMaps[references.LargeMapType][-1] = newLayeredMap(references.XLargeMapTiles, references.YLargeMapTiles, tileRefs, xTilesInMap, yTilesInMap, true)

	maps.SetInitialLargeMap(0, overworldRef)
	maps.SetInitialLargeMap(-1, underworldRef)

	return maps
}

func GetMapTypeByLocation(location references.Location) references.GeneralMapType {
	if location == references.Britannia_Underworld {
		return references.LargeMapType
	}
	return references.SmallMapType
}

func (l *LayeredMaps) GetLayeredMap(generalMapType references.GeneralMapType, floorNumber references.FloorNumber) *LayeredMap {
	return l.layeredMaps[generalMapType][floorNumber]
}

func (l *LayeredMaps) SetInitialLargeMap(nFloor references.FloorNumber, lrg *references.LargeMapReference) {
	for x := references.Coordinate(0); x < references.XLargeMapTiles; x++ {
		for y := references.Coordinate(0); y < references.YLargeMapTiles; y++ {
			l.layeredMaps[references.LargeMapType][nFloor].SetTileByLayer(MapLayer, &references.Position{X: x, Y: y}, lrg.GetSpriteIndex(x, y)) // slr.GetTileNumber(floor, x, y))
		}
	}
}

func (l *LayeredMaps) ResetAndCreateSmallMap(
	slr *references.SmallLocationReference,
	tileRefs *references.Tiles,
	xTilesInMap int,
	yTilesInMap int,
) {
	l.layeredMaps[references.SmallMapType] = make(map[references.FloorNumber]*LayeredMap)
	for _, floor := range slr.ListOfFloors {
		l.layeredMaps[references.SmallMapType][floor] = newLayeredMap(references.XSmallMapTiles, references.YSmallMapTiles, tileRefs, xTilesInMap, yTilesInMap, false)
		theFloor := l.layeredMaps[references.SmallMapType][floor]
		// populate the MapLayer immediately
		for x := references.Coordinate(0); x < references.XSmallMapTiles; x++ {
			for y := references.Coordinate(0); y < references.YSmallMapTiles; y++ {
				pos := references.Position{X: x, Y: y}
				theFloor.SetTileByLayer(MapLayer, &pos, slr.GetTileNumber(pos, floor))
			}
		}
	}
}

func (l *LayeredMaps) GetTileRefByPosition(mapType references.GeneralMapType, mapLayer LayerType, pos *references.Position, nFloor references.FloorNumber) *references.Tile {
	index := l.layeredMaps[mapType][nFloor].layers[mapLayer][pos.X][pos.Y]
	return l.layeredMaps[mapType][nFloor].tileRefs.GetTile(index)
}

func (l *LayeredMaps) GetTileTopMapOnlyTileByPosition(mapType references.GeneralMapType, pos *references.Position, nFloor references.FloorNumber) *references.Tile {
	return l.GetLayeredMap(mapType, nFloor).GetTileTopMapOnlyTile(pos)
}

func (l *LayeredMaps) GetSmallMapFloorKlimbOffset(position references.Position, nFloor references.FloorNumber) int {
	tile := l.GetTileTopMapOnlyTileByPosition(references.SmallMapType, &position, nFloor)
	if !tile.Index.IsStairs() {
		return 0
	}

	if l.HasLowerFloor(nFloor) {
		lowerTile := l.GetTileTopMapOnlyTileByPosition(references.SmallMapType, &position, nFloor-1)
		if lowerTile.Index.IsStairs() {
			return -1
		}
	}
	return 1
}

func (l *LayeredMaps) HasLowerFloor(currentFloor references.FloorNumber) bool {
	_, exists := l.layeredMaps[references.SmallMapType][currentFloor-1]
	return exists
}
