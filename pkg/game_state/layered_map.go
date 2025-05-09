package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const totalLayers = 6
const overflowTiles = references.Coordinate(10)

const (
	MapLayer LayerType = iota
	MapOverrideLayer
	AvatarAndPartyLayer
	MapUnitLayer
	EquipmentAndProvisionsLayer
	EffectLayer
)

type LayerType int

type Layer map[references.Coordinate]map[references.Coordinate]indexes.SpriteIndex

type LayeredMap struct {
	layers [totalLayers]Layer

	visibleFlags         VisibilityCoords
	testForVisibilityMap VisibilityCoords
	distanceMaskMaps     []DistanceMaskMap
	// lighting *Lighting

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

	for mapLayer := range layeredMap.layers {
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

func (l *LayeredMap) RecalculateVisibleTiles(avatarPos references.Position, lighting *Lighting, timeOfDay datetime.UltimaDate) {
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
	l.testForVisibilityMap.SetVisibilityCoordsRectangle(&l.topLeft, &l.bottomRight, l.xMax, l.yMax, l.bWrappingMap)
	l.visibleFlags.ResetVisibilityCoords(false)

	l.SetVisible(true, &avatarPos)

	l.setTilesAroundPositionVisible(&avatarPos)

	l.floodFillIfInside(&avatarPos, true)
	l.floodFillIfInside(avatarPos.GetPositionToLeft(), true)
	l.floodFillIfInside(avatarPos.GetPositionToRight(), true)
	l.floodFillIfInside(avatarPos.GetPositionDown(), true)
	l.floodFillIfInside(avatarPos.GetPositionUp(), true)

	// build lighting mask
	l.distanceMaskMaps = make([]DistanceMaskMap, 1)
	l.distanceMaskMaps[0] = lighting.BuildDistanceMap(avatarPos, 1)

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

func (l *LayeredMap) IsPositionVisible(pos *references.Position, timeOfDay datetime.UltimaDate, lighting *Lighting) bool {
	// note: some of this may feel like overkill - but it is setting up for an eventual alpha or gradient
	// based lighting that degrades as it gets further away

	// first check to see if line of sight allows it to be seen
	if !l.visibleFlags[pos.X][pos.Y] {
		return false
	}

	// next focus on checking the primary lighting
	xUp := l.xVisibleTiles/2 + 1
	nToShow := helpers.RoundUp(float32(xUp) * timeOfDay.GetVisibilityFactorWithoutTorch(0.1))
	bShow := l.distanceMaskMaps[0][*pos] <= int(nToShow)
	if lighting.HasTorchLit() && !bShow {
		return true
	}
	return bShow

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
	l.SetTileByLayer(layer, position, indexes.NoSprites)
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

func (l *LayeredMap) ClearMapUnitTiles() {
	for _, innerMap := range l.layers[MapUnitLayer] {
		for innerKey := range innerMap {
			innerMap[innerKey] = 0
		}
	}
}
