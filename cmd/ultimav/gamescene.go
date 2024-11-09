package main

import (
	"image"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"

	mainscreen2 "github.com/bradhannah/Ultima5ReduxGo/internal/ui/mainscreen"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ui/widgets"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const (
	borderWidthScaling = 601
	xTilesInMap        = 19
	yTilesInMap        = 13
)

const (
	defaultOutputFontPoint = 23
	defaultLineSpacing     = 25
	maxCharsPerLine        = 16
	maxLinesForOutput      = 11
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

	inputBox      *widgets.InputBox
	bShowInputBox bool

	debugMessage string

	gameState *game_state.GameState

	secondaryKeyState InputState

	borders gameBorders

	lastCheckedResolution config.ScreenResolution
}

func (g *GameScene) InvalidateResolution() {
	g.initializeResizeableVisualElements()
}

func (g *GameScene) GetUltimaConfiguration() *config.UltimaVConfiguration {
	return g.gameConfig
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
	err = gameScene.gameState.LoadLegacySaveGame(path.Join(gameScene.gameConfig.SavedConfigData.DataFilePath, "SAVED.GAM"), gameScene.gameReferences)

	if err != nil {
		log.Fatal(err)
	}

	gameScene.spriteSheet = sprites.NewSpriteSheet()
	gameScene.keyboard = input.NewKeyboard(keyPressDelay)
	gameScene.initializeResizeableVisualElements()

	// ebiten.SetTPS(120)
	ebiten.SetTPS(60)

	return &gameScene
}

func (g *GameScene) initializeResizeableVisualElements() {
	g.lastCheckedResolution = config.GetWindowResolutionFromEbiten()

	g.initializeBorders()
	g.ultimaFont = text.NewUltimaFont(text.GetScaledNumberToResolution(defaultOutputFontPoint))
	if g.output == nil {
		g.output = text.NewOutput(
			g.ultimaFont,
			text.GetScaledNumberToResolution(defaultLineSpacing),
			maxLinesForOutput,
			maxCharsPerLine)
	} else {
		g.output.SetFont(
			g.ultimaFont,
			text.GetScaledNumberToResolution(defaultLineSpacing),
		)
	}

	if g.debugConsole == nil {
		g.debugConsole = NewDebugConsole(g)
	} else {
		g.debugConsole.Refresh()
	}

	g.characterSummary = mainscreen2.NewCharacterSummary(g.spriteSheet)
	g.provisionSummary = mainscreen2.NewProvisionSummary(g.spriteSheet)
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

func (g *GameScene) DoModalInputBox(question string, textCommand *grammar.TextCommand) {
	g.inputBox = widgets.NewInputBox(question, textCommand)
	g.bShowInputBox = true
}
