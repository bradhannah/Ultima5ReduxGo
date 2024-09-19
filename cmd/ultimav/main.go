package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

// Layout sets the game's screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Set the width and height of the game window
	return 1024, 768
}

func main() {
	// Set the window title and size
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Ultima V Redux")

	game := &Game{
		currentScene: &IntroMenuScene{}, // Start with the main menu scene
	}
	// Run the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
