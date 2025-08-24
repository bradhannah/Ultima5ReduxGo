# Message Output System

This document describes the message output system that displays text on the right-hand side of the screen, showing movement results, command outcomes, current commands, and other game feedback.

## Architecture Overview

The message output system follows a dependency injection pattern to maintain separation between game logic (`internal/game_state`) and UI concerns (`cmd/ultimav`).

### Core Components

- **UI Layer**: `GameScene.output` (type `text.Output`) - The main message window
- **Game Logic Layer**: Action methods in `GameState` - Need to output messages but shouldn't depend on UI
- **Dependency Bridge**: Function pointer struct injected into `GameState`

## Current Implementation

### UI Layer (`cmd/ultimav/gamescene.go`)

```go
type GameScene struct {
    output text.Output  // Main message window on right side of screen
    // ... other fields
}
```

**Key Methods:**
- `addRowStr(message string)` - Adds a new line of text to the output window
- `appendToCurrentRowStr(message string)` - Appends text to the current line

**Usage Examples:**
```go
g.addRowStr("Push-")                    // Command prompt
g.addRowStr("Won't budge!")            // Command result
g.appendToCurrentRowStr("North")       // Direction completion
```

### Game Logic Layer (`internal/game_state`)

Action methods like `ActionPushSmallMap` need to output messages but currently have TODO comments:

```go
func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
    if g.ObjectPresentAt(pushableThingPosition) || !g.IsPushable(pushableThingTile) {
        // TODO: Show message when messaging system is available
        // g.ShowMessage("Won't budge!\n")
        return false
    }
    
    // Push logic...
    
    // TODO: Show message when messaging system is available
    // g.ShowMessage("Pushed!\n")
}
```

## Proposed Dependency Injection Solution

### Message Output Interface

Create a message output interface that the game logic can use:

```go
// internal/game_state/message_output.go
type MessageOutput interface {
    ShowMessage(message string)
    ShowMessageAppend(message string)
    ShowCommandPrompt(command string)
}
```

### Constructor-Based System Callbacks ✅ IMPLEMENTED

The system uses validated constructors and enum-based sound effects:

```go
// internal/game_state/system_callbacks.go
type SoundEffect int
const (
    SoundCannonFire SoundEffect = iota
    SoundPushObject
    SoundGlassBreak
    // ... other sound effects
)

type SystemCallbacks struct {
    Message MessageCallbacks  // Message output (addRowStr, etc.)
    Visual  VisualCallbacks   // Visual effects (kapow, missiles)
    Audio   AudioCallbacks    // Sound effects (enum-based)
    Screen  ScreenCallbacks   // UI updates (stats, inventory)
    Flow    FlowCallbacks     // Game flow (guards, time, turns)
}

// Created with validation
systemCallbacks, err := NewSystemCallbacks(messageCallbacks, visualCallbacks, ...)
```

### Constructor-Based UI Integration ✅ IMPLEMENTED

Use validated constructors during GameScene initialization:

```go
// cmd/ultimav/gamescene.go initialization
messageCallbacks, err := game_state.NewMessageCallbacks(
    gameScene.addRowStr,
    gameScene.appendToCurrentRowStr,
    gameScene.addRowStr, // command prompts
)
if err != nil {
    log.Fatalf("Failed to create MessageCallbacks: %v", err)
}

audioCallbacks := game_state.NewAudioCallbacks(
    func(effect game_state.SoundEffect) {
        // Handle sound effect by enum
    },
)

systemCallbacks, err := game_state.NewSystemCallbacks(
    messageCallbacks, visualCallbacks, audioCallbacks, screenCallbacks, flowCallbacks)
```

### Usage in Action Methods ✅ IMPLEMENTED

Direct function calls with guaranteed validity:

```go
func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
    if g.ObjectPresentAt(pushableThingPosition) || !g.IsPushable(pushableThingTile) {
        g.SystemCallbacks.Message.AddRowStr("Won't budge!")
        return false
    }
    
    // Push logic...
    g.pushIt(...)
    g.SystemCallbacks.Message.AddRowStr("Pushed!")
    g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
    
    return true
}
```

## Message Types and Patterns

### Command Prompts
- **Pattern**: `"Command-"` (e.g., `"Push-"`, `"Look-"`, `"Open-"`)
- **Usage**: Display immediately when command key pressed
- **Method**: `ShowCommandPrompt()` or `AddRowStr()`

### Command Results
- **Success Messages**: `"Pushed!"`, `"Opened!"`, `"Unlocked!"`
- **Failure Messages**: `"Won't budge!"`, `"Nothing to open!"`, `"Locked!"`
- **Usage**: Display after action completion
- **Method**: `AddRowStr()`

### Direction Completion
- **Pattern**: Append direction to existing command prompt
- **Example**: `"Push-"` + `"North"` = `"Push-North"`
- **Usage**: After direction key pressed in secondary input
- **Method**: `AppendToCurrentRowStr()`

### Movement Feedback
- **Pattern**: Direction names (`"North"`, `"South"`, `"East"`, `"West"`)
- **Usage**: After successful movement
- **Method**: `AddRowStr()`

### Multi-line Messages
- **Pattern**: Multiple `AddRowStr()` calls for complex messages
- **Example**: 
  ```go
  g.MessageCallbacks.AddRowStr("You see:")
  g.MessageCallbacks.AddRowStr("a wooden chest.")
  ```

## Testing Considerations

### Mock Message Output
For unit tests, inject a mock callback:

```go
func TestActionPush(t *testing.T) {
    var messages []string
    mockCallbacks := &MessageCallbacks{
        AddRowStr: func(msg string) { messages = append(messages, msg) },
    }
    
    gameState := NewGameState()
    gameState.MessageCallbacks = mockCallbacks
    
    result := gameState.ActionPushSmallMap(direction)
    
    assert.Contains(t, messages, "Won't budge!")
}
```

### Nil Safety
Always check for nil callbacks in game logic to prevent panics:

```go
if g.MessageCallbacks != nil && g.MessageCallbacks.AddRowStr != nil {
    g.MessageCallbacks.AddRowStr("Message")
}
```

## Implementation Benefits

### Constructor Validation
- **Required functions validated**: Prevents nil pointer panics at initialization
- **Early error detection**: Fails fast with clear error messages
- **No runtime checks needed**: Functions guaranteed to exist

### Type Safety
- **Sound effect enums**: Prevents typos and invalid sound names
- **Compile-time validation**: Invalid sound effects caught at build time
- **Easy extension**: Adding new sound effects is straightforward

### Separation of Concerns
- Game logic focuses on game rules and state
- UI layer handles presentation and formatting
- Clean boundaries between layers

### Testability  
- Game logic can be tested without UI dependencies
- Mock callbacks easily created with constructors
- Verify correct messages and sound effects are triggered

### Flexibility
- Easy to change implementations without affecting game logic
- Support for different UI implementations (console, GUI, etc.)
- No-op defaults allow partial implementation

## Integration with Command System

### Primary Input (Command Keys)
```go
// In UI layer - smallMapInputHandler
case ebiten.KeyP:
    g.addRowStr("Push-")  // Show command prompt
    g.secondaryKeyState = PushDirectionInput
```

### Secondary Input (Directions)
```go  
// In UI layer - smallMapPushSecondary
func (g *GameScene) smallMapPushSecondary(direction references.Direction) {
    success := g.gameState.ActionPushSmallMap(direction)
    if !success {
        g.addRowStr("Won't budge!")  // Fallback if action doesn't show message
    }
}
```

### Action Methods (Game Logic)
```go
// In game logic - ActionPushSmallMap
func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
    // Validation logic...
    
    if validPush {
        g.MessageCallbacks.AddRowStr("Pushed!")
        return true
    } else {
        g.MessageCallbacks.AddRowStr("Won't budge!")
        return false
    }
}
```

## Future Enhancements

### Message Categories
- Error messages (red text)
- Success messages (green text)  
- Information messages (default)
- Debug messages (filtered in release)

### Message History
- Scrollable message buffer
- Message timestamps
- Message filtering by category

### Internationalization
- Message key lookup instead of hardcoded strings
- Language-specific message formatting
- Cultural adaptations for text presentation

## Other System Dependencies Identified

Based on pseudocode analysis, these additional systems need dependency injection:

### Audio System ✅ ENUM-BASED
- **Unified interface**: `PlaySoundEffect(SoundEffect)` with enum parameter
- **Sound effects**: `SoundCannonFire`, `SoundGlassBreak`, `SoundPushObject`, etc.
- **Benefits**: Type safety, easy to extend, single callback function
- **Implementation**: AudioCallbacks with enum-based sound system

### Visual Effects System  
- **Impact effects**: `kapow_xy(x, y)` for explosions and hits
- **Projectile animations**: `show_missile_effect()` for arrows, cannonballs
- **Animation timing**: `delay_glide()` for visual pauses
- **Implementation**: VisualCallbacks with function pointers to graphics system
- **Note**: Screen redraw/LOS updates handled automatically by rendering system

### Game Flow System
- **Guard activation**: `activate_guards()` for town aggression response
- **Time management**: `addtime(minutes)` for clock advancement
- **Turn completion**: Integration with existing `FinishTurn()` method
- **Implementation**: FlowCallbacks with function pointers to game flow system

### UI State System
- **Stats updates**: `mark_stats_changed()`, `update_stats()`
- **Inventory refresh**: When items are added/removed
- **Modal dialogs**: Yes/no prompts, character selection
- **Implementation**: ScreenCallbacks with function pointers to UI system

## File Locations

- **✅ System Callbacks**: `internal/game_state/system_callbacks.go` (implemented)
- **✅ Message Callbacks**: `internal/game_state/message_callbacks.go` (implemented)
- **✅ UI Integration**: `cmd/ultimav/gamescene.go` (implemented)
- **✅ Action Methods**: `internal/game_state/action_*.go` (updated)
- **Tests**: `internal/game_state/*_test.go` (existing/new)
- **Documentation**: `docs/MESSAGE_OUTPUT_SYSTEM.md` (this file)
- **✅ TimeOfDay System**: `internal/datetime/ultima_date.go` (existing, integrated)
- **Time Documentation**: `docs/TIME_OF_DAY.md` (comprehensive time system guide)