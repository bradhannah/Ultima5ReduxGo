package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

const totalLayers = 4

type LayeredMaps struct {
	LargeMapLayers LayeredMap
	SmallMapLayers LayeredMap
}

type LayeredMap struct {
	Layers [totalLayers][][]int
}

func NewLayeredMaps() *LayeredMaps {
	maps := &LayeredMaps{}
	maps.LargeMapLayers = newLayeredMap(int(references.XLargeMapTiles), int(references.YLargeMapTiles))
	maps.SmallMapLayers = newLayeredMap(int(references.XSmallMapTiles), int(references.YSmallMapTiles))
	return maps
}

func newLayeredMap(xMax int, yMax int) LayeredMap {
	layeredMap := LayeredMap{}
	for i, _ := range layeredMap.Layers {
		layeredMap.Layers[i] = make([][]int, xMax)
		for row, _ := range layeredMap.Layers[i] {
			layeredMap.Layers[i][row] = make([]int, yMax)
		}
	}
	return layeredMap
}
