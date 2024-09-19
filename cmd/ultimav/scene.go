package main

import "github.com/hajimehoshi/ebiten/v2"

// Scene interface that requires Update and Draw methods
type Scene interface {
	Update(game *Game) error
	Draw(screen *ebiten.Image)
}

// Game struct that holds the current scene
type Game struct {
	currentScene Scene
}

// Update calls the current scene's Update method
func (g *Game) Update() error {
	if g.currentScene != nil {
		return g.currentScene.Update(g)
	}
	return nil
}

// Draw calls the current scene's Draw method
func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}
