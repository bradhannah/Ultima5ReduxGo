package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type IntroMenuScene struct {
	sprites *sprites.IntroSprites
}

type screenDimensions struct {
	x int
	y int
}

//var ScreenDimension = screenDimensions{x: 800, y: 600}

var ScreenDimension = screenDimensions{x: 1024, y: 768}

// Update method for the IntroMenuScene
func (m *IntroMenuScene) Update(game *Game) error {
	// Switch to the gameplay scene on keypress (e.g., pressing "Enter")
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// Replace this with code to switch to the game scene
		fmt.Println("Switching to Game Scene")
		game.currentScene = &GameScene{}
	}
	return nil
}

// Draw method for the IntroMenuScene
func (m *IntroMenuScene) Draw(screen *ebiten.Image) {
	op := sprites.GetCenteredXSprite(m.sprites.Ultima16Logo, 50, 100)

	screen.DrawImage(m.sprites.Ultima16Logo, op)

	// Render the main menu
	ebitenutil.DebugPrint(screen, "Main Menu: Press Enter to Start")
}

func CreateIntroMenuScene() *IntroMenuScene {
	return &IntroMenuScene{
		sprites: sprites.NewIntroSprites(),
	}
}
