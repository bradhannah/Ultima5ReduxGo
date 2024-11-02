package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ui/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

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
	debugConsole.border = widgets.NewBorder(
		sprites.PercentBasedPlacement{
			StartPercentX: 0 + debugPercentOffEdge,
			EndPercentX:   .75 + .01 - debugPercentOffEdge,
			StartPercentY: .7,
			EndPercentY:   1,
		},
		borderWidthScaling)
	debugConsole.font = text.NewUltimaFont(debugFontPoint)
	debugConsole.Output = text.NewOutput(debugConsole.font, debugFontLineSpacing, debugMaxLines, debugMaxCharsInput)
	debugConsole.TextInput = widgets.NewTextInput(
		sprites.PercentBasedPlacement{
			StartPercentX: 0 + debugPercentIntoBorder,
			EndPercentX:   .75 + .01 - debugPercentIntoBorder,
			StartPercentY: .955,
			EndPercentY:   1,
		},
		debugFontPoint,
		debugMaxCharsInput,
		debugConsole.createDebugFunctions(gameScene),
		widgets.TextInputCallbacks{
			AmbiguousAutoComplete: func(message string) {
				debugConsole.Output.AddRowStr(message)
			},
		})
	return &debugConsole
}

func (d *DebugConsole) update() {
	if ebiten.IsKeyPressed(ebiten.KeyBackquote) || ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if !d.gameScene.keyboard.TryToRegisterKeyPress(ebiten.KeyBackquote) {
			return
		}
		d.gameScene.bShowDebugConsole = false
	} else {
		d.TextInput.Update()
	}

	return
}

func (d *DebugConsole) drawDebugConsole(screen *ebiten.Image) {
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
	}, false)
	d.border.DrawBorder(screen)

	d.TextInput.Draw(screen)
}
