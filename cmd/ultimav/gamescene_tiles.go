package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func (g *GameScene) getSmallCalculatedAvatarTileIndex(ogSpriteIndex indexes.SpriteIndex) indexes.SpriteIndex {
	//pos := g.gameState.Position
	//if !g.gameState.IsAvatarAtPosition(pos) {
	//	return ogSpriteIndex
	//}

	switch ogSpriteIndex {
	case indexes.LeftBed:
		return indexes.AvatarSleepingInBed
	case indexes.ChairFacingRight, indexes.ChairFacingLeft, indexes.ChairFacingUp, indexes.ChairFacingDown:
		return g.getCorrectAvatarOnChairTile(ogSpriteIndex, &g.gameState.Position)
	case indexes.LadderUp:
		return indexes.AvatarOnLadderUp
	case indexes.LadderDown:
		return indexes.AvatarOnLadderDown
	}
	return indexes.Avatar_KeyIndex
}

func (g *GameScene) getSmallCalculatedTileIndex(ogSpriteIndex indexes.SpriteIndex, pos *references.Position) indexes.SpriteIndex {
	switch ogSpriteIndex {
	case indexes.Mirror:
		// is avatar in front of it
		pos := pos.GetPositionDown()
		if g.gameState.IsAvatarAtPosition(pos) {
			return indexes.MirrorAvatar
		}
	}
	return ogSpriteIndex
}

func (g *GameScene) getCorrectAvatarOnChairTile(spriteIndex indexes.SpriteIndex, position *references.Position) indexes.SpriteIndex {
	switch spriteIndex {
	case indexes.ChairFacingRight:
		return indexes.AvatarSittingFacingRight
	case indexes.ChairFacingLeft:
		return indexes.AvatarSittingFacingLeft
	case indexes.ChairFacingUp, indexes.ChairFacingDown:
		return g.getCorrectAvatarEatingInChairTile(spriteIndex, position)
	}
	return spriteIndex
}

func (g *GameScene) getCorrectAvatarEatingInChairTile(avatarChairTileIndex indexes.SpriteIndex, pos *references.Position) indexes.SpriteIndex {
	switch avatarChairTileIndex {
	case indexes.ChairFacingDown:
		downPos := pos.GetPositionDown()
		downPosTile := g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor).GetTopTile(downPos)
		if downPosTile.Index == indexes.TableFoodBoth || downPosTile.Index == indexes.TableFoodTop {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingDown)
		}
		return indexes.AvatarSittingFacingDown
	case indexes.ChairFacingUp:
		upPos := pos.GetPositionUp()
		upPosTile := g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor).GetTopTile(upPos)
		if upPosTile.Index == indexes.TableFoodBoth || upPosTile.Index == indexes.TableFoodBottom {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingUp)
		}
		return indexes.AvatarSittingFacingUp
	}

	return avatarChairTileIndex
}

func (g *GameScene) refreshMapLayerTiles() {
	layer := g.gameState.GetLayeredMapByCurrentLocation()
	layer.RecalculateVisibleTiles(g.gameState.Position)

	if g.unscaledMapImage == nil {
		g.unscaledMapImage = ebiten.NewImage(sprites.TileSize*xTilesInMap, sprites.TileSize*yTilesInMap)
	}

	g.unscaledMapImage.Fill(image.Black)
	mapType := g.gameState.Location.GetMapType()

	do := ebiten.DrawImageOptions{}

	xCenter := references.Coordinate(xTilesInMap / 2)
	yCenter := references.Coordinate(yTilesInMap / 2)

	var avatarPos references.Position
	var avatarDo ebiten.DrawImageOptions

	theMap := g.gameState.LayeredMaps.GetLayeredMap(mapType, g.gameState.Floor)

	var x, y references.Coordinate
	// get and set tiles
	for x = 0; x < xTilesInMap; x++ {
		for y = 0; y < yTilesInMap; y++ {
			pos := references.Position{X: x + g.gameState.Position.X - xCenter, Y: y + g.gameState.Position.Y - yCenter}
			if mapType == references.LargeMapType {
				pos = *pos.GetWrapped(references.XLargeMapTiles, references.YLargeMapTiles)
			}
			do.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))

			tile := theMap.GetTileTopMapOnlyTile(&pos)
			var spriteIndex indexes.SpriteIndex
			if tile == nil {
				if g.gameState.IsOutOfBounds(pos) {
					spriteIndex = 5
				} else {
					log.Fatal("bad index")
				}
			} else {
				spriteIndex = tile.Index
			}

			// get from the reference
			spriteIndex = g.getSmallCalculatedTileIndex(spriteIndex, &pos)

			if g.gameState.Position.Equals(pos) {
				avatarPos = pos
				avatarDo = do
			}

			if layer.IsPositionVisible(&pos) {
				g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(spriteIndex), &do)
			}
			do.GeoM.Reset()
		}
	}

	avatarSpriteIndex := theMap.GetTileTopMapOnlyTile(&avatarPos).Index
	avatarSpriteIndex = g.getSmallCalculatedAvatarTileIndex(avatarSpriteIndex)

	g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(avatarSpriteIndex), &avatarDo)
}
