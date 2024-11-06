package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

// Layout sets the game's screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	aspectRatio := 16.0 / 9.0
	if float64(outsideWidth)/float64(outsideHeight) > aspectRatio {
		return int(float64(outsideHeight) * aspectRatio), outsideHeight
	} else {
		return outsideWidth, int(float64(outsideWidth) / aspectRatio)
	}
}

func main() {
	// Set the window title and size
	ebiten.SetWindowTitle("Ultima V Redux")

	game := &Game{
		currentScene: CreateIntroMenuScene(), // Start with the main menu scene
	}

	// Run the game loop
	ebiten.SetFullscreen(false)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
