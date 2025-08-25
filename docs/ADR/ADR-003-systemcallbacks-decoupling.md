# ADR-003: SystemCallbacks for UI Decoupling

**Status**: Accepted  
**Date**: 2025-08-25  
**Context**: Command implementation pattern and clean architecture

---

## Context

During command implementation (Tasks 1, 4, 10) and architecture cleanup (Task 8), we needed a clean way for game logic to communicate with the UI layer:

1. **Separation of Concerns**: Game logic shouldn't directly call UI functions
2. **Testability**: Need to validate UI interactions without actual UI
3. **Architecture Boundaries**: Core logic must stay independent of presentation
4. **User Feedback**: Commands need to provide messages, sounds, visual effects

### Specific Requirements
- Game logic needs to display messages to users
- Actions should trigger sound effects  
- Time advancement needs to be communicated
- Dialog system needs to be integrated
- All of this must be testable without UI

---

## Decision

Implement **SystemCallbacks** as a dependency injection pattern for UI interactions:

### 1. **SystemCallbacks Interface Design**
```go
type SystemCallbacks struct {
    Message  MessageCallbacks   // User text feedback
    Audio    AudioCallbacks     // Sound effects  
    Visual   VisualCallbacks    // Screen effects
    Screen   ScreenCallbacks    // Display updates
    Flow     FlowCallbacks      // Time, state flow
    Talk     TalkCallbacks      // Dialog system
}
```

### 2. **GameState Integration Pattern**
```go
type GameState struct {
    SystemCallbacks *SystemCallbacks
    // ... other fields
}

func (gs *GameState) ActionJimmySmallMap(direction Direction) bool {
    // Game logic validation...
    if !gs.hasKeys() {
        gs.SystemCallbacks.Message.AddRowStr("No lock picks!")
        return false
    }
    
    // Core logic...
    gs.consumeKey()
    success := gs.attemptJimmy(direction)
    
    // User feedback via callbacks
    if success {
        gs.SystemCallbacks.Message.AddRowStr("Unlocked!")
        gs.SystemCallbacks.Audio.PlaySoundEffect(SoundUnlock)
        gs.SystemCallbacks.Flow.AdvanceTime(1)
        return true
    }
    
    gs.SystemCallbacks.Message.AddRowStr("Lock pick broke!")
    return false
}
```

### 3. **Testing Pattern with MockSystemCallbacks**
```go
type MockSystemCallbacks struct {
    Messages            []string
    SoundEffectsPlayed  []SoundEffect  
    TimeAdvanced        []int
    TalkDialogCalls     []NPCFriendly
    // ... other tracking fields
}

func TestJimmySystemCallbacks(t *testing.T) {
    gs, mockCallbacks := NewIntegrationTestBuilder(t).
        WithSystemCallbacks().
        Build()
        
    result := gs.ActionJimmySmallMap(references.Up)
    
    // Validate SystemCallbacks integration
    assert.True(t, len(mockCallbacks.Messages) > 0)
    assert.Equal(t, "Not lock!", mockCallbacks.Messages[0])
}
```

### 4. **UI Implementation Pattern**
```go
// GameScene implements the callbacks
func (g *GameScene) CreateMessageCallbacks() MessageCallbacks {
    return MessageCallbacks{
        AddRowStr: func(message string) {
            g.output.AddRowStrWithTrim(message)
        },
        ClearScreen: func() {
            g.output.ClearScreen()  
        },
    }
}
```

---

## Consequences

### ✅ **Positive Consequences**

1. **Clean Separation of Concerns**
   - Game logic focuses on game mechanics
   - UI logic handles presentation details
   - Clear interface between layers

2. **Excellent Testability**
   - MockSystemCallbacks tracks all UI interactions
   - Can validate messages, sounds, time advancement
   - Game logic testable without actual UI

3. **Flexible UI Implementation**
   - UI layer can implement callbacks however appropriate
   - Could support different UI frameworks
   - Easy to change presentation without touching game logic

4. **Comprehensive Integration Validation**
   - 47 integration tests validate SystemCallbacks usage
   - Ensures all user feedback is properly implemented
   - Catches missing UI interactions

### ⚠️ **Trade-offs**

1. **Additional Complexity**
   - Requires understanding of dependency injection
   - More setup code for creating callback implementations
   - Interface design needs to be thoughtful

2. **Potential for Callback Explosion**
   - Could lead to many small callback interfaces
   - Need to balance granularity vs simplicity
   - Maintenance overhead for callback changes

3. **Indirection**
   - One more layer between game logic and UI
   - Stack traces go through callback layer
   - Slightly harder to trace execution flow

### 🔄 **Neutral Consequences**

1. **Performance**: Negligible overhead from function pointers
2. **Memory**: Small increase for callback function storage  
3. **Debugging**: Can trace callback calls through mocks

---

## Implementation Details

### Callback Categories Implemented

```go
// Message system - user text feedback
type MessageCallbacks struct {
    AddRowStr    func(message string)
    ClearScreen  func()
}

// Audio system - sound effects  
type AudioCallbacks struct {
    PlaySoundEffect func(effect SoundEffect)
}

// Flow system - time and state management
type FlowCallbacks struct {
    AdvanceTime func(minutes int)
}

// Talk system - dialog management
type TalkCallbacks struct {
    CreateTalkDialog func(npc *NPCFriendly) TalkDialog
    PushDialog       func(dialog TalkDialog)  
}
```

### Integration Testing Support
```go
// MockSystemCallbacks provides comprehensive tracking
type MockSystemCallbacks struct {
    Messages            []string
    SoundEffectsPlayed  []SoundEffect
    TimeAdvanced        []int
    TalkDialogCalls     []NPCFriendly
    DialogsPushed       []TalkDialog
    // ... methods for validation
}

func (m *MockSystemCallbacks) Reset() {
    // Clear all tracking for next test
}

func (m *MockSystemCallbacks) ToSystemCallbacks() *SystemCallbacks {
    // Convert to actual SystemCallbacks interface
}
```

### Usage Patterns by Command Type

```go
// Simple success/failure pattern
func (gs *GameState) ActionLook() bool {
    result := gs.performLook()
    gs.SystemCallbacks.Message.AddRowStr(result.message)
    gs.SystemCallbacks.Flow.AdvanceTime(1)
    return result.success
}

// Resource consumption pattern  
func (gs *GameState) ActionJimmy() bool {
    if !gs.hasKeys() {
        gs.SystemCallbacks.Message.AddRowStr("No lock picks!")
        return false
    }
    
    gs.consumeKey() // Always consume
    if gs.attemptJimmy() {
        gs.SystemCallbacks.Message.AddRowStr("Unlocked!")
        gs.SystemCallbacks.Audio.PlaySoundEffect(SoundUnlock)
    } else {
        gs.SystemCallbacks.Message.AddRowStr("Lock pick broke!")
    }
    gs.SystemCallbacks.Flow.AdvanceTime(1)
    return success
}

// Dialog integration pattern
func (gs *GameState) ActionTalk() bool {
    npc := gs.findNPC()
    if npc == nil {
        gs.SystemCallbacks.Message.AddRowStr("No-one to talk to!")
        return false
    }
    
    dialog := gs.SystemCallbacks.Talk.CreateTalkDialog(npc)
    if dialog != nil {
        gs.SystemCallbacks.Talk.PushDialog(dialog)
        gs.SystemCallbacks.Flow.AdvanceTime(1)
        return true
    }
    return false
}
```

---

## Alternatives Considered

### Alternative 1: Direct UI Calls from Game Logic
**Rejected**: Violates separation of concerns, makes testing impossible, creates tight coupling.

### Alternative 2: Event/Observer Pattern
**Considered**: More decoupled but adds complexity. SystemCallbacks provides simpler direct communication.

### Alternative 3: Return Structs with UI Instructions
**Considered**: Game logic returns data structures describing UI changes. More complex, less direct than callbacks.

### Alternative 4: Global UI Service
**Rejected**: Creates global state dependencies, harder to test, violates dependency injection principles.

---

## References

- **Integration**: Used throughout Task 10's 47 integration tests
- **Commands**: Implements pattern from CODING_CONVENTIONS.md command input patterns
- **Architecture**: Supports clean boundaries from ADR-002 DisplayManager
- **Testing**: Enables comprehensive integration testing with real game data

---

## Validation

This pattern was validated through:

1. **47 Integration Tests**: All use SystemCallbacks with comprehensive validation
2. **Command Implementation**: Jimmy, Get, Talk, Push commands all use pattern successfully
3. **UI Integration**: GameScene successfully implements all callback categories  
4. **Mock Validation**: MockSystemCallbacks accurately tracks all UI interactions

### Key Success Metrics
- **100% SystemCallbacks coverage** in Action* methods
- **Zero direct UI calls** from game logic
- **Comprehensive test validation** of all callback types
- **Clean architecture boundaries** maintained

---

## Future Considerations

1. **Callback Evolution**: May need to add new callback categories as features are implemented
2. **Performance Optimization**: Could batch callbacks if performance becomes concern
3. **Error Handling**: May need error return values from callbacks for failure scenarios
4. **Async Operations**: Currently synchronous, may need async support for complex UI operations

The SystemCallbacks pattern successfully provides clean UI decoupling while enabling comprehensive testing and maintaining excellent separation of concerns.