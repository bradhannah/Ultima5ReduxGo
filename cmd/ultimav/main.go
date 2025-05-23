package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil)) // Exposes pprof at /debug/pprof
	}()

	// Set the window title and size
	ebiten.SetWindowTitle("Ultima V Redux")

	game := &Game{
		currentScene: CreateIntroMenuScene(), // Start with the main menu scene
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
