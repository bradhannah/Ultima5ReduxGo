package widgets

import "github.com/hajimehoshi/ebiten/v2"

type TextInput struct {
	textArea *TextArea
}

func NewTextInput(textArea *TextArea) *TextInput {
	textInput := &TextInput{textArea: NewTextArea("")}
	return textInput
}

func (t *TextInput) Draw(screen *ebiten.Image) {}

func (t *TextInput) GetText() string {
	return t.textArea.GetText()
}
