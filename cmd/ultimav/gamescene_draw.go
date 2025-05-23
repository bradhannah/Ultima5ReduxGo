package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
)

var gameScreenPercents = sprites.PercentBasedPlacement{
	StartPercentX: .015,
	EndPercentX:   0.75,
	StartPercentY: 0.02,
	EndPercentY:   0.98,
}

// Draw method for the GameScene
func (g *GameScene) Draw(screen *ebiten.Image) {
	if g.lastCheckedResolution != config.GetWindowResolutionFromEbiten() {
		g.initializeResizeableVisualElements()
	}

	mapWidth := sprites.TileSize * xTilesVisibleOnGameScreen
	mapHeight := sprites.TileSize * yTilesVisibleOnGameScreen

	if g.mapImage == nil {
		g.mapImage = ebiten.NewImage(mapWidth, mapHeight)
	}

	g.mapImage.Fill(image.Black)
	g.refreshAllMapLayerTiles()

	g.drawMap(g.mapImage)

	op := sprites.GetDrawOptionsFromPercentsForWholeScreen(g.mapImage,
		gameScreenPercents)

	screen.DrawImage(g.mapImage, op)
	g.drawBorders(screen)

	g.output.DrawRightSideOutput(screen)

	g.characterSummary.Draw(&g.gameState.PartyState, screen)
	g.provisionSummary.Draw(g.gameState, screen)

	// draw the dialogs - but stacked on top of each other
	g.drawDialogs(screen)

	// Render the game scene
	ebitenutil.DebugPrint(screen, g.debugMessage)
}

func (g *GameScene) drawDialogs(screen *ebiten.Image) {
	for _, val := range g.dialogStack.Dialogs {
		(*val).Draw(screen)
	}
}

// drawBorders
func (g *GameScene) drawBorders(screen *ebiten.Image) {
	screen.DrawImage(g.borders.outsideBorder, g.borders.outsideBorderOp)
	screen.DrawImage(g.borders.rightSideBorder, g.borders.rightSideBorderOp)
	screen.DrawImage(g.borders.midRightTextOutputBorder, g.borders.midRightTextOutputBorderOp)
	screen.DrawImage(g.borders.bottomRightProvisionsBorder, g.borders.bottomRightProvisionsBorderOp)
	screen.DrawImage(g.borders.topRightCharacterBorder, g.borders.topRightCharacterBorderOp)
}

// drawMap
func (g *GameScene) drawMap(screen *ebiten.Image) {
	screen.DrawImage(g.unscaledMapImage, &ebiten.DrawImageOptions{})
	g.debugMessage = fmt.Sprintf("%d, %d", g.gameState.MapState.PlayerLocation.Position.X, g.gameState.MapState.PlayerLocation.Position.Y)
}
