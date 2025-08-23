# Debug System Documentation

The Ultima 5 Redux Go project includes a comprehensive in-game debug console system designed to facilitate development and testing. This system provides real-time access to game state manipulation, testing commands, and development utilities without requiring code recompilation.

## System Architecture

### Core Components

1. **DebugConsole** (`cmd/ultimav/debug_console.go`)
   - Main UI component providing the interactive console interface
   - Text-based input with command parsing and autocomplete
   - Scrollable output display showing command results
   - Modal dialog overlay that doesn't interrupt game state

2. **DebugConsoleActions** (`cmd/ultimav/debug_console_actions.go`)
   - Contains all debug command implementations
   - Uses grammar-based command parsing for flexible input handling
   - Organized into logical command categories (teleport, time, NPCs, etc.)

3. **DebugOptions** (`internal/references/debug.go`)
   - Runtime configuration flags for debug behaviors
   - Persistent settings that affect gameplay during debugging
   - Currently includes: FreeMove, MonsterGen, ExperimentalConversationFeatures

4. **Debug Action Methods** (`cmd/ultimav/debug_console_actions.go`)
   - Low-level debug operations for position changes, floor navigation
   - Direct game state manipulation methods

## Activation and Usage

### Opening the Debug Console
- **Activation Key**: Backtick (`) key
- **Context**: Available in both small maps (towns, castles, etc.) and large map (overworld)
- **Toggle Behavior**: Pressing backtick again closes the console
- **Exit Methods**: Backtick (`) or Escape key

### Interface Layout
```
┌─────────────────────────────────────────┐
│  Debug Console Output (14 lines max)   │
│  > Command executed successfully        │
│  X=125, Y=87, Floor=0                   │
│                                         │
│  [Previous commands and results...]     │
│                                         │
│                                         │
├─────────────────────────────────────────┤
│  > [command input field]                │
└─────────────────────────────────────────┘
```

### Command Input Features
- **Autocomplete**: Tab completion for command names and parameters
- **Command History**: Up/Down arrows to navigate previous commands
- **Ambiguous Completion**: Shows available options when multiple matches exist
- **Parameter Validation**: Real-time validation of command parameters
- **Case Insensitive**: Commands work regardless of case

## Available Commands

### Position and Navigation

#### `teleport <X> <Y>`
Move player to specific coordinates on current map.
- **Parameters**: X (0-255), Y (0-255)
- **Example**: `teleport 100 50`
- **Output**: Shows new position coordinates

#### `fy <floor>`
Teleport to specific floor number.
- **Parameters**: Floor (-1 to 5, depending on location)
- **Example**: `fy 2`
- **Special**: Underworld uses floors -1 (surface) and 0 (underground)

#### `fu` / `fd`
Move up or down one floor if available.
- **fu**: Floor up
- **fd**: Floor down
- **Output**: Success/failure status

#### `gos <location>`
Instantly travel to any small map location.
- **Parameters**: Location name (autocompletes)
- **Example**: `gos britain`, `gos castle`, `gos dwelling`
- **Effect**: Immediately enters the specified building/location

### Game State Modification

#### `freemove`
Toggle movement boundary checking.
- **Effect**: When enabled, allows moving through walls and obstacles
- **Status**: Shows current FreeMove setting
- **Use Case**: Exploring inaccessible areas, testing map boundaries

### Time and Environment

#### `tsh <hour>`
Set game time to specific hour.
- **Parameters**: Hour (0-23)
- **Example**: `tsh 14` (sets to 2 PM)
- **Effect**: Immediately changes game time, affecting NPC schedules

#### `qt <time>`
Quick time setting using descriptive names.
- **Parameters**: morning, evening, midnight, noon, dusk
- **Example**: `qt midnight`
- **Effect**: Sets time to predetermined hours for each period

### Monster and Combat System

#### `mon-toggle-gen`
Toggle random monster generation.
- **Effect**: Enables/disables random encounters and monster spawning
- **Status**: Shows current MonsterGen setting
- **Use Case**: Testing without combat interruptions

#### `mon-change-odds <odds>`
Modify monster generation probability.
- **Parameters**: One-in-X odds (0-1000)
- **Example**: `mon-change-odds 100` (1 in 100 chance)
- **Default**: Usually 1 in 64 or similar
- **Use Case**: Increase encounters for combat testing or reduce for peaceful exploration

#### `mon-delall`
Remove all monsters from current map.
- **Effect**: Immediately clears all enemy NPCs
- **Use Case**: Clear hostiles for safe exploration

### NPC and Conversation Testing

#### `talk <location> <npc_id>`
Initiate conversation with specific NPC.
- **Parameters**: Location name, NPC dialog number (0-1000)
- **Example**: `talk castle 1`, `talk britain 5`
- **System**: Uses LinearConversationEngine
- **Use Case**: Test specific NPC conversations without navigating to them

### Vehicle and Transportation

#### `buyboat <type> <location>`
Spawn boats at dock locations.
- **Parameters**: 
  - Type: frigate, skiff
  - Location: Any location with docks, or "avatar" for current position
- **Example**: `buyboat frigate britain`, `buyboat skiff avatar`
- **Effect**: Creates fully functional vehicle with 1 skiff capacity

### Display and Technical

#### `ru` / `rd`
Adjust game resolution.
- **ru**: Increase resolution
- **rd**: Decrease resolution
- **Effect**: Real-time resolution change with UI scaling

#### `fs`
Toggle fullscreen mode.
- **Effect**: Switches between windowed and fullscreen display

## Command Grammar System

The debug console uses a sophisticated grammar-based parsing system that provides:

### Parameter Types
- **MatchString**: Exact string matches (case insensitive)
- **MatchInt**: Integer parameters with min/max validation
- **MatchStringList**: Selection from predefined lists with autocomplete
- **SingleCharacterInput**: Commands that accept single character responses

### Validation Features
- **Range Checking**: Automatic validation of numeric parameters
- **List Validation**: Parameters must match predefined options
- **Autocomplete Support**: Tab completion for all parameter types
- **Error Feedback**: Clear error messages for invalid input

### Example Command Definition
```go
grammar.NewTextCommand([]grammar.Match{
    grammar.MatchString{
        Str: "teleport",
        Description: "Move to an X, Y coordinate",
        CaseSensitive: false,
    },
    grammar.MatchInt{IntMin: 0, IntMax: 255, Description: "X coordinate"},
    grammar.MatchInt{IntMin: 0, IntMax: 255, Description: "Y coordinate"},
}, handlerFunction)
```

## Integration with Game Systems

### Game State Integration
- Direct manipulation of `GameScene.gameState`
- Real-time updates to player position, time, and flags
- Immediate effect on rendering and game logic

### NPC System Integration
- Access to NPCReferences for conversation testing
- Integration with LinearConversationEngine
- Monster/enemy management through NPCAIController

### Map System Integration
- Floor navigation respects map constraints
- Location teleportation validates available maps
- Coordinate validation within map boundaries

### Debug Options Integration
Debug flags affect real gameplay behavior:
- **FreeMove**: Bypasses collision detection
- **MonsterGen**: Controls random encounter generation
- **ExperimentalConversationFeatures**: Reserved for testing new conversation features

## Development and Testing Benefits

### Rapid Iteration
- Test new features without lengthy setup
- Quickly reproduce specific game states
- Skip time-consuming navigation for testing

### State Exploration
- Access any location instantly
- Test time-sensitive features (NPC schedules, lighting)
- Examine different floor configurations

### Content Validation
- Test all NPC conversations systematically
- Verify monster generation algorithms
- Validate map boundaries and transitions

### Bug Investigation
- Reproduce specific coordinate-based issues
- Test edge cases in floor transitions
- Isolate problems by controlling variables

## Technical Implementation Details

### Command Processing Flow
1. User types command in debug console input field
2. Grammar system parses input and validates parameters
3. Matching command handler function executes
4. Results displayed in console output area
5. Game state updated immediately (if applicable)

### Memory Management
- Debug console uses persistent text output buffer
- Limited to 14 lines of output history
- Automatic text wrapping for long output lines
- No persistent storage - resets when console closes

### Performance Considerations
- Debug commands execute immediately (no queuing)
- Minimal performance impact when console is closed
- UI rendering only when console is visible
- Grammar parsing optimized for real-time input

### Error Handling
- Parameter validation prevents invalid game states
- Graceful handling of edge cases (invalid floors, coordinates)
- Clear error messages for debugging command issues
- Fallback behaviors for system failures

## Best Practices for Use

### Testing Workflows
1. **Feature Testing**: Use `gos` to quickly access test locations
2. **Time Testing**: Use `qt` for quick time changes, `tsh` for precision
3. **Combat Testing**: Use `mon-change-odds` to control encounter frequency
4. **Conversation Testing**: Use `talk` to systematically test all NPCs

### Safety Considerations
- Debug commands can put game in invalid states
- Save game before extensive debug testing
- Some commands may affect save file compatibility
- FreeMove can lead to getting stuck in inaccessible areas

### Development Integration
- Add new debug commands for new game features
- Use debug flags to test experimental implementations
- Leverage command grammar system for consistent UX
- Document new commands in this file

This debug system significantly accelerates development and testing cycles by providing immediate access to game state manipulation and testing utilities, making it an essential tool for efficient Ultima 5 Redux development.