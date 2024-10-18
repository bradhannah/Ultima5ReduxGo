package widgets

import (
	_ "embed"
	u_color "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"log"
	"time"
)

import "golang.org/x/image/font"

// todo: different background for text?
// todo: accelerate text after a period of time
// todo: autocomplete?
// todo:

// /teleport/ , /0-255/ , /0-255/
// /tplace/ ,

const (
	keyPressDelay            = 165
	cursorWidth              = 10
	cursorBlinkRateInMs      = 1000
	pauseBlinkKeyPressedInMs = 750
)

var (
	//go:embed assets/ultima-v-warriors-of-destiny/ultima-v-warriors-of-destiny.ttf
	rawUltimaTTF []byte
)

type TextInputColors struct {
	DefaultColor          color.Color
	NoMatchesColor        color.Color
	OneMatchColor         color.Color
	MoreThanOneMatchColor color.Color
}

type TextInput struct {
	output     *text.Output
	ultimaFont *text.UltimaFont
	face       font.Face

	maxCharsPerLine int
	hasFocus        bool

	inputColors TextInputColors

	textCommands *grammar.TextCommands

	keyboard *input.Keyboard
}

const percentIntoBorder = 0.02

var mainTextPlacement = sprites.PercentBasedPlacement{
	StartPercentX: 0 + percentIntoBorder,
	EndPercentX:   .75 + .01 - percentIntoBorder,
	StartPercentY: .955,
	EndPercentY:   1,
}

var nonAlphaNumericBoundKeys = []ebiten.Key{ebiten.KeyDown,
	ebiten.KeyEnter,
	ebiten.KeySpace,
	ebiten.KeyBackspace,
	ebiten.KeyDelete,
	ebiten.KeyEscape,
	ebiten.KeyTab,
}

func NewTextInput(fontPointSize float64, maxCharsPerLine int, textCommands *grammar.TextCommands) *TextInput {
	textInput := &TextInput{}
	textInput.ultimaFont = text.NewUltimaFont(fontPointSize)
	textInput.maxCharsPerLine = maxCharsPerLine
	textInput.keyboard = input.NewKeyboard(keyPressDelay)
	textInput.textCommands = textCommands

	// NOTE: single line input only (for now?)
	textInput.output = text.NewOutput(textInput.ultimaFont, 0, 1, maxCharsPerLine)
	textInput.output.AddRowStr("")

	textInput.inputColors = TextInputColors{
		DefaultColor:          color.White,
		NoMatchesColor:        u_color.Red,
		OneMatchColor:         u_color.Green,
		MoreThanOneMatchColor: u_color.Yellow,
	}

	textInput.getAndSetTtf(fontPointSize)

	return textInput
}

//func (t *TextInput) SetColor(c color.Color) {
//	t.defaultColor = c
//	t.output.SetColor(t.defaultColor)
//}

func (t *TextInput) getAndSetTtf(fontPointSize float64) {
	tt, err := opentype.Parse(rawUltimaTTF)
	if err != nil {
		log.Fatal(err)
	}

	t.face, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontPointSize, // Set font size here
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (t *TextInput) Draw(screen *ebiten.Image) {
	textRect := sprites.GetRectangleFromPercents(mainTextPlacement)
	t.output.SetColor(t.getTextColor())
	t.output.DrawContinuousOutputTexOnXy(screen, image.Point{
		X: textRect.Min.X,
		Y: textRect.Min.Y,
	}, false)
	t.drawCursor(screen)
}

func (t *TextInput) drawCursor(screen *ebiten.Image) {
	if t.keyboard.GetMsSinceLastKeyPress() > pauseBlinkKeyPressedInMs {
		blink := time.Now().UnixMilli() % cursorBlinkRateInMs
		if blink <= cursorBlinkRateInMs/2 {
			return
		}
	}

	var cursorPlacement = sprites.PercentBasedPlacement{
		StartPercentX: 0 + percentIntoBorder,
		EndPercentX:   .75 + .01 - percentIntoBorder,
		StartPercentY: .952,
		EndPercentY:   .968,
	}
	textRect := sprites.GetRectangleFromPercents(cursorPlacement)
	width := t.CalculateTextWidth(t.output.GetOutputStr(false))
	vector.DrawFilledRect(screen,
		float32(textRect.Min.X)+float32(width),
		float32(textRect.Min.Y),
		cursorWidth,
		float32(textRect.Dy()),
		t.getTextColor(),
		false)
}

func (t *TextInput) getTextColor() color.Color {
	if t.textCommands == nil {
		return t.inputColors.DefaultColor
	}

	// before we draw the text input - first check if the text is valid
	if t.textCommands.OneOrMoreCommandsMatch(t.GetText()) {
		return t.inputColors.OneMatchColor
	} else {
		return t.inputColors.NoMatchesColor
	}

}

func (t *TextInput) GetText() string {
	return t.output.GetOutputStr(false)
}

func (t *TextInput) addCharacter(keyStr string) {
	if t.output.GetLengthOfCurrentRowStr() < t.maxCharsPerLine {
		t.output.AppendToCurrentRowStr(keyStr)
	}
}

func (t *TextInput) removeRightCharacter(keyStr string) {

}

func (t *TextInput) Update() {
	//if !t.hasFocus {
	//	return
	//}

	var boundKey *ebiten.Key
	key, keyStr := t.keyboard.GetAlphaNumericPressed()
	if key != nil {
		if !t.keyboard.TryToRegisterKeyPress(*key) {
			return
		}
		t.addCharacter(keyStr)
		return
	} else {
		// the idea here is if the user pushes, then lets go and repushes the same button it should
		// show both characters since it is not technically repeating
		boundKey = t.keyboard.GetBoundKeyPressed(&nonAlphaNumericBoundKeys)
		if boundKey == nil {
			t.keyboard.SetAllowKeyPressImmediately()
			return
		}
		if !t.keyboard.TryToRegisterKeyPress(*boundKey) {
			return
		}
	}

	switch *boundKey {
	case ebiten.KeyEnter:
		return
	case ebiten.KeySpace, ebiten.KeyTab:
		{
			if !t.hasTextCommands() {
				// fill in the autocomplete
				t.addCharacter(" ")
				break
			}
			outputStr := t.output.GetOutputStr(false)
			nMatch := t.textCommands.HowManyCommandsMatch(outputStr)
			if nMatch == 1 {
				// need to feed into autocomplete
				matches := t.textCommands.GetAllMatchingTextCommands(outputStr)
				if len(*matches) > 1 {
					log.Fatal("Should only be a single match")
				}

				t.output.AppendToCurrentRowStr((*matches)[0].GetAutoComplete(outputStr) + " ")
			}
			// do nothing - it must have a match
		}
	case ebiten.KeyBackspace:
		t.output.RemoveRightSideCharacter()
	}
	return
}

func (t *TextInput) hasTextCommands() bool {
	return len(*t.textCommands) > 0
}

// CalculateTextWidth returns the total width of the string in pixels
func (t *TextInput) CalculateTextWidth(s string) int {
	var width fixed.Int26_6
	for _, r := range s {
		awidth, _ := t.face.GlyphAdvance(r)
		width += awidth
	}
	return width.Ceil() // Convert from fixed.Int26_6 to int
}
