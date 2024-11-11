package widgets

import "github.com/hajimehoshi/ebiten/v2"

type Button struct {
	text     string
	callback func()
}

func NewButton(text string, callback func()) *Button {
	button := &Button{}
	button.text = text
	button.callback = callback
	return button
}

func (b Button) Draw(screen *ebiten.Image) {
}

func (b Button) Update() {
}
