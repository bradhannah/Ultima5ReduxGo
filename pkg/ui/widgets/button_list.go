package widgets

import "github.com/hajimehoshi/ebiten/v2"

// ButtonListModal GUI modal dialog that shows a number of buttons and lets you select
type ButtonListModal struct {
	buttons []*Button
}

func NewButtonListModal() *ButtonListModal {
	buttonListModal := &ButtonListModal{}
	return buttonListModal
}

func (b *ButtonListModal) Draw(screen *ebiten.Image) {
}

func (b *ButtonListModal) Update() {
}

func (b *ButtonListModal) AddButton(button *Button) {
	b.buttons = append(b.buttons, button)
}
