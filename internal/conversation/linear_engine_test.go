package conversation

import (
	"strings"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// TestActionCallbacks provides a mock implementation for testing
type TestActionCallbacks struct {
	joinPartyCalled        bool
	callGuardsCalled       bool
	increaseKarmaCalled    bool
	decreaseKarmaCalled    bool
	goToJailCalled         bool
	makeHorseCalled        bool
	payExtortionCalled     bool
	payHalfExtortionCalled bool

	userInputResponses []string
	userInputIndex     int
	avatarName         string
	karmaLevel         int
	hasMetNPC          bool

	outputShown   []string
	keypressWaits int
	errors        []error
}

func (t *TestActionCallbacks) JoinParty() error {
	t.joinPartyCalled = true
	return nil
}

func (t *TestActionCallbacks) CallGuards() error {
	t.callGuardsCalled = true
	return nil
}

func (t *TestActionCallbacks) IncreaseKarma() error {
	t.increaseKarmaCalled = true
	return nil
}

func (t *TestActionCallbacks) DecreaseKarma() error {
	t.decreaseKarmaCalled = true
	return nil
}

func (t *TestActionCallbacks) GoToJail() error {
	t.goToJailCalled = true
	return nil
}

func (t *TestActionCallbacks) MakeHorse() error {
	t.makeHorseCalled = true
	return nil
}

func (t *TestActionCallbacks) PayExtortion(amount int) error {
	t.payExtortionCalled = true
	return nil
}

func (t *TestActionCallbacks) PayHalfExtortion() error {
	t.payHalfExtortionCalled = true
	return nil
}

func (t *TestActionCallbacks) GetUserInput(prompt string) (string, error) {
	if t.userInputIndex < len(t.userInputResponses) {
		response := t.userInputResponses[t.userInputIndex]
		t.userInputIndex++
		return response, nil
	}
	return "", nil
}

func (t *TestActionCallbacks) AskPlayerName() (string, error) {
	return t.GetUserInput("What is thy name?")
}

func (t *TestActionCallbacks) GetGoldAmount(prompt string) (int, error) {
	response, err := t.GetUserInput(prompt)
	if err != nil {
		return 0, err
	}
	// Simple conversion for testing
	if response == "100" {
		return 100, nil
	}
	return 0, nil
}

func (t *TestActionCallbacks) ShowOutput(text string) {
	t.outputShown = append(t.outputShown, text)
}

func (t *TestActionCallbacks) WaitForKeypress() {
	t.keypressWaits++
}

func (t *TestActionCallbacks) HasMet(npcID int) bool {
	return t.hasMetNPC
}

func (t *TestActionCallbacks) GetAvatarName() string {
	return t.avatarName
}

func (t *TestActionCallbacks) GetKarmaLevel() int {
	return t.karmaLevel
}

func (t *TestActionCallbacks) OnError(err error) {
	t.errors = append(t.errors, err)
}

// createTestScript creates a basic TalkScript for testing
func createTestScript() *references.TalkScript {
	return &references.TalkScript{
		Lines: []references.ScriptLine{
			// Name (index 0)
			{
				{Cmd: references.PlainString, Str: "Treanna"},
			},
			// Description (index 1)
			{
				{Cmd: references.PlainString, Str: "a mysterious woman in robes"},
			},
			// Greeting (index 2)
			{
				{Cmd: references.PlainString, Str: "Welcome back, "},
				{Cmd: references.AvatarsName},
				{Cmd: references.PlainString, Str: "!"},
			},
			// Job (index 3)
			{
				{Cmd: references.PlainString, Str: "I am a keeper of ancient knowledge."},
			},
			// Bye (index 4)
			{
				{Cmd: references.PlainString, Str: "Farewell, brave "},
				{Cmd: references.AvatarsName},
				{Cmd: references.PlainString, Str: "."},
			},
		},
		QuestionGroups: []references.QuestionGroup{
			{
				Options: []string{"VIRTUE", "HONOR", "JUSTICE"},
				Script: references.ScriptLine{
					{Cmd: references.PlainString, Str: "The virtues guide our path through darkness."},
				},
			},
			{
				Options: []string{"MAGIC", "SPELL", "RUNE"},
				Script: references.ScriptLine{
					{Cmd: references.PlainString, Str: "Magic flows through all things, if one knows how to see."},
				},
			},
		},
	}
}

func TestLinearConversationEngine_NewEngine(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  false,
	}

	engine := NewLinearConversationEngine(script, callbacks)

	if engine == nil {
		t.Fatal("Expected engine to be created, got nil")
	}

	if engine.script != script {
		t.Error("Expected script to be set correctly")
	}

	if engine.callbacks == nil {
		t.Error("Expected callbacks to be set correctly")
	}

	if engine.isActive {
		t.Error("Expected engine to not be active initially")
	}
}

func TestLinearConversationEngine_StartConversation_FirstMeeting(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  false, // First meeting
	}

	engine := NewLinearConversationEngine(script, callbacks)
	response := engine.Start(1)

	if !engine.IsActive() {
		t.Error("Expected engine to be active after starting")
	}

	if response.IsComplete {
		t.Error("Expected conversation to not be complete yet")
	}

	if !response.NeedsInput {
		t.Error("Expected conversation to need input")
	}

	if response.InputPrompt != "Your interest?" {
		t.Errorf("Expected prompt 'Your interest?', got '%s'", response.InputPrompt)
	}

	// Should contain description and introduction
	output := response.Output
	if !strings.Contains(output, "a mysterious woman in robes") {
		t.Error("Expected output to contain NPC description")
	}

	if !strings.Contains(output, "I am called Treanna") {
		t.Error("Expected output to contain NPC introduction for first meeting")
	}
}

func TestLinearConversationEngine_StartConversation_HasMet(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true, // Already met
	}

	engine := NewLinearConversationEngine(script, callbacks)
	response := engine.Start(1)

	if !response.NeedsInput {
		t.Error("Expected conversation to need input")
	}

	// Should contain description and greeting
	output := response.Output
	if !strings.Contains(output, "a mysterious woman in robes") {
		t.Error("Expected output to contain NPC description")
	}

	if !strings.Contains(output, "Welcome back, TestHero!") {
		t.Error("Expected output to contain greeting with avatar name")
	}
}

func TestLinearConversationEngine_HandleNameKeyword(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	response := engine.ProcessInput("NAME")

	if response.IsComplete {
		t.Error("Expected conversation to continue after name inquiry")
	}

	if !strings.Contains(response.Output, "My name is Treanna") {
		t.Error("Expected response to contain NPC name")
	}

	if response.InputPrompt != "Your interest?" {
		t.Error("Expected to prompt for more input")
	}
}

func TestLinearConversationEngine_HandleJobKeyword(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	response := engine.ProcessInput("JOB")

	if response.IsComplete {
		t.Error("Expected conversation to continue after job inquiry")
	}

	if !strings.Contains(response.Output, "I am a keeper of ancient knowledge") {
		t.Error("Expected response to contain NPC job description")
	}
}

func TestLinearConversationEngine_HandleWorkKeyword(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	response := engine.ProcessInput("WORK")

	if response.IsComplete {
		t.Error("Expected conversation to continue after work inquiry")
	}

	if !strings.Contains(response.Output, "I am a keeper of ancient knowledge") {
		t.Error("Expected WORK to trigger job response")
	}
}

func TestLinearConversationEngine_HandleByeKeyword(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	response := engine.ProcessInput("BYE")

	if !response.IsComplete {
		t.Error("Expected conversation to be complete after bye")
	}

	if !strings.Contains(response.Output, "Farewell, brave TestHero") {
		t.Error("Expected farewell message with avatar name")
	}

	if engine.IsActive() {
		t.Error("Expected engine to be inactive after bye")
	}
}

func TestLinearConversationEngine_HandleEmptyInput(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	response := engine.ProcessInput("")

	if !response.IsComplete {
		t.Error("Expected conversation to be complete after empty input (bye)")
	}

	if !strings.Contains(response.Output, "Farewell, brave TestHero") {
		t.Error("Expected empty input to trigger bye response")
	}
}

func TestLinearConversationEngine_HandleQuestionGroupKeywords(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	// Test virtue keyword
	response := engine.ProcessInput("VIRTUE")

	if response.IsComplete {
		t.Error("Expected conversation to continue after virtue inquiry")
	}

	if !strings.Contains(response.Output, "The virtues guide our path through darkness") {
		t.Error("Expected virtue response")
	}

	// Test magic keyword
	response = engine.ProcessInput("MAGIC")

	if response.IsComplete {
		t.Error("Expected conversation to continue after magic inquiry")
	}

	if !strings.Contains(response.Output, "Magic flows through all things") {
		t.Error("Expected magic response")
	}
}

func TestLinearConversationEngine_HandleUnrecognizedInput(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	response := engine.ProcessInput("NONSENSE")

	if response.IsComplete {
		t.Error("Expected conversation to continue after unrecognized input")
	}

	if !strings.Contains(response.Output, "I cannot help thee with that") {
		t.Error("Expected default response for unrecognized input")
	}
}

func TestLinearConversationEngine_ProcessInputWhenNotActive(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	// Don't start the conversation

	response := engine.ProcessInput("NAME")

	if response.Error == nil {
		t.Error("Expected error when processing input on inactive engine")
	}

	if !strings.Contains(response.Error.Error(), "conversation not active") {
		t.Error("Expected 'conversation not active' error message")
	}
}

func TestLinearConversationEngine_ActionCallbacks(t *testing.T) {
	// Create script with action commands
	script := &references.TalkScript{
		Lines: []references.ScriptLine{
			// Name (index 0)
			{
				{Cmd: references.PlainString, Str: "Guard"},
			},
			// Description (index 1)
			{
				{Cmd: references.PlainString, Str: "a stern guard"},
			},
			// Greeting (index 2)
			{
				{Cmd: references.PlainString, Str: "State your business!"},
			},
			// Job (index 3)
			{
				{Cmd: references.PlainString, Str: "I protect this place."},
				{Cmd: references.KarmaPlusOne},
			},
			// Bye (index 4)
			{
				{Cmd: references.PlainString, Str: "Move along."},
			},
		},
		QuestionGroups: []references.QuestionGroup{
			{
				Options: []string{"HELP"},
				Script: references.ScriptLine{
					{Cmd: references.PlainString, Str: "I shall join thee!"},
					{Cmd: references.JoinParty},
				},
			},
		},
	}

	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	// Test karma increase callback
	response := engine.ProcessInput("JOB")
	if !callbacks.increaseKarmaCalled {
		t.Error("Expected IncreaseKarma callback to be called")
	}

	// Test join party callback
	response = engine.ProcessInput("HELP")
	if !callbacks.joinPartyCalled {
		t.Error("Expected JoinParty callback to be called")
	}

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}
}

func TestLinearConversationEngine_Stop(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	if !engine.IsActive() {
		t.Error("Expected engine to be active after start")
	}

	engine.Stop()

	if engine.IsActive() {
		t.Error("Expected engine to be inactive after stop")
	}
}

// Integration test simulating a full conversation
func TestLinearConversationEngine_FullConversation(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  false, // First meeting
	}

	engine := NewLinearConversationEngine(script, callbacks)

	// Start conversation
	response := engine.Start(1)
	if !strings.Contains(response.Output, "I am called Treanna") {
		t.Error("Expected introduction for first meeting")
	}

	// Ask about name
	response = engine.ProcessInput("NAME")
	if !strings.Contains(response.Output, "My name is Treanna") {
		t.Error("Expected name response")
	}

	// Ask about job
	response = engine.ProcessInput("JOB")
	if !strings.Contains(response.Output, "keeper of ancient knowledge") {
		t.Error("Expected job response")
	}

	// Ask about virtue
	response = engine.ProcessInput("VIRTUE")
	if !strings.Contains(response.Output, "virtues guide our path") {
		t.Error("Expected virtue response")
	}

	// Say goodbye
	response = engine.ProcessInput("BYE")
	if !response.IsComplete {
		t.Error("Expected conversation to end")
	}
	if !strings.Contains(response.Output, "Farewell, brave TestHero") {
		t.Error("Expected farewell message")
	}
}

func TestLinearConversationEngine_LabelNavigation(t *testing.T) {
	// Create script with labels and jumps
	script := &references.TalkScript{
		Lines: []references.ScriptLine{
			// Name (index 0)
			{
				{Cmd: references.PlainString, Str: "Wizard"},
			},
			// Description (index 1)
			{
				{Cmd: references.PlainString, Str: "a wise old wizard"},
			},
			// Greeting (index 2)
			{
				{Cmd: references.PlainString, Str: "Greetings, seeker."},
			},
			// Job (index 3)
			{
				{Cmd: references.PlainString, Str: "I study the mysteries of magic."},
				{Cmd: references.GotoLabel, Num: 1}, // Jump to Label1
			},
			// Bye (index 4)
			{
				{Cmd: references.PlainString, Str: "Farewell."},
			},
			// Label1 content (index 5)
			{
				{Cmd: references.DefineLabel},
				{Cmd: references.Label1},
				{Cmd: references.PlainString, Str: "You have been transported by magic!"},
			},
		},
		QuestionGroups: []references.QuestionGroup{},
	}

	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)

	// Check that label map was built
	if len(engine.labelMap) == 0 {
		t.Error("Expected label map to be built")
	}

	// Verify Label1 is mapped correctly
	if position, exists := engine.labelMap[references.Label1]; !exists {
		t.Error("Expected Label1 to be in label map")
	} else if position != 5 {
		t.Errorf("Expected Label1 to map to position 5, got %d", position)
	}

	engine.Start(1)

	// Ask about job which should trigger goto Label1
	response := engine.ProcessInput("JOB")

	if response.Error != nil {
		t.Errorf("Unexpected error: %v", response.Error)
	}

	if !strings.Contains(response.Output, "I study the mysteries of magic") {
		t.Error("Expected job response")
	}
}

func TestLinearConversationEngine_GotoLabelError(t *testing.T) {
	script := createTestScript()
	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true,
	}

	engine := NewLinearConversationEngine(script, callbacks)

	// Try to goto a non-existent label
	err := engine.gotoLabel(references.Label5)

	if err == nil {
		t.Error("Expected error when going to non-existent label")
	}

	if !strings.Contains(err.Error(), "label") && !strings.Contains(err.Error(), "not found") {
		t.Error("Expected 'label not found' error message")
	}
}

func TestLinearConversationEngine_IfElseKnowsName_HasMet(t *testing.T) {
	// Create script with IfElseKnowsName conditional
	script := &references.TalkScript{
		Lines: []references.ScriptLine{
			// Name (index 0)
			{
				{Cmd: references.PlainString, Str: "Guard"},
			},
			// Description (index 1)
			{
				{Cmd: references.PlainString, Str: "a town guard"},
			},
			// Greeting (index 2)
			{
				{Cmd: references.PlainString, Str: "Halt! "},
				{Cmd: references.IfElseKnowsName},
				{Cmd: references.PlainString, Str: "Good to see thee again, "},
				{Cmd: references.PlainString, Str: "State thy business, stranger."},
				{Cmd: references.AvatarsName},
			},
			// Job (index 3)
			{
				{Cmd: references.PlainString, Str: "I guard this town."},
			},
			// Bye (index 4)
			{
				{Cmd: references.PlainString, Str: "Move along."},
			},
		},
		QuestionGroups: []references.QuestionGroup{},
	}

	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  true, // Has met before
	}

	engine := NewLinearConversationEngine(script, callbacks)
	response := engine.Start(1)

	// Debug output
	t.Logf("Actual output: %q", response.Output)
	if response.Error != nil {
		t.Logf("Error: %v", response.Error)
	}

	// Check for errors first
	if response.Error != nil {
		t.Fatalf("Unexpected error: %v", response.Error)
	}

	// Should show the "has met" branch
	if !strings.Contains(response.Output, "Good to see thee again") {
		t.Error("Expected 'has met' branch output")
	}

	if !strings.Contains(response.Output, "TestHero") {
		t.Error("Expected avatar name in output")
	}
}

func TestLinearConversationEngine_IfElseKnowsName_FirstMeeting(t *testing.T) {
	// Test IfElseKnowsName conditional in QuestionGroup context
	script := &references.TalkScript{
		Lines: []references.ScriptLine{
			// Name (index 0)
			{
				{Cmd: references.PlainString, Str: "Guard"},
			},
			// Description (index 1)
			{
				{Cmd: references.PlainString, Str: "a town guard"},
			},
			// Greeting (index 2)
			{
				{Cmd: references.PlainString, Str: "Halt!"},
			},
			// Job (index 3)
			{
				{Cmd: references.PlainString, Str: "I guard this town."},
			},
			// Bye (index 4)
			{
				{Cmd: references.PlainString, Str: "Move along."},
			},
		},
		QuestionGroups: []references.QuestionGroup{
			{
				Options: []string{"TEST"},
				Script: references.ScriptLine{
					{Cmd: references.PlainString, Str: "Testing "},
					{Cmd: references.IfElseKnowsName},
					{Cmd: references.PlainString, Str: "met before"},
					{Cmd: references.PlainString, Str: "first time"},
				},
			},
		},
	}

	callbacks := &TestActionCallbacks{
		avatarName: "TestHero",
		hasMetNPC:  false, // First meeting
	}

	engine := NewLinearConversationEngine(script, callbacks)
	engine.Start(1)

	// Test the conditional in a keyword response
	response := engine.ProcessInput("TEST")

	// Should show the "first meeting" branch (first time)
	if !strings.Contains(response.Output, "first time") {
		t.Error("Expected 'first meeting' branch output")
	}

	// Should NOT show the "has met" text
	if strings.Contains(response.Output, "met before") {
		t.Error("Should not show 'has met' text for first meeting")
	}
}
