package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
)

type TextArea struct {
	text      string
	placement sprites.PercentBasedPlacement
}

// world space vs. local space
// if you pass in a whole screen (I think I will) then it's covered
// if you pass in an image then it is all skewed

// Design challenges I will encounter:
// * multiple colours in the text
// * world vs local space
// *

func NewTextArea(text string, placement sprites.PercentBasedPlacement) *TextArea {
	textArea := &TextArea{text: text}
	textArea.placement = placement

	return textArea
}

func (t *TextArea) Draw(screen *ebiten.Image) {
}

func (t *TextArea) SetText(text string) {
	t.text = text
}

func (t *TextArea) GetText() string {
	return t.text
}
