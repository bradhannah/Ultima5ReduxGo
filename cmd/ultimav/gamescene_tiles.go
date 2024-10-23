package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *GameScene) getSmallCalculatedTileIndex(ogSpriteIndex indexes.SpriteIndex, pos *references.Position) indexes.SpriteIndex {
	switch ogSpriteIndex {
	case indexes.Mirror:
		// is avatar in front of it
		pos := pos.GetPositionDown()
		if g.gameState.IsAvatarAtPosition(&pos) {
			return indexes.MirrorAvatar
		}
	}

	if !g.gameState.IsAvatarAtPosition(pos) {
		return ogSpriteIndex
	}

	switch ogSpriteIndex {
	case indexes.LeftBed:
		return indexes.AvatarSleepingInBed
	case indexes.ChairFacingRight, indexes.ChairFacingLeft, indexes.ChairFacingUp, indexes.ChairFacingDown:
		return g.getCorrectAvatarOnChairTile(ogSpriteIndex, pos)
	case indexes.LadderUp:
		return indexes.AvatarOnLadderUp
	case indexes.LadderDown:
		return indexes.AvatarOnLadderDown
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
		downPosTileIndex := g.gameReferences.SingleMapReferences.GetLocationReference(g.gameState.Location).GetTileNumberWithAnimation(int(g.gameState.Floor), &downPos)
		if downPosTileIndex == indexes.TableFoodBoth || downPosTileIndex == indexes.TableFoodTop {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingDown)
		}
		return indexes.AvatarSittingFacingDown
	case indexes.ChairFacingUp:
		upPos := pos.GetPositionUp()
		upPosTileIndex := g.gameReferences.SingleMapReferences.GetLocationReference(g.gameState.Location).GetTileNumberWithAnimation(int(g.gameState.Floor), &upPos)
		if upPosTileIndex == indexes.TableFoodBoth || upPosTileIndex == indexes.TableFoodBottom {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingUp)
		}
		return indexes.AvatarSittingFacingUp
	}

	return avatarChairTileIndex
}

func (g *GameScene) refreshMapLayerTiles() {
	if g.unscaledMapImage == nil {
		g.unscaledMapImage = ebiten.NewImage(sprites.TileSize*xTilesInMap, sprites.TileSize*yTilesInMap)
	}

	do := ebiten.DrawImageOptions{}

	xCenter := int16(xTilesInMap / 2)
	yCenter := int16(yTilesInMap / 2)
	var x, y int16
	for x = 0; x < xTilesInMap; x++ {
		for y = 0; y < yTilesInMap; y++ {
			do.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))

			pos := references.Position{X: x + g.gameState.Position.X - xCenter, Y: y + g.gameState.Position.Y - yCenter}
			spriteIndex := g.GetSpriteIndex(&pos)

			if g.gameState.Location == references.Britannia_Underworld { // Large Map
				g.gameState.LayeredMaps.GetLayeredMap(game_state.LargeMap, g.gameState.Floor).Layers[game_state.MapLayer][int(pos.X)][int(pos.Y)] = spriteIndex
			} else { // Small Map
				g.gameState.LayeredMaps.GetLayeredMap(game_state.SmallMap, g.gameState.Floor).Layers[game_state.MapLayer][int(pos.X)][int(pos.Y)] = spriteIndex
				// always favour the Avatar sprite if it is the actual map tile
				if spriteIndex != indexes.Avatar_KeyIndex {
					spriteIndex = g.gameState.LayeredMaps.GetLayeredMap(game_state.SmallMap, g.gameState.Floor).GetTopTile(&pos).Index
				}
			}
			g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(spriteIndex), &do)
			do.GeoM.Reset()
		}
	}

	return
}
