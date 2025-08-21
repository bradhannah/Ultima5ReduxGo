package conversation

import (
	"fmt"
	"log"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// LinearConversationEngine implements a straightforward pointer-based conversation system
// that processes TalkScript commands sequentially with linear navigation
type LinearConversationEngine struct {
	script        *references.TalkScript
	pointer       int // Current position in the script
	callbacks     ActionCallbacks
	currentOutput strings.Builder
	inputBuffer   string
	hasMet        bool
	isActive      bool
}

// ActionCallbacks defines the interface for handling conversation actions
type ActionCallbacks interface {
	// Core callbacks for game actions
	JoinParty() error
	CallGuards() error
	IncreaseKarma() error
	DecreaseKarma() error
	GoToJail() error
	MakeHorse() error
	PayExtortion(amount int) error
	PayHalfExtortion() error

	// Player interaction callbacks
	GetUserInput(prompt string) (string, error)
	AskPlayerName() (string, error)
	GetGoldAmount(prompt string) (int, error)
	ShowOutput(text string)
	WaitForKeypress()

	// Game state queries
	HasMet(npcID int) bool
	GetAvatarName() string
	GetKarmaLevel() int

	// Error handling
	OnError(err error)
}

// ConversationResponse represents the result of processing a conversation step
type ConversationResponse struct {
	Output      string
	NeedsInput  bool
	InputPrompt string
	IsComplete  bool
	Error       error
}

// NewLinearConversationEngine creates a new linear conversation engine
func NewLinearConversationEngine(script *references.TalkScript, callbacks ActionCallbacks) *LinearConversationEngine {
	return &LinearConversationEngine{
		script:    script,
		pointer:   0,
		callbacks: callbacks,
		isActive:  false,
	}
}

// Start begins the conversation with the NPC introduction
func (e *LinearConversationEngine) Start(npcID int) *ConversationResponse {
	e.isActive = true
	e.pointer = 0
	e.currentOutput.Reset()

	// Check if player has met this NPC before
	e.hasMet = e.callbacks.HasMet(npcID)

	// Bootstrap procedure: NPC introduction
	return e.performBootstrap()
}

// ProcessInput handles user input and continues the conversation
func (e *LinearConversationEngine) ProcessInput(input string) *ConversationResponse {
	if !e.isActive {
		return &ConversationResponse{
			Error: fmt.Errorf("conversation not active"),
		}
	}

	e.inputBuffer = strings.TrimSpace(strings.ToUpper(input))
	return e.processNextCommand()
}

// performBootstrap handles the initial NPC introduction sequence
func (e *LinearConversationEngine) performBootstrap() *ConversationResponse {
	e.currentOutput.Reset()

	// Show NPC description (fixed entry 1)
	if len(e.script.Lines) > references.TalkScriptConstantsDescription {
		descLine := e.script.Lines[references.TalkScriptConstantsDescription]
		e.processScriptLine(descLine)
	}

	// Determine greeting based on whether player has met NPC
	var greetingLine references.ScriptLine
	if e.hasMet && len(e.script.Lines) > references.TalkScriptConstantsGreeting {
		// Use greeting for known NPCs
		greetingLine = e.script.Lines[references.TalkScriptConstantsGreeting]
	} else if len(e.script.Lines) > references.TalkScriptConstantsName {
		// Introduce themselves if first meeting
		e.currentOutput.WriteString("\"I am called ")
		nameLine := e.script.Lines[references.TalkScriptConstantsName]
		e.processScriptLine(nameLine)
		e.currentOutput.WriteString("\"\n\n")
		return e.promptForInput("Your interest?")
	}

	if len(greetingLine) > 0 {
		e.currentOutput.WriteString("\"")
		e.processScriptLine(greetingLine)
		e.currentOutput.WriteString("\"\n\n")
	}

	return e.promptForInput("Your interest?")
}

// processNextCommand continues processing the conversation
func (e *LinearConversationEngine) processNextCommand() *ConversationResponse {
	// Handle empty input as BYE
	if e.inputBuffer == "" {
		return e.handleBye()
	}

	// Try to match input against keywords and handle accordingly
	response := e.handleKeywordMatch()
	if response != nil {
		return response
	}

	// Try to find keyword in script data
	response = e.searchScriptKeywords()
	if response != nil {
		return response
	}

	// Default response for unrecognized input
	return e.handleUnrecognizedInput()
}

// handleKeywordMatch processes standard conversation keywords
func (e *LinearConversationEngine) handleKeywordMatch() *ConversationResponse {
	switch {
	case strings.Contains(e.inputBuffer, "NAME"):
		return e.handleName()
	case strings.Contains(e.inputBuffer, "JOB") || strings.Contains(e.inputBuffer, "WORK"):
		return e.handleJob()
	case strings.Contains(e.inputBuffer, "BYE") || strings.Contains(e.inputBuffer, "THANK"):
		return e.handleBye()
	}
	return nil
}

// handleName responds with NPC's name
func (e *LinearConversationEngine) handleName() *ConversationResponse {
	e.currentOutput.Reset()
	e.currentOutput.WriteString("\"My name is ")

	if len(e.script.Lines) > references.TalkScriptConstantsName {
		nameLine := e.script.Lines[references.TalkScriptConstantsName]
		e.processScriptLine(nameLine)
	}

	e.currentOutput.WriteString("\"\n\n")
	return e.promptForInput("Your interest?")
}

// handleJob responds with NPC's job
func (e *LinearConversationEngine) handleJob() *ConversationResponse {
	e.currentOutput.Reset()
	e.currentOutput.WriteString("\"")

	if len(e.script.Lines) > references.TalkScriptConstantsJob {
		jobLine := e.script.Lines[references.TalkScriptConstantsJob]
		e.processScriptLine(jobLine)
	}

	e.currentOutput.WriteString("\"\n\n")
	return e.promptForInput("Your interest?")
}

// handleBye ends the conversation
func (e *LinearConversationEngine) handleBye() *ConversationResponse {
	e.currentOutput.Reset()
	e.currentOutput.WriteString("\"")

	if len(e.script.Lines) > references.TalkScriptConstantsBye {
		byeLine := e.script.Lines[references.TalkScriptConstantsBye]
		e.processScriptLine(byeLine)
	}

	e.currentOutput.WriteString("\"\n")
	e.isActive = false

	return &ConversationResponse{
		Output:     e.currentOutput.String(),
		IsComplete: true,
	}
}

// searchScriptKeywords searches for keywords in the script's question groups
func (e *LinearConversationEngine) searchScriptKeywords() *ConversationResponse {
	for _, group := range e.script.QuestionGroups {
		for _, option := range group.Options {
			if strings.Contains(e.inputBuffer, strings.ToUpper(option)) {
				e.currentOutput.Reset()
				e.currentOutput.WriteString("\"")
				e.processScriptLine(group.Script)
				e.currentOutput.WriteString("\"\n\n")
				return e.promptForInput("Your interest?")
			}
		}
	}
	return nil
}

// handleUnrecognizedInput provides default response for unmatched input
func (e *LinearConversationEngine) handleUnrecognizedInput() *ConversationResponse {
	e.currentOutput.Reset()
	e.currentOutput.WriteString("\"I cannot help thee with that.\"\n\n")
	return e.promptForInput("Your interest?")
}

// processScriptLine processes a single script line command by command
func (e *LinearConversationEngine) processScriptLine(line references.ScriptLine) error {
	for _, item := range line {
		if err := e.processScriptItem(item); err != nil {
			return err
		}
	}
	return nil
}

// processScriptItem processes a single script item
func (e *LinearConversationEngine) processScriptItem(item references.ScriptItem) error {
	switch item.Cmd {
	case references.PlainString:
		e.currentOutput.WriteString(item.Str)

	case references.AvatarsName:
		e.currentOutput.WriteString(e.callbacks.GetAvatarName())

	case references.NewLine:
		e.currentOutput.WriteString("\n")

	case references.JoinParty:
		return e.callbacks.JoinParty()

	case references.CallGuards:
		return e.callbacks.CallGuards()

	case references.KarmaPlusOne:
		return e.callbacks.IncreaseKarma()

	case references.KarmaMinusOne:
		return e.callbacks.DecreaseKarma()

	case references.Pause:
		e.callbacks.WaitForKeypress()

	case references.KeyWait:
		e.callbacks.WaitForKeypress()

	case references.GoToJail:
		return e.callbacks.GoToJail()

	case references.MakeAHorse:
		return e.callbacks.MakeHorse()

	case references.EndConversation:
		e.isActive = false

	default:
		// Log unknown commands but don't fail
		log.Printf("Unknown talk command: %s (0x%02X)", item.Cmd, byte(item.Cmd))
	}

	return nil
}

// promptForInput creates a response that requests user input
func (e *LinearConversationEngine) promptForInput(prompt string) *ConversationResponse {
	return &ConversationResponse{
		Output:      e.currentOutput.String(),
		NeedsInput:  true,
		InputPrompt: prompt,
	}
}

// GetCurrentOutput returns the current conversation output
func (e *LinearConversationEngine) GetCurrentOutput() string {
	return e.currentOutput.String()
}

// IsActive returns whether the conversation is currently active
func (e *LinearConversationEngine) IsActive() bool {
	return e.isActive
}

// Stop forcefully ends the conversation
func (e *LinearConversationEngine) Stop() {
	e.isActive = false
}
