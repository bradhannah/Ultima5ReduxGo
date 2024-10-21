package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

const totalLayers = 4

type LayeredMaps struct {
	LayeredMaps map[GeneralMapType]*LayeredMap
}

type GeneralMapType int

const (
	LargeMap GeneralMapType = iota
	SmallMap
)

type Layer int

func NewLayeredMaps(tileRefs *references.Tiles) *LayeredMaps {
	maps := &LayeredMaps{}
	maps.LayeredMaps = make(map[GeneralMapType]*LayeredMap)
	maps.LayeredMaps[LargeMap] = newLayeredMap(int(references.XLargeMapTiles), int(references.YLargeMapTiles), tileRefs)
	maps.LayeredMaps[SmallMap] = newLayeredMap(int(references.XSmallMapTiles), int(references.YSmallMapTiles), tileRefs)

	return maps
}

func GetMapTypeByLocation(location references.Location) GeneralMapType {
	if location == references.Britannia_Underworld {
		return LargeMap
	}
	return SmallMap
}

func (l *LayeredMaps) GetTileRefByPosition(mapType GeneralMapType, mapLayer Layer, pos *references.Position) *references.Tile {
	index := l.LayeredMaps[mapType].Layers[mapLayer][int(pos.X)][int(pos.Y)]

	return l.LayeredMaps[mapType].tileRefs.GetTile(index)
}
