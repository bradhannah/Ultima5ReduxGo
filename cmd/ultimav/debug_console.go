package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	e_text "github.com/hajimehoshi/ebiten/v2/text/v2"

	u_color "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ui/widgets"
)

var _ widgets.Widget = &DebugConsole{}

const (
	debugPercentOffEdge    = 0.0
	debugPercentIntoBorder = 0.02

	debugFontPoint       = 13
	debugFontLineSpacing = 17
	debugMaxLines        = 14
	debugMaxCharsInput   = 90
)

type DebugConsole struct {
	border *widgets.Border

	font   *text.UltimaFont
	Output *text.Output

	TextInput *widgets.TextInput

	gameScene *GameScene
}

func NewDebugConsole(gameScene *GameScene) *DebugConsole {
	debugConsole := DebugConsole{}
	debugConsole.gameScene = gameScene
	debugConsole.initializeResizeableVisualElements()

	return &debugConsole
}

func (d *DebugConsole) Refresh() {
	d.initializeResizeableVisualElements()
}

func (d *DebugConsole) initializeResizeableVisualElements() {
	d.border = widgets.NewBorder(
		sprites.PercentBasedPlacement{
			StartPercentX: 0 + debugPercentOffEdge,
			EndPercentX:   .75 + .01 - debugPercentOffEdge,
			StartPercentY: .7,
			EndPercentY:   1,
		},
		borderWidthScaling,
		u_color.Black)

	d.font = text.NewUltimaFont(text.GetScaledNumberToResolution(debugFontPoint))

	if d.Output == nil {
		d.Output = text.NewOutput(d.font,
			text.GetScaledNumberToResolution(debugFontLineSpacing),
			debugMaxLines,
			debugMaxCharsInput)
	} else {
		d.Output.SetFont(
			d.font,
			text.GetScaledNumberToResolution(debugFontLineSpacing),
		)
	}
	if d.TextInput == nil {
		d.TextInput = widgets.NewTextInput(
			sprites.PercentBasedPlacement{
				StartPercentX: 0 + debugPercentIntoBorder,
				EndPercentX:   .75 + .01 - debugPercentIntoBorder,
				StartPercentY: .955,
				EndPercentY:   1,
			},
			text.GetScaledNumberToResolution(debugFontPoint),
			debugMaxCharsInput,
			d.createDebugFunctions(d.gameScene),
			widgets.TextInputCallbacks{
				AmbiguousAutoComplete: func(message string) {
					d.Output.AddRowStr(message)
				},
			},
			d.gameScene.keyboard)
	} else {
		d.TextInput.SetFontPoint(text.GetScaledNumberToResolution(debugFontPoint))
	}
}

func (d *DebugConsole) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyBackquote) || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if !d.gameScene.keyboard.TryToRegisterKeyPress(ebiten.KeyBackquote) {
			return
		}
		d.gameScene.dialogStack.PopModalDialog()
	} else {
		d.TextInput.Update()
	}
}

func (d *DebugConsole) Draw(screen *ebiten.Image) {
	d.border.DrawBackground(screen)

	const percentIntoBorder = 0.02
	textRect := sprites.GetRectangleFromPercents(sprites.PercentBasedPlacement{
		StartPercentX: 0 + percentIntoBorder,
		EndPercentX:   .75 + .01 - percentIntoBorder,
		StartPercentY: .7 + 0.03,
		EndPercentY:   1,
	})
	d.Output.DrawContinuousOutputTexOnXy(screen, image.Point{
		X: textRect.Min.X,
		Y: textRect.Min.Y,
	}, false, e_text.AlignStart, e_text.AlignStart)
	d.border.DrawBorder(screen)

	d.TextInput.Draw(screen)
}
