# Linear Conversation System Migration Guide

This document describes the migration from the channel-based conversation system to the new LinearConversationEngine, providing a roadmap for testing and integration.

## Migration Overview

The linear conversation system provides a synchronous, callback-based alternative to the asynchronous channel-based conversation system. Both systems can coexist during the migration period.

### Key Components Created

1. **LinearConversationEngine** (`internal/conversation/linear_engine.go`) - New synchronous conversation processor
2. **LinearTalkDialog** (`cmd/ultimav/linear_talk_dialog.go`) - UI component using linear engine
3. **GameActionCallbacks** (`cmd/ultimav/game_action_callbacks.go`) - Implementation of ActionCallbacks interface
4. **Debug Toggle** - Runtime switching between conversation systems

## Manual Testing Instructions

### Prerequisites
1. Build the application: `go build ./cmd/ultimav`
2. Ensure you have NPCs available for conversation testing
3. Have the debug console available for commands

### Basic Testing Steps

#### 1. Toggle Between Systems
Use the debug command to switch conversation systems:
```
linear-talk-toggle
```
- This toggles the `UseLinearConversationSystem` debug flag
- Default is `false` (channel-based system)
- Set to `true` to use the linear system

#### 2. Test Basic Conversation
1. Start game and approach an NPC
2. Press 'T' then arrow key to initiate talk
3. Verify conversation dialog appears
4. Test basic interactions:
   - Say "name" 
   - Say "job"
   - Say "bye"

#### 3. Test Debug Talk Command
Use the debug talk command to test specific NPCs:
```
talk castle 1
```
- Replace "castle" with location (castle, towne, dwelling, keep)
- Replace "1" with NPC dialog number
- This allows testing specific NPCs without navigation

#### 4. Compare Systems
1. Test conversation with channel system (`linear-talk-toggle` to disable)
2. Test same conversation with linear system (`linear-talk-toggle` to enable)
3. Compare behavior, output, and responsiveness

### Expected Behaviors

#### Linear System Features
- **Synchronous Processing**: No goroutines, immediate response processing
- **Callback Integration**: Game actions handled through GameActionCallbacks
- **Same UI**: Identical visual appearance to channel system
- **Debug Output**: ActionCallback placeholders show "[ACTION - Not yet implemented]"

#### Known Limitations (Current Implementation)
- **Stub Callbacks**: Most game actions show placeholder text
- **Basic Commands Only**: Limited TalkCommand support compared to full linear engine
- **No Save Integration**: HasMet status may not persist properly

### Testing Scenarios

#### Scenario 1: Basic NPC Interaction
1. Toggle to linear system
2. Talk to a friendly NPC
3. Verify:
   - Description appears ("You see...")
   - Name/Job/Bye responses work
   - Dialog closes properly on "bye"

#### Scenario 2: Action Commands (Placeholders)
1. Find NPCs with special actions (join party, karma changes, etc.)
2. Trigger action commands
3. Verify:
   - Placeholder text appears
   - No crashes occur
   - Dialog continues properly

#### Scenario 3: System Switching
1. Start conversation with channel system
2. Close dialog
3. Toggle to linear system
4. Start same conversation
5. Verify:
   - Different system is used
   - Behavior is similar
   - No crashes during switching

## Troubleshooting

### Common Issues

#### Compilation Errors
- Ensure all imports are correct
- Check that ActionCallbacks interface matches implementation
- Verify NPCReference field usage (no Name field)

#### Runtime Errors
- Check debug flag state with debug commands
- Verify TalkScript data is loaded properly
- Monitor console logs for error messages

#### UI Issues  
- Ensure LinearTalkDialog inherits proper UI constants
- Check that text output is displaying correctly
- Verify input handling works as expected

### Debugging Commands

```bash
# Toggle conversation systems
linear-talk-toggle

# Test specific NPC (examples)
talk castle 1          # Test NPC in castle
talk towne 5           # Test NPC in town  
talk dwelling 2        # Test NPC in dwelling

# Monitor debug output
# Check console logs for ActionCallback calls
# Look for placeholder text in dialog
```

## Implementation Status

### ✅ Completed
- LinearTalkDialog UI component
- GameActionCallbacks implementation  
- Debug toggle integration
- Basic conversation flow
- Compilation and integration

### 🚧 In Progress  
- Manual testing and iteration
- ActionCallback implementation refinement
- Error handling improvements

### ❌ Not Yet Implemented
- Full ActionCallback functionality (join party, guards, etc.)
- Complex conversation flow testing
- Save game integration for HasMet status
- Performance optimization

## Next Steps

1. **Manual Testing Phase**: Complete testing scenarios outlined above
2. **Iterate on Issues**: Fix bugs found during testing
3. **Enhance Callbacks**: Implement full ActionCallback functionality  
4. **Documentation**: Document final solution
5. **Migration**: Plan removal of channel-based system
6. **Cleanup**: Remove obsolete code

## Migration Path Forward

### Short Term (Testing Phase)
- Focus on verifying basic conversation functionality
- Identify and fix critical issues
- Ensure UI consistency between systems

### Medium Term (Feature Completion)  
- Implement full ActionCallback functionality
- Add proper error handling and edge cases
- Integrate with save game system

### Long Term (Migration Completion)
- Default to linear system
- Remove channel-based system 
- Clean up obsolete code
- Update documentation

This migration approach allows for safe, incremental adoption of the linear conversation system while maintaining the existing functionality during the transition period.