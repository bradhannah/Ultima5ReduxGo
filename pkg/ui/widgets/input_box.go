package widgets

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/hajimehoshi/ebiten/v2"
)

const inputBoxPercentOffEdge = 0.2
const inputBoxBorderWidthScaling = 601
const inputPercentIntoBorder = 0.02

type InputBox struct {
	textInput *TextInput
	Question  string

	border *Border

	font *text.UltimaFont
}

func (i *InputBox) Draw(screen *ebiten.Image) {
	i.border.Draw(screen)
	i.textInput.Draw(screen)
}

func (i *InputBox) Update() {
	i.textInput.Update()
}

func (i *InputBox) GetText() string {
	return i.textInput.GetText()
}

func NewInputBox(question string, textCommand *grammar.TextCommand) *InputBox {
	inputBox := &InputBox{}
	inputBox.Question = question

	textCommands := &grammar.TextCommands{
		*textCommand,
	}

	inputBox.textInput = NewTextInput(
		sprites.PercentBasedPlacement{
			StartPercentX: 0 + inputPercentIntoBorder,
			EndPercentX:   .75 + .01 - inputPercentIntoBorder,
			StartPercentY: .955,
			EndPercentY:   1,
		},
		20, 20,
		textCommands,
		TextInputCallbacks{
			AmbiguousAutoComplete: func(message string) {
			},
		})

	inputBox.border = NewBorder(
		sprites.PercentBasedPlacement{
			StartPercentX: 0 + inputBoxPercentOffEdge,
			EndPercentX:   .75 + .01 - inputBoxPercentOffEdge,
			StartPercentY: .65,
			EndPercentY:   .9,
		},
		inputBoxBorderWidthScaling,
	)
	return inputBox
}
