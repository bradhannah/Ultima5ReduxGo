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
	script           *references.TalkScript
	pointer          int // Current position in the script
	callbacks        ActionCallbacks
	currentOutput    strings.Builder
	inputBuffer      string
	hasMet           bool
	isActive         bool
	labelMap         map[references.TalkCommand]int // Maps label commands to script positions
	currentLabel     int                            // Current label for question mode (-1 = not in question mode)
	waitingForName   bool                           // True when waiting for name input from AskName command
	waitingForPause  bool                           // True when waiting for keypress from Pause command
	pausedScriptLine references.ScriptLine          // Script line being processed when pause occurred
	pausedItemIndex  int                            // Index in script line where pause occurred
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
		script:          script,
		pointer:         0,
		callbacks:       callbacks,
		isActive:        false,
		labelMap:        make(map[references.TalkCommand]int),
		currentLabel:    -1,    // Not in question mode
		waitingForName:  false, // Not waiting for name input
		waitingForPause: false, // Not waiting for pause input
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

	// If we're waiting for name input from AskName command, handle it
	if e.waitingForName {
		return e.processNameInput()
	}

	// If we're waiting for pause keypress, handle it
	if e.waitingForPause {
		return e.processPauseInput()
	}

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
		e.currentOutput.WriteString("You see ")
		descLine := e.script.Lines[references.TalkScriptConstantsDescription]
		if err := e.processScriptLine(descLine); err != nil {
			return &ConversationResponse{Error: err}
		}
	}

	// Determine greeting based on whether player has met NPC
	if e.hasMet && len(e.script.Lines) > references.TalkScriptConstantsGreeting {
		// Use greeting for known NPCs - this may contain commands like <Goto Label 0>
		greetingLine := e.script.Lines[references.TalkScriptConstantsGreeting]
		e.currentOutput.WriteString("\"")
		if err := e.processScriptLine(greetingLine); err != nil {
			return &ConversationResponse{Error: err}
		}
		e.currentOutput.WriteString("\"\n\n")
	}
	// For first meeting (HasMet=false), just show description and wait for input
	// Don't automatically process the name line - let player ask for "name" explicitly

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

		// Check if AskName command was encountered and we're now waiting for name input
		if e.waitingForName {
			e.currentOutput.WriteString("\"\n\nWhat is thy name?")
			return &ConversationResponse{
				Output:      e.currentOutput.String(),
				NeedsInput:  true,
				InputPrompt: "You respond:",
			}
		}
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

				// Check if we encountered a pause during processing
				if e.waitingForPause {
					// Show output so far and wait for keypress
					return &ConversationResponse{
						Output:      e.currentOutput.String() + "[PAUSED, press enter]",
						NeedsInput:  true,
						InputPrompt: "[Press Enter to continue...]",
					}
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

		// Check if we encountered a pause and need to stop processing
		if e.waitingForPause {
			// Pause state is already set by the inner processing (e.g., in processQuestion)
			// Don't override it here unless it hasn't been set
			if e.pausedScriptLine == nil {
				e.pausedScriptLine = line
				e.pausedItemIndex = i + 1 // Resume from next item
			}
			return nil
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
		e.waitingForPause = true

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

	case references.AskName:
		return e.processAskName()

	case references.StartNewSection:
		// StartNewSection (0xA2) - formatting/organizational marker, no action needed

	case references.DoNothingSection:
		// DoNothingSection (0xFF) - explicitly does nothing

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

				// Check if we encountered a pause during processing
				if e.waitingForPause {
					// Make sure the pause state is set correctly for the label content
					// The pause occurred in the label content, not the calling script
					e.pausedScriptLine = labelData.Initial
					// Find where the pause occurred by looking for the Pause command
					for i, item := range labelData.Initial {
						if item.Cmd == references.Pause {
							e.pausedItemIndex = i + 1 // Resume after the pause
							break
						}
					}
					log.Printf("DEBUG: Pause in label %d, resuming from item %d of %d", labelNum, e.pausedItemIndex, len(e.pausedScriptLine))
					return nil
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

		// First try exact match
		if qa, exists := labelData.QA[inputKey]; exists {
			log.Printf("DEBUG: Found exact QA mapping for '%s'", inputKey)

			// Store original label before processing answer
			originalLabel := e.currentLabel

			e.currentOutput.Reset()
			e.currentOutput.WriteString("\"")

			// Process the answer - it may contain navigation commands
			if err := e.processScriptLine(qa.Answer); err != nil {
				return &ConversationResponse{Error: err}
			}

			// Only exit question mode if we haven't navigated to another question
			if e.currentLabel == originalLabel {
				e.currentLabel = -1 // Exit question mode after answering
			}

			e.currentOutput.WriteString("\"\n\n")
			return e.promptForInput("Your interest?")
		}

		// Try intelligent matching for common yes/no variations
		if inputKey == "yes" || inputKey == "yeah" || inputKey == "yep" || inputKey == "yea" {
			if qa, exists := labelData.QA["y"]; exists {
				log.Printf("DEBUG: Found QA mapping for '%s' via intelligent matching (matched 'y')", inputKey)

				// Store original label before processing answer
				originalLabel := e.currentLabel

				e.currentOutput.Reset()
				e.currentOutput.WriteString("\"")

				// Process the answer - it may contain navigation commands
				if err := e.processScriptLine(qa.Answer); err != nil {
					return &ConversationResponse{Error: err}
				}

				// Only exit question mode if we haven't navigated to another question
				if e.currentLabel == originalLabel {
					e.currentLabel = -1 // Exit question mode after answering
				}

				e.currentOutput.WriteString("\"\n\n")
				return e.promptForInput("Your interest?")
			}
		}

		if inputKey == "no" || inputKey == "nope" || inputKey == "nay" {
			if qa, exists := labelData.QA["n"]; exists {
				log.Printf("DEBUG: Found QA mapping for '%s' via intelligent matching (matched 'n')", inputKey)

				// Store original label before processing answer
				originalLabel := e.currentLabel

				e.currentOutput.Reset()
				e.currentOutput.WriteString("\"")

				// Process the answer - it may contain navigation commands
				if err := e.processScriptLine(qa.Answer); err != nil {
					return &ConversationResponse{Error: err}
				}

				// Only exit question mode if we haven't navigated to another question
				if e.currentLabel == originalLabel {
					e.currentLabel = -1 // Exit question mode after answering
				}

				e.currentOutput.WriteString("\"\n\n")
				return e.promptForInput("Your interest?")
			}
		}
	}

	// No match found, check for default answers
	log.Printf("DEBUG: No QA mapping found, checking for default answers")

	if len(labelData.DefaultAnswers) > 0 {
		// Process the first default answer
		defaultAnswer := labelData.DefaultAnswers[0]
		log.Printf("DEBUG: Processing default answer with %d items", len(defaultAnswer))

		// Store original label before processing default answer
		originalLabel := e.currentLabel

		e.currentOutput.Reset()
		e.currentOutput.WriteString("\"")

		// Process the default answer - it may contain navigation commands
		if err := e.processScriptLine(defaultAnswer); err != nil {
			return &ConversationResponse{Error: err}
		}

		// Only exit question mode if we haven't navigated to another question
		if e.currentLabel == originalLabel {
			e.currentLabel = -1 // Exit question mode after answering
		}

		e.currentOutput.WriteString("\"\n\n")
		return e.promptForInput("Your interest?")
	}

	// No default answer either, exit question mode and provide fallback response
	log.Printf("DEBUG: No default answer found, exiting question mode")
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

// processAskName implements the AskName (0x88) command
// Based on the original askname() function in TALKNPC.C lines 701-729
// This returns a ConversationResponse that requests name input
func (e *LinearConversationEngine) processAskName() error {
	// Mark that we're waiting for name input
	e.waitingForName = true
	return nil // The actual prompting will be handled by the response system
}

// processPauseInput handles keypress after Pause command
func (e *LinearConversationEngine) processPauseInput() *ConversationResponse {
	// No longer waiting for pause input
	e.waitingForPause = false

	log.Printf("DEBUG: Resuming from pause. PausedScriptLine length: %d, PausedItemIndex: %d",
		len(e.pausedScriptLine), e.pausedItemIndex)

	// Continue processing from where we left off
	if e.pausedScriptLine != nil && e.pausedItemIndex < len(e.pausedScriptLine) {
		// Resume processing the remaining items in the paused script line
		remainingItems := e.pausedScriptLine[e.pausedItemIndex:]
		log.Printf("DEBUG: Processing %d remaining items after pause", len(remainingItems))

		if err := e.processScriptLine(remainingItems); err != nil {
			return &ConversationResponse{Error: err}
		}

		// Clear the paused state
		e.pausedScriptLine = nil
		e.pausedItemIndex = 0

		// After resuming from pause and finishing the script, we should exit question mode
		// The label content has been fully processed
		if e.currentLabel >= 0 {
			e.currentLabel = -1 // Exit question mode
		}

		// Close the quote and return to normal conversation flow
		e.currentOutput.WriteString("\"\n\n")
		return e.promptForInput("Your interest?")
	}

	log.Printf("DEBUG: No paused script to resume")
	// Return to normal conversation flow
	return e.promptForInput("Your interest?")
}

// processNameInput handles the response to AskName command
func (e *LinearConversationEngine) processNameInput() *ConversationResponse {
	// No longer waiting for name input
	e.waitingForName = false

	e.currentOutput.Reset()
	e.currentOutput.WriteString("\"")

	// Clean up the input (trim spaces, convert to uppercase for comparison)
	nameInput := strings.TrimSpace(strings.ToUpper(e.inputBuffer))

	// If empty input, respond with "If you say so..."
	if nameInput == "" {
		e.currentOutput.WriteString("If you say so...\"")
		return e.promptForInput("Your interest?")
	}

	// Check if the name matches any party member's name
	// We don't have direct access to party data, so we'll use the callback system
	// and check if the input matches the avatar name (simplified for now)
	avatarName := strings.ToUpper(e.callbacks.GetAvatarName())

	// Check if input contains the avatar's name (allowing partial matches)
	if strings.Contains(nameInput, avatarName) || strings.Contains(avatarName, nameInput) {
		// Name recognized - this should mark the NPC as "met"
		// In the original game, this would call setmet(talknum)
		// For now, we'll just respond positively
		e.currentOutput.WriteString("A pleasure!\"")
		return e.promptForInput("Your interest?")
	}

	// Name not recognized
	e.currentOutput.WriteString("If you say so...\"")
	return e.promptForInput("Your interest?")
}
