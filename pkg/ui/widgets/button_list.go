package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"

	u_color "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
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

	currentSelectionIndex int
}

func NewButtonListModal(closeNoActionCallback func(),
	keyboard *input.Keyboard,
	gameScreenPercents *sprites.PercentBasedPlacement) *ButtonListModal {
	buttonListModal := &ButtonListModal{}
	buttonListModal.keyboard = keyboard
	buttonListModal.closeNoActionCallback = closeNoActionCallback
	buttonListModal.gameScreenPercents = gameScreenPercents
	buttonListModal.gameScreenCenter = gameScreenPercents.GetCenterPoint()

	buttonListModal.currentSelectionIndex = -1

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
		b.closeNoActionCallback()
	case ebiten.KeyDown, ebiten.KeyRight:
		b.GoDown()
	case ebiten.KeyUp, ebiten.KeyLeft:
		b.GoUp()
	}
}

func (b *ButtonListModal) GoDown() {
	b.currentSelectionIndex = (b.currentSelectionIndex + 1) % len(b.buttons)
	b.updateButtons()
}

func (b *ButtonListModal) GoUp() {
	if b.currentSelectionIndex == 0 {
		b.currentSelectionIndex = len(b.buttons) - 1
	} else {
		b.currentSelectionIndex = helpers.Max(b.currentSelectionIndex-1, 0)
	}
	b.updateButtons()
}

func (b *ButtonListModal) AddButton(buttonText string, onClickCallback func()) {
	nButton := len(b.buttons)
	button := NewButton(buttonText, onClickCallback, b.getCenterPoint(nButton), MediumButton)
	button.SetButtonStatus(Selected)
	b.buttons = append(b.buttons, button)
	b.initializeBorder()

	if b.currentSelectionIndex == -1 {
		b.currentSelectionIndex = nButton
	}
	b.updateButtons()
}

func (b *ButtonListModal) updateButtons() {
	for i, button := range b.buttons {
		if i == b.currentSelectionIndex {
			button.SetButtonStatus(Selected)
		} else {
			button.SetButtonStatus(NotSelected)
		}
	}
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

	b.border = NewBorder(p, 601, u_color.Black)
}
