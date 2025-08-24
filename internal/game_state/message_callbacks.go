package game_state

import "fmt"

// MessageCallbacks provides dependency injection for UI message output functions.
// This allows game logic to output messages without depending on UI implementation.
// Use NewMessageCallbacks() to create with validation.
type MessageCallbacks struct {
	// AddRowStr adds a new line of text to the message output window
	AddRowStr func(message string)

	// AppendToCurrentRowStr appends text to the current line in the message output window
	AppendToCurrentRowStr func(message string)

	// ShowCommandPrompt displays a command prompt (e.g., "Push-", "Look-")
	ShowCommandPrompt func(command string)
}

// NewMessageCallbacks creates MessageCallbacks with required function validation
func NewMessageCallbacks(addRowStr, appendToCurrentRowStr, showCommandPrompt func(string)) (MessageCallbacks, error) {
	if addRowStr == nil {
		return MessageCallbacks{}, fmt.Errorf("AddRowStr function is required")
	}
	if appendToCurrentRowStr == nil {
		return MessageCallbacks{}, fmt.Errorf("AppendToCurrentRowStr function is required")
	}
	if showCommandPrompt == nil {
		showCommandPrompt = addRowStr // Default to AddRowStr if not provided
	}

	return MessageCallbacks{
		AddRowStr:             addRowStr,
		AppendToCurrentRowStr: appendToCurrentRowStr,
		ShowCommandPrompt:     showCommandPrompt,
	}, nil
}
