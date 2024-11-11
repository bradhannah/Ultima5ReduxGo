package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
)

const borderWidthScaling = 1202

type ButtonSize int

const (
	SmallButton ButtonSize = iota
	MediumButton
	LargeButton
)

var mediumButtonPercentWidth = 0.1
var mediumButtonPercentHeight = 0.05

type Button struct {
	text        string
	callback    func()
	border      *Border
	size        ButtonSize
	centerPoint sprites.PercentBasedCenterPoint
	pbp         sprites.PercentBasedPlacement
}

func NewButton(
	text string,
	callback func(),
	centerPoint sprites.PercentBasedCenterPoint,
	size ButtonSize,
) *Button {
	button := &Button{}
	button.text = text
	button.callback = callback
	button.size = size
	button.setPosition(centerPoint)
	button.initializeBorder()

	return button
}

func (b *Button) setPosition(centerPoint sprites.PercentBasedCenterPoint) {
	b.centerPoint = centerPoint
	switch b.size {
	case MediumButton:
		b.pbp = sprites.PercentBasedPlacement{
			StartPercentX: centerPoint.X - mediumButtonPercentWidth/2,
			EndPercentX:   centerPoint.X + mediumButtonPercentWidth/2,
			StartPercentY: centerPoint.Y - mediumButtonPercentHeight/2,
			EndPercentY:   centerPoint.Y + mediumButtonPercentHeight/2,
		}
	default:
		panic("unhandled default case")
	}
	b.initializeBorder()
}

func (b *Button) initializeBorder() {
	b.border = NewBorder(
		b.pbp,
		borderWidthScaling)
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.border.Draw(screen)
}

func (b *Button) Update() {
}
