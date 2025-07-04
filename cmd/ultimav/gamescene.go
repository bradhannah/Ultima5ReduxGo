package main

import (
	// "image"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/internal/text"
	mainscreen2 "github.com/bradhannah/Ultima5ReduxGo/internal/ui/mainscreen"
	"github.com/bradhannah/Ultima5ReduxGo/internal/ui/widgets"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
)

const (
	borderWidthScaling        = 601
	xTilesVisibleOnGameScreen = 19
	yTilesVisibleOnGameScreen = 13
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
	TalkDirectionInput
)

// GameScene is another scene (e.g., the actual game)
type GameScene struct {
	gameConfig       *config.UltimaVConfiguration
	gameReferences   *references2.GameReferences
	spriteSheet      *sprites.SpriteSheet
	keyboard         *input.Keyboard
	output           *text.Output
	ultimaFont       *text.UltimaFont
	mapImage         *ebiten.Image
	unscaledMapImage *ebiten.Image
	// rightSideImage      *ebiten.Image
	// debugWindowImage    *ebiten.Image
	// debugWindowSizeRect *image.Rectangle
	// debugWindowPosRect  *image.Rectangle
	characterSummary *mainscreen2.CharacterSummary
	provisionSummary *mainscreen2.ProvisionSummary

	dialogStack widgets.DialogStack

	debugConsole *DebugConsole

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
	gameScene.gameReferences, err = references2.NewGameReferences(gameConfig)
	if err != nil {
		log.Fatal(err)
	}

	saveFilePath := path.Join(gameScene.gameConfig.SavedConfigData.DataFilePath, "SAVED.GAM")

	gameScene.gameState = game_state.NewGameStateFromLegacySaveFile(saveFilePath,
		gameConfig,
		gameScene.gameReferences,
		xTilesVisibleOnGameScreen,
		yTilesVisibleOnGameScreen)

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
	g.output.AddRowStrWithTrim(str)
	g.debugConsole.Output.AddRowStrWithTrim(str)
}

func (g *GameScene) GetCurrentLocationReference() *references2.SmallLocationReference {
	return g.gameReferences.LocationReferences.GetLocationReference(g.gameState.MapState.PlayerLocation.Location)
}

func (g *GameScene) IsDebugDialogOpen() bool {
	return g.dialogStack.GetIndexOfDialogType((*DebugConsole)(nil)) != -1
}

func (g *GameScene) toggleDebug() {
	if !g.IsDebugDialogOpen() {
		g.dialogStack.PushModalDialog(g.debugConsole)

		return
	}

	nIndex := g.dialogStack.GetIndexOfDialogType((*DebugConsole)(nil))

	if nIndex == -1 {
		log.Fatal("Unexpected - should find debug dialog index")
	}
	// g.dialogs = append(g.dialogs[:nIndex], g.dialogs[nIndex+1:]...)
	dc := (*DebugConsole)(nil)
	g.dialogStack.RemoveWidget(dc)
}
