package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
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
	// Use DisplayManager to handle resolution changes
	if g.currentScene != nil && g.currentScene.GetUltimaConfiguration() != nil {
		dm := g.currentScene.GetUltimaConfiguration().DisplayManager
		if dm != nil {
			// DisplayManager handles all resolution change detection internally
			dm.Update()
		}

		// Handle tracked resolution changes for window resizing
		trackedRes := g.currentScene.GetUltimaConfiguration().GetCurrentTrackedWindowResolution()
		if lastResolution.X == 0 || lastResolution.Y == 0 || lastResolution != trackedRes {
			lastResolution = trackedRes
			if !dm.IsFullscreen() {
				dm.SetWindowSize(trackedRes.X, trackedRes.Y)
			}
			g.currentScene.InvalidateResolution()
		}
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
