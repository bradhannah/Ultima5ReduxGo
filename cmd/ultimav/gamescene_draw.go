package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Draw method for the GameScene
func (g *GameScene) Draw(screen *ebiten.Image) {
	const widthRatio = 16
	const heightRatio = 9

	mapWidth := sprites.TileSize * xTilesInMap
	mapHeight := sprites.TileSize * yTilesInMap

	if g.mapImage == nil {
		g.mapImage = ebiten.NewImage(mapWidth, mapHeight)
	}

	g.drawMap(g.mapImage)
	g.drawMapUnits(g.mapImage)

	op := sprites.GetDrawOptionsFromPercentsForWholeScreen(g.mapImage, sprites.PercentBasedPlacement{
		StartPercentX: .015,
		EndPercentX:   0.75,
		StartPercentY: 0.02,
		EndPercentY:   0.98,
	})

	screen.DrawImage(g.mapImage, op)
	g.drawBorders(screen)

	rightSideWidth := mapWidth / 4
	rightSideHeight := mapHeight
	if g.rightSideImage == nil {
		g.rightSideImage = ebiten.NewImage(rightSideWidth*5, int(float64(rightSideHeight)*3.7))
	} else {
		g.rightSideImage.Clear()
	}
	op = sprites.GetDrawOptionsFromPercentsForWholeScreen(g.rightSideImage, sprites.PercentBasedPlacement{
		StartPercentX: .751,
		EndPercentX:   1.0,
		StartPercentY: 0.02,
		EndPercentY:   0.98,
	})

	g.output.DrawContinuousOutputText(g.rightSideImage)

	screen.DrawImage(g.rightSideImage, op)
	g.characterSummary.Draw(g.gameState, screen)
	g.provisionSummary.Draw(g.gameState, screen)

	// Render the game scene
	ebitenutil.DebugPrint(screen, g.debugMessage)
}

// drawBorders
func (g *GameScene) drawBorders(screen *ebiten.Image) {
	screen.DrawImage(g.borders.outsideBorder, g.borders.outsideBorderOp)
	screen.DrawImage(g.borders.rightSideBorder, g.borders.rightSideBorderOp)
	screen.DrawImage(g.borders.midRightTextOutputBorder, g.borders.midRightTextOutputBorderOp)
	screen.DrawImage(g.borders.bottomRightProvisionsBorder, g.borders.bottomRightProvisionsBorderOp)
	screen.DrawImage(g.borders.topRightCharacterBorder, g.borders.topRightCharacterBorderOp)
}

// drawMapUnits
func (g *GameScene) drawMapUnits(screen *ebiten.Image) {
	do := ebiten.DrawImageOptions{}

	do.GeoM.Translate(float64(sprites.TileSize*(xTilesInMap/2)), float64(sprites.TileSize*(yTilesInMap/2)))

	screen.DrawImage(g.spriteSheet.GetSprite(indexes.Avatar), &do)
}

// drawMap
func (g *GameScene) drawMap(screen *ebiten.Image) {
	do := ebiten.DrawImageOptions{}

	if g.unscaledMapImage == nil {
		g.unscaledMapImage = ebiten.NewImage(sprites.TileSize*xTilesInMap, sprites.TileSize*yTilesInMap)
	}

	xCenter := int16(xTilesInMap / 2)
	yCenter := int16(yTilesInMap / 2)
	var x, y int16
	for x = 0; x < xTilesInMap; x++ {
		for y = 0; y < yTilesInMap; y++ {
			do.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))
			var tileNumber int
			if g.gameState.Location == references.Britannia_Underworld {
				tileNumber = g.gameReferences.OverworldLargeMapReference.GetTileNumber(x+g.gameState.Position.X-xCenter, y+g.gameState.Position.Y-yCenter)
			} else {
				// small map for now
				pos := references.Position{
					X: x - xCenter + g.gameState.Position.X,
					Y: y - yCenter + g.gameState.Position.Y,
				}
				if pos.X < 0 || pos.X >= references.XSmallMapTiles || pos.Y < 0 || pos.Y >= references.YSmallMapTiles {
					tileNumber = 5
				} else {
					tileNumber = g.gameReferences.SingleMapReferences.GetSingleMapReference(g.gameState.Location,
						int(g.gameState.Floor)).GetTileNumber(&pos)
				}
			}
			g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(tileNumber), &do)
			do.GeoM.Reset()
		}
	}
	screen.DrawImage(g.unscaledMapImage, &do)

	g.debugMessage = fmt.Sprintf("%d, %d", g.gameState.Position.X, g.gameState.Position.Y)
}
