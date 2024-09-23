package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type gameBorders struct {
	outsideBorder   *ebiten.Image
	outsideBorderOp *ebiten.DrawImageOptions
}

// GameScene is another scene (e.g., the actual game)
type GameScene struct {
	gameConfig     *config.UltimaVConfiguration
	gameReferences *references.GameReferences
	spriteSheet    *sprites.SpriteSheet

	borders gameBorders
}

// Update method for the GameScene
func (g *GameScene) Update(game *Game) error {
	// Handle gameplay logic here
	return nil
}

// Draw method for the GameScene
func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.borders.outsideBorder, g.borders.outsideBorderOp)

	g.drawMap(screen)

	// Render the game scene
	ebitenutil.DebugPrint(screen, "Game Scene: Enjoy the Game!")

}

func CreateGameScene(gameConfig *config.UltimaVConfiguration) *GameScene {
	gameScene := GameScene{gameConfig: gameConfig}

	// load the files man
	var err error
	gameScene.gameReferences, err = references.NewGameReferences(gameConfig)
	if err != nil {
		log.Fatal(err)
	}

	border := sprites.NewBorderSprites()
	image, op := border.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(601, sprites.PercentBasedPlacement{
		StartPercentX: 0,
		StartPercentY: 0,
		EndPercentY:   1,
		EndtPercentX:  1,
	})
	gameScene.borders.outsideBorder = image
	gameScene.borders.outsideBorderOp = op

	gameScene.spriteSheet = sprites.NewSpriteSheet()

	return &gameScene
}

func (g *GameScene) drawMap(screen *ebiten.Image) {
	do := ebiten.DrawImageOptions{}
	do.GeoM.Translate(0, 0)

	//screen.DrawImage(g.spriteSheet.SpriteImage, &do)
	screen.DrawImage(g.spriteSheet.GetSprite(183), &do)
}
