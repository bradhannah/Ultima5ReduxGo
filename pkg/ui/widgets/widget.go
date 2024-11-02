package widgets

import "github.com/hajimehoshi/ebiten/v2"

type Widget interface {
	Draw(screen *ebiten.Image)
	Update()
}
