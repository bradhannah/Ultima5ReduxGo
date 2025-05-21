package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
)

var lastResolution config.ScreenResolution

// Scene interface that requires Update and Draw methods
type Scene interface {
	Update(game *Game) error
	Draw(screen *ebiten.Image)
	GetUltimaConfiguration() *config.UltimaVConfiguration
	InvalidateResolution()
}

// Game struct that holds the current scene
type Game struct {
	currentScene Scene
}

// Update calls the current scene's Update method
func (g *Game) Update() error {
	// if lastResolution.X == 0 || lastResolution.Y == 0 || lastResolution != config.GetWindowResolutionFromEbiten() {

	if lastResolution.X == 0 || lastResolution.Y == 0 || lastResolution != g.currentScene.GetUltimaConfiguration().GetCurrentTrackedWindowResolution() || lastResolution != config.GetWindowResolutionFromEbiten() {

		lastResolution = g.currentScene.GetUltimaConfiguration().GetCurrentTrackedWindowResolution()

		// lastResolution = config.GetWindowResolutionFromEbiten()

		config.SetWindowSize(config.ScreenResolution{X: lastResolution.X, Y: lastResolution.Y})
		g.currentScene.InvalidateResolution()
	}

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
