package main

import (
	"fmt"
	mainscreen2 "github.com/bradhannah/Ultima5ReduxGo/internal/ui/mainscreen"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
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

var boundKeysGame = []ebiten.Key{ebiten.KeyDown, ebiten.KeyUp, ebiten.KeyEnter, ebiten.KeyLeft, ebiten.KeyRight, ebiten.KeyE, ebiten.KeyX, ebiten.KeyO}

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
	characterSummary *mainscreen2.CharacterSummary
	provisionSummary *mainscreen2.ProvisionSummary

	debugMessage string

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

	// TODO: add a New function to GameState
	gameScene.gameState = &game_state.GameState{}
	err = gameScene.gameState.LoadLegacySaveGame(path.Join(gameScene.gameConfig.DataFilePath, "SAVED.GAM"), gameScene.gameReferences)
	if err != nil {
		log.Fatal(err)
	}

	gameScene.initializeBorders()

	gameScene.spriteSheet = sprites.NewSpriteSheet()
	gameScene.ultimaFont = text.NewUltimaFont(text.OutputFontPoint)
	gameScene.output = text.NewOutput(gameScene.ultimaFont, 20)

	gameScene.keyboard = &input.Keyboard{MillisecondDelayBetweenKeyPresses: keyPressDelay}

	ebiten.SetTPS(120)

	gameScene.characterSummary = mainscreen2.NewCharacterSummary(gameScene.spriteSheet)
	gameScene.provisionSummary = mainscreen2.NewProvisionSummary(gameScene.spriteSheet)

	return &gameScene
}

// Update method for the GameScene
func (g *GameScene) Update(game *Game) error {
	if g.gameState.SecondaryKeyState != game_state.PrimaryInput {
		g.handleSecondaryInput()
		return nil
	}

	// Handle gameplay logic here
	if !g.keyboard.IsBoundKeyPressed(boundKeysGame) {
		return nil
	}
	if !g.keyboard.TryToRegisterKeyPress() {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.debugMessage = "enter"
		g.output.AddToContinuousOutput("Enter")
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.handleMovement("North", ebiten.KeyUp)
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.handleMovement("South", ebiten.KeyDown)
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.handleMovement("West", ebiten.KeyLeft)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.handleMovement("East", ebiten.KeyRight)
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.gameState.Location = references.Britannia_Underworld
		g.gameState.Floor = 0
		g.gameState.Position = g.gameState.LastLargeMapPosition
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
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
			g.output.AddToContinuousOutput(fmt.Sprintf("%s",
				g.gameReferences.SingleMapReferences.GetSingleMapReference(newLocation).EnteringText))
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyO) {
		g.debugMessage = "Open"
		g.output.AddToContinuousOutput("Open-")
		if g.gameState.Location == references.Britannia_Underworld {
			g.output.AppendToOutput("Cannot")
			return nil
		}
		g.gameState.SecondaryKeyState = game_state.OpenDirectionInput
	}

	// only process end of turn if the turn is actually done.
	if g.gameState.SecondaryKeyState == game_state.PrimaryInput {
		g.gameState.ProcessEndOfTurn()
	}
	return nil
}
