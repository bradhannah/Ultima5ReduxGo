package main

import (
	"fmt"
	"image"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/internal/text"
	"github.com/bradhannah/Ultima5ReduxGo/internal/ui/widgets"
	ucolor "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/hajimehoshi/ebiten/v2"
	etext "github.com/hajimehoshi/ebiten/v2/text/v2"
)

var _ widgets.Widget = &TalkDialog{}

const (
	talkFontPoint       = 20
	talkFontLineSpacing = talkFontPoint + 6
	talkMaxLines        = 11
	talkMaxCharsInput   = 45
	//talkMaxCharsInput   = 15
)

type TalkDialog struct {
	border *widgets.Border

	font   *text.UltimaFont
	Output *text.Output

	TextInput *widgets.TextInput

	gameScene *GameScene

	friendly     *map_units.NPCFriendly
	conversation *party_state.Conversation
}

func NewTalkDialog(gameScene *GameScene, friendly *map_units.NPCFriendly) *TalkDialog {
	talkDialog := TalkDialog{}
	talkDialog.gameScene = gameScene
	talkDialog.friendly = friendly
	talk := gameScene.gameState.GameReferences.TalkReferences.GetTalkScriptByNpcIndex(references.Castle, 2)
	talkDialog.conversation = party_state.NewConversation(2, &gameScene.gameState.PartyState, talk)
	talkDialog.conversation.Start()
	talkDialog.initializeResizeableVisualElements()

	return &talkDialog
}

func (d *TalkDialog) Refresh() {
	d.initializeResizeableVisualElements()
}

const (
	talkPercentOffEdge          = 0.1
	percentTextIndentFromBorder = 0.017
)

const (
	borderStartPercentX = talkPercentOffEdge
	borderEndPercentX   = .75 + .01 - talkPercentOffEdge
	borderStartPercentY = .605
	borderEndPercentY   = 0.95
)

const (
	textRectStartPercentX = talkPercentOffEdge + percentTextIndentFromBorder
	textRectEndPercentX   = .75 + .01 - talkPercentOffEdge
	textRectStartPercentY = borderStartPercentY + 0.03
	textRectEndPercentY   = 0.9
)
const (
	textInputStartPercentX = textRectStartPercentX
	textInputEndPercentX   = textRectEndPercentX
	textInputStartPercentY = 0.90
	textInputEndPercentY   = textInputStartPercentY + .05
)

func (d *TalkDialog) initializeResizeableVisualElements() {
	d.border = widgets.NewBorder(
		sprites.PercentBasedPlacement{
			StartPercentX: borderStartPercentX,
			EndPercentX:   borderEndPercentX,
			StartPercentY: borderStartPercentY,
			EndPercentY:   borderEndPercentY,
		},
		borderWidthScaling,
		ucolor.Black)

	d.font = text.NewUltimaFont(text.GetScaledNumberToResolution(talkFontPoint))

	if d.Output == nil {
		d.Output = text.NewOutput(d.font,
			text.GetScaledNumberToResolution(talkFontLineSpacing),
			talkMaxLines,
			talkMaxCharsInput)
	} else {
		d.Output.SetFont(
			d.font,
			text.GetScaledNumberToResolution(talkFontLineSpacing),
		)
	}

	if d.TextInput == nil {
		d.TextInput = widgets.NewTextInput(
			sprites.PercentBasedPlacement{
				StartPercentX: textInputStartPercentX,
				EndPercentX:   textInputEndPercentX,
				StartPercentY: textInputStartPercentY,
				EndPercentY:   textInputEndPercentY,
			},
			text.GetScaledNumberToResolution(talkFontPoint),
			talkMaxCharsInput,
			d.createTalkFunctions(d.gameScene),
			widgets.TextInputCallbacks{
				AmbiguousAutoComplete: func(message string) {
					d.Output.AddRowStr(message, false)
				},
				OnEnter: d.onEnter,
			},
			d.gameScene.keyboard)
		d.TextInput.SetInputColors(widgets.TextInputColors{
			DefaultColor:          ucolor.Green,
			NoMatchesColor:        ucolor.Green,
			OneMatchColor:         ucolor.Green,
			MoreThanOneMatchColor: ucolor.Green,
		})
	} else {
		d.TextInput.SetFontPoint(text.GetScaledNumberToResolution(talkFontPoint))
	}
}

var emptyLastTime = false

func (d *TalkDialog) onEnter() {
	str := d.TextInput.GetText()
	select {
	case d.conversation.In() <- str: // succeeds only if ch has room
		emptyLastTime = false
		d.Output.AddRowStr(fmt.Sprintf("--- %s\n", str), false)
	default: // immediately chosen if ch is full
	}

	d.TextInput.ClearText()
}

func (d *TalkDialog) handleConversationOut() bool {
	select {
	case v := <-d.conversation.Out(): // succeeds only if ch has a value
		switch v.Cmd {
		case references.AvatarsName:
			break
		case references.EndConversation:
			d.conversation.Stop()

			d.gameScene.dialogStack.PopModalDialog()
		}
		d.Output.AppendToCurrentRowStr(v.Str)
	default: // immediately chosen if ch is empty
		if !emptyLastTime {
			emptyLastTime = true
			//d.Output.AddRowStr("\n", false)
		}
	}
	return true
}

func (d *TalkDialog) Update() {
	d.handleConversationOut()
	d.TextInput.Update()
}

func (d *TalkDialog) Draw(screen *ebiten.Image) {
	d.border.DrawBackground(screen)

	textRect := sprites.GetRectangleFromPercents(sprites.PercentBasedPlacement{
		StartPercentX: textRectStartPercentX,
		EndPercentX:   textRectEndPercentX,
		StartPercentY: textRectStartPercentY,
		EndPercentY:   textRectEndPercentY,
	})
	d.Output.DrawContinuousOutputTexOnXy(screen, image.Point{
		X: textRect.Min.X,
		Y: textRect.Min.Y,
	}, false, etext.AlignStart, etext.AlignStart)
	d.border.DrawBorder(screen)

	d.TextInput.Draw(screen)
}

func (d *TalkDialog) createTalkFunctions(gameScene *GameScene) *grammar.TextCommands {
	textCommands := make(grammar.TextCommands, 0)

	// Add each command by calling helper functions
	//textCommands = append(textCommands, *d.createTeleportCommand(gameScene))

	return &textCommands
}

func (d *TalkDialog) AddTestTest() {
	//d.Output.AddRowStrWithTrim("dsdsadsadsadsadsa\ndsdsadsadsadsadsa\ndsdsadsadsadsadsa\ndsdsadsadsadsadsa\ndsdsadsa\ndsadsadsadsadsadsadsadsa\nsdsadsadsadsa\ntesting 123\nand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some more...")
}
