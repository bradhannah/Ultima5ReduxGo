package game_state

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const totalLayers = 5

const (
	MapLayer LayerType = iota
	MapOverrideLayer
	AvatarAndPartyLayer
	MapUnitLayer
	EffectLayer
)

type LayerType int

type Layer map[references.Coordinate]map[references.Coordinate]indexes.SpriteIndex

type LayeredMap struct {
	layers   [totalLayers]map[references.Coordinate]map[references.Coordinate]indexes.SpriteIndex
	tileRefs *references.Tiles
}

func newLayeredMap(xMax references.Coordinate, yMax references.Coordinate, tileRefs *references.Tiles) *LayeredMap {
	const overflowTiles = references.Coordinate(10)
	layeredMap := LayeredMap{}
	layeredMap.tileRefs = tileRefs
	for mapLayer, _ := range layeredMap.layers {
		layeredMap.layers[mapLayer] = make(map[references.Coordinate]map[references.Coordinate]indexes.SpriteIndex)
		for yRow := -overflowTiles; yRow < yMax+10; yRow++ {
			layeredMap.layers[mapLayer][yRow] = make(map[references.Coordinate]indexes.SpriteIndex)
		}
	}
	return &layeredMap
}

func (l *LayeredMap) GetTopTile(position *references.Position) *references.Tile {
	if position == nil {
		fmt.Sprintf("oof")
	}
	for i := EffectLayer; i >= MapLayer; i-- {
		tile := l.GetTileByLayer(i, position)
		if tile.Index <= 0 {
			continue
		}
		return tile
	}
	return nil
}

func (l *LayeredMap) GetTileTopMapOnlyTile(position *references.Position) *references.Tile {
	for i := MapOverrideLayer; i >= MapLayer; i-- {
		tile := l.GetTileByLayer(i, position)
		if tile.Index <= 0 {
			continue
		}
		return tile
	}
	return nil
}

func (l *LayeredMap) SetTileByLayer(layer LayerType, position *references.Position, nIndex indexes.SpriteIndex) {
	l.layers[layer][position.X][position.Y] = nIndex
}

func (l *LayeredMap) UnSetTileByLayer(layer LayerType, position *references.Position) {
	l.SetTileByLayer(layer, position, -1)
}

func (l *LayeredMap) GetTileByLayer(layer LayerType, position *references.Position) *references.Tile {
	return l.tileRefs.GetTile(l.layers[layer][position.X][position.Y])
}

func (l *LayeredMap) SwapTiles(pos1 *references.Position, pos2 *references.Position) {
	tile1 := l.GetTileTopMapOnlyTile(pos1)
	tile2 := l.GetTileTopMapOnlyTile(pos2)
	l.SetTileByLayer(MapLayer, pos1, tile2.Index)
	l.SetTileByLayer(MapLayer, pos2, tile1.Index)
}
