package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
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

type NPCInfo struct {
	Name       string `json:"name"`
	TLKFile    string `json:"tlk_file"`
	TLKIndex   int    `json:"tlk_index"`
	NPCFile    string `json:"npc_file"`
	NPCIndex   int    `json:"npc_index"`
	Location   string `json:"location"`
	Occupation string `json:"occupation"`
}

// loadNPCList loads available NPCs for selection
func loadNPCList() []NPCInfo {
	return []NPCInfo{
		{Name: "Alistair", TLKFile: "CASTLE.TLK", TLKIndex: 1, NPCFile: "CASTLE.NPC", NPCIndex: 1, Location: "Castle", Occupation: "Bard"},
		{Name: "Treanna", TLKFile: "CASTLE.TLK", TLKIndex: 3, NPCFile: "CASTLE.NPC", NPCIndex: 3, Location: "Castle", Occupation: "Girl"},
		{Name: "Blackthorn", TLKFile: "CASTLE.TLK", TLKIndex: 0, NPCFile: "CASTLE.NPC", NPCIndex: 0, Location: "Castle", Occupation: "King"},
		{Name: "Margaret", TLKFile: "CASTLE.TLK", TLKIndex: 2, NPCFile: "CASTLE.NPC", NPCIndex: 2, Location: "Castle", Occupation: "Cook"},
		{Name: "Chuckles", TLKFile: "CASTLE.TLK", TLKIndex: 4, NPCFile: "CASTLE.NPC", NPCIndex: 4, Location: "Castle", Occupation: "Jester"},
	}
}

// loadTalkScript loads a TalkScript from a TLK file
func loadTalkScript(npcInfo NPCInfo) (*references.TalkScript, error) {
	cfg := config.NewUltimaVConfiguration()
	tlkPath := filepath.Join(cfg.SavedConfigData.DataFilePath, npcInfo.TLKFile)

	// Load TLK file
	talkData, err := references.LoadFile(tlkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLK file: %v", err)
	}

	// Get NPC data
	npcData, exists := talkData[npcInfo.TLKIndex]
	if !exists {
		return nil, fmt.Errorf("NPC data not found at index %d in %s", npcInfo.TLKIndex, npcInfo.TLKFile)
	}

	// Load proper word dictionary from DATA.OVL
	dataOvl := references.NewDataOvl(cfg)
	wordDict := references.NewWordDict(dataOvl.CompressedWords)

	// Parse NPC's blob into a TalkScript
	script, err := references.ParseNPCBlob(npcData, wordDict)
	if err != nil {
		return nil, fmt.Errorf("failed to parse NPC data: %v", err)
	}

	return script, nil
}

func selectNPC() NPCInfo {
	npcs := loadNPCList()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Available NPCs:")
	for i, npc := range npcs {
		fmt.Printf("%d. %s (%s - %s)\n", i+1, npc.Name, npc.Location, npc.Occupation)
	}

	fmt.Print("\nSelect an NPC (1-5): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(npcs) {
		fmt.Println("Invalid choice, defaulting to Alistair")
		return npcs[0]
	}

	return npcs[choice-1]
}

func runConversation(npcInfo NPCInfo, avatarName string, hasMet bool) error {
	// Load the real TalkScript
	script, err := loadTalkScript(npcInfo)
	if err != nil {
		return fmt.Errorf("failed to load script for %s: %v", npcInfo.Name, err)
	}

	// Create callbacks
	callbacks := NewDemoCallbacks(avatarName)
	if hasMet {
		callbacks.SetMet(npcInfo.TLKIndex)
	}

	// Create conversation engine
	engine := conversation.NewLinearConversationEngine(script, callbacks)

	// Start conversation
	response := engine.Start(npcInfo.TLKIndex)
	reader := bufio.NewReader(os.Stdin)

	// Main conversation loop
	for engine.IsActive() && !response.IsComplete {
		if response.Error != nil {
			return fmt.Errorf("conversation error: %v", response.Error)
		}

		// Display output
		fmt.Print(response.Output)

		// Get input if needed
		if response.NeedsInput {
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading input: %v", err)
			}
			input = strings.TrimSpace(input)

			response = engine.ProcessInput(input)
		}
	}

	// Final output
	fmt.Print(response.Output)
	return nil
}

func main() {
	fmt.Println("=== Linear Conversation System Demo (Real TLK Data) ===")
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

	// Select NPC
	npcInfo := selectNPC()
	fmt.Printf("\nYou have chosen to speak with %s.\n\n", npcInfo.Name)

	// First meeting
	fmt.Printf("--- First Meeting with %s ---\n", npcInfo.Name)
	fmt.Printf("You approach %s...\n\n", npcInfo.Name)

	if err := runConversation(npcInfo, avatarName, false); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Ask if they want to try a return visit
	fmt.Print("\n\nWould you like to try a return visit? (y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" {
		// Second meeting (HasMet=true)
		fmt.Printf("\n--- Return Visit to %s ---\n", npcInfo.Name)
		fmt.Printf("You approach %s again...\n\n", npcInfo.Name)

		if err := runConversation(npcInfo, avatarName, true); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	fmt.Println("\n\n=== Demo Complete ===")
	fmt.Println("Try different keywords like: NAME, JOB, BYE")
	fmt.Printf("For %s, also try keywords from the test cases like: MUSI, SMIT, VAL\n", npcInfo.Name)
}
