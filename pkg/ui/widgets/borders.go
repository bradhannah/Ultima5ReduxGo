package widgets

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
)

type Border struct {
	border            *ebiten.Image
	borderDrawOptions *ebiten.DrawImageOptions

	background            *ebiten.Image
	backgroundDrawOptions *ebiten.DrawImageOptions

	bordersWithScaling int

	interiorColor color.Color
	interiorAlpha float32
	interiorRect  *image.Rectangle
}

func (b *Border) Draw(screen *ebiten.Image) {
	b.DrawBackground(screen)
	b.DrawBorder(screen)
}

func (b *Border) DrawBackground(screen *ebiten.Image) {
	screen.DrawImage(b.background, b.backgroundDrawOptions)
}

func (b *Border) DrawBorder(screen *ebiten.Image) {
	screen.DrawImage(b.border, b.borderDrawOptions)
}

func (b *Border) Update() {
}

func NewBorder(percentBasedPlacement sprites.PercentBasedPlacement, bordersWithScaling int, interiorColor color.Color) *Border {
	b := &Border{}
	b.bordersWithScaling = bordersWithScaling

	b.interiorColor = interiorColor
	b.interiorAlpha = 0.85

	b.initializeBorders(percentBasedPlacement)
	b.initializeInterior(percentBasedPlacement)
	b.fillInternalRect()
	return b
}

func (b *Border) SetInteriorColor(color color.Color) {

}

func (b *Border) initializeBorders(percentBasedPlacement sprites.PercentBasedPlacement) {
	// todo: generalize this function for multiple uses - it's the same as the debug console
	mainBorder := sprites.NewBorderSprites()
	b.border, b.borderDrawOptions = mainBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(b.bordersWithScaling, percentBasedPlacement)
}

func (b *Border) initializeInterior(percentBasedPlacement sprites.PercentBasedPlacement) {
	b.backgroundDrawOptions = &ebiten.DrawImageOptions{}
	*b.backgroundDrawOptions = *b.borderDrawOptions
	backgroundPercents := percentBasedPlacement
	b.interiorRect = sprites.GetRectangleFromPercents(backgroundPercents)
	b.background = ebiten.NewImageWithOptions(*b.interiorRect, &ebiten.NewImageOptions{})
}

func (b *Border) fillInternalRect() {
	xDiff := float32(b.interiorRect.Dx()) * 0.01
	yDiff := float32(b.interiorRect.Dy()) * 0.01
	vector.DrawFilledRect(b.background,
		float32(b.interiorRect.Min.X)+xDiff,
		float32(b.interiorRect.Min.Y)+yDiff,
		float32(b.interiorRect.Dx())-xDiff*2,
		float32(b.interiorRect.Dy())-yDiff*2,
		b.interiorColor,
		// u_color.Black,
		false)

	b.backgroundDrawOptions.ColorScale.ScaleAlpha(b.interiorAlpha)
}
