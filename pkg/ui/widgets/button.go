package widgets

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	e_text "github.com/hajimehoshi/ebiten/v2/text/v2"

	u_color "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
)

const borderWidthScaling = 1202

type ButtonSize int

const (
	SmallButton ButtonSize = iota
	MediumButton
	LargeButton
)

const (
	SmallButtonFontPoint  = 15
	MediumButtonFontPoint = 20
	LargeButtonFontPoint  = 30
)

const (
	SmallButtonLineSpacing  = SmallButtonFontPoint
	MediumButtonLineSpacing = MediumButtonFontPoint
	LargeButtonLineSpacing  = LargeButtonFontPoint
)

type ButtonStatus int

const (
	NotSelected ButtonStatus = iota
	Selected
	Disabled
)

type Button struct {
	text         string
	callback     func()
	border       *Border
	size         ButtonSize
	centerPoint  sprites.PercentBasedCenterPoint
	pbp          sprites.PercentBasedPlacement
	buttonStatus ButtonStatus

	font   *text.UltimaFont
	Output *text.Output

	interiorColor color.Color
}

func NewButton(
	buttonText string,
	callback func(),
	centerPoint sprites.PercentBasedCenterPoint,
	size ButtonSize,
) *Button {
	button := &Button{}
	button.text = buttonText
	button.callback = callback
	button.size = size
	button.interiorColor = color.Black
	button.setPosition(centerPoint)
	button.initializeBorder(button.interiorColor)

	button.font = text.NewUltimaFont(text.GetScaledNumberToResolution(MediumButtonFontPoint))
	button.Output = text.NewOutput(button.font, text.GetScaledNumberToResolution(MediumButtonFontPoint), 1, 20)
	button.Output.SetFont(
		button.font,
		text.GetScaledNumberToResolution(MediumButtonLineSpacing),
	)

	button.Output.AddRowStr(buttonText)
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
	b.initializeBorder(b.interiorColor)
}

func (b *Button) initializeBorder(color color.Color) {
	b.border = NewBorder(
		b.pbp,
		borderWidthScaling,
		color,
	)
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.border.Draw(screen)

	const percentIntoBorder = 0.02
	textRect := sprites.GetRectangleFromPercents(
		// sprites.PercentBasedPlacement{
		b.pbp,
		// StartPercentX: 0 + percentIntoBorder,
		// EndPercentX:   .75 + .01 - percentIntoBorder,
		// StartPercentY: .7 + 0.03,
		// EndPercentY:   1,
		// }
	)

	b.Output.DrawContinuousOutputTexOnXy(screen, image.Point{
		X: textRect.Min.X + (textRect.Dx() / 2),
		Y: textRect.Min.Y + (textRect.Dy() / 2),
		// X: textRect.Min.X,
		// Y: textRect.Min.Y,
	}, false, e_text.AlignCenter, e_text.AlignCenter)
}

func (b *Button) Update() {
}

func (b *Button) SetButtonStatus(buttonStatus ButtonStatus) {
	b.buttonStatus = buttonStatus
	switch buttonStatus {
	case Selected:
		b.interiorColor = u_color.UltimaBlue
	case NotSelected:
		b.interiorColor = u_color.LighterBlackSemi
	}
	b.initializeBorder(b.interiorColor)
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
