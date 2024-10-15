package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ui/widgets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
)

const percentOffEdge = 0.0
const debugFontPoint = 13
const debugFontLineSpacing = 17
const maxDebugLines = 14
const maxCharsForInput = 90

type DebugConsole struct {
	border            *ebiten.Image
	borderDrawOptions *ebiten.DrawImageOptions

	background            *ebiten.Image
	backgroundDrawOptions *ebiten.DrawImageOptions

	font   *text.UltimaFont
	Output *text.Output

	TextInput *widgets.TextInput

	gameScene *GameScene
}

func NewDebugConsole(gameScene *GameScene) *DebugConsole {
	debugConsole := DebugConsole{}
	debugConsole.gameScene = gameScene
	debugConsole.initializeDebugBorders()
	debugConsole.font = text.NewUltimaFont(debugFontPoint)
	debugConsole.Output = text.NewOutput(debugConsole.font, debugFontLineSpacing, maxDebugLines, maxCharsForInput)
	debugConsole.TextInput = widgets.NewTextInput(debugFontPoint, maxCharsForInput)

	return &debugConsole
}

func (d *DebugConsole) update() {
	if ebiten.IsKeyPressed(ebiten.KeyBackquote) {
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
	screen.DrawImage(d.background, d.backgroundDrawOptions)
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
	screen.DrawImage(d.border, d.borderDrawOptions)
	d.TextInput.Draw(screen)
}

func (d *DebugConsole) initializeDebugBorders() {
	mainBorder := sprites.NewBorderSprites()
	percentBased := sprites.PercentBasedPlacement{
		StartPercentX: 0 + percentOffEdge,
		EndPercentX:   .75 + .01 - percentOffEdge,
		StartPercentY: .7,
		EndPercentY:   1,
	}
	d.border, d.borderDrawOptions = mainBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, percentBased)

	d.backgroundDrawOptions = &ebiten.DrawImageOptions{}
	*d.backgroundDrawOptions = *d.borderDrawOptions

	backgroundPercents := percentBased

	rect := sprites.GetRectangleFromPercents(backgroundPercents)

	d.background = ebiten.NewImageWithOptions(*rect, &ebiten.NewImageOptions{})
	xDiff := float32(rect.Dx()) * 0.01
	yDiff := float32(rect.Dy()) * 0.01
	vector.DrawFilledRect(d.background,
		float32(rect.Min.X)+xDiff,
		float32(rect.Min.Y)+yDiff,
		float32(rect.Dx())-xDiff*2,
		float32(rect.Dy())-yDiff*2,
		color.Black,
		false)

	d.backgroundDrawOptions.ColorScale.ScaleAlpha(.85)
}
