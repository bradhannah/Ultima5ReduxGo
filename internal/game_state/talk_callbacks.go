package game_state

import "github.com/bradhannah/Ultima5ReduxGo/internal/map_units"

// TalkDialog represents a dialog that can be pushed to the UI
type TalkDialog interface {
	// Interface for talk dialogs - implementation will be in UI layer
}

// TalkCallbacks defines the interface for talk-related UI operations
type TalkCallbacks struct {
	// CreateTalkDialog creates a new talk dialog for the given NPC
	CreateTalkDialog func(npc *map_units.NPCFriendly) TalkDialog

	// PushDialog pushes a dialog to the UI dialog stack
	PushDialog func(dialog TalkDialog)
}

// NewTalkCallbacks creates TalkCallbacks with required function validation
func NewTalkCallbacks(createTalkDialog func(*map_units.NPCFriendly) TalkDialog, pushDialog func(TalkDialog)) TalkCallbacks {
	if createTalkDialog == nil {
		createTalkDialog = func(*map_units.NPCFriendly) TalkDialog { return nil } // No-op default
	}
	if pushDialog == nil {
		pushDialog = func(TalkDialog) {} // No-op default
	}

	return TalkCallbacks{
		CreateTalkDialog: createTalkDialog,
		PushDialog:       pushDialog,
	}
}
