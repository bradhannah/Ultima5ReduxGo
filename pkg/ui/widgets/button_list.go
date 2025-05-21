package widgets

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	text2 "github.com/hajimehoshi/ebiten/v2/text/v2"

	ucolor "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
)

const buttonListModalStartYPercent = 0.35
const buttonSize = MediumButton
const buttonSpacingPercent = 1.15

const borderStartYPercent = buttonListModalStartYPercent - 0.1
const titleStartYPercent = borderStartYPercent + 0.0455
const fontPoint = 22
const borderPercentEachSide = 0.125

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

	ultimaFont *text.UltimaFont
	titleText  string
	// titleTextScreenPercents *sprites.PercentBasedPlacement
	titleTextPoint  *image.Point
	titleTextOutput *text.Output
}

func NewButtonListModal(
	titleText string,
	closeNoActionCallback func(),
	keyboard *input.Keyboard,
	gameScreenPercents *sprites.PercentBasedPlacement,
) *ButtonListModal {
	buttonListModal := &ButtonListModal{}
	buttonListModal.keyboard = keyboard
	buttonListModal.closeNoActionCallback = closeNoActionCallback
	buttonListModal.gameScreenPercents = gameScreenPercents
	buttonListModal.gameScreenCenter = gameScreenPercents.GetCenterPoint()

	buttonListModal.ultimaFont = text.NewUltimaFont(text.GetScaledNumberToResolution(fontPoint))
	buttonListModal.titleTextOutput = text.NewOutput(buttonListModal.ultimaFont, fontPoint, 1, 30)

	textRect := sprites.GetRectangleFromPercents(sprites.PercentBasedPlacement{
		StartPercentX: buttonListModal.gameScreenCenter.X,
		EndPercentX:   buttonListModal.gameScreenCenter.X,
		StartPercentY: titleStartYPercent,
		EndPercentY:   1,
	})
	buttonListModal.titleTextPoint = &image.Point{
		X: textRect.Min.X + (textRect.Dx() / 2),
		Y: textRect.Min.Y,
	}
	buttonListModal.titleText = titleText
	buttonListModal.titleTextOutput.AddRowStr(buttonListModal.titleText)

	buttonListModal.currentSelectionIndex = -1

	return buttonListModal
}

func (b *ButtonListModal) Draw(screen *ebiten.Image) {
	b.border.Draw(screen)
	for _, b := range b.buttons {
		b.Draw(screen)
	}
	b.titleTextOutput.DrawContinuousOutputTexOnXy(screen, *b.titleTextPoint, false, text2.AlignCenter, text2.AlignCenter)
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
	case ebiten.KeyEnter:
		b.buttons[b.currentSelectionIndex].callback()
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
		Y: buttonListModalStartYPercent + (float64(nButton) * GetButtonHeightPercent(buttonSize) * buttonSpacingPercent),
	}
}

func (b *ButtonListModal) initializeBorder() {
	p := sprites.PercentBasedPlacement{
		StartPercentX: b.gameScreenCenter.X - borderPercentEachSide,
		EndPercentX:   b.gameScreenCenter.X + borderPercentEachSide,
		StartPercentY: borderStartYPercent,
		EndPercentY:   buttonListModalStartYPercent + (float64(len(b.buttons)) * GetButtonHeightPercent(buttonSize) * buttonSpacingPercent),
	}

	b.border = NewBorder(p, 601, ucolor.Black)
}
