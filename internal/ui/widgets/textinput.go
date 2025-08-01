package widgets

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	etext "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/internal/text"
	ucolor "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
)

// todo: different background for text?
// todo: accelerate text after a period of time
// todo: autocomplete?

// /teleport/ , /0-255/ , /0-255/
// /tplace/ ,

const (
	keyPressDelay            = 300
	cursorWidth              = 10
	cursorBlinkRateInMs      = 1000
	pauseBlinkKeyPressedInMs = 750
	percentIntoBorder        = 0.02
)

//go:embed assets/ultima-v-warriors-of-destiny/ultima-v-warriors-of-destiny.ttf
var rawUltimaTTF []byte

type TextInputColors struct {
	DefaultColor          color.Color
	NoMatchesColor        color.Color
	OneMatchColor         color.Color
	MoreThanOneMatchColor color.Color
}

type TextInputCallbacks struct {
	AmbiguousAutoComplete func(string)
	OnEnter               func()
}

type TextInput struct {
	output     *text.Output
	ultimaFont *text.UltimaFont
	face       font.Face

	maxCharsPerLine int
	hasFocus        bool

	inputColors TextInputColors

	textCommands *grammar.TextCommands

	mainTextPlacement sprites.PercentBasedPlacement

	keyboard *input.Keyboard

	TextInputCallbacks TextInputCallbacks
}

var nonAlphaNumericBoundKeys = []ebiten.Key{
	ebiten.KeyEnter,
	ebiten.KeySpace,
	ebiten.KeyBackspace,
	ebiten.KeyDelete,
	ebiten.KeyEscape,
	ebiten.KeyTab,
	ebiten.KeyMinus,
	ebiten.KeyDown,
	ebiten.KeyLeft,
	ebiten.KeyRight,
	ebiten.KeyUp,
}

func NewTextInput(
	mainTextPlacement sprites.PercentBasedPlacement,
	fontPointSize float64,
	maxCharsPerLine int,
	textCommands *grammar.TextCommands,
	callbacks TextInputCallbacks,
	keyboard *input.Keyboard,
) *TextInput {
	textInput := &TextInput{}
	textInput.maxCharsPerLine = maxCharsPerLine
	textInput.keyboard = keyboard
	textInput.textCommands = textCommands
	textInput.TextInputCallbacks = callbacks
	textInput.mainTextPlacement = mainTextPlacement

	// NOTE: single line input only (for now?)
	textInput.output = text.NewOutput(textInput.ultimaFont, 0, 1, maxCharsPerLine)
	textInput.SetFontPoint(fontPointSize)

	textInput.output.AddRowStrWithTrim("")

	textInput.inputColors = TextInputColors{
		DefaultColor:          color.White,
		NoMatchesColor:        ucolor.Red,
		OneMatchColor:         ucolor.Green,
		MoreThanOneMatchColor: ucolor.Yellow,
	}

	return textInput
}

func (t *TextInput) SetInputColors(inputColors TextInputColors) {
	t.inputColors = inputColors
}

func (t *TextInput) SetFontPoint(fontPointSize float64) {
	t.ultimaFont = text.NewUltimaFont(fontPointSize)
	t.getAndSetTtf(fontPointSize)
	t.output.SetFont(t.ultimaFont, fontPointSize)
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
	textRect := sprites.GetRectangleFromPercents(t.mainTextPlacement)
	t.output.SetColor(t.getTextColor())
	t.output.DrawContinuousOutputTexOnXy(screen, image.Point{
		X: textRect.Min.X,
		Y: textRect.Min.Y,
	}, false, etext.AlignStart, etext.AlignStart)
	t.drawCursor(screen)
}

func (t *TextInput) createCursorPlacement() sprites.PercentBasedPlacement {
	cursorPlacement := sprites.PercentBasedPlacement{
		StartPercentX: t.mainTextPlacement.StartPercentX,         // 0 + debugPercentOffEdge
		EndPercentX:   t.mainTextPlacement.EndPercentY,           // .75 + .01 - percentIntoBorder,
		StartPercentY: t.mainTextPlacement.StartPercentY - 0.003, // .955 - .952,
		EndPercentY:   t.mainTextPlacement.EndPercentY - 0.032,   // .968,
	}
	return cursorPlacement
}

func (t *TextInput) drawCursor(screen *ebiten.Image) {
	if t.keyboard.GetMsSinceLastKeyPress() > pauseBlinkKeyPressedInMs {
		blink := time.Now().UnixMilli() % cursorBlinkRateInMs
		if blink <= cursorBlinkRateInMs/2 {
			return
		}
	}

	cursorPlacement := t.createCursorPlacement()
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

func (t *TextInput) ClearText() {
	t.output.Clear()
}

func (t *TextInput) addCharacter(keyStr string) {
	if t.output.GetLengthOfCurrentRowStr() < t.maxCharsPerLine {
		t.output.AppendToCurrentRowStr(keyStr)
	}
}

// func (t *TextInput) removeRightCharacter(keyStr string) {
//
// }

func (t *TextInput) isPerfectMatchWithAutoFillFirstCharacter() bool {
	perfectMatches := *t.textCommands.GetAllPerfectMatches(t.output.GetOutputStr(false))
	if len(perfectMatches) == 1 {
		if perfectMatches[0].Matches[0].ShouldAutofillWithFirstCharacter() {
			return true
		}
	}
	return false
}

func (t *TextInput) ForceWaitForNextKeyPress() {
	t.keyboard.SetLastKeyPressedNowPlusMs(500)
	t.keyboard.ForceLastKeyPressed(ebiten.KeyEnter)
}

func (t *TextInput) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if ebiten.IsKeyPressed(ebiten.KeyU) {
			t.output.Clear()
			t.keyboard.SetForceWaitAnyKey(500)
			//t.keyboard.ForceLastKeyPressed(ebiten.KeyU)
		}
		return
	}

	var boundKey *ebiten.Key
	key, keyStr := t.keyboard.GetAlphaNumericPressed()
	if key != nil {
		if !t.keyboard.TryToRegisterKeyPress(*key) {
			return
		}

		if t.isPerfectMatchWithAutoFillFirstCharacter() {
			// don't add letters if the match is perfect already
			return
		}

		t.addCharacter(keyStr)

		// special block for single character inputs that expand to a full string
		if t.isPerfectMatchWithAutoFillFirstCharacter() {
			t.tryToAutoComplete(false)
		}
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

	switch *boundKey { //nolint:exhaustive
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
			t.keyboard.SetLastKeyPressedNowPlusMs(keyPressDelay)
		}
		// if it does then run execute

		if t.TextInputCallbacks.OnEnter != nil {
			t.TextInputCallbacks.OnEnter()
		}

		return
	case ebiten.KeySpace, ebiten.KeyTab:
		t.tryToAutoComplete(true)
	case ebiten.KeyBackspace:
		if t.isPerfectMatchWithAutoFillFirstCharacter() {
			t.output.Clear()
		} else {
			t.output.RemoveRightSideCharacter()
		}
	}
	return
}

func (t *TextInput) tryToAutoComplete(bAddSpaceAtEndIfPerfectMatch bool) {
	outputStr := t.output.GetOutputStr(false)

	if outputStr == "" {
		t.TextInputCallbacks.AmbiguousAutoComplete(fmt.Sprintf("Available Commands: %s",
			t.textCommands.GetFriendlyListOfAllCommands()))

		return
	}

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
			t.TextInputCallbacks.AmbiguousAutoComplete(fmt.Sprintf("Ambigious - %s", strings.Join(autoCompleteMatches, ", ")))
		}

		return
	}

	if nPrimaryCommandMatch == 1 {
		autoCompleteMatches := (*primaryCommandMatches)[0].GetAutoComplete(outputStr)
		if len(autoCompleteMatches) == 1 {
			t.output.AppendToCurrentRowStr(strings.TrimSpace((autoCompleteMatches)[0]))
			if bAddSpaceAtEndIfPerfectMatch {
				t.output.AppendToCurrentRowStr(" ")
			}
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

func (t *TextInput) CalculateTextHeight(s string) int {
	if s == "" {
		return 0
	}
	rect, _, _ := t.face.GlyphBounds(rune(s[0]))

	return rect.Min.Y.Ceil()
}
