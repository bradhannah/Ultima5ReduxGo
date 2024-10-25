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
		downPosTileIndex := g.gameReferences.LocationReferences.GetLocationReference(g.gameState.Location).GetTileNumberWithAnimation(int(g.gameState.Floor), &downPos)
		if downPosTileIndex == indexes.TableFoodBoth || downPosTileIndex == indexes.TableFoodTop {
			return sprites.GetSpriteIndexWithAnimationBySpriteIndex(indexes.AvatarSittingAndEatingFacingDown)
		}
		return indexes.AvatarSittingFacingDown
	case indexes.ChairFacingUp:
		upPos := pos.GetPositionUp()
		upPosTileIndex := g.gameReferences.LocationReferences.GetLocationReference(g.gameState.Location).GetTileNumberWithAnimation(int(g.gameState.Floor), &upPos)
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
	mapType := g.gameState.Location.GetMapType()
	// remove it so it doesn't stick around for a single frame
	if mapType == references.SmallMapType {
		g.gameState.WipeOldAvatarPosition()
	}

	xCenter := references.Coordinate(xTilesInMap / 2)
	yCenter := references.Coordinate(yTilesInMap / 2)
	var x, y references.Coordinate
	for x = 0; x < xTilesInMap; x++ {
		for y = 0; y < yTilesInMap; y++ {
			do.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))

			pos := references.Position{X: x + g.gameState.Position.X - xCenter, Y: y + g.gameState.Position.Y - yCenter}
			spriteIndex := g.GetSpriteIndex(&pos)

			g.gameState.LayeredMaps.GetLayeredMap(mapType, g.gameState.Floor).SetTile(game_state.MapLayer, &pos, spriteIndex)
			switch mapType {
			case references.SmallMapType:
				if g.gameState.Position.Equals(pos) {
					// the avatar is on this tile
					g.gameState.SetNewAvatarPosition(&pos)
					spriteIndex = indexes.Avatar_KeyIndex
				} else {
					spriteIndex = g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor).GetTopTile(&pos).Index
				}
			case references.LargeMapType:
				_ = ""
			default:
				panic("unhandled default case")
			}

			g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(spriteIndex), &do)
			do.GeoM.Reset()
		}
	}
	// get the previous avatar position - wipe it

	// re-add avatar to the avatar layer

	return
}
