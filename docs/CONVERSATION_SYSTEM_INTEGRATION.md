# Conversation System Integration Solution

This document provides a complete reference for the conversation system integration, including both the legacy channel-based system and the new linear system implementation.

## Architecture Overview

### Legacy Channel-Based System
```
GameScene → TalkDialog → Conversation (goroutine) → TalkScript
    ↓           ↓              ↓                      ↓
  Input → TextInput.OnEnter → channels (In/Out) → Script Processing
```

### New Linear System with Dependency Injection (2025)
```
GameScene → ActionTalkSmallMap → SystemCallbacks.Talk → LinearTalkDialog
    ↓              ↓                      ↓                    ↓
  Input → GameState Logic → TalkCallbacks Interface → UI Dialog Stack
                                  ↓                    ↓
                        CreateTalkDialog()     ProcessInput()
                        PushDialog()              ↓
                                  ↓        ActionCallbacks (GameActionCallbacks)
                              No Circular         ↓
                              Dependencies   GameState Updates
```

## Component Details

### 1. LinearTalkDialog (`cmd/ultimav/linear_talk_dialog.go`)

**Purpose**: UI component that provides the same visual interface as the original TalkDialog but uses the LinearConversationEngine instead of channels.

**Key Features**:
- Identical UI layout and styling to original TalkDialog
- Synchronous conversation processing  
- Direct integration with LinearConversationEngine
- Response-based state management

**Critical Methods**:
- `startConversation()`: Initiates conversation with engine
- `processCurrentResponse()`: Handles engine responses  
- `onEnter()`: Processes user input through engine
- `Update()`: Displays accumulated callback output

### 2. GameActionCallbacks (`cmd/ultimav/game_action_callbacks.go`)

**Purpose**: Implementation of the ActionCallbacks interface that bridges conversation commands with game state modifications.

**Key Features**:
- Complete ActionCallbacks interface implementation
- Game state integration through GameScene
- Output accumulation for UI display
- Error handling and logging

**Critical Methods**:
- **Core Actions**: JoinParty, CallGuards, IncreaseKarma, etc.
- **Player Interactions**: GetUserInput, AskPlayerName, ShowOutput
- **Game Queries**: HasMet, GetAvatarName, GetKarmaLevel
- **Helpers**: ExtractGoldAmountFromText for GoldPrompt parsing

### 3. TalkCallbacks Interface System (2025)

**Purpose**: Provides dependency injection for talk dialog operations without circular dependencies.

**Architecture**:
```go
// Interface definition (in game_state package)
type TalkCallbacks struct {
    CreateTalkDialog func(npc *map_units.NPCFriendly) TalkDialog
    PushDialog      func(dialog TalkDialog)
}

// Implementation (in GameScene)
func (g *GameScene) CreateTalkDialog(npc *map_units.NPCFriendly) game_state.TalkDialog {
    return NewLinearTalkDialog(g, npc.NPCReference)
}

func (g *GameScene) PushDialog(dialog game_state.TalkDialog) {
    if linearDialog, ok := dialog.(*LinearTalkDialog); ok {
        g.dialogStack.PushModalDialog(linearDialog)
    }
}
```

**Integration with ActionTalkSmallMap**:
```go
func (g *GameState) ActionTalkSmallMap(direction references.Direction) bool {
    // Find NPC and validate...
    if friendly, ok := (*npc).(*map_units.NPCFriendly); ok {
        // Create and push dialog using dependency injection
        dialog := g.SystemCallbacks.Talk.CreateTalkDialog(friendly)
        if dialog != nil {
            g.SystemCallbacks.Talk.PushDialog(dialog)
            g.SystemCallbacks.Flow.AdvanceTime(1) // Time advances after dialog setup
            return true
        }
    }
    // Handle errors with SystemCallbacks.Message.AddRowStr()...
}
```

**Key Benefits**:
- **No Circular Dependencies**: GameState only knows about interfaces, not concrete GameScene
- **Proper Timing**: AdvanceTime() happens after dialog setup, not before
- **Clean Separation**: Dialog creation/pushing is UI concern, game logic stays in GameState
- **Testable**: Interface allows easy mocking for unit tests

### 4. Integration Points

**GameScene Integration**:
- `smallMapTalkSecondary()`: Modified to support both systems via debug flag
- `talkWithNpc()`: Debug command updated to support both systems  
- Debug toggle command: `linear-talk-toggle`

**Debug System Integration**:
- Added `UseLinearConversationSystem` flag to DebugOptions
- Runtime switching between conversation systems
- Debug talk command supports both systems

## Configuration and Usage

### Debug Commands

```bash
# Toggle between conversation systems
linear-talk-toggle

# Test specific NPCs (works with both systems)  
talk castle 1      # Talk to NPC 1 in castle
talk towne 5       # Talk to NPC 5 in town
talk dwelling 2    # Talk to NPC 2 in dwelling
```

### Runtime Behavior

**Channel System (Default)**:
- `UseLinearConversationSystem = false`
- Uses goroutines and channels
- Asynchronous processing
- Original conversation.go engine

**Linear System**:  
- `UseLinearConversationSystem = true`
- Direct method calls
- Synchronous processing
- LinearConversationEngine

## Implementation Status

### ✅ Fully Implemented

#### Core Integration
- [x] LinearTalkDialog UI component
- [x] GameActionCallbacks implementation
- [x] **TalkCallbacks dependency injection system (2025)**
- [x] **ActionTalkSmallMap refactored with SystemCallbacks (2025)**
- [x] **Circular dependency resolution (2025)**
- [x] Debug flag integration
- [x] Runtime system switching
- [x] Compilation and basic functionality

#### UI Components  
- [x] Identical visual layout to original
- [x] Text input handling
- [x] Output display and formatting
- [x] Border and styling consistency

#### Talk System Refactoring (2025)
- [x] **Removed TalkResult struct** - no longer needed with dependency injection
- [x] **Simplified ActionTalkSmallMap** - now returns bool instead of complex result
- [x] **Unified messaging** - all talk messages use SystemCallbacks.Message
- [x] **Proper time advancement** - AdvanceTime() called after dialog setup

#### Callback System
- [x] All ActionCallbacks interface methods
- [x] Game state query methods (HasMet, GetAvatarName, etc.)
- [x] Output accumulation and display
- [x] Error handling infrastructure

### 🚧 Partially Implemented

#### Game Actions (Placeholder Status)
- [⚠️] JoinParty: Basic logic implemented, needs testing
- [⚠️] CallGuards: Placeholder only
- [⚠️] GoToJail: Placeholder only  
- [⚠️] MakeHorse: Placeholder only
- [⚠️] PayExtortion/PayHalfExtortion: Basic logic implemented
- [⚠️] GiveItem: Placeholder only

#### Advanced Features
- [⚠️] Complex conversation flows
- [⚠️] Save game integration for HasMet
- [⚠️] Performance optimization

### ❌ Not Yet Implemented

- [ ] Full ActionCallback functionality (beyond placeholders)
- [ ] Integration testing with complex conversations
- [ ] Error recovery and edge case handling
- [ ] Migration of all TalkCommand types

## Testing Strategy

### Manual Testing Phases

1. **Phase 1: Basic Functionality**
   - Toggle between systems
   - Basic name/job/bye interactions
   - UI consistency verification

2. **Phase 2: Action Commands**  
   - Test placeholder actions
   - Verify no crashes with complex commands
   - Check output formatting

3. **Phase 3: Integration Testing**
   - Test with real NPCs from game data
   - Verify save game compatibility
   - Performance comparison

### Automated Testing

The linear conversation engine already has comprehensive tests in `internal/conversation/linear_engine_test.go`. Integration tests should focus on:

- GameActionCallbacks method coverage
- UI component interaction
- System switching reliability

## Migration Path

### Current State: Dual System Support
- Both systems coexist
- Runtime switching via debug flag
- No breaking changes to existing functionality

### Phase 1: Testing and Refinement (Current)
- Manual testing of basic functionality
- Bug fixes and UI refinement
- ActionCallback implementation completion

### Phase 2: Feature Parity  
- Full ActionCallback implementation
- Complex conversation flow support
- Performance optimization

### Phase 3: Migration
- Default to linear system
- Deprecation warnings for channel system
- Migration of existing integrations

### Phase 4: Cleanup
- Removal of channel-based system
- Code cleanup and documentation updates
- Final optimization

## Code Files Reference

### New Files Created
```
cmd/ultimav/linear_talk_dialog.go       # Linear UI component
cmd/ultimav/game_action_callbacks.go    # ActionCallbacks implementation  
docs/LINEAR_CONVERSATION_MIGRATION.md   # Migration guide
docs/CONVERSATION_SYSTEM_INTEGRATION.md # This document
```

### Modified Files
```
cmd/ultimav/gamescene_input_smallmap.go     # Added system toggle logic
cmd/ultimav/gamescene_debug_commands.go     # Added toggle command
internal/references/debug.go                # Added debug flag
```

### Existing Files (Preserved)
```
cmd/ultimav/talk_dialog.go                  # Original channel-based UI
internal/conversation/conversation.go       # Original channel engine
internal/conversation/linear_engine.go      # New linear engine (existing)
```

## Best Practices for Continued Development

### When Adding New Features

1. **UI Changes**: Modify both TalkDialog and LinearTalkDialog to maintain consistency
2. **New Actions**: Add to ActionCallbacks interface first, then implement in GameActionCallbacks
3. **Testing**: Test both conversation systems when making changes
4. **Documentation**: Update both system documentation

### When Fixing Bugs

1. **Root Cause**: Determine if bug affects one or both systems
2. **Consistent Fixes**: Apply equivalent fixes to both systems during dual-support phase  
3. **Regression Testing**: Test both conversation systems after fixes

### Performance Considerations

- Linear system should be more performant (no goroutines/channels)
- Monitor memory usage during conversations
- Profile conversation-heavy scenarios

## Troubleshooting Guide

### Common Issues

**"Method not found" errors**:
- Check ActionCallbacks interface implementation
- Verify all methods are implemented in GameActionCallbacks

**UI inconsistencies**:  
- Compare constants between TalkDialog and LinearTalkDialog
- Check UI element positioning and styling

**Conversation not starting**:
- Verify TalkScript loading for NPC
- Check debug flag state
- Monitor console for LinearConversationEngine errors

**System switching not working**:
- Check DebugOptions.UseLinearConversationSystem flag
- Verify debug command registration
- Test debug console functionality

### Debug Information

**Logging**:
- ActionCallback calls are logged to console
- Conversation errors include NPC position information  
- System switching is logged via debug console

**State Inspection**:
- Use debug talk command to test specific NPCs
- Check conversation engine active state
- Monitor UI state changes

This integration provides a solid foundation for migrating from the channel-based conversation system to the linear system while maintaining full backward compatibility during the transition period.