# Input Handling System

This document describes the input handling architecture for player commands and actions in the Ultima 5 Redux Go project.

## Architecture Overview

The input system uses a two-phase approach:
1. **Primary Input** - Initial key press detection and command identification
2. **Secondary Input** - Directional input for commands that require targeting

This design mirrors the original Ultima V command structure where players press a command key (like "L" for Look) followed by a direction arrow.

## Input State Management

### InputState Enum
```go
type InputState int

const (
    PrimaryInput InputState = iota
    OpenDirectionInput
    JimmyDoorDirectionInput
    KlimbDirectionInput
    PushDirectionInput
    GetDirectionInput
    LookDirectionInput
    TalkDirectionInput
    SearchDirectionInput
    AttackDirectionInput
    UseDirectionInput
    YellDirectionInput
)
```

### State Flow
```
[Key Press] → [Primary Input Handler] → [Secondary Input State]
                                      ↓
[Direction Key] → [Secondary Handler] → [Action Execution] → [Primary Input]
```

## File Organization

### Input Handlers
- **`cmd/ultimav/gamescene.go`** - InputState definitions and GameScene struct
- **`cmd/ultimav/gamescene_input_smallmap.go`** - Small map input handlers
- **`cmd/ultimav/gamescene_input_largemap.go`** - Large map input handlers  
- **`cmd/ultimav/gamescene_input_common.go`** - Shared input handlers

### Action Methods
- **`internal/game_state/action_*.go`** - GameState action method implementations
- Follow pattern: `Action[Command][SmallMap|LargeMap](direction references.Direction) bool`

## Primary Input Handling

Primary input handlers detect the initial command key and transition to the appropriate secondary input state.

### Example - Look Command
```go
case ebiten.KeyL:
    g.addRowStr("Look-")
    g.secondaryKeyState = LookDirectionInput
```

### Command Key Mappings
| Key | Command | Secondary State | Description |
|-----|---------|----------------|-------------|
| L   | Look    | LookDirectionInput | Examine adjacent tile |
| G   | Get     | GetDirectionInput | Pick up items |
| P   | Push    | PushDirectionInput | Move objects |
| O   | Open    | OpenDirectionInput | Open doors/chests |
| J   | Jimmy   | JimmyDoorDirectionInput | Pick locks |
| K   | Klimb   | KlimbDirectionInput | Use ladders/climb |
| T   | Talk    | TalkDirectionInput | Converse with NPCs |
| S   | Search  | SearchDirectionInput | Find hidden items |
| A   | Attack  | AttackDirectionInput | Combat actions |
| U   | Use     | UseDirectionInput | Use items |
| Y   | Yell    | YellDirectionInput | Shout commands |

### Non-Directional Commands
Some commands execute immediately without secondary input:
- **I** - Ignite Torch (`ActionIgnite()`)
- **B** - Board Vehicle  
- **X** - Exit/Leave
- **E** - Enter Building

## Secondary Input Handling

Secondary input handlers wait for directional input and execute the corresponding action.

### Direction Detection
```go
if g.isDirectionKeyValidAndOutput() {
    direction := getCurrentPressedArrowKeyAsDirection()
    g.smallMapActionSecondary(direction)
    g.secondaryKeyState = PrimaryInput
}
```

### Secondary Handler Pattern
```go
func (g *GameScene) smallMap[Command]Secondary(direction references.Direction) {
    success := g.gameState.Action[Command]SmallMap(direction)
    if !success {
        g.addRowStr("Appropriate error message!")
    }
}
```

## Action Method Patterns

### Directional Actions
All directional actions follow this signature:
```go
func (g *GameState) Action[Command]SmallMap(direction references.Direction) bool
func (g *GameState) Action[Command]LargeMap(direction references.Direction) bool
```

#### Examples:
- `ActionLookSmallMap(direction references.Direction) bool`
- `ActionGetLargeMap(direction references.Direction) bool`
- `ActionPushSmallMap(direction references.Direction) bool`

### Non-Directional Actions
Non-directional actions omit the direction parameter:
```go
func (g *GameState) ActionIgnite() bool
func (g *GameState) ActionEnter(location *SmallLocationReference) bool
```

### Return Values
- **`true`** - Action succeeded
- **`false`** - Action failed (invalid target, insufficient resources, etc.)

## Map-Specific Behavior

### Small Map vs Large Map
- **Small Map** - Towns, buildings, dungeons - detailed tile-by-tile interaction
- **Large Map** - Overworld - broader movement and interaction

### Context Validation
Actions may behave differently or be disabled based on:
- Current map type (small vs large)
- Player location (town, dungeon, overworld)
- Player state (on foot, mounted, in vehicle)
- Environmental conditions (light level, terrain)

## Error Handling

### Input Validation
- Invalid directions are rejected with appropriate messages
- Commands unavailable in current context show "Not here!" or similar
- Resource requirements (keys, torches, etc.) are validated before execution

### User Feedback
- Primary input shows command prompt: "Look-", "Get-", etc.
- Secondary input success/failure provides specific messages
- Failed actions return to PrimaryInput state immediately

## Integration Points

### Game State
- Action methods in `GameState` handle core game logic
- UI handlers in `GameScene` manage display and user feedback
- Clear separation between input handling and game logic

### Turn Management  
- Most actions call `g.gameState.FinishTurn()` after execution
- Some actions (failed attempts, cancelled inputs) don't advance turn
- Combat actions may have special turn handling

## Future Extensions

### Adding New Commands
1. Add InputState constant for secondary input (if directional)
2. Add key mapping in primary input handler
3. Add secondary input handler (if directional)
4. Create `Action[Command]SmallMap/LargeMap` methods in GameState
5. Update this documentation

### Command Categories
Consider organizing commands by category:
- **Movement** - Klimb, Push
- **Interaction** - Look, Get, Open, Talk
- **Combat** - Attack, Use (weapons)
- **Utility** - Search, Ignite, Yell

## Legacy Compatibility

### Deprecated Methods
Some methods maintain backwards compatibility:
- `IgniteTorch()` → `ActionIgnite()`
- `EnterBuilding()` → `ActionEnter()`
- `JimmyDoor()` → `ActionJimmySmallMap()`

These will be removed once all callers are updated to use the standardized Action patterns.