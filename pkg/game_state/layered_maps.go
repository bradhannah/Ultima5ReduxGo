package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

const totalLayers = 4

type LayeredMaps struct {
	layeredMaps map[GeneralMapType]map[references.FloorNumber]*LayeredMap
}

type GeneralMapType int

const (
	LargeMap GeneralMapType = iota
	SmallMap
)

type Layer int

// TODO: ADD floors - this will be ideal since it will prevent us from getting locked out
// of certain areas when we run out of keys

func NewLayeredMaps(tileRefs *references.Tiles) *LayeredMaps {
	maps := &LayeredMaps{}
	maps.layeredMaps = make(map[GeneralMapType]map[references.FloorNumber]*LayeredMap)
	maps.layeredMaps[LargeMap] = make(map[references.FloorNumber]*LayeredMap)
	maps.layeredMaps[LargeMap][0] = newLayeredMap(int(references.XLargeMapTiles), int(references.YLargeMapTiles), tileRefs)
	maps.layeredMaps[LargeMap][-1] = newLayeredMap(int(references.XLargeMapTiles), int(references.YLargeMapTiles), tileRefs)

	return maps
}

func GetMapTypeByLocation(location references.Location) GeneralMapType {
	if location == references.Britannia_Underworld {
		return LargeMap
	}
	return SmallMap
}

func (l *LayeredMaps) GetLayeredMap(generalMapType GeneralMapType, floorNumber references.FloorNumber) *LayeredMap {
	return l.layeredMaps[generalMapType][floorNumber]
}

func (l *LayeredMaps) ResetAndCreateSmallMap(slr *references.SmallLocationReference, tileRefs *references.Tiles) {
	l.layeredMaps[SmallMap] = make(map[references.FloorNumber]*LayeredMap)
	for _, floor := range slr.GetListOfFloors() {
		l.layeredMaps[SmallMap][floor] = newLayeredMap(int(references.XSmallMapTiles), int(references.YSmallMapTiles), tileRefs)
	}
}

func (l *LayeredMaps) GetTileRefByPosition(mapType GeneralMapType, mapLayer Layer, pos *references.Position, nFloor references.FloorNumber) *references.Tile {
	index := l.layeredMaps[mapType][nFloor].Layers[mapLayer][int(pos.X)][int(pos.Y)]

	return l.layeredMaps[mapType][nFloor].tileRefs.GetTile(index)
}
