package text

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/color"
)

var IntroChoices = []string{"Journey Onward", "Select Save Game", "Import from Legacy", "Set Data Directory", "Create Character", "Introduction", "Acknowledgments"}

func (f *UltimaFont) DrawIntroChoices(screen *ebiten.Image, nSelection int) {
	// nChoices := len(IntroChoices)
	const lineSpacing = 48
	const percentBetweenLines = 0.0435
	const startXFontPercent = 0.655
	const endPercentYDiff = 0.035

	for i, choice := range IntroChoices {
		w, h := text.Measure(choice, f.TextFace, lineSpacing)
		startYPercent := startXFontPercent + (percentBetweenLines * float64(i))
		op, _ := sprites.GetYSpriteWithPercents(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(w), Y: int(h)},
		}, sprites.PercentYBasedPlacement{
			StartPercentY: startYPercent,
			EndPercentY:   startYPercent + endPercentYDiff,
		})

		bIsSelected := i == nSelection

		if bIsSelected {
			op.ColorScale.ScaleWithColor(color.Black)
		} else {
			op.ColorScale.ScaleWithColor(color.White)
		}

		dop := text.DrawOptions{
			DrawImageOptions: *op,
			LayoutOptions: text.LayoutOptions{
				LineSpacing:  lineSpacing,
				PrimaryAlign: text.AlignCenter,
			},
		}

		if bIsSelected {
			text.Measure(choice, f.TextFace, lineSpacing)
			const startXHighlightPercent = 0.2
			const endXHighlightPercent = 1 - startXHighlightPercent
			const extraYHeightPercent = 0.003
			rect := sprites.GetRectangleFromPercents(sprites.PercentBasedPlacement{
				StartPercentX: startXHighlightPercent,
				EndPercentX:   endXHighlightPercent,
				StartPercentY: startYPercent - extraYHeightPercent,
				EndPercentY:   startYPercent + endPercentYDiff + extraYHeightPercent,
			})
			vector.DrawFilledRect(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), color.White, false)
		}

		text.Draw(screen, choice, f.TextFace, &dop)
		op.GeoM.Reset()
	}
}
