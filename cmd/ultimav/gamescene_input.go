package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

const (
	keyPressDelay = 165
)

var boundKeysGame = []ebiten.Key{
	ebiten.KeySpace,
	ebiten.KeyDown,
	ebiten.KeyUp,
	ebiten.KeyEnter,
	ebiten.KeyLeft,
	ebiten.KeyRight,
	ebiten.KeyB,
	ebiten.KeyE,
	ebiten.KeyX,
	ebiten.KeyO,
	ebiten.KeyJ,
	ebiten.KeyK,
	ebiten.KeyP,
	ebiten.KeyG,
	ebiten.KeyL,
	ebiten.KeyI,
	ebiten.KeyT,
	ebiten.KeySlash,
	ebiten.KeyBackquote,
	ebiten.KeyEscape,
}

// Update method for the GameScene
func (g *GameScene) Update(_ *Game) error {
	mapType := g.gameState.MapState.PlayerLocation.Location.GetMapType()

	if g.updateDialogs() {
		// only run a single update routine at a time
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

	switch g.gameState.MapState.PlayerLocation.Location.GetMapType() {
	case references.LargeMapType:
		g.largeMapInputHandler(*boundKey)
	case references.SmallMapType:
		g.smallMapInputHandler(*boundKey)
	case references.CombatMapType:
	case references.DungeonMapType:
	}

	return nil
}

func (g *GameScene) updateDialogs() bool {
	if !g.dialogStack.HasOpenDialog() {
		return false
	}
	(*g.dialogStack.PeekTopModalDialog()).Update()
	return true
}
