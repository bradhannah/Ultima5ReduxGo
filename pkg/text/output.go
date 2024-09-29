package text

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Output struct {
	Font            *UltimaFont
	lines           [maxLines]string
	nextLineToIndex int
	lineSpacing     float64
}

const (
	maxCharsPerLine = 30
	maxLines        = 10
	OutputFontPoint = 20
)

func NewOutput(font *UltimaFont, lineSpacing float64) *Output {
	output := &Output{}
	output.Font = font
	output.lineSpacing = lineSpacing
	return output
}

func (o *Output) DrawText(screen *ebiten.Image, textStr string, op *ebiten.DrawImageOptions) {
	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing:  o.lineSpacing,
			PrimaryAlign: text.AlignStart,
		},
	}

	text.Draw(screen, textStr, o.Font.textFace, &dop)
}

func (o *Output) DrawTextCenter(screen *ebiten.Image, textStr string, op *ebiten.DrawImageOptions) {
	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing:  o.lineSpacing,
			PrimaryAlign: text.AlignCenter,
		},
	}

	text.Draw(screen, textStr, o.Font.textFace, &dop)
}

func (o *Output) DrawTextRightToLeft(screen *ebiten.Image, textStr string, op *ebiten.DrawImageOptions) {
	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing:  o.lineSpacing,
			PrimaryAlign: text.AlignEnd,
		},
	}

	text.Draw(screen, textStr, o.Font.textFace, &dop)
}

func (o *Output) AddToContinuousOutput(outputStr string) {
	o.lines[o.nextLineToIndex] = outputStr
	o.nextLineToIndex = (o.nextLineToIndex + 1) % maxLines
}

func (o *Output) DrawContinuousOutputText(screen *ebiten.Image) {
	const lineSpacing = 20

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(25, float64(screen.Bounds().Dy())*.52)

	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing: lineSpacing,
		},
	}

	text.Draw(screen, o.getOutputStr(), o.Font.textFace, &dop)
}

func (o *Output) getOutputStr() string {
	var outputStr string
	for i := 0; i < len(o.lines); i++ {
		lineToAdd := o.lines[(i+o.nextLineToIndex)%maxLines]
		if i != 0 {
			outputStr += " \n"
		}
		outputStr += lineToAdd
	}
	return outputStr
}
