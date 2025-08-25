package game_state

import (
	"strings"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
)

// MockSystemCallbacks provides comprehensive mock implementations for all SystemCallbacks
// interfaces, with built-in assertion helpers and call tracking.
type MockSystemCallbacks struct {
	t *testing.T

	// Message tracking
	Messages          []string
	CurrentRowContent string
	MessageCallCount  int

	// Audio tracking
	SoundEffectsPlayed []SoundEffect
	AudioCallCount     int

	// Visual tracking
	KapowCalls         []KapowCall
	MissileEffectCalls []MissileEffectCall
	DelayGlideCalls    int
	VisualCallCount    int

	// Screen tracking
	StatsMarkedChanged      int
	StatsDisplayUpdated     int
	InventoryDisplayRefresh int
	ModalDialogCalls        []ModalDialogCall
	YesNoPromptCalls        []YesNoPromptCall
	ScreenCallCount         int

	// Flow tracking
	TurnsFinished    int
	GuardActivations int
	TimeAdvanced     []int // Minutes advanced in each call
	TimeOfDayChanges []datetime.TimeOfDay
	DelayFxCalls     int
	CheckUpdateCalls int
	FlowCallCount    int

	// Talk tracking
	TalkDialogCalls []TalkDialogCall
	DialogsPushed   []TalkDialog
	TalkCallCount   int

	// Test configuration
	YesNoPromptResponse bool   // Default response for yes/no prompts
	ModalDialogResponse string // Default response for modal dialogs
}

// KapowCall tracks calls to visual KapowAt function
type KapowCall struct {
	X, Y int
}

// MissileEffectCall tracks calls to missile effect function
type MissileEffectCall struct {
	FromX, FromY   int
	ToX, ToY       int
	ProjectileType string
}

// ModalDialogCall tracks modal dialog requests
type ModalDialogCall struct {
	Title   string
	Message string
}

// YesNoPromptCall tracks yes/no prompt requests
type YesNoPromptCall struct {
	Question string
}

// TalkDialogCall tracks talk dialog creation requests
type TalkDialogCall struct {
	NPC *map_units.NPCFriendly
}

// NewMockSystemCallbacks creates a new mock system callbacks instance
func NewMockSystemCallbacks(t *testing.T) *MockSystemCallbacks {
	return &MockSystemCallbacks{
		t:                   t,
		Messages:            make([]string, 0),
		SoundEffectsPlayed:  make([]SoundEffect, 0),
		KapowCalls:          make([]KapowCall, 0),
		MissileEffectCalls:  make([]MissileEffectCall, 0),
		ModalDialogCalls:    make([]ModalDialogCall, 0),
		YesNoPromptCalls:    make([]YesNoPromptCall, 0),
		TalkDialogCalls:     make([]TalkDialogCall, 0),
		DialogsPushed:       make([]TalkDialog, 0),
		TimeAdvanced:        make([]int, 0),
		TimeOfDayChanges:    make([]datetime.TimeOfDay, 0),
		YesNoPromptResponse: true, // Default to "yes"
		ModalDialogResponse: "OK", // Default response
	}
}

// ToSystemCallbacks converts the mock to actual SystemCallbacks interface
func (m *MockSystemCallbacks) ToSystemCallbacks() *SystemCallbacks {
	messageCallbacks, err := NewMessageCallbacks(
		m.addRowStr,
		m.appendToCurrentRowStr,
		m.addRowStr, // Use same for command prompts
	)
	if err != nil {
		m.t.Fatalf("Failed to create MessageCallbacks: %v", err)
	}

	visualCallbacks := NewVisualCallbacks(
		m.kapowAt,
		m.showMissileEffect,
		m.delayGlide,
	)

	audioCallbacks := NewAudioCallbacks(
		m.playSoundEffect,
	)

	screenCallbacks := NewScreenCallbacks(
		m.markStatsChanged,
		m.updateStatsDisplay,
		m.refreshInventoryDisplay,
		m.showModalDialog,
		m.promptYesNo,
	)

	flowCallbacks := NewFlowCallbacks(
		m.finishTurn,
		m.activateGuards,
		m.advanceTime,
		m.setTimeOfDay,
		m.delayFx,
		m.checkUpdate,
	)

	talkCallbacks := NewTalkCallbacks(
		m.createTalkDialog,
		m.pushDialog,
	)

	systemCallbacks, err := NewSystemCallbacks(
		messageCallbacks,
		visualCallbacks,
		audioCallbacks,
		screenCallbacks,
		flowCallbacks,
		talkCallbacks,
	)

	if err != nil {
		m.t.Fatalf("Failed to create SystemCallbacks: %v", err)
	}

	return systemCallbacks
}

// Message callback implementations
func (m *MockSystemCallbacks) addRowStr(str string) {
	m.Messages = append(m.Messages, str)
	m.MessageCallCount++
	m.t.Logf("📝 Message: %s", str)
}

func (m *MockSystemCallbacks) appendToCurrentRowStr(str string) {
	m.CurrentRowContent += str
	m.MessageCallCount++
	m.t.Logf("📝 Append: %s", str)
}

// Visual callback implementations
func (m *MockSystemCallbacks) kapowAt(x, y int) {
	m.KapowCalls = append(m.KapowCalls, KapowCall{X: x, Y: y})
	m.VisualCallCount++
	m.t.Logf("💥 Kapow at (%d,%d)", x, y)
}

func (m *MockSystemCallbacks) showMissileEffect(fromX, fromY, toX, toY int, projectileType string) {
	m.MissileEffectCalls = append(m.MissileEffectCalls, MissileEffectCall{
		FromX: fromX, FromY: fromY,
		ToX: toX, ToY: toY,
		ProjectileType: projectileType,
	})
	m.VisualCallCount++
	m.t.Logf("🏹 Missile: %s from (%d,%d) to (%d,%d)", projectileType, fromX, fromY, toX, toY)
}

func (m *MockSystemCallbacks) delayGlide() {
	m.DelayGlideCalls++
	m.VisualCallCount++
	m.t.Logf("⏸️ Delay glide")
}

// Audio callback implementations
func (m *MockSystemCallbacks) playSoundEffect(effect SoundEffect) {
	m.SoundEffectsPlayed = append(m.SoundEffectsPlayed, effect)
	m.AudioCallCount++
	m.t.Logf("🔊 Sound: %v", effect)
}

// Screen callback implementations
func (m *MockSystemCallbacks) markStatsChanged() {
	m.StatsMarkedChanged++
	m.ScreenCallCount++
	m.t.Logf("📊 Stats marked changed")
}

func (m *MockSystemCallbacks) updateStatsDisplay() {
	m.StatsDisplayUpdated++
	m.ScreenCallCount++
	m.t.Logf("📊 Stats display updated")
}

func (m *MockSystemCallbacks) refreshInventoryDisplay() {
	m.InventoryDisplayRefresh++
	m.ScreenCallCount++
	m.t.Logf("🎒 Inventory refreshed")
}

func (m *MockSystemCallbacks) showModalDialog(title, message string) string {
	m.ModalDialogCalls = append(m.ModalDialogCalls, ModalDialogCall{
		Title:   title,
		Message: message,
	})
	m.ScreenCallCount++
	m.t.Logf("💬 Modal: %s - %s", title, message)
	return m.ModalDialogResponse
}

func (m *MockSystemCallbacks) promptYesNo(question string) bool {
	m.YesNoPromptCalls = append(m.YesNoPromptCalls, YesNoPromptCall{
		Question: question,
	})
	m.ScreenCallCount++
	m.t.Logf("❓ Yes/No: %s (returning %v)", question, m.YesNoPromptResponse)
	return m.YesNoPromptResponse
}

// Flow callback implementations
func (m *MockSystemCallbacks) finishTurn() {
	m.TurnsFinished++
	m.FlowCallCount++
	m.t.Logf("🔄 Turn finished (total: %d)", m.TurnsFinished)
}

func (m *MockSystemCallbacks) activateGuards() {
	m.GuardActivations++
	m.FlowCallCount++
	m.t.Logf("🛡️ Guards activated")
}

func (m *MockSystemCallbacks) advanceTime(minutes int) {
	m.TimeAdvanced = append(m.TimeAdvanced, minutes)
	m.FlowCallCount++
	m.t.Logf("⏰ Time advanced: %d minutes", minutes)
}

func (m *MockSystemCallbacks) setTimeOfDay(timeOfDay datetime.TimeOfDay) {
	m.TimeOfDayChanges = append(m.TimeOfDayChanges, timeOfDay)
	m.FlowCallCount++
	m.t.Logf("🌅 Time of day set: %v", timeOfDay)
}

func (m *MockSystemCallbacks) delayFx() {
	m.DelayFxCalls++
	m.FlowCallCount++
	m.t.Logf("⏳ Delay FX")
}

func (m *MockSystemCallbacks) checkUpdate() {
	m.CheckUpdateCalls++
	m.FlowCallCount++
	m.t.Logf("🔄 Check update")
}

// Talk callback implementations
func (m *MockSystemCallbacks) createTalkDialog(npc *map_units.NPCFriendly) TalkDialog {
	m.TalkDialogCalls = append(m.TalkDialogCalls, TalkDialogCall{NPC: npc})
	m.TalkCallCount++
	m.t.Logf("💭 Create talk dialog for NPC")
	// Return nil for now - can be extended with mock dialog if needed
	return nil
}

func (m *MockSystemCallbacks) pushDialog(dialog TalkDialog) {
	m.DialogsPushed = append(m.DialogsPushed, dialog)
	m.TalkCallCount++
	m.t.Logf("📋 Dialog pushed")
}

// Test assertion helpers
func (m *MockSystemCallbacks) AssertMessageContains(expected string) {
	m.t.Helper()
	for _, msg := range m.Messages {
		if strings.Contains(msg, expected) {
			return
		}
	}
	m.t.Errorf("Expected message containing '%s', but not found in: %v", expected, m.Messages)
}

func (m *MockSystemCallbacks) AssertLastMessage(expected string) {
	m.t.Helper()
	if len(m.Messages) == 0 {
		m.t.Errorf("Expected last message '%s', but no messages recorded", expected)
		return
	}
	if m.Messages[len(m.Messages)-1] != expected {
		m.t.Errorf("Expected last message '%s', got '%s'", expected, m.Messages[len(m.Messages)-1])
	}
}

func (m *MockSystemCallbacks) AssertSoundEffectPlayed(expected SoundEffect) {
	m.t.Helper()
	for _, effect := range m.SoundEffectsPlayed {
		if effect == expected {
			return
		}
	}
	m.t.Errorf("Expected sound effect %v, but not found in: %v", expected, m.SoundEffectsPlayed)
}

func (m *MockSystemCallbacks) AssertTimeAdvanced(expectedMinutes int) {
	m.t.Helper()
	totalAdvanced := 0
	for _, minutes := range m.TimeAdvanced {
		totalAdvanced += minutes
	}
	if totalAdvanced != expectedMinutes {
		m.t.Errorf("Expected %d minutes advanced, got %d", expectedMinutes, totalAdvanced)
	}
}

func (m *MockSystemCallbacks) AssertNoMessages() {
	m.t.Helper()
	if len(m.Messages) > 0 {
		m.t.Errorf("Expected no messages, but got: %v", m.Messages)
	}
}

func (m *MockSystemCallbacks) AssertNoSoundEffects() {
	m.t.Helper()
	if len(m.SoundEffectsPlayed) > 0 {
		m.t.Errorf("Expected no sound effects, but got: %v", m.SoundEffectsPlayed)
	}
}

// Reset clears all recorded calls for test reuse
func (m *MockSystemCallbacks) Reset() {
	m.Messages = m.Messages[:0]
	m.CurrentRowContent = ""
	m.MessageCallCount = 0

	m.SoundEffectsPlayed = m.SoundEffectsPlayed[:0]
	m.AudioCallCount = 0

	m.KapowCalls = m.KapowCalls[:0]
	m.MissileEffectCalls = m.MissileEffectCalls[:0]
	m.DelayGlideCalls = 0
	m.VisualCallCount = 0

	m.StatsMarkedChanged = 0
	m.StatsDisplayUpdated = 0
	m.InventoryDisplayRefresh = 0
	m.ModalDialogCalls = m.ModalDialogCalls[:0]
	m.YesNoPromptCalls = m.YesNoPromptCalls[:0]
	m.ScreenCallCount = 0

	m.TurnsFinished = 0
	m.GuardActivations = 0
	m.TimeAdvanced = m.TimeAdvanced[:0]
	m.TimeOfDayChanges = m.TimeOfDayChanges[:0]
	m.DelayFxCalls = 0
	m.CheckUpdateCalls = 0
	m.FlowCallCount = 0

	m.TalkDialogCalls = m.TalkDialogCalls[:0]
	m.DialogsPushed = m.DialogsPushed[:0]
	m.TalkCallCount = 0
}
