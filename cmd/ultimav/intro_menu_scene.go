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
	op := &ebiten.DrawImageOptions{}
	//op.GeoM.Translate(100, 100) // Position the sprite
	//op.GeoM.Scale(float64(screenDimension.x), float64(screenDimension.y))
	//op.GeoM.Translate(float64(screenDimension.x)+100, float64(screenDimension.y)+100) // Position the sprite
	const xDiff = 100
	targetWidth := float64(ScreenDimension.x - xDiff)
	//targetHeight := float64(screenDimension.y - 100)

	imgWidth := m.sprites.Ultima16Logo.Bounds().Dx()
	//imgHeight := m.sprites.Ultima16Logo.Bounds().Dy()

	scaleX := targetWidth / float64(imgWidth)
	scaledWidth := scaleX * float64(imgWidth)

	// Scale the height (Y-axis) by the same factor to maintain the aspect ratio
	scaleY := scaleX
	//targetHeight := float64(imgHeight) * scaleY
	//offsetY := (float64(screenDimension.y) - targetHeight) / 2

	// 320 across, 61 down
	// 800 /2 = 400
	// 400 - (320/2) = 400 - 160 = 240

	//op.GeoM.Translate(0, 0)
	op.GeoM.Scale(scaleX, scaleY)
	//op.GeoM.Scale(1, 1)
	//= 700 / 2 + 50 = 400
	// 800 - 700 + 50
	// 1024 - 924 + 50 = 150
	topLeftX := (float64(ScreenDimension.x) - float64(scaledWidth)) / 2
	//+ (xDiff / 2)
	op.GeoM.Translate(topLeftX, 100)
	//op.GeoM.Translate(240-50, 100)
	//op.GeoM.Translate(50, offsetY)

	screen.DrawImage(m.sprites.Ultima16Logo, op)

	// Render the main menu
	ebitenutil.DebugPrint(screen, "Main Menu: Press Enter to Start")
}

func CreateIntroMenuScene() *IntroMenuScene {
	return &IntroMenuScene{
		sprites: sprites.NewIntroSprites(),
	}
}
