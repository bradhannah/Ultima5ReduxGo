package widgets

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
)

var _ Widget = &InputBox{}

const (
	inputBoxPercentOffEdge     = 0.2
	inputBoxBorderWidthScaling = 601
	inputPercentIntoBorder     = 0.02
	inputBoxFontPoint          = 20
	inputBoxDefaultLineSpacing = 20
	inputBoxMaxCharsPerLine    = 32
	inputBoxMaxLinesForOutput  = 10

	inputBoxStartPercentY = .70
)

type InputBox struct {
	textInput    *TextInput
	textQuestion *text.Output
	Question     string

	border *Border

	font *text.UltimaFont

	inputBoxPercents  sprites.PercentBasedPlacement
	borderBoxPercents sprites.PercentBasedPlacement
}

func NewInputBox(question string, textCommand *grammar.TextCommand, keyboard *input.Keyboard) *InputBox {
	inputBox := &InputBox{}
	inputBox.font = text.NewUltimaFont(text.GetScaledNumberToResolution(inputBoxFontPoint))

	inputBox.Question = question

	textCommands := &grammar.TextCommands{
		*textCommand,
	}

	inputBox.textQuestion = text.NewOutput(inputBox.font,
		inputBoxDefaultLineSpacing,
		inputBoxMaxLinesForOutput,
		inputBoxMaxCharsPerLine)

	inputBox.inputBoxPercents = sprites.PercentBasedPlacement{
		StartPercentX: 0 + inputPercentIntoBorder + inputBoxPercentOffEdge,
		EndPercentX:   .75 + .01 - inputPercentIntoBorder,
		StartPercentY: inputBoxStartPercentY + .05,
		EndPercentY:   inputBoxStartPercentY + .05 + 0.05,
	}

	inputBox.textQuestion.AddRowStr(inputBox.Question)
	questionStr := inputBox.textQuestion.GetOutputStr(false)
	nQuestionRows := strings.Count(questionStr, "\n") + 1

	inputBox.textInput = NewTextInput(
		inputBox.inputBoxPercents,
		text.GetScaledNumberToResolution(inputBoxFontPoint), 20,
		textCommands,
		TextInputCallbacks{
			AmbiguousAutoComplete: func(message string) {
			},
		},
		keyboard)

	heightOfText := -float64(inputBoxDefaultLineSpacing*(nQuestionRows)) + (inputBoxDefaultLineSpacing * 0.5)
	screenResolution := config.GetWindowResolutionFromEbiten()
	percentTextHeight := heightOfText / float64(screenResolution.Y)

	inputBox.borderBoxPercents = sprites.PercentBasedPlacement{
		StartPercentX: 0 + inputBoxPercentOffEdge,
		EndPercentX:   .75 + .01 - inputBoxPercentOffEdge,
		StartPercentY: inputBoxStartPercentY + percentTextHeight,
		EndPercentY:   inputBoxStartPercentY + .1,
	}

	inputBox.border = NewBorder(
		inputBox.borderBoxPercents,
		inputBoxBorderWidthScaling,
	)

	return inputBox
}

func (i *InputBox) Draw(screen *ebiten.Image) {

	i.border.Draw(screen)
	i.textInput.Draw(screen)
	i.textQuestion.DrawText(
		screen,
		i.textQuestion.GetOutputStr(false),
		i.getTextQuestionDrawOptions())
}

func (i *InputBox) Update() {
	i.textInput.Update()
}

func (i *InputBox) GetText() string {
	return i.textInput.GetText()
}

func (i *InputBox) getTextQuestionDrawOptions() *ebiten.DrawImageOptions {
	dop := &ebiten.DrawImageOptions{}

	leftTextX, leftTextY := sprites.GetTranslateXYByPercent(
		sprites.PercentBasedCenterPoint{X: i.borderBoxPercents.StartPercentX + inputPercentIntoBorder,
			Y: i.borderBoxPercents.StartPercentY + inputPercentIntoBorder + 0.01})
	dop.GeoM.Translate(leftTextX, leftTextY)

	return dop
}
