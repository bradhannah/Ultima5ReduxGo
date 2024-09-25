package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	borderWidthScaling = 601
	xTilesInMap        = 19
	yTilesInMap        = 13
	keyPressDelay      = 175
)

// GameScene is another scene (e.g., the actual game)
type GameScene struct {
	gameConfig     *config.UltimaVConfiguration
	gameReferences *references.GameReferences
	spriteSheet    *sprites.SpriteSheet
	keyboard       *input.Keyboard
	output         *text.Output
	ultimaFont     *text.UltimaFont

	debugMessage string

	avatarX, avatarY int

	borders gameBorders
}

func NewGameScene(gameConfig *config.UltimaVConfiguration) *GameScene {
	gameScene := GameScene{gameConfig: gameConfig}

	// load the files man
	var err error
	gameScene.gameReferences, err = references.NewGameReferences(gameConfig)
	if err != nil {
		log.Fatal(err)
	}

	gameScene.initializeBorders()

	gameScene.spriteSheet = sprites.NewSpriteSheet()
	gameScene.ultimaFont = text.NewUltimaFont(text.OutputFontPoint)
	gameScene.output = text.NewOutput(gameScene.ultimaFont)

	gameScene.keyboard = &input.Keyboard{MillisecondDelayBetweenKeyPresses: keyPressDelay}
	gameScene.avatarX = 75
	gameScene.avatarY = 75

	return &gameScene
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
		g.output.AddToOutput("Enter")
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.debugMessage = "up"
		g.output.AddToOutput("> North")
		g.avatarY = helpers.Max(g.avatarY-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.debugMessage = "down"
		g.output.AddToOutput("> South")
		g.avatarY = (g.avatarY + 1) % references.YTiles
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.debugMessage = "left"
		g.output.AddToOutput("> West")
		g.avatarX = helpers.Max(g.avatarX-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.debugMessage = "right"
		g.output.AddToOutput("> East")
		g.avatarX = (g.avatarX + 1) % references.XTiles
	}
	return nil
}
