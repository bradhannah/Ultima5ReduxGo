package widgets

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

const (
	keyPressDelay = 165
)

type TextInput struct {
	output     *text.Output
	ultimaFont *text.UltimaFont

	maxCharsPerLine int
	hasFocus        bool

	keyboard *input.Keyboard
}

var nonAlphaNumericBoundKeys = []ebiten.Key{ebiten.KeyDown,
	ebiten.KeyEnter,
	ebiten.KeySpace,
	ebiten.KeyBackspace,
	ebiten.KeyDelete,
	ebiten.KeyEscape,
}

func NewTextInput(fontPointSize float64, maxCharsPerLine int) *TextInput {
	textInput := &TextInput{}
	textInput.ultimaFont = text.NewUltimaFont(fontPointSize)
	textInput.maxCharsPerLine = maxCharsPerLine
	textInput.keyboard = input.NewKeyboard(keyPressDelay)

	// NOTE: single line input only (for now?)
	textInput.output = text.NewOutput(textInput.ultimaFont, 0, 1, maxCharsPerLine)
	textInput.output.AddRowStr("TESTINPUT_oof")
	textInput.output.SetColor(color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	})
	return textInput
}

func (t *TextInput) Draw(screen *ebiten.Image) {
	const percentIntoBorder = 0.02
	textRect := sprites.GetRectangleFromPercents(sprites.PercentBasedPlacement{
		StartPercentX: 0 + percentIntoBorder,
		EndPercentX:   .75 + .01 - percentIntoBorder,
		StartPercentY: .955,
		EndPercentY:   1,
	})
	t.output.DrawContinuousOutputTexOnXy(screen, image.Point{
		X: textRect.Min.X,
		Y: textRect.Min.Y,
	}, false)
}

func (t *TextInput) GetText() string {
	return ""
	//return t.textArea.GetText()
}

func (t *TextInput) addCharacter(keyStr string) {
	if t.output.GetLengthOfCurrentRowStr() < t.maxCharsPerLine {
		t.output.AppendToCurrentRowStr(keyStr)
	}
}

func (t *TextInput) removeRightCharacter(keyStr string) {

}

func (t *TextInput) Update() {
	//if !t.hasFocus {
	//	return
	//}

	key, keyStr := t.keyboard.GetAlphaNumericPressed()
	if key != nil {
		if !t.keyboard.TryToRegisterKeyPress(*key) {
			return
		}
		t.addCharacter(keyStr)
		return
	}

	boundKey := t.keyboard.GetBoundKeyPressed(&nonAlphaNumericBoundKeys)
	if boundKey == nil {
		return
	}
	if !t.keyboard.TryToRegisterKeyPress(*boundKey) {
		return
	}

	switch *boundKey {
	case ebiten.KeyEnter:
		return
	case ebiten.KeySpace:
		t.addCharacter(" ")
	case ebiten.KeyBackspace:
		t.output.RemoveRightSideCharacter()
	}
	return
}
