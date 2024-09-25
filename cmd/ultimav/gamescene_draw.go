package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Draw method for the GameScene
func (g *GameScene) Draw(screen *ebiten.Image) {

	mapImage := ebiten.NewImage(sprites.TileSize*xTilesInMap, sprites.TileSize*yTilesInMap)
	g.drawMap(mapImage)
	g.drawMapUnits(mapImage)

	op := sprites.GetDrawOptionsFromPercents(mapImage, sprites.PercentBasedPlacement{
		StartPercentX: .015,
		EndPercentX:   0.75,
		StartPercentY: 0.02,
		EndPercentY:   0.98,
	})

	screen.DrawImage(mapImage, op)
	g.drawBorders(screen)
	g.output.Draw(screen)

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

	mapImage := ebiten.NewImage(sprites.TileSize*xTilesInMap, sprites.TileSize*yTilesInMap)

	for x := 0; x < xTilesInMap; x++ {
		for y := 0; y < yTilesInMap; y++ {
			do.GeoM.Translate(float64(x*sprites.TileSize), float64(y*sprites.TileSize))
			mapImage.DrawImage(g.spriteSheet.GetSprite(g.gameReferences.OverworldLargeMapReference.GetTileNumber(x+g.avatarX, y+g.avatarY)), &do)
			do.GeoM.Reset()
		}
	}
	screen.DrawImage(mapImage, &do)
}
