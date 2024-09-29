package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var oof = 0

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

//func (g *GameScene) drawCharacterSummary(screen *ebiten.Image) {
//	g.characterSummary.DrawContinuousOutputText(g.gameState)
//	dop := ebiten.DrawImageOptions{}
//	dop.GeoM.Translate(15, 15)
//	dop.GeoM.Scale(5, 5)
//	screen.DrawImage(g.characterSummary.FullSummaryImage, &dop)
//}

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

	for x := 0; x < xTilesInMap; x++ {
		for y := 0; y < yTilesInMap; y++ {
			do.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))
			tileNumber := g.gameReferences.OverworldLargeMapReference.GetTileNumber(x+g.avatarX, y+g.avatarY)
			g.unscaledMapImage.DrawImage(g.spriteSheet.GetSprite(tileNumber), &do)
			do.GeoM.Reset()
		}
	}
	screen.DrawImage(g.unscaledMapImage, &do)
}
