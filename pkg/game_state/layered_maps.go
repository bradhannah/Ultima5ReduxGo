package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

type LayeredMaps struct {
	layeredMaps map[references.GeneralMapType]map[references.FloorNumber]*LayeredMap
}

func NewLayeredMaps(tileRefs *references.Tiles) *LayeredMaps {
	maps := &LayeredMaps{}
	maps.layeredMaps = make(map[references.GeneralMapType]map[references.FloorNumber]*LayeredMap)
	maps.layeredMaps[references.LargeMapType] = make(map[references.FloorNumber]*LayeredMap)
	maps.layeredMaps[references.LargeMapType][0] = newLayeredMap(references.XLargeMapTiles, references.YLargeMapTiles, tileRefs)
	maps.layeredMaps[references.LargeMapType][-1] = newLayeredMap(references.XLargeMapTiles, references.YLargeMapTiles, tileRefs)

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

func (l *LayeredMaps) ResetAndCreateSmallMap(slr *references.SmallLocationReference, tileRefs *references.Tiles) {
	l.layeredMaps[references.SmallMapType] = make(map[references.FloorNumber]*LayeredMap)
	for _, floor := range slr.ListOfFloors {
		l.layeredMaps[references.SmallMapType][floor] = newLayeredMap(references.XSmallMapTiles, references.YSmallMapTiles, tileRefs)
	}
}

func (l *LayeredMaps) GetTileRefByPosition(mapType references.GeneralMapType, mapLayer LayerType, pos *references.Position, nFloor references.FloorNumber) *references.Tile {
	index := l.layeredMaps[mapType][nFloor].layers[mapLayer][pos.X][pos.Y]

	return l.layeredMaps[mapType][nFloor].tileRefs.GetTile(index)
}
