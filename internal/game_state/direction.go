package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func GetKeyAsDirection(key ebiten.Key) references.Direction {
	switch key {
	case ebiten.KeyUp:
		return references.Up
	case ebiten.KeyDown:
		return references.Down
	case ebiten.KeyLeft:
		return references.Left
	case ebiten.KeyRight:
		return references.Right
	}
	return references.NoneDirection
}
