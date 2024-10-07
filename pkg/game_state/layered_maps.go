package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

const totalLayers = 4

type LayeredMaps struct {
	LayeredMaps map[GeneralMapType]*LayeredMap
	//SmallMapLayers LayeredMap
}

type GeneralMapType int

const (
	LargeMap GeneralMapType = iota
	SmallMap
)

type Layer int

const (
	MapLayer Layer = iota
	MapOverrideLayer
	MapUnitLayer
	EffectLayer
)

type LayeredMap struct {
	Layers   [totalLayers]map[int]map[int]int
	tileRefs *references.Tiles
}

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

func (l *LayeredMap) SetTile(layer Layer, position *references.Position, nIndex int) {
	l.Layers[layer][int(position.X)][int(position.Y)] = nIndex
}

func (l *LayeredMap) UnSetTile(layer Layer, position *references.Position) {
	l.SetTile(layer, position, -1)
}
