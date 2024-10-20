package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
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

	g.refreshMapLayerTiles()
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

	g.output.DrawRightSideOutput(g.rightSideImage)

	screen.DrawImage(g.rightSideImage, op)
	g.characterSummary.Draw(g.gameState, screen)
	g.provisionSummary.Draw(g.gameState, screen)

	if g.bShowDebugConsole {
		g.debugConsole.drawDebugConsole(screen)
	}

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
	//screen.DrawImage(g.spriteSheet.GetSprite(indexes.Avatar), &do)
}

// drawMap
func (g *GameScene) drawMap(screen *ebiten.Image) {
	screen.DrawImage(g.unscaledMapImage, &ebiten.DrawImageOptions{})
	g.debugMessage = fmt.Sprintf("%d, %d", g.gameState.Position.X, g.gameState.Position.Y)
}
