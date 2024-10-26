package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	keyPressDelay = 165
)

var boundKeysGame = []ebiten.Key{
	ebiten.KeyDown,
	ebiten.KeyUp,
	ebiten.KeyEnter,
	ebiten.KeyLeft,
	ebiten.KeyRight,
	ebiten.KeyE,
	ebiten.KeyX,
	ebiten.KeyO,
	ebiten.KeyJ,g.gameState.Karma
	ebiten.KeyK,
	ebiten.KeyP,
	ebiten.KeyG,
	ebiten.KeySlash,
	ebiten.KeyBackquote,
}

// Update method for the GameScene
func (g *GameScene) Update(_ *Game) error {
	mapType := g.gameState.Location.GetMapType()

	if g.bShowDebugConsole {
		// if debug console is showing, then we eat all the input
		g.debugConsole.update()
		return nil
	}

	if g.secondaryKeyState != PrimaryInput {
		switch mapType {
		case references.LargeMapType:
			g.largeMapHandleSecondaryInput()
		case references.SmallMapType:
			g.smallMapHandleSecondaryInput()
		case references.CombatMapType:
		case references.DungeonMapType:
		}
		return nil
	}

	// Handle gameplay logic here
	boundKey := g.keyboard.GetBoundKeyPressed(&boundKeysGame)
	if boundKey == nil {
		// allow quick tapping
		g.keyboard.SetAllowKeyPressImmediately()
		return nil
	}

	if !g.keyboard.TryToRegisterKeyPress(*boundKey) {
		return nil
	}

	switch g.gameState.Location.GetMapType() {
	case references.LargeMapType:
		g.largeMapInputHandler(*boundKey)
	case references.SmallMapType:
		g.smallMapInputHandler(*boundKey)
	case references.CombatMapType:
	case references.DungeonMapType:
	}

	return nil
}
