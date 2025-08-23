package main

import (
	"fmt"
	"image"

	"github.com/bradhannah/Ultima5ReduxGo/internal/conversation"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/internal/text"
	"github.com/bradhannah/Ultima5ReduxGo/internal/ui/widgets"
	ucolor "github.com/bradhannah/Ultima5ReduxGo/pkg/color"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/hajimehoshi/ebiten/v2"
	etext "github.com/hajimehoshi/ebiten/v2/text/v2"
)

var _ widgets.Widget = &LinearTalkDialog{}

// LinearTalkDialog implements the same UI as TalkDialog but uses the LinearConversationEngine
type LinearTalkDialog struct {
	border *widgets.Border

	font   *text.UltimaFont
	Output *text.Output

	TextInput *widgets.TextInput

	gameScene *GameScene

	npcReference references.NPCReference
	engine       *conversation.LinearConversationEngine
	callbacks    *GameActionCallbacks

	// State management for the linear system
	waitingForInput    bool
	waitingForKeypress bool
	currentResponse    *conversation.ConversationResponse
}

func NewLinearTalkDialog(gameScene *GameScene, npcRef references.NPCReference) *LinearTalkDialog {
	dialog := &LinearTalkDialog{}
	dialog.gameScene = gameScene
	dialog.npcReference = npcRef

	// Get the TalkScript for this NPC
	talkScript := gameScene.gameState.GameReferences.TalkReferences.
		GetTalkScriptByNpcIndex(references.Castle, int(npcRef.DialogNumber)-1)

	// Create callbacks and engine
	dialog.callbacks = NewGameActionCallbacks(gameScene, npcRef)
	dialog.engine = conversation.NewLinearConversationEngine(talkScript, dialog.callbacks)

	// Initialize UI elements (borrowed from original TalkDialog)
	dialog.initializeResizeableVisualElements()

	// Start the conversation
	dialog.startConversation()

	return dialog
}

func (d *LinearTalkDialog) Refresh() {
	d.initializeResizeableVisualElements()
}

// initializeResizeableVisualElements - exact copy from original TalkDialog
func (d *LinearTalkDialog) initializeResizeableVisualElements() {
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

// startConversation initiates the conversation and displays initial output
func (d *LinearTalkDialog) startConversation() {
	// Convert the NPC reference dialog number to an NPC ID for the engine
	npcID := int(d.npcReference.DialogNumber)
	d.currentResponse = d.engine.Start(npcID)
	d.processCurrentResponse()
}

// processCurrentResponse handles the current response from the engine
func (d *LinearTalkDialog) processCurrentResponse() {
	if d.currentResponse == nil {
		return
	}

	// Display any output
	if d.currentResponse.Output != "" {
		d.Output.AppendToCurrentRowStr(d.currentResponse.Output)
	}

	// Handle errors
	if d.currentResponse.Error != nil {
		d.Output.AddRowStr(fmt.Sprintf("[Error: %v]\n", d.currentResponse.Error), false)
		return
	}

	// Check if conversation is complete
	if d.currentResponse.IsComplete {
		d.gameScene.dialogStack.PopModalDialog()
		return
	}

	// Update state flags
	d.waitingForInput = d.currentResponse.NeedsInput
	d.waitingForKeypress = false // TODO: detect keypress needs from response

	// Display input prompt if needed
	if d.currentResponse.NeedsInput && d.currentResponse.InputPrompt != "" {
		d.Output.AppendToCurrentRowStr(d.currentResponse.InputPrompt)
	}
}

// onEnter handles user input (modified from original to use linear engine)
func (d *LinearTalkDialog) onEnter() {
	str := d.TextInput.GetText()

	// Display the user input (same as original)
	d.Output.AddRowStr(fmt.Sprintf("--- %s\n", str), false)

	// Process input through linear engine
	if d.waitingForInput {
		d.currentResponse = d.engine.ProcessInput(str)
		d.processCurrentResponse()
	}

	d.TextInput.ClearText()
}

// Update method (simplified from original since no channel handling needed)
func (d *LinearTalkDialog) Update() {
	d.TextInput.Update()

	// Check if engine has any pending output to display
	if d.callbacks != nil {
		accumulatedOutput := d.callbacks.GetAccumulatedOutput()
		if accumulatedOutput != "" {
			d.Output.AppendToCurrentRowStr(accumulatedOutput)
			d.callbacks.ClearAccumulatedOutput()
		}
	}
}

// Draw method - exact copy from original TalkDialog
func (d *LinearTalkDialog) Draw(screen *ebiten.Image) {
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

// createTalkFunctions - exact copy from original TalkDialog
func (d *LinearTalkDialog) createTalkFunctions(gameScene *GameScene) *grammar.TextCommands {
	textCommands := make(grammar.TextCommands, 0)

	// Add each command by calling helper functions
	//textCommands = append(textCommands, *d.createTeleportCommand(gameScene))

	return &textCommands
}

// AddTestTest - exact copy from original TalkDialog
func (d *LinearTalkDialog) AddTestTest() {
	//d.Output.AddRowStrWithTrim("dsdsadsadsadsadsa\\ndsdsadsadsadsadsa\\ndsdsadsadsadsadsa\\ndsdsadsadsadsadsa\\ndsdsadsa\\ndsadsadsadsadsadsadsadsa\\nsdsadsadsadsa\\ntesting 123\\nand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some moreand then some more...")
}
