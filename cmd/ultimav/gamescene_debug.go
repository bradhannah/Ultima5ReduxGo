package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ui/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const percentOffEdge = 0.0
const debugFontPoint = 13
const debugFontLineSpacing = 17
const maxDebugLines = 14
const maxCharsForInput = 90

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
	debugConsole.border = widgets.NewBorder(sprites.PercentBasedPlacement{
		StartPercentX: 0 + percentOffEdge,
		EndPercentX:   .75 + .01 - percentOffEdge,
		StartPercentY: .7,
		EndPercentY:   1,
	},
		borderWidthScaling)
	debugConsole.font = text.NewUltimaFont(debugFontPoint)
	debugConsole.Output = text.NewOutput(debugConsole.font, debugFontLineSpacing, maxDebugLines, maxCharsForInput)
	debugConsole.TextInput = widgets.NewTextInput(
		debugFontPoint,
		maxCharsForInput,
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
