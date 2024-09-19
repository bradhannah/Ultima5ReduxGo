package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// GameScene is another scene (e.g., the actual game)
type GameScene struct{}

// Update method for the GameScene
func (g *GameScene) Update(game *Game) error {
	// Handle gameplay logic here
	return nil
}

// Draw method for the GameScene
func (g *GameScene) Draw(screen *ebiten.Image) {
	// Render the game scene
	ebitenutil.DebugPrint(screen, "Game Scene: Enjoy the Game!")
}
