package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"

	u_color "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
)

const buttonListModalStartYPercent = 0.35
const buttonSize = MediumButton

// const buttonHeight = 0.05

// ButtonListModal GUI modal dialog that shows a number of buttons and lets you select
type ButtonListModal struct {
	border                *Border
	buttons               []*Button
	keyboard              *input.Keyboard
	closeNoActionCallback func()
	gameScreenPercents    *sprites.PercentBasedPlacement
	gameScreenCenter      *sprites.PercentBasedCenterPoint
}

func NewButtonListModal(closeNoActionCallback func(),
	keyboard *input.Keyboard,
	gameScreenPercents *sprites.PercentBasedPlacement) *ButtonListModal {
	buttonListModal := &ButtonListModal{}
	buttonListModal.keyboard = keyboard
	buttonListModal.closeNoActionCallback = closeNoActionCallback
	buttonListModal.gameScreenPercents = gameScreenPercents
	buttonListModal.gameScreenCenter = gameScreenPercents.GetCenterPoint()

	return buttonListModal
}

func (b *ButtonListModal) Draw(screen *ebiten.Image) {
	b.border.Draw(screen)
	for _, b := range b.buttons {
		b.Draw(screen)
	}
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

func (b *ButtonListModal) AddButton(buttonText string, onClickCallback func()) {
	nButton := len(b.buttons)
	button := NewButton(buttonText, onClickCallback, b.getCenterPoint(nButton), MediumButton)
	b.buttons = append(b.buttons, button)
	b.initializeBorder()
}

func (b *ButtonListModal) getCenterPoint(nButton int) sprites.PercentBasedCenterPoint {

	return sprites.PercentBasedCenterPoint{
		X: b.gameScreenPercents.GetCenterPoint().X,
		Y: buttonListModalStartYPercent + (float64(nButton) * GetButtonHeightPercent(buttonSize) * 1.25),
	}
}

func (b *ButtonListModal) initializeBorder() {
	p := sprites.PercentBasedPlacement{
		StartPercentX: b.gameScreenCenter.X - .15,
		EndPercentX:   b.gameScreenCenter.X + .15,
		StartPercentY: buttonListModalStartYPercent - 0.1,
		EndPercentY:   buttonListModalStartYPercent + (float64(len(b.buttons)) * GetButtonHeightPercent(buttonSize) * 1.25),
	}

	b.border = NewBorder(p, 601, u_color.LighterBlueSemi)
}
