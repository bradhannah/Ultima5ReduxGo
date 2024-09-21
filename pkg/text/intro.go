package text

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
)

var introChoices = []string{"Journey Onward", "Select Save Game", "Import from Legacy", "Set Data Directory", "Create Character", "Introduction", "Acknowledgments"}

func (f *UltimaFont) DrawIntroChoices(screen *ebiten.Image, nSelection int) {
	//nChoices := len(introChoices)
	//gray := color.RGBA{0x80, 0x80, 0x80, 0xff}
	const lineSpacing = 48
	const percentBetweenLines = 0.04
	const startPercent = 0.63

	for i, choice := range introChoices {
		w, h := text.Measure(choice, f.textFace, lineSpacing)
		startY := startPercent + (percentBetweenLines * float64(i))
		op := sprites.GetYSpriteWithPercents(image.Rectangle{
			Min: image.Point{X: 0, Y: 0},
			Max: image.Point{X: int(w), Y: int(h)},
		}, sprites.PercentYBasedPlacement{
			StartPercentX: startPercent,
			StartPercentY: startY,
			EndPercentY:   startY + .03,
		})
		dop := text.DrawOptions{
			DrawImageOptions: *op,
			LayoutOptions: text.LayoutOptions{
				LineSpacing:  lineSpacing,
				PrimaryAlign: text.AlignCenter,
				//SecondaryAlign: text.AlignCenter,
			},
		}

		if i == nSelection {
			//vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), gray, false)
		}

		text.Draw(screen, choice, f.textFace, &dop)
		op.GeoM.Reset()
	}
	//screenWidth := screen.Bounds().Dx()

	//gray := color.RGBA{0x80, 0x80, 0x80, 0xff}

	{
		const testText = "THIS IS my Ultima V Test 12349!#$[]"

	}
}
