package text

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"strings"
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

//func (o *Output) AddToContinuousOutput(outputStr string) {
//	o.lines[o.nextLineToIndex] = outputStr
//	o.nextLineToIndex = (o.nextLineToIndex + 1) % maxLines
//}

func (o *Output) AddToContinuousOutput(outputStr string) {
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
			o.nextLineToIndex = (o.nextLineToIndex + 1) % maxLines
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

func (o *Output) AppendToOutput(outputStr string) {
	lastLineAdded := (o.nextLineToIndex - 1) % maxLines
	if lastLineAdded < 0 {
		lastLineAdded = maxLines - 1
	}
	currentStr := o.lines[lastLineAdded]

	o.nextLineToIndex = lastLineAdded
	o.AddToContinuousOutput(currentStr + outputStr)
}
