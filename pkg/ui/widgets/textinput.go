package widgets

import (
	_ "embed"
	"fmt"
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
	"strings"
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

type TextInputCallbacks struct {
	AmbiguousAutoComplete func(string)
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

	TextInputCallbacks TextInputCallbacks
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
	ebiten.KeyMinus,
}

func NewTextInput(fontPointSize float64, maxCharsPerLine int, textCommands *grammar.TextCommands, callbacks TextInputCallbacks) *TextInput {
	textInput := &TextInput{}
	textInput.ultimaFont = text.NewUltimaFont(fontPointSize)
	textInput.maxCharsPerLine = maxCharsPerLine
	textInput.keyboard = input.NewKeyboard(keyPressDelay)
	textInput.textCommands = textCommands
	textInput.TextInputCallbacks = callbacks

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
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if ebiten.IsKeyPressed(ebiten.KeyU) {
			t.output.Clear()
			t.keyboard.SetLastKeyPressedNowPlusMs(500)
			t.keyboard.ForceLastKeyPressed(ebiten.KeyU)
		}
		return
	}

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
	case ebiten.KeyMinus:
		t.output.AppendToCurrentRowStr("-")
	case ebiten.KeyEnter:
		// check if it has single full match
		command := t.output.GetOutputStr(false)
		matches := t.textCommands.GetAllPerfectMatches(command)
		if len(*matches) == 1 {
			// call the callback function
			firstMatch := (*matches)[0]
			if firstMatch.ExecuteCommand == nil {
				return
			}
			firstMatch.ExecuteCommand(command, &firstMatch)
		}
		// if it does then run execute
		return
	case ebiten.KeySpace, ebiten.KeyTab:
		t.tryToAutoComplete()
	case ebiten.KeyBackspace:
		t.output.RemoveRightSideCharacter()
	}
	return
}

func (t *TextInput) tryToAutoComplete() {
	outputStr := t.output.GetOutputStr(false)
	nPrimaryCommandMatch := t.textCommands.HowManyCommandsMatch(outputStr)
	primaryCommandMatches := t.textCommands.GetAllMatchingTextCommands(outputStr)

	splitCommandStrings := strings.Split(outputStr, " ")

	if nPrimaryCommandMatch == 1 && len(splitCommandStrings) > (*primaryCommandMatches)[0].GetNumberOfParameters() {
		return
	}
	if nPrimaryCommandMatch == 0 {
		return
	}

	// if the command ends with a space - then we show all available details
	if outputStr[len(outputStr)-1:] == " " {
		autoCompleteMatches := (*primaryCommandMatches)[0].GetAutoComplete(outputStr)
		if len(autoCompleteMatches) > 0 {
			t.TextInputCallbacks.AmbiguousAutoComplete(fmt.Sprintf("Ambigious - %s", strings.Join(autoCompleteMatches, ",")))
		}
		return
	}

	if nPrimaryCommandMatch == 1 {
		autoCompleteMatches := (*primaryCommandMatches)[0].GetAutoComplete(outputStr)
		if len(autoCompleteMatches) == 1 {
			t.output.AppendToCurrentRowStr(strings.TrimSpace((autoCompleteMatches)[0]) + " ")
			return
		}

		commandNames := strings.Join(autoCompleteMatches, ",")
		t.TextInputCallbacks.AmbiguousAutoComplete(fmt.Sprintf("Ambigious - %s", commandNames))
		return
	}
	var commandNames string
	for i, m := range *primaryCommandMatches {
		if i > 0 {
			commandNames += ", "
		}
		commandNames += m.Matches[0].GetString()
	}
	t.TextInputCallbacks.AmbiguousAutoComplete(fmt.Sprintf("Ambigious - %s", commandNames))
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
