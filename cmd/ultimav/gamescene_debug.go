package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DebugConsole struct {
	border            *ebiten.Image
	borderDrawOptions *ebiten.DrawImageOptions

	background            *ebiten.Image
	backgroundDrawOptions *ebiten.DrawImageOptions

	gameState *game_state.GameState
}

func NewDebugConsole(gameState *game_state.GameState) *DebugConsole {
	debugConsole := DebugConsole{}
	debugConsole.gameState = gameState
	debugConsole.initializeDebugBorders()
	return &debugConsole
}

func (d *DebugConsole) drawDebugConsole(screen *ebiten.Image) {
	screen.DrawImage(d.background, d.backgroundDrawOptions)
	screen.DrawImage(d.border, d.borderDrawOptions)
}

func (d *DebugConsole) initializeDebugBorders() {
	mainBorder := sprites.NewBorderSprites()
	//const percentOffEdge = 0.04
	const percentOffEdge = 0.0
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
