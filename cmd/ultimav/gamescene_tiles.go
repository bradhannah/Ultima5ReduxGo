package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *GameScene) getSmallCalculatedTileIndex(ogTileIndex int, pos *references.Position) int {
	switch ogTileIndex {
	case indexes.Mirror:
		// is avatar in front of it
	}

	if !g.gameState.IsAvatarAtPosition(pos) {
		return ogTileIndex
	}

	switch ogTileIndex {
	case indexes.LeftBed:
		// is avatar on it?
		return indexes.AvatarSleepingInBed
	case indexes.ChairFacingRight, indexes.ChairFacingLeft, indexes.ChairFacingUp, indexes.ChairFacingDown:
		return g.getCorrectAvatarOnChairTile(ogTileIndex, pos)
	case indexes.LadderUp:
		return indexes.AvatarOnLadderUp
	case indexes.LadderDown:
		return indexes.AvatarOnLadderDown
	}
	return ogTileIndex
}

func (g *GameScene) getCorrectAvatarOnChairTile(tileIndex int, position *references.Position) int {
	switch tileIndex {
	case indexes.ChairFacingRight:
		return indexes.AvatarSittingFacingRight
	case indexes.ChairFacingLeft:
		return indexes.AvatarSittingFacingLeft
	case indexes.ChairFacingUp, indexes.ChairFacingDown:
		return g.getCorrectAvatarEatingInChairTile(tileIndex, position)
	}
	return tileIndex
}

func (g *GameScene) getCorrectAvatarEatingInChairTile(avatarChairTileIndex int, pos *references.Position) int {
	//func isFoodTable(tileIndex) bool {
	//	return avatarChairTileIndex != indexes.TableFoodTop &&
	//		avatarChairTileIndex != indexes.TableFoodBottom &&
	//		avatarChairTileIndex != indexes.TableFoodBoth
	//
	//}

	switch avatarChairTileIndex {
	case indexes.ChairFacingDown:
		downPos := pos.GetPositionDown()
		downPosTileIndex := g.gameReferences.SingleMapReferences.GetLocationReference(g.gameState.Location).GetTileNumberWithAnimation(int(g.gameState.Floor), &downPos)
		if downPosTileIndex == indexes.TableFoodBoth || downPosTileIndex == indexes.TableFoodTop {
			return sprites.GetTileNumberWithAnimationByTile(indexes.AvatarSittingAndEatingFacingDown)
		}
		return indexes.AvatarSittingFacingUp
	case indexes.ChairFacingUp:
		upPos := pos.GetPositionUp()
		upPosTileIndex := g.gameReferences.SingleMapReferences.GetLocationReference(g.gameState.Location).GetTileNumberWithAnimation(int(g.gameState.Floor), &upPos)
		if upPosTileIndex == indexes.TableFoodBoth || upPosTileIndex == indexes.TableFoodBottom {
			return sprites.GetTileNumberWithAnimationByTile(indexes.AvatarSittingAndEatingFacingUp)
		}
		return indexes.AvatarSittingFacingDown
	}

	return avatarChairTileIndex
}

func (g *GameScene) refreshMap() {
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
			tileNumber := g.GetTileIndex(&pos)

			if g.gameState.Location == references.Britannia_Underworld { // Large Map
				g.gameState.LayeredMaps.LayeredMaps[game_state.LargeMap].Layers[game_state.MapLayer][int(pos.X)][int(pos.Y)] = tileNumber
			} else { // Small Map
				g.gameState.LayeredMaps.LayeredMaps[game_state.SmallMap].Layers[game_state.MapLayer][int(pos.X)][int(pos.Y)] = tileNumber
				// always favour the Avatar sprite if it is the actual map tile
				if tileNumber != indexes.Avatar_KeyIndex {
					tileNumber = g.gameState.LayeredMaps.LayeredMaps[game_state.SmallMap].GetTopTile(&pos).Index
				}
			}
			g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(tileNumber), &do)
			do.GeoM.Reset()
		}
	}

	return
}
