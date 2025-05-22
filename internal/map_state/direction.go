package map_state

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
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
