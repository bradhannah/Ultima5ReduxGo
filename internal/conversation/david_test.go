package conversation

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// DavidTestCallbacks implements ActionCallbacks for David's conversation testing
type DavidTestCallbacks struct {
	avatarName    string
	metNPCs       map[int]bool
	karmaLevel    int
	itemsReceived []int // Track items received
}

func NewDavidTestCallbacks(avatarName string) *DavidTestCallbacks {
	return &DavidTestCallbacks{
		avatarName:    avatarName,
		metNPCs:       make(map[int]bool),
		karmaLevel:    50,
		itemsReceived: make([]int, 0),
	}
}

func (t *DavidTestCallbacks) JoinParty() error                           { return nil }
func (t *DavidTestCallbacks) CallGuards() error                          { return nil }
func (t *DavidTestCallbacks) IncreaseKarma() error                       { t.karmaLevel++; return nil }
func (t *DavidTestCallbacks) DecreaseKarma() error                       { t.karmaLevel--; return nil }
func (t *DavidTestCallbacks) GoToJail() error                            { return nil }
func (t *DavidTestCallbacks) MakeHorse() error                           { return nil }
func (t *DavidTestCallbacks) PayExtortion(amount int) error              { return nil }
func (t *DavidTestCallbacks) PayHalfExtortion() error                    { return nil }
func (t *DavidTestCallbacks) GetUserInput(prompt string) (string, error) { return "", nil }
func (t *DavidTestCallbacks) AskPlayerName() (string, error)             { return t.avatarName, nil }
func (t *DavidTestCallbacks) GetGoldAmount(prompt string) (int, error)   { return 0, nil }
func (t *DavidTestCallbacks) ShowOutput(text string)                     {}
func (t *DavidTestCallbacks) WaitForKeypress()                           {}
func (t *DavidTestCallbacks) TimedPause()                                { time.Sleep(1 * time.Millisecond) }
func (t *DavidTestCallbacks) HasMet(npcID int) bool                      { return t.metNPCs[npcID] }
func (t *DavidTestCallbacks) GetAvatarName() string                      { return t.avatarName }
func (t *DavidTestCallbacks) GetKarmaLevel() int                         { return t.karmaLevel }
func (t *DavidTestCallbacks) OnError(err error)                          {}
func (t *DavidTestCallbacks) SetMet(npcID int)                           { t.metNPCs[npcID] = true }

// GiveItem tracks when items are given to the player
func (t *DavidTestCallbacks) GiveItem(itemID int) error {
	t.itemsReceived = append(t.itemsReceived, itemID)
	return nil
}

// HasReceivedItem checks if player received a specific item
func (t *DavidTestCallbacks) HasReceivedItem(itemID int) bool {
	for _, receivedItem := range t.itemsReceived {
		if receivedItem == itemID {
			return true
		}
	}
	return false
}

// loadDavidScript loads David's conversation script for testing
// We'll need to find David's actual TLK index - trying a few common indices
func loadDavidScript(t *testing.T) (*references.TalkScript, int) {
	cfg := config.NewUltimaVConfiguration()
	tlkPath := filepath.Join(cfg.SavedConfigData.DataFilePath, "CASTLE.TLK")

	talkData, err := references.LoadFile(tlkPath)
	if err != nil {
		t.Fatalf("Failed to load TLK file: %v", err)
	}

	dataOvl := references.NewDataOvl(cfg)
	wordDict := references.NewWordDict(dataOvl.CompressedWords)

	// Try different indices to find David - we'll check for "crotchety" or "sext" keywords
	possibleIndices := []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	for _, index := range possibleIndices {
		if davidData, exists := talkData[index]; exists {
			script, err := references.ParseNPCBlob(davidData, wordDict)
			if err != nil {
				continue // Try next index
			}

			// Check if this looks like David by looking for "sext" keyword
			for _, group := range script.QuestionGroups {
				for _, option := range group.Options {
					if strings.Contains(strings.ToLower(option), "sext") {
						t.Logf("Found David at TLK index %d", index)
						return script, index
					}
				}
			}
		}
	}

	t.Skip("Could not find David in CASTLE.TLK - may need to check other TLK files or add David data")
	return nil, -1
}

// TestGiveItemCallback tests that the GiveItem callback is properly called
func TestGiveItemCallback(t *testing.T) {
	t.Run("GiveItemCallback", func(t *testing.T) {
		callbacks := NewDavidTestCallbacks("TestHero")

		// Directly test the callback functionality
		err := callbacks.GiveItem(72)
		if err != nil {
			t.Errorf("GiveItem callback failed: %v", err)
		}

		// Check that we received Item 72
		if !callbacks.HasReceivedItem(72) {
			t.Errorf("Expected to receive Item 72, but didn't. Items received: %v", callbacks.itemsReceived)
		}

		// Test multiple items
		callbacks.GiveItem(100)
		callbacks.GiveItem(200)

		if len(callbacks.itemsReceived) != 3 {
			t.Errorf("Expected 3 items received, got: %v", callbacks.itemsReceived)
		}

		expectedItems := []int{72, 100, 200}
		for _, expectedItem := range expectedItems {
			if !callbacks.HasReceivedItem(expectedItem) {
				t.Errorf("Expected to have received item %d, but didn't. Items received: %v", expectedItem, callbacks.itemsReceived)
			}
		}

		t.Logf("Successfully tested GiveItem callback functionality")
		t.Logf("Items received: %v", callbacks.itemsReceived)
	})
}
