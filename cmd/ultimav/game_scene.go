package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

const (
	borderWidthScaling = 601
	xTilesInMap        = 19
	yTilesInMap        = 14
)

type gameBorders struct {
	outsideBorder     *ebiten.Image
	outsideBorderOp   *ebiten.DrawImageOptions
	rightSideBorder   *ebiten.Image
	rightSideBorderOp *ebiten.DrawImageOptions
}

// GameScene is another scene (e.g., the actual game)
type GameScene struct {
	gameConfig     *config.UltimaVConfiguration
	gameReferences *references.GameReferences
	spriteSheet    *sprites.SpriteSheet
	keyboard       *input.Keyboard

	debugMessage string

	avatarX, avatarY int

	borders gameBorders
}

// Update method for the GameScene
func (g *GameScene) Update(game *Game) error {
	// Handle gameplay logic here
	if !g.keyboard.IsBoundKeyPressed(boundKeys) {
		return nil
	}
	if !g.keyboard.TryToRegisterKeyPress() {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.debugMessage = "enter"
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.debugMessage = "up"
		g.avatarY = helpers.Max(g.avatarY-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.debugMessage = "down"
		g.avatarY = (g.avatarY + 1) % references.YTiles
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.debugMessage = "left"
		g.avatarX = helpers.Max(g.avatarX-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.debugMessage = "right"
		g.avatarX = (g.avatarX + 1) % references.XTiles
	}
	return nil
}

// Draw method for the GameScene
func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.borders.rightSideBorder, g.borders.rightSideBorderOp)
	screen.DrawImage(g.borders.outsideBorder, g.borders.outsideBorderOp)

	// make unscaled map image
	//do := ebiten.DrawImageOptions{}

	//const xTiles = 20
	//const yTiles = 14

	mapImage := ebiten.NewImage(sprites.TileSize*xTilesInMap, sprites.TileSize*yTilesInMap)

	g.drawMap(mapImage)
	g.drawMapUnits(mapImage)

	op := sprites.GetDrawOptionsFromPercents(mapImage, sprites.PercentBasedPlacement{
		StartPercentX: .015,
		EndPercentX:   0.75, //(.98 - 0.015) / (16.0 / 9.0),
		StartPercentY: 0.02,
		EndPercentY:   0.98,
	})

	screen.DrawImage(mapImage, op)

	// Render the game scene
	ebitenutil.DebugPrint(screen, g.debugMessage)

}

func NewGameScene(gameConfig *config.UltimaVConfiguration) *GameScene {
	gameScene := GameScene{gameConfig: gameConfig}

	// load the files man
	var err error
	gameScene.gameReferences, err = references.NewGameReferences(gameConfig)
	if err != nil {
		log.Fatal(err)
	}

	// initialize borders
	mainBorder := sprites.NewBorderSprites()
	image, op := mainBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, sprites.PercentBasedPlacement{
		StartPercentX: 0,
		StartPercentY: 0,
		EndPercentY:   1,
		EndPercentX:   1,
	})
	gameScene.borders.outsideBorder = image
	gameScene.borders.outsideBorderOp = op
	rightSideBorder := sprites.NewBorderSprites()
	image, op = rightSideBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, sprites.PercentBasedPlacement{
		StartPercentX: .75,
		EndPercentX:   1,
		StartPercentY: 0,
		EndPercentY:   1,
	})
	gameScene.borders.rightSideBorder = image
	gameScene.borders.rightSideBorderOp = op

	gameScene.spriteSheet = sprites.NewSpriteSheet()

	gameScene.keyboard = &input.Keyboard{MillisecondDelayBetweenKeyPresses: 100}
	gameScene.avatarX = 75
	gameScene.avatarY = 75

	return &gameScene
}

// drawMapUnits
func (g *GameScene) drawMapUnits(screen *ebiten.Image) {
	do := ebiten.DrawImageOptions{}

	do.GeoM.Translate(float64(sprites.TileSize*(xTilesInMap/2)), float64(sprites.TileSize*(yTilesInMap/2)))

	//avatarImage := ebiten.NewImage(sprites.TileSize, sprites.TileSize)
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
