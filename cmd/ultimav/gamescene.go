package main

import (
	mainscreen2 "github.com/bradhannah/Ultima5ReduxGo/internal/ui/mainscreen"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"path"
)

const (
	borderWidthScaling = 601
	xTilesInMap        = 19
	yTilesInMap        = 13
)

const (
	defaultOutputFontPoint = 20
	defaultLineSpacing     = 20
	maxCharsPerLine        = 16
	maxLinesForOutput      = 10
)

type InputState int

const (
	PrimaryInput InputState = iota
	OpenDirectionInput
	JimmyDoorDirectionInput
	KlimbDirectionInput
	PushDirectionInput
	GetDirectionInput
	LookDirectionInput
)

// GameScene is another scene (e.g., the actual game)
type GameScene struct {
	gameConfig          *config.UltimaVConfiguration
	gameReferences      *references.GameReferences
	spriteSheet         *sprites.SpriteSheet
	keyboard            *input.Keyboard
	output              *text.Output
	ultimaFont          *text.UltimaFont
	mapImage            *ebiten.Image
	unscaledMapImage    *ebiten.Image
	rightSideImage      *ebiten.Image
	debugWindowImage    *ebiten.Image
	debugWindowSizeRect *image.Rectangle
	debugWindowPosRect  *image.Rectangle
	characterSummary    *mainscreen2.CharacterSummary
	provisionSummary    *mainscreen2.ProvisionSummary

	debugConsole      *DebugConsole
	bShowDebugConsole bool

	debugMessage string

	gameState *game_state.GameState

	secondaryKeyState InputState

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

	// TODO: add a New function to GameState
	gameScene.gameState = &game_state.GameState{
		XTilesInMap:    xTilesInMap,
		YTilesInMap:    yTilesInMap,
		GameReferences: gameScene.gameReferences,
	}
	err = gameScene.gameState.LoadLegacySaveGame(path.Join(gameScene.gameConfig.DataFilePath, "SAVED.GAM"), gameScene.gameReferences)

	if err != nil {
		log.Fatal(err)
	}

	gameScene.initializeBorders()

	gameScene.spriteSheet = sprites.NewSpriteSheet()
	gameScene.ultimaFont = text.NewUltimaFont(defaultOutputFontPoint)
	gameScene.output = text.NewOutput(gameScene.ultimaFont, defaultLineSpacing, maxLinesForOutput, maxCharsPerLine)

	gameScene.keyboard = input.NewKeyboard(keyPressDelay)

	//ebiten.SetTPS(120)
	ebiten.SetTPS(60)

	gameScene.characterSummary = mainscreen2.NewCharacterSummary(gameScene.spriteSheet)
	gameScene.provisionSummary = mainscreen2.NewProvisionSummary(gameScene.spriteSheet)

	gameScene.debugConsole = NewDebugConsole(&gameScene)

	return &gameScene
}

func (g *GameScene) appendToCurrentRowStr(str string) {
	g.output.AppendToCurrentRowStr(str)
	g.debugConsole.Output.AppendToCurrentRowStr(str)
}

func (g *GameScene) addRowStr(str string) {
	g.output.AddRowStr(str)
	g.debugConsole.Output.AddRowStr(str)
}

func (g *GameScene) GetCurrentLocationReference() *references.SmallLocationReference {
	return g.gameReferences.LocationReferences.GetLocationReference(g.gameState.Location)
}
