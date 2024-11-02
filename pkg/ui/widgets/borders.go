package widgets

import (
	u_color "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Border struct {
	border            *ebiten.Image
	borderDrawOptions *ebiten.DrawImageOptions

	background            *ebiten.Image
	backgroundDrawOptions *ebiten.DrawImageOptions

	bordersWithScaling int
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

func NewBorder(percentBasedPlacement sprites.PercentBasedPlacement, bordersWithScaling int) *Border {
	b := &Border{}
	b.bordersWithScaling = bordersWithScaling
	b.initializeDebugBorders(percentBasedPlacement)
	return b
}

func (b *Border) initializeDebugBorders(percentBasedPlacement sprites.PercentBasedPlacement) {
	// todo: generalize this function for multiple uses - it's the same as the debug console
	mainBorder := sprites.NewBorderSprites()
	b.border, b.borderDrawOptions = mainBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(b.bordersWithScaling, percentBasedPlacement)

	b.backgroundDrawOptions = &ebiten.DrawImageOptions{}
	*b.backgroundDrawOptions = *b.borderDrawOptions

	backgroundPercents := percentBasedPlacement

	rect := sprites.GetRectangleFromPercents(backgroundPercents)

	b.background = ebiten.NewImageWithOptions(*rect, &ebiten.NewImageOptions{})
	xDiff := float32(rect.Dx()) * 0.01
	yDiff := float32(rect.Dy()) * 0.01
	vector.DrawFilledRect(b.background,
		float32(rect.Min.X)+xDiff,
		float32(rect.Min.Y)+yDiff,
		float32(rect.Dx())-xDiff*2,
		float32(rect.Dy())-yDiff*2,
		u_color.Black,
		false)

	b.backgroundDrawOptions.ColorScale.ScaleAlpha(.85)
}
