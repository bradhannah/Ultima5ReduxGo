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
			StartPercentX: centerPoint.X - b.GetWidthPercent()/2,
			EndPercentX:   centerPoint.X + b.GetWidthPercent()/2,
			StartPercentY: centerPoint.Y - b.GetHeightPercent()/2,
			EndPercentY:   centerPoint.Y + b.GetHeightPercent()/2,
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

func (b *Button) GetHeightPercent() float64 {
	return GetButtonHeightPercent(b.size)
}

func (b *Button) GetWidthPercent() float64 {
	return GetButtonWidthPercent(b.size)
}

func GetButtonWidthPercent(size ButtonSize) float64 {
	var mediumButtonPercentWidth = 0.2

	switch size {
	case SmallButton:
		return 0
	case MediumButton:
		return mediumButtonPercentWidth
	case LargeButton:
		return 0
	}
	return -1
}

func GetButtonHeightPercent(size ButtonSize) float64 {
	var mediumButtonPercentHeight = 0.06
	switch size {
	case SmallButton:
		return 0
	case MediumButton:
		return mediumButtonPercentHeight
	case LargeButton:
		return 0
	}
	return -1
}
