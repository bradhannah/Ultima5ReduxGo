package main

import (
	"fmt"
	mainscreen2 "github.com/bradhannah/Ultima5ReduxGo/internal/ui/mainscreen"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
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
	keyPressDelay      = 165
)

var boundKeysGame = []ebiten.Key{ebiten.KeyDown,
	ebiten.KeyUp,
	ebiten.KeyEnter,
	ebiten.KeyLeft,
	ebiten.KeyRight,
	ebiten.KeyE,
	ebiten.KeyX,
	ebiten.KeyO,
	ebiten.KeyJ,
	ebiten.KeySlash,
	ebiten.KeyBackquote,
}

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

	borders gameBorders
}

const (
	defaultOutputFontPoint = 20
	defaultLineSpacing     = 20
	maxCharsPerLine        = 16
	maxLinesForOutput      = 10
)

func NewGameScene(gameConfig *config.UltimaVConfiguration) *GameScene {
	gameScene := GameScene{gameConfig: gameConfig}

	// load the files man
	var err error
	gameScene.gameReferences, err = references.NewGameReferences(gameConfig)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: add a New function to GameState
	gameScene.gameState = &game_state.GameState{}
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

// Update method for the GameScene
func (g *GameScene) Update(_ *Game) error {
	if g.bShowDebugConsole {
		// if debug console is showing, then we eat all the input
		g.debugConsole.update()
		return nil
	}

	if g.gameState.SecondaryKeyState != game_state.PrimaryInput {
		g.handleSecondaryInput()
		return nil
	}

	// Handle gameplay logic here
	boundKey := g.keyboard.GetBoundKeyPressed(&boundKeysGame)
	if boundKey == nil {
		// allow quick tapping
		g.keyboard.SetAllowKeyPressImmediately()
		return nil
	}

	if !g.keyboard.TryToRegisterKeyPress(*boundKey) {
		return nil
	}

	switch *boundKey {
	case ebiten.KeyBackquote:
		g.bShowDebugConsole = !g.bShowDebugConsole
	case ebiten.KeyEnter:
		g.addRowStr("Enter")
	case ebiten.KeyUp:
		g.handleMovement(game_state.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(game_state.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(game_state.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(game_state.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyX:
		g.gameState.Location = references.Britannia_Underworld
		g.gameState.Floor = 0
		g.gameState.Position = g.gameState.LastLargeMapPosition
	case ebiten.KeyE:
		g.debugMessage = "Enter a place"
		newLocation := g.gameReferences.SingleMapReferences.WorldLocations.GetLocationByPosition(g.gameState.Position)
		if newLocation != references.EmptyLocation {
			g.gameState.LastLargeMapPosition = g.gameState.Position
			g.gameState.Position = references.Position{
				X: 15,
				Y: 30,
			}
			g.gameState.Location = newLocation
			g.gameState.Floor = 0
			g.addRowStr(fmt.Sprintf("%s",
				g.gameReferences.SingleMapReferences.GetLocationReference(newLocation).EnteringText))
		}
	case ebiten.KeyO:
		g.debugConsole.Output.AddRowStr("Open")
		g.addRowStr("Open-")
		if g.gameState.Location == references.Britannia_Underworld {
			g.appendToCurrentRowStr("Cannot")
			return nil
		}
		g.gameState.SecondaryKeyState = game_state.OpenDirectionInput
		// we don't want the delay, it feels unnatural
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyJ:
		g.debugMessage = "Jimmy"
		g.addRowStr("Jimmy-")
		if g.gameState.Location == references.Britannia_Underworld {
			g.appendToCurrentRowStr("Cannot")
			return nil
		}

		g.gameState.SecondaryKeyState = game_state.JimmyDoorDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	}

	// only process end of turn if the turn is actually done.
	if g.gameState.SecondaryKeyState == game_state.PrimaryInput {
		g.gameState.ProcessEndOfTurn()
	}
	return nil
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
	return g.gameReferences.SingleMapReferences.GetLocationReference(g.gameState.Location)
}

func (g *GameScene) GetSpriteIndex(position *references.Position) indexes.SpriteIndex {
	var spriteIndex = indexes.SpriteIndex(0)

	if g.gameState.Location == references.Britannia_Underworld {
		if g.gameState.IsAvatarAtPosition(position) {
			return indexes.Avatar_KeyIndex
		}
		if g.gameState.Floor == 0 {
			spriteIndex = g.gameReferences.OverworldLargeMapReference.GetSpriteIndex(position.X, position.Y)
		} else {
			spriteIndex = g.gameReferences.UnderworldLargeMapReference.GetSpriteIndex(position.X, position.Y)
		}
		return spriteIndex
	}

	if position.X < 0 || position.X >= references.XSmallMapTiles || position.Y < 0 || position.Y >= references.YSmallMapTiles {
		return g.GetCurrentLocationReference().GetOuterTile()
	}

	spriteIndex = g.gameReferences.SingleMapReferences.GetLocationReference(g.gameState.Location).GetTileNumberWithAnimation(int(g.gameState.Floor), position)
	// get an adjustment to tile if one is warranted
	updatedTileNumber := g.getSmallCalculatedTileIndex(spriteIndex, position)

	// if there is an overridden tile, then we will always favour that
	if spriteIndex != updatedTileNumber {
		return updatedTileNumber
	}

	// if it's the avatar's position and there wasn't an override then we favour showing the avatar
	if g.gameState.IsAvatarAtPosition(position) {
		return indexes.Avatar_KeyIndex
	}

	// stick to the original tile if no exceptions are made
	return spriteIndex
}
