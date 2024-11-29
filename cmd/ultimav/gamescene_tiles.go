package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameScene) getSmallCalculatedAvatarTileIndex(ogSpriteIndex indexes.SpriteIndex) indexes.SpriteIndex {
	return g.getSmallCalculatedNPCTileIndex(ogSpriteIndex, indexes.Avatar_KeyIndex, g.gameState.Position)

}

func (g *GameScene) getSmallCalculatedNPCTileIndex(ogSpriteIndex indexes.SpriteIndex, npcIndex indexes.SpriteIndex, spritePosition references.Position) indexes.SpriteIndex {
	switch ogSpriteIndex {
	case indexes.LeftBed:
		return indexes.AvatarSleepingInBed
	case indexes.ChairFacingRight, indexes.ChairFacingLeft, indexes.ChairFacingUp, indexes.ChairFacingDown:
		return g.getCorrectAvatarOnChairTile(ogSpriteIndex, &spritePosition)
	case indexes.LadderUp:
		return indexes.AvatarOnLadderUp
	case indexes.LadderDown:
		return indexes.AvatarOnLadderDown
	}
	return npcIndex
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
	return sprites.GetSpriteIndexWithAnimationBySpriteIndex(ogSpriteIndex, pos.GetHash())
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

func (g *GameScene) setDrawBridge(theMap *game_state.LayeredMap, pos *references.Position, spriteIndex indexes.SpriteIndex) bool {
	const (
		leftXDrawBridge   = 14
		rightXDrawBridge  = 16
		topYDrawBridge    = 28
		bottomYDrawBridge = 29
	)

	if (pos.X >= leftXDrawBridge && pos.X <= rightXDrawBridge) && (pos.Y >= topYDrawBridge && pos.Y <= bottomYDrawBridge) {
		theMap.SetTileByLayer(game_state.MapOverrideLayer, pos, g.gameState.GetDrawBridgeWaterByTime(spriteIndex))
		return true
	}
	return false
}

func (g *GameScene) getCorrectAvatarEatingInChairTile(avatarChairTileIndex indexes.SpriteIndex, pos *references.Position) indexes.SpriteIndex {
	switch avatarChairTileIndex {
	case indexes.ChairFacingDown:
		downPos := pos.GetPositionDown()
		downPosTile := g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor).GetTopTile(downPos)
		if downPosTile.Index == indexes.TableFoodBoth || downPosTile.Index == indexes.TableFoodTop {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingDown, 0)
		}
		return indexes.AvatarSittingFacingDown
	case indexes.ChairFacingUp:
		upPos := pos.GetPositionUp()
		upPosTile := g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor).GetTopTile(upPos)
		if upPosTile.Index == indexes.TableFoodBoth || upPosTile.Index == indexes.TableFoodBottom {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingUp, 0)
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

	// START SPECIAL EXCEPTIONS
	for x = 0; x < xTilesInMap; x++ {
		for y = 0; y < yTilesInMap; y++ {
			pos := references.Position{X: x + g.gameState.Position.X - xCenter, Y: y + g.gameState.Position.Y - yCenter}
			tile := theMap.GetTileTopMapOnlyTile(&pos)
			if tile == nil {
				continue
			}
			switch tile.Index {
			case indexes.Portcullis, indexes.BrickWallArchway:
				theMap.SetTileByLayer(game_state.MapOverrideLayer, &pos, g.gameState.GetArchwayPortcullisSpriteByTime())
			case indexes.WoodenPlankVert1Floor, indexes.WoodenPlankVert2Floor:
				g.setDrawBridge(theMap, &pos, tile.Index)
			}
		}
	}
	// END SPECIAL EXCEPTIONS

	// STATIC MAP
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
					spriteIndex = g.gameState.GetCurrentSmallLocationReference().GetOuterTile()
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
	// END STATIC MAP
	uf := text.NewUltimaFont(10)

	// BEGIN MAPUNITS
	g.gameState.Chests = make(map[references.Position]references.Chest)
	for x = 0; x < xTilesInMap; x++ {
		for y = 0; y < yTilesInMap; y++ {
			pos := references.Position{X: x + g.gameState.Position.X - xCenter, Y: y + g.gameState.Position.Y - yCenter}
			if mapType == references.LargeMapType {
				pos = *pos.GetWrapped(references.XLargeMapTiles, references.YLargeMapTiles)
			}
			do.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))

			mapUnitTile := theMap.GetTileByLayer(game_state.MapUnitLayer, &pos)
			underTile := theMap.GetTileTopMapOnlyTile(&pos)
			if mapUnitTile == nil || mapUnitTile.Index == 0 {
				do.GeoM.Reset()
				continue
			}
			if layer.IsPositionVisible(&pos) {
				tileIndex := g.getSmallCalculatedNPCTileIndex(underTile.Index, mapUnitTile.Index, pos)
				tileIndex = g.getSmallCalculatedTileIndex(tileIndex, &pos)
				o := text.NewOutput(uf, 20, 1, 10)

				g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(tileIndex), &do)

				o.DrawText(g.unscaledMapImage, fmt.Sprintf("x=%d y=%d", pos.X, pos.Y), &do)
			}
			do.GeoM.Reset()
		}
	}
	// END MAPUNITS

	avatarSpriteIndex := theMap.GetTileTopMapOnlyTile(&avatarPos).Index
	avatarSpriteIndex = g.getSmallCalculatedAvatarTileIndex(avatarSpriteIndex)

	g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(avatarSpriteIndex), &avatarDo)
}
