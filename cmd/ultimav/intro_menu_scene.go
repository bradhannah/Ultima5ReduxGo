package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type IntroMenuScene struct{}

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
	// Render the main menu
	ebitenutil.DebugPrint(screen, "Main Menu: Press Enter to Start")
}
