package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const (
	MapLayer Layer = iota
	MapOverrideLayer
	MapUnitLayer
	EffectLayer
)

type LayeredMap struct {
	Layers   [totalLayers]map[int]map[int]indexes.SpriteIndex
	tileRefs *references.Tiles
}

func newLayeredMap(xMax int, yMax int, tileRefs *references.Tiles) *LayeredMap {
	const overflowTiles = 10
	layeredMap := LayeredMap{}
	layeredMap.tileRefs = tileRefs
	for i, _ := range layeredMap.Layers {
		layeredMap.Layers[i] = make(map[int]map[int]indexes.SpriteIndex)
		for j := -overflowTiles; j < yMax+10; j++ {
			layeredMap.Layers[i][j] = make(map[int]indexes.SpriteIndex)
		}
	}
	return &layeredMap
}

func (l *LayeredMap) GetTopTile(position *references.Position) *references.Tile {
	var tileIndex indexes.SpriteIndex
	for i := EffectLayer; i >= MapLayer; i-- {
		tileIndex = l.Layers[i][int(position.X)][int(position.Y)]
		if tileIndex <= 0 {
			continue
		}
		return l.tileRefs.GetTile(tileIndex)
	}
	return nil
}

func (l *LayeredMap) GetTileTopMapOnlyTile(position *references.Position) *references.Tile {
	var tileValue indexes.SpriteIndex
	for i := MapOverrideLayer; i >= MapLayer; i-- {
		tileValue = l.Layers[i][int(position.X)][int(position.Y)]
		if tileValue <= 0 {
			continue
		}
		return l.tileRefs.GetTile(tileValue)
	}
	return nil
}

func (l *LayeredMap) SetTile(layer Layer, position *references.Position, nIndex indexes.SpriteIndex) {
	l.Layers[layer][int(position.X)][int(position.Y)] = nIndex
}

func (l *LayeredMap) UnSetTile(layer Layer, position *references.Position) {
	l.SetTile(layer, position, -1)
}
