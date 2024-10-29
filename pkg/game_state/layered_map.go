package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"log"
)

const totalLayers = 5
const overflowTiles = references.Coordinate(10)

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
	layers [totalLayers]Layer

	visibleFlags         VisibilityCoords
	testForVisibilityMap VisibilityCoords

	tileRefs *references.Tiles

	xVisibleTiles, yVisibleTiles int
	topLeft, bottomRight         references.Position
	xMax, yMax                   references.Coordinate

	bWrappingMap bool
}

func newLayeredMap(xMax references.Coordinate, yMax references.Coordinate, tileRefs *references.Tiles, xVisibleTiles int, yVisibleTiles int, bWrappingMap bool) *LayeredMap {
	layeredMap := LayeredMap{}
	layeredMap.xVisibleTiles = xVisibleTiles
	layeredMap.yVisibleTiles = yVisibleTiles
	layeredMap.xMax = xMax
	layeredMap.yMax = yMax
	layeredMap.bWrappingMap = bWrappingMap
	layeredMap.tileRefs = tileRefs

	for mapLayer, _ := range layeredMap.layers {
		layeredMap.visibleFlags = make(map[references.Coordinate]map[references.Coordinate]bool)
		layeredMap.testForVisibilityMap = make(map[references.Coordinate]map[references.Coordinate]bool)
		layeredMap.layers[mapLayer] = make(map[references.Coordinate]map[references.Coordinate]indexes.SpriteIndex)
		for yRow := -overflowTiles - 1; yRow < yMax+overflowTiles; yRow++ {
			layeredMap.visibleFlags[yRow] = make(map[references.Coordinate]bool)
			layeredMap.testForVisibilityMap[yRow] = make(map[references.Coordinate]bool)
			layeredMap.layers[mapLayer][yRow] = make(map[references.Coordinate]indexes.SpriteIndex)
		}
	}

	return &layeredMap
}

func (l *LayeredMap) RecalculateVisibleTiles(avatarPos references.Position) {
	l.topLeft = references.Position{
		X: avatarPos.X - overflowTiles,
		Y: avatarPos.Y - overflowTiles,
	}
	l.bottomRight = references.Position{
		X: avatarPos.X + overflowTiles,
		Y: avatarPos.Y + overflowTiles,
	}

	// TODO: it's lazy to make both of these calls since it could do in one pass
	l.testForVisibilityMap.ResetVisibilityCoords(false)
	l.testForVisibilityMap.SetVisibilityCoordsRectangle(l.topLeft, l.bottomRight, l.xMax, l.yMax, l.bWrappingMap)
	l.visibleFlags.ResetVisibilityCoords(false)

	l.SetVisible(true, &avatarPos)

	l.setTilesAroundPositionVisible(&avatarPos)

	l.floodFillIfInside(&avatarPos, true)
	l.floodFillIfInside(avatarPos.GetPositionToLeft(), true)
	l.floodFillIfInside(avatarPos.GetPositionToRight(), true)
	l.floodFillIfInside(avatarPos.GetPositionDown(), true)
	l.floodFillIfInside(avatarPos.GetPositionUp(), true)
}

func (l *LayeredMap) floodFillIfInside(pos *references.Position, bForce bool) {
	if l.bWrappingMap {
		pos = pos.GetWrapped(l.xMax, l.yMax)
	}

	if !bForce && !l.testForVisibilityMap[pos.X][pos.Y] {
		// don't reprocess and existing tile if it has been set
		return
	}

	l.SetVisible(true, pos)
	l.testForVisibilityMap[pos.X][pos.Y] = false

	tile := l.GetTileTopMapOnlyTile(pos)
	if tile != nil && tile.BlocksLight && !tile.IsWindow {
		return
	}

	l.setTilesAroundPositionVisible(pos)

	l.floodFillIfInside(pos.GetPositionToLeft(), false)
	l.floodFillIfInside(pos.GetPositionToRight(), false)
	l.floodFillIfInside(pos.GetPositionDown(), false)
	l.floodFillIfInside(pos.GetPositionUp(), false)
}

func (l *LayeredMap) setTilesAroundPositionVisible(pos *references.Position) {
	l.SetVisible(true, pos.GetPositionToLeft())
	l.SetVisible(true, pos.GetPositionToRight())
	l.SetVisible(true, pos.GetPositionDown())
	l.SetVisible(true, pos.GetPositionUp())

	// do the diagonal on the initial call
	l.SetVisible(true, pos.GetPositionUp().GetPositionToLeft())
	l.SetVisible(true, pos.GetPositionUp().GetPositionToRight())
	l.SetVisible(true, pos.GetPositionDown().GetPositionToLeft())
	l.SetVisible(true, pos.GetPositionDown().GetPositionToRight())
}

func (l *LayeredMap) SetVisible(bVisible bool, pos *references.Position) {
	l.visibleFlags[pos.X][pos.Y] = bVisible
}

func (l *LayeredMap) IsPositionVisible(pos *references.Position) bool {
	return l.visibleFlags[pos.X][pos.Y]
}

func (l *LayeredMap) GetTopTile(position *references.Position) *references.Tile {
	if position == nil {
		log.Fatal("entered a nil position so I can't get a tile")
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
