package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/conversation"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// GameActionCallbacks implements the ActionCallbacks interface to integrate
// the LinearConversationEngine with the game's state and UI systems
type GameActionCallbacks struct {
	gameScene *GameScene
	npcRef    references.NPCReference
	output    strings.Builder
}

// NewGameActionCallbacks creates a new callback implementation for the given game scene and NPC
func NewGameActionCallbacks(gameScene *GameScene, npcRef references.NPCReference) *GameActionCallbacks {
	return &GameActionCallbacks{
		gameScene: gameScene,
		npcRef:    npcRef,
	}
}

// GetAccumulatedOutput returns all text that has been accumulated via ShowOutput
func (g *GameActionCallbacks) GetAccumulatedOutput() string {
	return g.output.String()
}

// ClearAccumulatedOutput resets the accumulated output
func (g *GameActionCallbacks) ClearAccumulatedOutput() {
	g.output.Reset()
}

// === Core Action Callbacks ===

func (g *GameActionCallbacks) JoinParty() error {
	if !g.gameScene.gameState.PartyState.HasRoom() {
		g.ShowOutput("My party is full.\n")
		return nil
	}

	if err := g.gameScene.gameState.PartyState.JoinNPC(g.npcRef); err != nil {
		g.ShowOutput(fmt.Sprintf("%v\n", err))
		return err
	}

	// Note: NPC name will be provided by the calling TalkCommand, not NPC reference
	g.ShowOutput("NPC has joined thee!\n")
	return nil
}

func (g *GameActionCallbacks) CallGuards() error {
	// TODO: Implement guard calling logic
	g.ShowOutput("[GUARDS CALLED - Not yet implemented]\n")
	log.Printf("CallGuards: NPC at %+v called guards", g.npcRef.Position)
	return nil
}

func (g *GameActionCallbacks) IncreaseKarma() error {
	g.gameScene.gameState.PartyState.Karma.IncreaseKarma(1)
	g.ShowOutput("[KARMA INCREASED]\n")
	return nil
}

func (g *GameActionCallbacks) DecreaseKarma() error {
	g.gameScene.gameState.PartyState.Karma.DecreaseKarma(1)
	g.ShowOutput("[KARMA DECREASED]\n")
	return nil
}

func (g *GameActionCallbacks) GoToJail() error {
	// TODO: Implement jail logic
	g.ShowOutput("[SENT TO JAIL - Not yet implemented]\n")
	log.Printf("GoToJail: Avatar sent to jail by NPC at %+v", g.npcRef.Position)
	return nil
}

func (g *GameActionCallbacks) MakeHorse() error {
	// TODO: Implement horse creation logic
	g.ShowOutput("[HORSE CREATED - Not yet implemented]\n")
	log.Printf("MakeHorse: Horse created by NPC at %+v", g.npcRef.Position)
	return nil
}

func (g *GameActionCallbacks) PayExtortion(amount int) error {
	currentGold := int(g.gameScene.gameState.PartyState.Inventory.Gold.Get())
	if currentGold < amount {
		g.ShowOutput(fmt.Sprintf("Thou dost not have %d gold pieces!\n", amount))
		return fmt.Errorf("insufficient gold: have %d, need %d", currentGold, amount)
	}

	g.gameScene.gameState.PartyState.Inventory.Gold.DecrementBy(uint16(amount))
	g.ShowOutput(fmt.Sprintf("Thou dost pay %d gold pieces.\n", amount))
	return nil
}

func (g *GameActionCallbacks) PayHalfExtortion() error {
	currentGold := int(g.gameScene.gameState.PartyState.Inventory.Gold.Get())
	halfGold := currentGold / 2

	if halfGold == 0 {
		g.ShowOutput("Thou dost not have any gold!\n")
		return fmt.Errorf("no gold available")
	}

	g.gameScene.gameState.PartyState.Inventory.Gold.DecrementBy(uint16(halfGold))
	g.ShowOutput(fmt.Sprintf("Thou dost pay %d gold pieces (half thy gold).\n", halfGold))
	return nil
}

func (g *GameActionCallbacks) GiveItem(itemID int) error {
	// TODO: Implement item giving logic based on itemID
	g.ShowOutput(fmt.Sprintf("[ITEM %d GIVEN - Not yet implemented]\n", itemID))
	log.Printf("GiveItem: NPC at %+v gave item %d", g.npcRef.Position, itemID)
	return nil
}

// === Player Interaction Callbacks ===

func (g *GameActionCallbacks) GetUserInput(prompt string) (string, error) {
	// This is a stub - in practice, this would need to be handled by the UI layer
	// The LinearTalkDialog will handle this by setting response flags
	g.ShowOutput(prompt)
	return "", fmt.Errorf("GetUserInput should be handled by UI layer")
}

func (g *GameActionCallbacks) AskPlayerName() (string, error) {
	// This is a stub - similar to GetUserInput, handled by UI
	g.ShowOutput("What is thy name? ")
	return "", fmt.Errorf("AskPlayerName should be handled by UI layer")
}

func (g *GameActionCallbacks) GetGoldAmount(prompt string) (int, error) {
	// This is a stub - handled by UI layer
	g.ShowOutput(prompt)
	return 0, fmt.Errorf("GetGoldAmount should be handled by UI layer")
}

func (g *GameActionCallbacks) ShowOutput(text string) {
	g.output.WriteString(text)
}

func (g *GameActionCallbacks) WaitForKeypress() {
	// This is handled by the UI layer setting appropriate response flags
	g.ShowOutput("[Press any key to continue]")
}

func (g *GameActionCallbacks) TimedPause() {
	// This is handled by the UI layer with a 3-second pause
	g.ShowOutput("[Pausing...]")
}

// === Game State Query Callbacks ===

func (g *GameActionCallbacks) HasMet(npcID int) bool {
	return g.gameScene.gameState.PartyState.HasMet(int32(npcID))
}

func (g *GameActionCallbacks) GetAvatarName() string {
	return g.gameScene.gameState.PartyState.AvatarName()
}

func (g *GameActionCallbacks) GetKarmaLevel() int {
	// Return current karma value
	return g.gameScene.gameState.PartyState.Karma.Value
}

func (g *GameActionCallbacks) MatchPartyMemberName(inputName string) bool {
	inputLower := strings.ToLower(strings.TrimSpace(inputName))

	// Check Avatar name
	if strings.EqualFold(inputLower, g.GetAvatarName()) {
		return true
	}

	// Check party member names
	for _, character := range g.gameScene.gameState.PartyState.Characters {
		characterName := strings.TrimRight(string(character.Name[:]), string(rune(0)))
		if strings.EqualFold(inputLower, characterName) {
			return true
		}
	}

	return false
}

func (g *GameActionCallbacks) SetMet(npcID int) {
	g.gameScene.gameState.PartyState.SetMet(g.npcRef.Location, npcID)
}

// === Error Handling ===

func (g *GameActionCallbacks) OnError(err error) {
	g.ShowOutput(fmt.Sprintf("[ERROR: %v]\n", err))
	log.Printf("Conversation error with NPC at %+v: %v", g.npcRef.Position, err)
}

// === Helper Methods ===

// ExtractGoldAmountFromText extracts a numeric prefix from text (for GoldPrompt handling)
func (g *GameActionCallbacks) ExtractGoldAmountFromText(text string) (int, string, error) {
	if text == "" {
		return 0, "", fmt.Errorf("empty text")
	}

	// Find the first non-digit character
	digitEnd := 0
	for i, r := range text {
		if r < '0' || r > '9' {
			digitEnd = i
			break
		}
	}

	if digitEnd == 0 {
		return 0, text, fmt.Errorf("no numeric prefix found")
	}

	goldAmount, err := strconv.Atoi(text[:digitEnd])
	if err != nil {
		return 0, text, fmt.Errorf("failed to parse gold amount: %v", err)
	}

	remainingText := text[digitEnd:]
	return goldAmount, remainingText, nil
}

// Ensure GameActionCallbacks implements the ActionCallbacks interface
var _ conversation.ActionCallbacks = (*GameActionCallbacks)(nil)
