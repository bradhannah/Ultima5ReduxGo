# Linear Conversation System

This document describes the new linear conversation system implementation for Ultima 5 Redux Go, designed to be more straightforward and easier to understand than the previous implementation.

## Overview

The Linear Conversation Engine processes TalkScript commands sequentially using a simple pointer-based navigation system. Unlike the previous implementation that organized conversations into complex buckets, this system follows a left-to-right linear approach with a single active command pointer.

## Key Design Principles

1. **Stateless Operations**: Each conversation interaction is independent, with no retained state between conversations
2. **Linear Processing**: Commands are processed sequentially with simple pointer navigation
3. **Callback-Based Actions**: Game actions are handled through injected callbacks rather than channels
4. **Bootstrap Procedure**: Standardized NPC introduction sequence
5. **Simple Navigation**: Single pointer tracks current position, moves linearly through script

## Architecture

### Core Components

- **LinearConversationEngine**: Main engine that processes conversations
- **ActionCallbacks**: Interface for handling game actions and player interactions
- **ConversationResponse**: Response object containing output and state information
- **TalkScript**: Existing AST structure from `internal/references/tlk_script_ast.go`

### File Structure

```
internal/conversation/
├── linear_engine.go      # Main conversation engine implementation
└── linear_engine_test.go # Comprehensive test suite
```

## Usage

### Basic Setup

```go
import "github.com/bradhannah/Ultima5ReduxGo/internal/conversation"

// Implement the ActionCallbacks interface
type GameCallbacks struct {
    // ... implementation
}

func (g *GameCallbacks) HasMet(npcID int) bool {
    // Check if player has met this NPC
}

func (g *GameCallbacks) GetAvatarName() string {
    // Return player's name
}

// ... implement other required methods

// Create engine
script := loadTalkScript() // Load your TalkScript
callbacks := &GameCallbacks{}
engine := conversation.NewLinearConversationEngine(script, callbacks)
```

### Starting a Conversation

```go
// Start conversation with NPC ID 1
response := engine.Start(1)

// Display output to player
fmt.Print(response.Output)

// Check if input is needed
if response.NeedsInput {
    fmt.Print(response.InputPrompt)
    userInput := getUserInput()
    response = engine.ProcessInput(userInput)
}
```

### Processing Conversation Flow

```go
for engine.IsActive() && !response.IsComplete {
    if response.Error != nil {
        log.Printf("Conversation error: %v", response.Error)
        break
    }
    
    if response.NeedsInput {
        fmt.Print(response.InputPrompt)
        userInput := getUserInput()
        response = engine.ProcessInput(userInput)
    }
    
    fmt.Print(response.Output)
}
```

## Conversation Flow

### Bootstrap Procedure

When a conversation starts, the engine follows this sequence:

1. **Reset State**: Clear output buffer and reset pointer to 0
2. **Check Meeting Status**: Use `HasMet(npcID)` callback to determine if player has met NPC
3. **Display Description**: Show NPC description (fixed entry 1)
4. **Show Greeting**:
   - If first meeting: "I am called [Name]"
   - If already met: Use greeting text (fixed entry 2)
5. **Prompt for Input**: "Your interest?"

### Input Processing

The engine processes user input in this order:

1. **Empty Input**: Treated as "BYE" command
2. **Standard Keywords**: NAME, JOB/WORK, BYE/THANK
3. **Script Keywords**: Search through QuestionGroups for matches
4. **Default Response**: "I cannot help thee with that."

### Command Processing

Commands from the TalkScript are processed linearly:

- **PlainString**: Output text directly
- **AvatarsName**: Substitute player's name
- **NewLine**: Add line break
- **Action Commands**: Call appropriate callbacks (JoinParty, CallGuards, etc.)
- **EndConversation**: Mark conversation as complete

## ActionCallbacks Interface

All game interactions are handled through the ActionCallbacks interface:

### Game Actions
- `JoinParty()`: NPC joins the party
- `CallGuards()`: Summon guards
- `IncreaseKarma()` / `DecreaseKarma()`: Modify karma
- `GoToJail()`: Send player to jail
- `MakeHorse()`: Create a horse
- `PayExtortion(amount)` / `PayHalfExtortion()`: Handle extortion

### Player Interactions
- `GetUserInput(prompt)`: Get text input from player
- `AskPlayerName()`: Request player's name
- `GetGoldAmount(prompt)`: Request gold amount
- `ShowOutput(text)`: Display text to player
- `WaitForKeypress()`: Pause for keypress

### Game State Queries
- `HasMet(npcID)`: Check if player has met NPC
- `GetAvatarName()`: Get player's name
- `GetKarmaLevel()`: Get current karma level

### Error Handling
- `OnError(err)`: Handle errors during conversation

## Fixed Script Entries

The engine expects these fixed entries in the TalkScript.Lines array:

| Index | Constant | Purpose |
|-------|----------|---------|
| 0 | TalkScriptConstantsName | NPC's name |
| 1 | TalkScriptConstantsDescription | NPC's description |
| 2 | TalkScriptConstantsGreeting | Greeting for known players |
| 3 | TalkScriptConstantsJob | NPC's job/occupation |
| 4 | TalkScriptConstantsBye | Farewell message |

## Supported TalkCommands

The engine currently supports these commands from the TalkScript AST:

### Basic Output
- `PlainString`: Output literal text
- `AvatarsName`: Substitute player's name
- `NewLine`: Insert line break

### Flow Control
- `EndConversation`: End the conversation
- `Pause`: Pause for user keypress
- `KeyWait`: Wait for keypress

### Actions
- `JoinParty`: NPC joins party
- `CallGuards`: Call guards
- `KarmaPlusOne`: Increase karma
- `KarmaMinusOne`: Decrease karma
- `GoToJail`: Send to jail
- `MakeAHorse`: Create horse

## Testing

The system includes comprehensive tests covering:

- Engine initialization
- Bootstrap procedures for first meetings and return visits
- Standard keyword handling (NAME, JOB, BYE)
- Custom keyword processing through QuestionGroups
- Action callback execution
- Error handling and edge cases
- Full conversation flows

Run tests with:
```bash
go test ./internal/conversation/ -v
```

## Error Handling

The engine handles errors gracefully:

- **Invalid State**: Returns error if processing input on inactive engine
- **Callback Errors**: Action callback errors are propagated to caller
- **Unknown Commands**: Logged but don't stop conversation
- **Missing Data**: Gracefully handles missing script entries

## Future Extensions

The linear design makes it easy to add:

1. **Label Navigation**: Support for GotoLabel and DefineLabel commands
2. **Conditional Branches**: IfElseKnowsName and similar conditionals
3. **Question Handling**: Complex question/answer sequences
4. **Script Variables**: State tracking within conversations
5. **Extended Actions**: Additional game actions and commands

## Comparison with Old Implementation

### Advantages of Linear System

- **Simplicity**: Single pointer, linear processing
- **Debuggability**: Easy to trace conversation flow
- **Testability**: Straightforward unit testing
- **Maintainability**: Clear separation of concerns
- **Extensibility**: Easy to add new features incrementally

### Key Differences

- **No Channels**: Uses direct function calls instead of Go channels
- **Stateless**: No retained conversation state between interactions
- **Interface-Based**: Callbacks through interface rather than concrete types
- **Linear Navigation**: Simple pointer movement vs. complex bucket organization
- **Bootstrap Standard**: Consistent NPC introduction procedure

## Integration Notes

The linear conversation engine is designed to:

- Work with existing TalkScript AST definitions
- Integrate with the main game loop through callbacks
- Support the existing save game format (for HasMet functionality)
- Maintain compatibility with Ultima V conversation expectations
- Provide foundation for incremental feature additions

This implementation provides a solid, testable foundation for the conversation system while maintaining the flexibility to add more complex features as needed.