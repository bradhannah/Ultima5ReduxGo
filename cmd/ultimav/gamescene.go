package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ui/mainscreen"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultima_v_save/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"path"
)

const (
	borderWidthScaling = 601
	xTilesInMap        = 19
	yTilesInMap        = 13
	keyPressDelay      = 115
)

// GameScene is another scene (e.g., the actual game)
type GameScene struct {
	gameConfig       *config.UltimaVConfiguration
	gameReferences   *references.GameReferences
	spriteSheet      *sprites.SpriteSheet
	keyboard         *input.Keyboard
	output           *text.Output
	ultimaFont       *text.UltimaFont
	mapImage         *ebiten.Image
	unscaledMapImage *ebiten.Image
	rightSideImage   *ebiten.Image
	characterSummary *mainscreen.CharacterSummary
	provisionSummary *mainscreen.ProvisionSummary

	debugMessage string

	avatarX, avatarY int

	gameState *game_state.GameState

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
	gameScene.output = text.NewOutput(gameScene.ultimaFont, 20)

	gameScene.keyboard = &input.Keyboard{MillisecondDelayBetweenKeyPresses: keyPressDelay}
	gameScene.avatarX = 75
	gameScene.avatarY = 75

	ebiten.SetTPS(120)

	// TODO: add a New function to GameState
	gameScene.gameState = &game_state.GameState{}
	err = gameScene.gameState.LoadSaveGame(path.Join(gameScene.gameConfig.DataFilePath, "SAVED.GAM"))
	if err != nil {
		log.Fatal(err)
	}

	gameScene.characterSummary = mainscreen.NewCharacterSummary(gameScene.spriteSheet)
	gameScene.provisionSummary = mainscreen.NewProvisionSummary(gameScene.spriteSheet)

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
		g.output.AddToContinuousOutput("Enter")
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.debugMessage = "up"
		g.output.AddToContinuousOutput("> North")
		g.avatarY = helpers.Max(g.avatarY-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.debugMessage = "down"
		g.output.AddToContinuousOutput("> South")
		g.avatarY = (g.avatarY + 1) % references.YTiles
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.debugMessage = "left"
		g.output.AddToContinuousOutput("> West")
		g.avatarX = helpers.Max(g.avatarX-1, 0)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.debugMessage = "right"
		g.output.AddToContinuousOutput("> East")
		g.avatarX = (g.avatarX + 1) % references.XTiles
	}
	return nil
}
