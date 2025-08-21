package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/conversation"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// DemoCallbacks implements the ActionCallbacks interface for demonstration
type DemoCallbacks struct {
	avatarName string
	metNPCs    map[int]bool
	karmaLevel int
}

func NewDemoCallbacks(avatarName string) *DemoCallbacks {
	return &DemoCallbacks{
		avatarName: avatarName,
		metNPCs:    make(map[int]bool),
		karmaLevel: 50,
	}
}

// Game action callbacks
func (d *DemoCallbacks) JoinParty() error {
	fmt.Println("\n[Game Action: NPC joins your party!]")
	return nil
}

func (d *DemoCallbacks) CallGuards() error {
	fmt.Println("\n[Game Action: Guards have been called!]")
	return nil
}

func (d *DemoCallbacks) IncreaseKarma() error {
	d.karmaLevel++
	fmt.Printf("\n[Game Action: Karma increased to %d]", d.karmaLevel)
	return nil
}

func (d *DemoCallbacks) DecreaseKarma() error {
	d.karmaLevel--
	fmt.Printf("\n[Game Action: Karma decreased to %d]", d.karmaLevel)
	return nil
}

func (d *DemoCallbacks) GoToJail() error {
	fmt.Println("\n[Game Action: You have been sent to jail!]")
	return nil
}

func (d *DemoCallbacks) MakeHorse() error {
	fmt.Println("\n[Game Action: A horse appears!]")
	return nil
}

func (d *DemoCallbacks) PayExtortion(amount int) error {
	fmt.Printf("\n[Game Action: You pay %d gold in extortion]", amount)
	return nil
}

func (d *DemoCallbacks) PayHalfExtortion() error {
	fmt.Println("\n[Game Action: You pay half your gold in extortion]")
	return nil
}

// Player interaction callbacks
func (d *DemoCallbacks) GetUserInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func (d *DemoCallbacks) AskPlayerName() (string, error) {
	return d.GetUserInput("What is thy name? ")
}

func (d *DemoCallbacks) GetGoldAmount(prompt string) (int, error) {
	input, err := d.GetUserInput(prompt)
	if err != nil {
		return 0, err
	}
	// Simple conversion for demo
	if input == "100" {
		return 100, nil
	}
	return 0, nil
}

func (d *DemoCallbacks) ShowOutput(text string) {
	fmt.Print(text)
}

func (d *DemoCallbacks) WaitForKeypress() {
	fmt.Print("\n[Press Enter to continue...]")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// Game state queries
func (d *DemoCallbacks) HasMet(npcID int) bool {
	return d.metNPCs[npcID]
}

func (d *DemoCallbacks) GetAvatarName() string {
	return d.avatarName
}

func (d *DemoCallbacks) GetKarmaLevel() int {
	return d.karmaLevel
}

func (d *DemoCallbacks) OnError(err error) {
	fmt.Printf("\n[Error: %v]", err)
}

// Mark NPC as met
func (d *DemoCallbacks) SetMet(npcID int) {
	d.metNPCs[npcID] = true
}

// createDemoScript creates a sample TalkScript for the demo
func createDemoScript() *references.TalkScript {
	return &references.TalkScript{
		Lines: []references.ScriptLine{
			// Name (index 0)
			{
				{Cmd: references.PlainString, Str: "Treanna"},
			},
			// Description (index 1)
			{
				{Cmd: references.PlainString, Str: "a mysterious woman in flowing robes"},
			},
			// Greeting (index 2)
			{
				{Cmd: references.PlainString, Str: "Greetings again, "},
				{Cmd: references.AvatarsName},
				{Cmd: references.PlainString, Str: ". I have been expecting thee."},
			},
			// Job (index 3)
			{
				{Cmd: references.PlainString, Str: "I am a keeper of ancient knowledge and seeker of the Eight Virtues."},
				{Cmd: references.NewLine},
				{Cmd: references.PlainString, Str: "My studies have revealed many secrets."},
			},
			// Bye (index 4)
			{
				{Cmd: references.PlainString, Str: "May the virtues guide thy path, "},
				{Cmd: references.AvatarsName},
				{Cmd: references.PlainString, Str: ". Farewell."},
			},
		},
		QuestionGroups: []references.QuestionGroup{
			{
				Options: []string{"VIRTUE", "VIRTUES", "HONOR", "JUSTICE", "COMPASSION"},
				Script: references.ScriptLine{
					{Cmd: references.PlainString, Str: "The Eight Virtues are the foundation of all that is good."},
					{Cmd: references.NewLine},
					{Cmd: references.PlainString, Str: "Through Honesty, Compassion, Valor, Justice, Sacrifice, Honor, Spirituality, and Humility, one achieves enlightenment."},
					{Cmd: references.KarmaPlusOne},
				},
			},
			{
				Options: []string{"MAGIC", "SPELL", "SPELLS", "RUNE", "RUNES"},
				Script: references.ScriptLine{
					{Cmd: references.PlainString, Str: "Magic flows through all things in Britannia, if one knows how to perceive it."},
					{Cmd: references.NewLine},
					{Cmd: references.PlainString, Str: "The runes hold power beyond mortal understanding."},
				},
			},
			{
				Options: []string{"JOIN", "PARTY", "HELP", "TRAVEL"},
				Script: references.ScriptLine{
					{Cmd: references.PlainString, Str: "Thy quest calls to me, Avatar. I shall join thee in thy noble cause!"},
					{Cmd: references.JoinParty},
				},
			},
			{
				Options: []string{"KNOWLEDGE", "SECRETS", "ANCIENT", "LORE"},
				Script: references.ScriptLine{
					{Cmd: references.PlainString, Str: "Knowledge is power, but with power comes responsibility."},
					{Cmd: references.NewLine},
					{Cmd: references.PlainString, Str: "I sense great wisdom within thee. This knowledge may aid thy quest."},
					{Cmd: references.Pause},
					{Cmd: references.PlainString, Str: "But remember - some secrets are better left undisturbed."},
				},
			},
		},
	}
}

func main() {
	fmt.Println("=== Linear Conversation System Demo ===")
	fmt.Println()

	// Get player name
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter thy name, brave adventurer: ")
	avatarName, _ := reader.ReadString('\n')
	avatarName = strings.TrimSpace(avatarName)
	if avatarName == "" {
		avatarName = "Hero"
	}

	fmt.Printf("\nWelcome, %s!\n\n", avatarName)

	// Create callbacks and script
	callbacks := NewDemoCallbacks(avatarName)
	script := createDemoScript()

	// Create conversation engine
	engine := conversation.NewLinearConversationEngine(script, callbacks)

	// Demo: First meeting
	fmt.Println("--- First Meeting ---")
	fmt.Println("You see a mysterious woman in flowing robes...")
	fmt.Println()

	npcID := 1
	response := engine.Start(npcID)

	// Main conversation loop
	for engine.IsActive() && !response.IsComplete {
		if response.Error != nil {
			fmt.Printf("Error: %v\n", response.Error)
			break
		}

		// Display output
		fmt.Print(response.Output)

		// Get input if needed
		if response.NeedsInput {
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				break
			}
			input = strings.TrimSpace(input)

			response = engine.ProcessInput(input)
		}
	}

	// Final output
	fmt.Print(response.Output)

	// Mark as met for second demo
	callbacks.SetMet(npcID)

	// Demo: Second meeting
	fmt.Println("\n\n--- Return Visit ---")
	fmt.Println("You approach the woman again...")
	fmt.Println()

	engine2 := conversation.NewLinearConversationEngine(script, callbacks)
	response = engine2.Start(npcID)

	// Second conversation loop
	for engine2.IsActive() && !response.IsComplete {
		if response.Error != nil {
			fmt.Printf("Error: %v\n", response.Error)
			break
		}

		// Display output
		fmt.Print(response.Output)

		// Get input if needed
		if response.NeedsInput {
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				break
			}
			input = strings.TrimSpace(input)

			response = engine2.ProcessInput(input)
		}
	}

	// Final output
	fmt.Print(response.Output)

	fmt.Println("\n\n=== Demo Complete ===")
	fmt.Println("Try different keywords like: NAME, JOB, VIRTUE, MAGIC, JOIN, KNOWLEDGE, BYE")
}
