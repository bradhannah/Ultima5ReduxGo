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
	labelMap      map[references.TalkCommand]int // Maps label commands to script positions
	currentLabel  int                            // Current label for question mode (-1 = not in question mode)
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
	engine := &LinearConversationEngine{
		script:       script,
		pointer:      0,
		callbacks:    callbacks,
		isActive:     false,
		labelMap:     make(map[references.TalkCommand]int),
		currentLabel: -1, // Not in question mode
	}

	// Build label map for fast navigation
	engine.buildLabelMap()

	return engine
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

	// If we're in question mode, handle differently
	if e.currentLabel >= 0 {
		return e.processQuestionAnswer()
	}

	return e.processNextCommand()
}

// performBootstrap handles the initial NPC introduction sequence
func (e *LinearConversationEngine) performBootstrap() *ConversationResponse {
	e.currentOutput.Reset()

	// Show NPC description (fixed entry 1)
	if len(e.script.Lines) > references.TalkScriptConstantsDescription {
		descLine := e.script.Lines[references.TalkScriptConstantsDescription]
		if err := e.processScriptLine(descLine); err != nil {
			return &ConversationResponse{Error: err}
		}
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
		if err := e.processScriptLine(nameLine); err != nil {
			return &ConversationResponse{Error: err}
		}
		e.currentOutput.WriteString("\"\n\n")
		return e.promptForInput("Your interest?")
	}

	if len(greetingLine) > 0 {
		e.currentOutput.WriteString("\"")
		if err := e.processScriptLine(greetingLine); err != nil {
			return &ConversationResponse{Error: err}
		}
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
				if err := e.processScriptLine(group.Script); err != nil {
					return &ConversationResponse{Error: err}
				}
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
	for i := 0; i < len(line); i++ {
		item := line[i]

		if item.Cmd == references.IfElseKnowsName {
			// Handle IfElseKnowsName inline with context
			var targetIndex int
			if e.hasMet {
				// Use the next item (+1) - they DO know the Avatar
				targetIndex = i + 1
			} else {
				// Use the item after that (+2) - they do NOT know the Avatar
				targetIndex = i + 2
			}

			if targetIndex < len(line) {
				if err := e.processScriptItem(line[targetIndex]); err != nil {
					return err
				}
			}

			// Skip the two conditional items and continue from after them
			i += 2
			continue
		}

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

	case references.GotoLabel:
		// Look for the label in the Num field
		if item.Num >= 1 && item.Num <= 10 {
			labelCmd := references.TalkCommand(int(references.Label1) + item.Num - 1)
			return e.gotoLabel(labelCmd)
		}

	case references.DefineLabel:
		// Label definitions don't execute, just mark positions

	case references.Label1, references.Label2, references.Label3, references.Label4, references.Label5,
		references.Label6, references.Label7, references.Label8, references.Label9, references.Label10:
		// When we encounter a label in a response, navigate to that label's content
		if err := e.processQuestion(item.Cmd); err != nil {
			return err
		}

	case references.StartLabelDef:
		// Label section markers don't execute

	case references.IfElseKnowsName:
		// This should be handled in processScriptLine, not here
		// If we reach here, something is wrong
		return fmt.Errorf("IfElseKnowsName should be handled at line level")

	default:
		// Handle question labels (Label1 through EndScript range for questions)
		if item.Cmd >= references.Label1 && item.Cmd <= references.EndScript {
			// This is a question - transition to question mode
			if err := e.processQuestion(item.Cmd); err != nil {
				return err
			}
			return nil
		}

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

// buildLabelMap scans the script and builds a map of labels to their positions
func (e *LinearConversationEngine) buildLabelMap() {
	for lineIndex, line := range e.script.Lines {
		for itemIndex, item := range line {
			// Look for DefineLabel followed by a label
			if item.Cmd == references.DefineLabel && itemIndex+1 < len(line) {
				nextItem := line[itemIndex+1]
				if nextItem.Cmd >= references.Label1 && nextItem.Cmd <= references.Label10 {
					e.labelMap[nextItem.Cmd] = lineIndex
				}
			}
		}
	}

	// Also check labels defined in script.Labels if available
	if e.script.Labels != nil {
		for labelNum := range e.script.Labels {
			if labelNum >= 1 && labelNum <= int(references.Label10-references.Label1+1) {
				labelCmd := references.TalkCommand(int(references.Label1) + labelNum - 1)
				// Find the label in the script lines
				for lineIndex, line := range e.script.Lines {
					for _, item := range line {
						if item.Cmd == labelCmd {
							e.labelMap[labelCmd] = lineIndex
							break
						}
					}
				}
			}
		}
	}
}

// gotoLabel jumps to the specified label in the script
func (e *LinearConversationEngine) gotoLabel(label references.TalkCommand) error {
	if position, exists := e.labelMap[label]; exists {
		e.pointer = position
		return nil
	}
	return fmt.Errorf("label %s not found", label.String())
}

// processQuestion handles question processing and waits for user response
func (e *LinearConversationEngine) processQuestion(questionCmd references.TalkCommand) error {
	// Find the question text from the script.Labels if available
	if e.script.Labels != nil {
		// Calculate label number: Label1=0x91 -> 0, Label2=0x92 -> 1, etc.
		// But the TLK data seems to be off by one, so Label5 (0x95) should map to Label 4
		labelNum := int(questionCmd - references.Label1)

		if labelData, exists := e.script.Labels[labelNum]; exists {
			// Enter question mode for this label
			e.currentLabel = labelNum
			log.Printf("DEBUG: Entering question mode for label %d", labelNum)

			// Skip the label definition header (StartLabelDef and the label itself)
			// Start processing from the actual content
			contentStart := 0
			for i, item := range labelData.Initial {
				if item.Cmd == references.StartLabelDef ||
					(item.Cmd >= references.Label1 && item.Cmd <= references.Label10) {
					contentStart = i + 1
				} else {
					break
				}
			}

			if contentStart < len(labelData.Initial) {
				contentItems := labelData.Initial[contentStart:]
				if err := e.processScriptLine(contentItems); err != nil {
					return err
				}
			}
		}
	}

	// For basic implementation, just output a generic question prompt
	if e.currentOutput.String() == "\"" {
		e.currentOutput.WriteString("I have a question for thee.")
	}

	return nil
}

// processQuestionAnswer handles input when in question mode for a specific label
func (e *LinearConversationEngine) processQuestionAnswer() *ConversationResponse {
	if e.script.Labels == nil || e.currentLabel < 0 {
		// Exit question mode and fall back to normal processing
		e.currentLabel = -1
		return e.processNextCommand()
	}

	labelData, exists := e.script.Labels[e.currentLabel]
	if !exists {
		// Exit question mode and fall back to normal processing
		e.currentLabel = -1
		return e.processNextCommand()
	}

	log.Printf("DEBUG: Processing question answer for label %d, input: '%s'", e.currentLabel, e.inputBuffer)

	// Check if input matches any QA mappings for this label
	if labelData.QA != nil {
		inputKey := strings.ToLower(e.inputBuffer)
		if qa, exists := labelData.QA[inputKey]; exists {
			log.Printf("DEBUG: Found QA mapping for '%s'", inputKey)
			e.currentOutput.Reset()
			e.currentOutput.WriteString("\"")

			// Store original label before processing answer
			originalLabel := e.currentLabel

			// Process the answer - it may contain navigation commands
			if err := e.processScriptLine(qa.Answer); err != nil {
				return &ConversationResponse{Error: err}
			}

			// Only exit question mode if we haven't navigated to another question
			// If processScriptLine triggered a label navigation, we'll stay in question mode
			if e.currentLabel == originalLabel {
				e.currentLabel = -1 // Exit question mode after answering
			}

			e.currentOutput.WriteString("\"\n\n")
			return e.promptForInput("Your interest?")
		}
	}

	// No match found, exit question mode and provide default response
	log.Printf("DEBUG: No QA mapping found, exiting question mode")
	e.currentLabel = -1
	return e.handleUnrecognizedInput()
}

// processIfElseKnowsName handles conditional branching based on whether NPC knows Avatar's name
func (e *LinearConversationEngine) processIfElseKnowsName() error {
	// This method is called when processing a script line that contains IfElseKnowsName
	// We need to find the current context (what line/item we're processing)

	// Since we're processing within processScriptLine, we need to look at the current processing context
	// For now, let's implement a simpler approach that works during bootstrap

	// Find the line that contains IfElseKnowsName
	var currentLine references.ScriptLine
	var currentItemIndex int = -1

	// First try the current pointer position
	if e.pointer < len(e.script.Lines) {
		for i, item := range e.script.Lines[e.pointer] {
			if item.Cmd == references.IfElseKnowsName {
				currentLine = e.script.Lines[e.pointer]
				currentItemIndex = i
				break
			}
		}
	}

	// If not found at current pointer, search all lines (for bootstrap phase)
	if currentItemIndex == -1 {
		for lineIdx, line := range e.script.Lines {
			for itemIdx, item := range line {
				if item.Cmd == references.IfElseKnowsName {
					currentLine = line
					currentItemIndex = itemIdx
					e.pointer = lineIdx // Update pointer for consistency
					break
				}
			}
			if currentItemIndex != -1 {
				break
			}
		}
	}

	if currentItemIndex == -1 {
		return fmt.Errorf("IfElseKnowsName command not found")
	}

	// According to the documentation:
	// The next script item (+1) will be what happens if they DO know the Avatar (HasMet),
	// the one after that (+2) will be what happens if they do NOT know the Avatar.

	var targetItemIndex int
	if e.hasMet {
		// Use the next item (+1) - they DO know the Avatar
		targetItemIndex = currentItemIndex + 1
	} else {
		// Use the item after that (+2) - they do NOT know the Avatar
		targetItemIndex = currentItemIndex + 2
	}

	if targetItemIndex < len(currentLine) {
		// Process the target item
		return e.processScriptItem(currentLine[targetItemIndex])
	}

	return nil
}
