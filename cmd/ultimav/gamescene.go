package main

import (
	// "image"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/clock"
	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
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
	SearchDirectionInput
	AttackDirectionInput
	UseDirectionInput
	YellDirectionInput
	FireDirectionInput
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

	clk *clock.GameClock

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
		log.Fatal(err) // Critical game references required for operation
	}

	saveFilePath := path.Join(gameScene.gameConfig.SavedConfigData.DataFilePath, "SAVED.GAM")

	gameScene.gameState = game_state.NewGameStateFromLegacySaveFile(saveFilePath,
		gameConfig,
		gameScene.gameReferences,
		xTilesVisibleOnGameScreen,
		yTilesVisibleOnGameScreen)

	// Wire up system callbacks dependency injection using constructors
	messageCallbacks, err := game_state.NewMessageCallbacks(
		gameScene.addRowStr,
		gameScene.appendToCurrentRowStr,
		gameScene.addRowStr, // Use addRowStr for command prompts
	)
	if err != nil {
		log.Fatalf("Failed to create MessageCallbacks: %v", err) // Critical system callbacks required for operation
	}

	visualCallbacks := game_state.NewVisualCallbacks(
		nil, // KapowAt - TODO: implement when visual system is ready
		nil, // ShowMissileEffect - TODO: implement when visual system is ready
		nil, // DelayGlide - TODO: implement when visual system is ready
	)

	audioCallbacks := game_state.NewAudioCallbacks(
		func(effect game_state.SoundEffect) {
			// TODO: Implement sound system - for now just log
			// log.Printf("Playing sound effect: %v", effect)
		},
	)

	screenCallbacks := game_state.NewScreenCallbacks(
		nil, // MarkStatsChanged - TODO: implement when stats UI is ready
		nil, // UpdateStatsDisplay - TODO: implement when stats UI is ready
		nil, // RefreshInventoryDisplay - TODO: implement when inventory UI is ready
		nil, // ShowModalDialog - TODO: implement when modal system is ready
		nil, // PromptYesNo - TODO: implement when dialog system is ready
	)

	flowCallbacks := game_state.NewFlowCallbacks(
		gameScene.gameState.FinishTurn, // Wire to existing FinishTurn method
		nil,                            // ActivateGuards - TODO: implement when guard system is ready
		func(minutes int) { // Wire to existing UltimaDate time system
			gameScene.gameState.DateTime.Advance(minutes)
		},
		func(timeOfDay datetime.TimeOfDay) { // Wire to existing TimeOfDay system
			gameScene.gameState.DateTime.SetTimeOfDay(timeOfDay)
		},
		nil, // DelayFx - TODO: implement when timing system is ready
		nil, // CheckUpdate - TODO: implement when update system is ready
	)

	talkCallbacks := game_state.NewTalkCallbacks(
		gameScene.CreateTalkDialog,
		gameScene.PushDialog,
	)

	systemCallbacks, err := game_state.NewSystemCallbacks(messageCallbacks, visualCallbacks, audioCallbacks, screenCallbacks, flowCallbacks, talkCallbacks)
	if err != nil {
		log.Fatalf("Failed to create SystemCallbacks: %v", err) // Critical system callbacks required for operation
	}

	gameScene.gameState.SystemCallbacks = systemCallbacks

	gameScene.spriteSheet = sprites.NewSpriteSheet()
	gameScene.keyboard = input.NewKeyboard(keyPressDelay)
	gameScene.initializeResizeableVisualElements()

	// ebiten.SetTPS(120)
	ebiten.SetTPS(60)

	gameScene.clk = clock.NewGameClock(16) // ~60 fixed updates per second

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
		log.Fatal("Unexpected - should find debug dialog index") // TODO: CONVERT TO SOFT ERROR - UI recovery possible, should not crash game
	}
	// g.dialogs = append(g.dialogs[:nIndex], g.dialogs[nIndex+1:]...)
	dc := (*DebugConsole)(nil)
	g.dialogStack.RemoveWidget(dc)
}

// TalkCallbacks interface implementation
func (g *GameScene) CreateTalkDialog(npc *map_units.NPCFriendly) game_state.TalkDialog {
	return NewLinearTalkDialog(g, npc.NPCReference)
}

func (g *GameScene) PushDialog(dialog game_state.TalkDialog) {
	if linearDialog, ok := dialog.(*LinearTalkDialog); ok {
		g.dialogStack.PushModalDialog(linearDialog)
	}
}
