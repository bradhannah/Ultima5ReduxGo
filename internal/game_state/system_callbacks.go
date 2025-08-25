package game_state

import (
	"fmt"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
)

// SoundEffect represents different types of sound effects
type SoundEffect int

const (
	SoundCannonFire SoundEffect = iota
	SoundBroadside
	SoundGlassBreak
	SoundEarthquake
	SoundSummon
	SoundAbsorb
	SoundPushObject
	SoundOpenDoor
	SoundUnlock
	SoundTrapTrigger
	SoundStepOnTrap
)

// SystemCallbacks provides comprehensive dependency injection for all external systems
// that game logic needs to interact with, maintaining separation of concerns.
// Use NewSystemCallbacks() to create with validation.
type SystemCallbacks struct {
	// Message Output System
	Message MessageCallbacks

	// Visual Effects System
	Visual VisualCallbacks

	// Audio System
	Audio AudioCallbacks

	// Screen Update System
	Screen ScreenCallbacks

	// Game Flow System
	Flow FlowCallbacks

	// Talk Dialog System
	Talk TalkCallbacks
}

// VisualCallbacks handles visual effects
type VisualCallbacks struct {
	// KapowAt shows explosion/impact effect at screen coordinates
	KapowAt func(screenX, screenY int)

	// ShowMissileEffect displays projectile animation
	ShowMissileEffect func(fromX, fromY, toX, toY int, projectileType string)

	// DelayGlide adds a short visual pause for animations
	DelayGlide func()
}

// NewVisualCallbacks creates VisualCallbacks with required function validation
func NewVisualCallbacks(kapowAt func(int, int), showMissileEffect func(int, int, int, int, string), delayGlide func()) VisualCallbacks {
	if kapowAt == nil {
		kapowAt = func(int, int) {} // No-op default
	}
	if showMissileEffect == nil {
		showMissileEffect = func(int, int, int, int, string) {} // No-op default
	}
	if delayGlide == nil {
		delayGlide = func() {} // No-op default
	}

	return VisualCallbacks{
		KapowAt:           kapowAt,
		ShowMissileEffect: showMissileEffect,
		DelayGlide:        delayGlide,
	}
}

// AudioCallbacks handles sound effects using enum-based system
type AudioCallbacks struct {
	// PlaySoundEffect plays a sound effect by enum type
	PlaySoundEffect func(effect SoundEffect)
}

// NewAudioCallbacks creates AudioCallbacks with required function validation
func NewAudioCallbacks(playSoundEffect func(SoundEffect)) AudioCallbacks {
	if playSoundEffect == nil {
		playSoundEffect = func(SoundEffect) {} // No-op default
	}

	return AudioCallbacks{
		PlaySoundEffect: playSoundEffect,
	}
}

// ScreenCallbacks handles screen management and UI updates
type ScreenCallbacks struct {
	// MarkStatsChanged flags that stats display needs refresh
	MarkStatsChanged func()

	// UpdateStatsDisplay immediately refreshes character stats
	UpdateStatsDisplay func()

	// RefreshInventoryDisplay updates inventory UI
	RefreshInventoryDisplay func()

	// ShowModalDialog displays a modal dialog box
	ShowModalDialog func(title, message string) string

	// PromptYesNo asks yes/no question and returns true for yes
	PromptYesNo func(question string) bool
}

// NewScreenCallbacks creates ScreenCallbacks with required function validation
func NewScreenCallbacks(markStatsChanged, updateStatsDisplay, refreshInventoryDisplay func(),
	showModalDialog func(string, string) string, promptYesNo func(string) bool) ScreenCallbacks {

	if markStatsChanged == nil {
		markStatsChanged = func() {} // No-op default
	}
	if updateStatsDisplay == nil {
		updateStatsDisplay = func() {} // No-op default
	}
	if refreshInventoryDisplay == nil {
		refreshInventoryDisplay = func() {} // No-op default
	}
	if showModalDialog == nil {
		showModalDialog = func(string, string) string { return "" } // No-op default
	}
	if promptYesNo == nil {
		promptYesNo = func(string) bool { return false } // No-op default
	}

	return ScreenCallbacks{
		MarkStatsChanged:        markStatsChanged,
		UpdateStatsDisplay:      updateStatsDisplay,
		RefreshInventoryDisplay: refreshInventoryDisplay,
		ShowModalDialog:         showModalDialog,
		PromptYesNo:             promptYesNo,
	}
}

// FlowCallbacks handles game flow and timing
type FlowCallbacks struct {
	// FinishTurn completes current turn and advances game state
	FinishTurn func()

	// ActivateGuards triggers town guard response to aggression
	ActivateGuards func()

	// AdvanceTime advances UltimaDate by specified minutes (integrates with existing TimeOfDay system)
	AdvanceTime func(minutes int)

	// SetTimeOfDay jumps to specific time period using TimeOfDay enum
	SetTimeOfDay func(timeOfDay datetime.TimeOfDay)

	// DelayFx adds a pause for effect timing
	DelayFx func()

	// CheckUpdate processes pending updates
	CheckUpdate func()
}

// NewFlowCallbacks creates FlowCallbacks with required function validation
func NewFlowCallbacks(finishTurn, activateGuards func(), advanceTime func(int), setTimeOfDay func(datetime.TimeOfDay), delayFx, checkUpdate func()) FlowCallbacks {
	if finishTurn == nil {
		finishTurn = func() {} // No-op default
	}
	if activateGuards == nil {
		activateGuards = func() {} // No-op default
	}
	if advanceTime == nil {
		advanceTime = func(int) {} // No-op default
	}
	if setTimeOfDay == nil {
		setTimeOfDay = func(datetime.TimeOfDay) {} // No-op default
	}
	if delayFx == nil {
		delayFx = func() {} // No-op default
	}
	if checkUpdate == nil {
		checkUpdate = func() {} // No-op default
	}

	return FlowCallbacks{
		FinishTurn:     finishTurn,
		ActivateGuards: activateGuards,
		AdvanceTime:    advanceTime,
		SetTimeOfDay:   setTimeOfDay,
		DelayFx:        delayFx,
		CheckUpdate:    checkUpdate,
	}
}

// NewSystemCallbacks creates SystemCallbacks with validation for all subsystems
func NewSystemCallbacks(message MessageCallbacks, visual VisualCallbacks, audio AudioCallbacks, screen ScreenCallbacks, flow FlowCallbacks, talk TalkCallbacks) (*SystemCallbacks, error) {
	// Validate that all callback structs have required functions
	if message.AddRowStr == nil {
		return nil, fmt.Errorf("MessageCallbacks.AddRowStr is required")
	}
	if message.AppendToCurrentRowStr == nil {
		return nil, fmt.Errorf("MessageCallbacks.AppendToCurrentRowStr is required")
	}
	if visual.KapowAt == nil {
		return nil, fmt.Errorf("VisualCallbacks.KapowAt is required")
	}
	if audio.PlaySoundEffect == nil {
		return nil, fmt.Errorf("AudioCallbacks.PlaySoundEffect is required")
	}

	return &SystemCallbacks{
		Message: message,
		Visual:  visual,
		Audio:   audio,
		Screen:  screen,
		Flow:    flow,
		Talk:    talk,
	}, nil
}
