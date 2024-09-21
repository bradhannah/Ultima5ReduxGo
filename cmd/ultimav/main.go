package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

// Layout sets the game's screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Set the width and height of the game window
	return ScreenDimension.x, ScreenDimension.y
}

func main() {
	// Set the window title and size
	ebiten.SetWindowSize(config.WindowWidth, config.WindowHeight)
	ebiten.SetWindowTitle("Ultima V Redux")

	game := &Game{
		currentScene: CreateIntroMenuScene(), // Start with the main menu scene
	}
	// Run the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
