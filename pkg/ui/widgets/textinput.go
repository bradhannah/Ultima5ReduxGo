package widgets

import (
	_ "embed"
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
)

import "golang.org/x/image/font"

// todo: cursor blink to show position
// todo: different background for text?
// todo: autocomplete?
// todo:

const (
	keyPressDelay = 165
	cursorWidth   = 10
)

var (
	//go:embed assets/ultima-v-warriors-of-destiny/ultima-v-warriors-of-destiny.ttf
	rawUltimaTTF []byte
)

type TextInput struct {
	output     *text.Output
	ultimaFont *text.UltimaFont
	face       font.Face

	maxCharsPerLine int
	hasFocus        bool

	keyboard *input.Keyboard
}

const percentIntoBorder = 0.02

var mainTextPlacement = sprites.PercentBasedPlacement{
	StartPercentX: 0 + percentIntoBorder,
	EndPercentX:   .75 + .01 - percentIntoBorder,
	StartPercentY: .955,
	EndPercentY:   1,
}

var textColor = color.RGBA{
	R: 0,
	G: 255,
	B: 0,
	A: 255,
}

var nonAlphaNumericBoundKeys = []ebiten.Key{ebiten.KeyDown,
	ebiten.KeyEnter,
	ebiten.KeySpace,
	ebiten.KeyBackspace,
	ebiten.KeyDelete,
	ebiten.KeyEscape,
}

func NewTextInput(fontPointSize float64, maxCharsPerLine int) *TextInput {
	textInput := &TextInput{}
	textInput.ultimaFont = text.NewUltimaFont(fontPointSize)
	textInput.maxCharsPerLine = maxCharsPerLine
	textInput.keyboard = input.NewKeyboard(keyPressDelay)

	// NOTE: single line input only (for now?)
	textInput.output = text.NewOutput(textInput.ultimaFont, 0, 1, maxCharsPerLine)
	textInput.output.AddRowStr("")
	textInput.output.SetColor(textColor)

	textInput.getTtf()

	return textInput
}

func (t *TextInput) getTtf() {
	// Load the font
	//ttf, err := ioutil.ReadFile("path/to/your/fontfile.ttf")
	//if err != nil {
	//	log.Fatal(err)
	//}

	tt, err := opentype.Parse(rawUltimaTTF)
	if err != nil {
		log.Fatal(err)
	}

	t.face, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14, // Set font size here
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (t *TextInput) Draw(screen *ebiten.Image) {
	textRect := sprites.GetRectangleFromPercents(mainTextPlacement)
	t.output.DrawContinuousOutputTexOnXy(screen, image.Point{
		X: textRect.Min.X,
		Y: textRect.Min.Y,
	}, false)
	t.drawCursor(screen)
}

func (t *TextInput) drawCursor(screen *ebiten.Image) {
	var cursorPlacement = sprites.PercentBasedPlacement{
		StartPercentX: 0 + percentIntoBorder,
		EndPercentX:   .75 + .01 - percentIntoBorder,
		StartPercentY: .955,
		EndPercentY:   .97,
	}
	textRect := sprites.GetRectangleFromPercents(cursorPlacement)
	width := t.CalculateTextWidth(t.output.GetOutputStr(false))
	vector.DrawFilledRect(screen,
		float32(textRect.Min.X)+float32(width),
		float32(textRect.Min.Y),
		cursorWidth,
		float32(textRect.Dy()),
		textColor,
		false)
}

func (t *TextInput) GetText() string {
	return ""
	//return t.textArea.GetText()
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
	case ebiten.KeySpace:
		t.addCharacter(" ")
	case ebiten.KeyBackspace:
		t.output.RemoveRightSideCharacter()
	}
	return
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
