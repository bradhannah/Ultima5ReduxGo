package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
)

// ButtonListModal GUI modal dialog that shows a number of buttons and lets you select
type ButtonListModal struct {
	buttons               []*Button
	keyboard              *input.Keyboard
	closeNoActionCallback func()
}

func NewButtonListModal(closeNoActionCallback func(), keyboard *input.Keyboard) *ButtonListModal {
	buttonListModal := &ButtonListModal{}
	buttonListModal.keyboard = keyboard
	buttonListModal.closeNoActionCallback = closeNoActionCallback

	return buttonListModal
}

func (b *ButtonListModal) Draw(screen *ebiten.Image) {
}

func (b *ButtonListModal) Update() {
	var boundKey *ebiten.Key
	// the idea here is if the user pushes, then lets go and repushes the same button it should
	// show both characters since it is not technically repeating
	boundKey = b.keyboard.GetBoundKeyPressed(&nonAlphaNumericBoundKeys)
	if boundKey == nil {
		b.keyboard.SetAllowKeyPressImmediately()
		return
	}
	if !b.keyboard.TryToRegisterKeyPress(*boundKey) {
		return
	}

	switch *boundKey {
	case ebiten.KeyEscape:

		// t.output.AppendToCurrentRowStr("-")
		b.closeNoActionCallback()
	}
}

func (b *ButtonListModal) AddButton(button *Button) {
	b.buttons = append(b.buttons, button)
}
