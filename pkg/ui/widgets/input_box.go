package widgets

import (
	u_color "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const percentOffEdge = 0.2
const debugFontPoint = 13
const debugFontLineSpacing = 17
const maxDebugLines = 14
const maxCharsForInput = 90
const borderWidthScaling = 601

type InputBox struct {
	textInput *TextInput
	Question  string

	border            *ebiten.Image
	borderDrawOptions *ebiten.DrawImageOptions

	background            *ebiten.Image
	backgroundDrawOptions *ebiten.DrawImageOptions

	font *text.UltimaFont
}

func (i *InputBox) Draw(screen *ebiten.Image) {
	screen.DrawImage(i.background, i.backgroundDrawOptions)
	const percentIntoBorder = 0.02
	//textRect := sprites.GetRectangleFromPercents(sprites.PercentBasedPlacement{
	//	StartPercentX: 0 + percentIntoBorder,
	//	EndPercentX:   .75 + .01 - percentIntoBorder,
	//	StartPercentY: .7 + 0.03,
	//	EndPercentY:   1,
	//})
	//i.Output.DrawContinuousOutputTexOnXy(screen, image.Point{
	//	X: textRect.Min.X,
	//	Y: textRect.Min.Y,
	//}, false)
	screen.DrawImage(i.border, i.borderDrawOptions)

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

	inputBox.textInput = NewTextInput(20, 20,
		//debugConsole.createDebugFunctions(gameScene),
		textCommands,
		TextInputCallbacks{
			AmbiguousAutoComplete: func(message string) {
				//debugConsole.Output.AddRowStr(message)
			},
		})

	inputBox.initializeDebugBorders()
	return inputBox
}

func (i *InputBox) initializeDebugBorders() {
	// todo: generalize this function for multiple uses - it's the same as the debug console
	mainBorder := sprites.NewBorderSprites()
	percentBased := sprites.PercentBasedPlacement{
		StartPercentX: 0 + percentOffEdge,
		EndPercentX:   .75 + .01 - percentOffEdge,
		StartPercentY: .65,
		EndPercentY:   .9,
	}
	i.border, i.borderDrawOptions = mainBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, percentBased)

	i.backgroundDrawOptions = &ebiten.DrawImageOptions{}
	*i.backgroundDrawOptions = *i.borderDrawOptions

	backgroundPercents := percentBased

	rect := sprites.GetRectangleFromPercents(backgroundPercents)

	i.background = ebiten.NewImageWithOptions(*rect, &ebiten.NewImageOptions{})
	xDiff := float32(rect.Dx()) * 0.01
	yDiff := float32(rect.Dy()) * 0.01
	vector.DrawFilledRect(i.background,
		float32(rect.Min.X)+xDiff,
		float32(rect.Min.Y)+yDiff,
		float32(rect.Dx())-xDiff*2,
		float32(rect.Dy())-yDiff*2,
		u_color.Black,
		false)

	i.backgroundDrawOptions.ColorScale.ScaleAlpha(.85)
}
