package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

// Draw method for the GameScene
func (g *GameScene) Draw(screen *ebiten.Image) {
	mapWidth := sprites.TileSize * xTilesInMap
	mapHeight := sprites.TileSize * yTilesInMap

	if g.mapImage == nil {
		g.mapImage = ebiten.NewImage(mapWidth, mapHeight)
	}

	g.mapImage.Fill(image.Black)
	if g.gameState.Location == references.Britannia_Underworld {
		g.refreshMapLayerTiles()
	} else {
		g.refreshMapLayerTiles()

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

	g.output.DrawRightSideOutput(screen)

	g.characterSummary.Draw(g.gameState, screen)
	g.provisionSummary.Draw(g.gameState, screen)

	if g.bShowDebugConsole {
		g.debugConsole.drawDebugConsole(screen)
	}
	if g.bShowInputBox {
		g.inputBox.Draw(screen)
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
}

// drawMap
func (g *GameScene) drawMap(screen *ebiten.Image) {
	screen.DrawImage(g.unscaledMapImage, &ebiten.DrawImageOptions{})
	g.debugMessage = fmt.Sprintf("%d, %d", g.gameState.Position.X, g.gameState.Position.Y)
}
