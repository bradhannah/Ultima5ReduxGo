package map_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

const (
	totalLayers   = 6
	overflowTiles = references.Coordinate(10)
)

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

type GameDimensions interface {
	GetTilesVisibleOnScreen() (int, int)
	GetTopLeftExtent() references.Position
	GetBottomRightExtent() references.Position
	GetBottomRightWithoutOverflow() references.Position
	IsWrappedMap() bool
}

// LayeredMap represents a map broken into separate layers. These layers are essentially "floors"
type LayeredMap struct {
	layers [totalLayers]Layer

	visibleFlags            VisibilityCoords
	testForVisibilityMap    VisibilityCoords
	primaryDistanceMaskMap  DistanceMap
	lightSourcesDistanceMap DistanceMap

	tileRefs *references.Tiles

	xVisibleTiles, yVisibleTiles     int
	TopLeft, BottomRight             references.Position
	XMaxTilesPerMap, YMaxTilesPerMap references.Coordinate

	bWrappingMap bool
}

func newLayeredMap(xMax references.Coordinate,
	yMax references.Coordinate,
	tileRefs *references.Tiles,
	xVisibleTiles int,
	yVisibleTiles int,
	bWrappingMap bool,
) *LayeredMap { //nolint:exhaustruct

	layeredMap := LayeredMap{ //nolint:exhaustruct
		xVisibleTiles:   xVisibleTiles,
		yVisibleTiles:   yVisibleTiles,
		XMaxTilesPerMap: xMax,
		YMaxTilesPerMap: yMax,
		bWrappingMap:    bWrappingMap,
		tileRefs:        tileRefs,
	}

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

func (l *LayeredMap) RecalculateVisibleTiles(avatarPos references.Position, lighting *Lighting) {
	l.TopLeft = references.Position{
		X: avatarPos.X - overflowTiles,
		Y: avatarPos.Y - overflowTiles,
	}
	l.BottomRight = references.Position{
		X: avatarPos.X + overflowTiles,
		Y: avatarPos.Y + overflowTiles,
	}

	// TODO: it's lazy to make both of these calls since it could do in one pass
	l.testForVisibilityMap.ResetVisibilityCoords(false)
	l.testForVisibilityMap.SetVisibilityCoordsRectangle(
		&l.TopLeft,
		&l.BottomRight,
		l.XMaxTilesPerMap,
		l.YMaxTilesPerMap,
		l.bWrappingMap,
	)
	l.visibleFlags.ResetVisibilityCoords(false)

	l.SetVisible(true, &avatarPos)

	l.setTilesAroundPositionVisible(&avatarPos)

	l.floodFillIfInside(&avatarPos, true)
	l.floodFillIfInside(avatarPos.GetPositionToLeft(), true)
	l.floodFillIfInside(avatarPos.GetPositionToRight(), true)
	l.floodFillIfInside(avatarPos.GetPositionDown(), true)
	l.floodFillIfInside(avatarPos.GetPositionUp(), true)

	// get a full game screen lighting map
	l.primaryDistanceMaskMap = lighting.BuildGameScreenDistanceMap(avatarPos)

	// get lighting sources and overlay them of a light map
	lightSources := l.getAllLightSourcesInRange(avatarPos)
	l.lightSourcesDistanceMap = lighting.BuildLightSourceDistanceMap(lightSources,
		l.visibleFlags,
		lighting.turnsToExtinguishTorch > 0,
		avatarPos,
	)
}

func (l *LayeredMap) getAllLightSourcesInRange(centerPosition references.Position) LightSources {
	lightSources := make(LightSources, 0)

	for dX := -l.XMaxTilesPerMap; dX < l.XMaxTilesPerMap; dX++ {
		for dY := -l.YMaxTilesPerMap; dY < l.YMaxTilesPerMap; dY++ {
			pos := references.Position{X: centerPosition.X + dX, Y: centerPosition.Y + dY}

			tile := l.GetTileTopMapOnlyTile(&pos)
			if tile == nil {
				continue
			}

			if tile.LightEmission > 0 {
				lightSources = append(lightSources, lightSource{
					Tile: tile,
					Pos:  pos,
				})
			}
		}
	}

	return lightSources
}

func (l *LayeredMap) floodFillIfInside(pos *references.Position, bForce bool) {
	if l.bWrappingMap {
		pos = pos.GetWrapped(l.XMaxTilesPerMap, l.YMaxTilesPerMap)
	}

	if !bForce && !l.testForVisibilityMap[pos.X][pos.Y] {
		// don't reprocess and existing tile if it has been set
		return
	}

	l.SetVisible(true, pos)
	l.testForVisibilityMap[pos.X][pos.Y] = false

	tile := l.GetTileTopMapOnlyTile(pos)

	// basically - if they are next to a window, then they can see through - otherwise we treat
	// the window as opaque, and nothing passes through
	if tile != nil {
		if tile.IsWindow && bForce {
		} else if tile.BlocksLight || tile.IsWindow {
			return
		}
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

// SetVisible sets the visibility of a tile to true or false.
func (l *LayeredMap) SetVisible(bVisible bool, pos *references.Position) {
	l.visibleFlags[pos.X][pos.Y] = bVisible
}

// IsPositionVisible returns true if the position is visible based on the time of day.
func (l *LayeredMap) IsPositionVisible(pos *references.Position, timeOfDay datetime.UltimaDate, bIsBasement bool) bool {
	// note: some of this may feel like overkill - but it is setting up for an eventual alpha- or gradient-
	// based lighting that degrades as it gets further away

	// first check to see if the line of sight allows it to be seen
	if !l.visibleFlags[pos.X][pos.Y] {
		return false
	}

	// next focus on checking the primary lighting
	xUp := l.xVisibleTiles/2 + 1

	var nToShow int

	if bIsBasement {
		nToShow = 1 // helpers.RoundUp(float32(xUp) * timeOfDay.GetVisibilityFactorWithoutTorch(0.1))
	} else {
		nToShow = helpers.RoundUp(float32(xUp) * timeOfDay.GetVisibilityFactorWithoutTorch(0.1))
	}

	// TODO: hey brad - the nToShow is small because of the Avatar - so this won't quite work.
	// figure out how to distinguish these two - i think the second map that is more "factual" is better
	bShow := l.primaryDistanceMaskMap[*pos] <= nToShow

	if bShow {
		return true
	}

	if _, ok := l.lightSourcesDistanceMap[*pos]; !ok {
		return false
	}

	return true
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

func (l *LayeredMap) SwapTiles(pos1, pos2 *references.Position) {
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

func (l *LayeredMap) GetTilesVisibleOnScreen() (int, int) {
	return l.xVisibleTiles, l.yVisibleTiles
}

func (l *LayeredMap) GetTopLeftExtent() references.Position {
	return l.TopLeft
}

func (l *LayeredMap) GetBottomRightExtent() references.Position {
	return l.BottomRight
}

func (l *LayeredMap) IsWrappedMap() bool {
	return l.bWrappingMap
}

func (l *LayeredMap) GetBottomRightWithoutOverflow() references.Position {
	return references.Position{X: l.XMaxTilesPerMap, Y: l.YMaxTilesPerMap}
}
