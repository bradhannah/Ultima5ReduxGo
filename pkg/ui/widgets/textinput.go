package widgets

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type TextInput struct {
	textArea *TextArea

	output     *text.Output
	ultimaFont *text.UltimaFont
}

func NewTextInput(fontPointSize float64) *TextInput {
	textInput := &TextInput{}
	textInput.ultimaFont = text.NewUltimaFont(fontPointSize)
	// NOTE: single line input only (for now?)
	textInput.output = text.NewOutput(textInput.ultimaFont, 0, 1)
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
	return t.textArea.GetText()
}
