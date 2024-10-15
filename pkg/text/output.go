package text

import (
	u_color "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"strings"
)

type Output struct {
	Font            *UltimaFont
	lines           []string
	nextLineToIndex int
	lineSpacing     float64
	maxLines        int
	color           color.Color
}

const (
// maxCharsPerLine = 30
// maxLines        = 10
// OutputFontPoint = 20
)

func NewOutput(font *UltimaFont, lineSpacing float64, maxLines int) *Output {
	output := &Output{}
	output.Font = font
	output.lineSpacing = lineSpacing
	output.maxLines = maxLines
	output.lines = make([]string, maxLines)
	output.SetColor(u_color.White)
	return output
}

func (o *Output) SetColor(color color.Color) {
	o.color = color
}

func (o *Output) DrawText(screen *ebiten.Image, textStr string, op *ebiten.DrawImageOptions) {
	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing:  o.lineSpacing,
			PrimaryAlign: text.AlignStart,
		},
	}

	text.Draw(screen, textStr, o.Font.TextFace, &dop)
}

func (o *Output) DrawTextCenter(screen *ebiten.Image, textStr string, op *ebiten.DrawImageOptions) {
	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing:  o.lineSpacing,
			PrimaryAlign: text.AlignCenter,
		},
	}

	text.Draw(screen, textStr, o.Font.TextFace, &dop)
}

func (o *Output) DrawTextRightToLeft(screen *ebiten.Image, textStr string, op *ebiten.DrawImageOptions) {
	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing:  o.lineSpacing,
			PrimaryAlign: text.AlignEnd,
		},
	}

	text.Draw(screen, textStr, o.Font.TextFace, &dop)
}

//func (o *Output) AddRowStr(outputStr string) {
//	o.lines[o.nextLineToIndex] = outputStr
//	o.nextLineToIndex = (o.nextLineToIndex + 1) % maxLines
//}

func (o *Output) AddRowStr(outputStr string) {
	const maxCharsPerLine = 16

	// Process the string line-by-line, splitting by '\n'
	lines := splitByNewline(outputStr)

	for _, line := range lines {
		// Process each line, splitting by 16 characters or nearest space
		for len(line) > 0 {
			// Find the position to split the string, favoring a space before 16 characters
			splitAt := maxCharsPerLine
			if len(line) < maxCharsPerLine {
				splitAt = len(line) // If the line is shorter, take the entire line
			} else {
				// Try to find the last space before the 16th character
				spaceIndex := -1
				for i := 0; i < maxCharsPerLine; i++ {
					if line[i] == ' ' {
						spaceIndex = i
					}
				}

				// If a space is found, split there
				if spaceIndex != -1 {
					splitAt = spaceIndex
				}
			}

			// Extract the chunk and remove leading spaces from the remaining line
			chunk := line[:splitAt]
			line = line[splitAt:]
			line = trimLeadingSpaces(line) // Helper function to remove leading spaces

			// Add hyphen if the chunk has the maximum length and there is more to process
			if splitAt == maxCharsPerLine && len(line) > 0 {
				chunk += "-"
			}

			// Add the chunk to the output slice
			o.lines[o.nextLineToIndex] = chunk
			o.nextLineToIndex = (o.nextLineToIndex + 1) % o.maxLines
		}
	}
}

// Helper function to split the input string by '\n' while keeping track of empty lines.
func splitByNewline(input string) []string {
	return strings.Split(input, "\n")
}

// Helper function to trim leading spaces
func trimLeadingSpaces(s string) string {
	for len(s) > 0 && s[0] == ' ' {
		s = s[1:]
	}
	return s
}

// DrawRightSideOutput
// Draws to the correct x and y image position for the right side output panel
func (o *Output) DrawRightSideOutput(screen *ebiten.Image) {
	const lineSpacing = 20

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(25, float64(screen.Bounds().Dy())*.52)

	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing: lineSpacing,
		},
	}

	text.Draw(screen, o.getOutputStr(true), o.Font.TextFace, &dop)
}

func (o *Output) DrawContinuousOutputTexOnXy(screen *ebiten.Image, point image.Point, bShowEmptyNewLines bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(point.X), float64(point.Y))

	dop := text.DrawOptions{
		DrawImageOptions: *op,
		LayoutOptions: text.LayoutOptions{
			LineSpacing:  o.lineSpacing,
			PrimaryAlign: text.AlignStart,
		},
	}

	dop.ColorScale.ScaleWithColor(o.color)
	text.Draw(screen, o.getOutputStr(bShowEmptyNewLines), o.Font.TextFace, &dop)
}

func (o *Output) getOutputStr(bShowEmptyNewLines bool) string {
	var outputStr string
	nTotalLines := len(o.lines)
	for i := 0; i < nTotalLines; i++ {
		lineToAdd := o.lines[(i+o.nextLineToIndex)%o.maxLines]

		outputStr += lineToAdd
		if (bShowEmptyNewLines || lineToAdd != "") && i < nTotalLines-1 {
			outputStr += " \n"
		}
	}
	return outputStr
}

func (o *Output) AppendToCurrentRowStr(outputStr string) {
	lastLineAdded := (o.nextLineToIndex - 1) % o.maxLines
	if lastLineAdded < 0 {
		lastLineAdded = o.maxLines - 1
	}
	currentStr := o.lines[lastLineAdded]

	o.nextLineToIndex = lastLineAdded
	o.AddRowStr(currentStr + outputStr)
}
