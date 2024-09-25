package text

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Output struct {
	Font  *UltimaFont
	lines [maxLines]string
	//totalLines      int
	nextLineToIndex int
}

const (
	maxCharsPerLine = 30
	maxLines        = 11
	OutputFontPoint = 20
)

func NewOutput(font *UltimaFont) *Output {
	output := &Output{}
	output.Font = font
	return output
}

func (o *Output) AddToOutput(outputStr string) {
	o.lines[o.nextLineToIndex] = outputStr
	o.nextLineToIndex = (o.nextLineToIndex + 1) % maxLines
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

func (o *Output) Draw(screen *ebiten.Image) {
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
