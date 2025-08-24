# Coding Conventions

This document defines the coding standards for the project. It balances idiomatic Go, public Google Go style guidance, and the project’s practical needs (classic game fidelity, clarity, and incremental development).

Goals:
- Keep core gameplay logic clear and simple.
- Favor readability over cleverness.
- Prefer conventions that make the code testable, deterministic, and reusable.

## Language and Tooling

- Go version: the repository’s go.mod defines the version; keep code compatible with it.
- Formatting: gofmt (via goimports). CI should fail if formatting is not applied.
- Linting: follow common static checks (e.g., golangci-lint defaults) where reasonable.
- Build tags: use sparingly (e.g., debug vs release).

## Package Structure and Boundaries

- Core logic should be independent of the rendering/input backend. Adapter layers (e.g., Ebiten) call into core, not vice versa.
- Keep packages focused and small; avoid cyclic dependencies.
- Exported types and functions should be stable and usable by external consumers where reasonable.

## Naming

- Exported identifiers: PascalCase. Unexported: camelCase.
- Initialisms: follow Go’s convention (ID, URL, HTTP). Example: UserID, not UserId.
- File names: lowercase with underscores only when it improves clarity; prefer consistent grouping by feature/domain.

## Imports and Dependencies

- Standard library first, external second, internal last; keep groups separated by a blank line.
- Do not import the adapter layer (e.g., rendering/input libraries) into core logic. Inject via interfaces if needed.
- Keep external dependencies minimal and pinned via go.mod/go.sum.

## Errors and Logging

- Return errors instead of panicking in normal control flow.
- Wrap underlying errors with context: `fmt.Errorf("loading sprite: %w", err)`.
- Temporary hard exits:
  - During early development, a critical path may use `log.Fatal` as a stub for “must not happen” states to catch regressions quickly.
  - Always add a comment `// TODO: soften to recoverable error` so we can convert to softer handling later.
- Logging should be concise and actionable. Avoid spamming logs in hot paths; prefer rate-limiting if necessary.

## Time, Randomness, and Ticks

- Do not call `time.Now()` inside core logic or animation helpers that affect gameplay determinism.
- Use a central game clock (ticks/elapsedMs) to drive:
  - Animations (waterfalls, clocks, idle cycles).
  - NPC schedules and time-of-day effects.
  - AI cadence and spawn checks.
- Randomness:
  - Centralize PRNG state and seed per session.
  - For deterministic tests, allow fixed seeds.

## Loops and Collections

- Prefer `range` loops for readability and safety when index is not needed:
  - Good: `for _, v := range items { ... }`
  - Use index-based loops when:
    - You need the index for math or spatial mapping.
    - You need to mutate by index efficiently.
    - You’re iterating over a fixed-size array where index is semantic.
- Preallocate slices with capacity when the size is known (`make(T, 0, n)` or `make(T, n)`).
- Avoid unnecessary allocations in hot paths; reuse buffers where appropriate.

## Data-Driven vs Hardcoded Logic

- Favor data-driven definitions for tile attributes, schedules, and references.
- It is acceptable to hardcode small, well-named rules when it materially improves clarity (e.g., special-case mirror behavior, chair/ladder substitutions).
- When hardcoding, isolate rules behind small helpers so they can be replaced by data later if needed.

## Rendering and Animation

- All animations should read from the tick/elapsed time (not `time.Now`) to keep visuals in sync with gameplay.
- Prefer simple, deterministic animation helpers that accept `(spriteIndex, positionHash, elapsedMs)`.
- Keep per-frame allocations to a minimum (reuse images/options where possible).

## AI and Pathfinding

- Cache computed paths (A*) on the unit and consume steps over subsequent ticks.
- Recompute on cooldowns or when blocked; fall back to simple greedy movement if needed.
- Respect terrain passability by agent type (avatar, vehicle, land/water enemy).
- Do not hard-exit on missing paths; use a clear fallback and log once with a TODO to soften behavior later.

## Concurrency

- Prefer single-threaded game-state mutation during the update step.
- If background loading is needed, communicate via channels or synchronized buffers, and apply results during the main update.

## Comments and Documentation

- Use doc comments for exported identifiers; explain the why, not the obvious what.
- Keep comments up-to-date with behavior; incorrect comments are worse than none.
- Add package-level docs for domains with non-trivial rules (e.g., AI, schedule resolution).

## Testing

- Unit tests for deterministic pieces: animation frame selection, schedule resolution, pathfinding.
- Use fixed seeds and a mock clock for deterministic tests.
- Avoid brittle pixel tests; prefer snapshot structures that capture tile IDs, frame indices, and positions.

## Performance and Allocation

- Avoid per-frame heap churn in tight loops.
- Use small structs and pass pointers only when mutation is required or copying is expensive.
- Profile with pprof before optimizing; optimize only bottlenecks.

## Code Review Checklist

- Readability: Is the code straightforward and consistent with these conventions?
- Safety: Are errors handled appropriately? Are temporary hard exits marked with TODO to soften?
- Determinism: Are time and randomness driven by the central clock/PRNG?
- Style: Are loops using `range` where reasonable? Are names idiomatic?
- Boundaries: Does core avoid importing renderer/input packages?
- Tests: Are new behaviors covered by tests or at least easy to test?
- Docs: Are exported symbols documented and comments accurate?

## Command Input Patterns

Player commands (Look, Push, Get, Klimb, Open, Jimmy, etc.) should follow a consistent separation of concerns between UI layer and GameState logic.

### UI Layer Command Methods
```go
// Small map UI commands - handle input, validation, call core logic
func (g *GameScene) smallMap[Command]Secondary(direction references.Direction)
func (g *GameScene) smallMap[Command]()

// Large map UI commands
func (g *GameScene) largeMap[Command]Secondary(direction references.Direction) 
func (g *GameScene) largeMap[Command]Primary()

// Combat map UI commands
func (g *GameScene) combatMap[Command]Secondary(direction references.Direction)
```

### GameState Action Methods
```go
// GameState logic methods - contain core game logic, use injected callbacks
func (g *GameState) Action[Command]SmallMap(direction references.Direction) bool
func (g *GameState) Action[Command]LargeMap(direction references.Direction) bool
func (g *GameState) Action[Command]CombatMap(direction references.Direction) bool
```

### Separation of Concerns

**UI Layer Responsibilities** (`GameScene` methods):
- Input handling and direction gathering
- Early validation that can be done without game logic
- Calling appropriate GameState action methods
- Handling cases where GameState methods aren't needed

**GameState Responsibilities** (`Action*` methods):
- Core game logic and validation
- State modifications (player position, inventory, doors, etc.)
- Time advancement via injected callbacks
- User feedback via injected callbacks
- Sound effects via injected callbacks
- Return success/failure status

### SystemCallbacks - Injected Function Usage

GameState uses dependency injection to communicate with the outer system without knowing implementation details. All UI interactions should go through `SystemCallbacks`:

#### Message System
```go
// Basic user message
g.SystemCallbacks.Message.AddRowStr("Won't budge!")

// Multi-line messages
g.SystemCallbacks.Message.AddRowStr("Found:")
g.SystemCallbacks.Message.AddRowStr("Gold coins, Magic sword")

// Critical system messages
g.SystemCallbacks.Message.AddRowStr("The door slams shut!")
```

#### Audio System
```go
// Action-specific sound effects
g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
g.SystemCallbacks.Audio.PlaySoundEffect(SoundOpenDoor)
g.SystemCallbacks.Audio.PlaySoundEffect(SoundKeyBreak)

// Combat sounds (when implemented)
g.SystemCallbacks.Audio.PlaySoundEffect(SoundHit)
g.SystemCallbacks.Audio.PlaySoundEffect(SoundMiss)
```

#### Time and Flow Control
```go
// Advance game time (in minutes)
g.SystemCallbacks.Flow.AdvanceTime(1)   // Most actions take 1 minute
g.SystemCallbacks.Flow.AdvanceTime(5)   // Longer actions (searching, etc.)

// Special time events
g.SystemCallbacks.Flow.AdvanceTime(15)  // Complex actions
```

#### Screen and Display Updates
```go
// Trigger screen refresh when needed
g.SystemCallbacks.Screen.UpdateDisplay()

// Map updates are usually automatic, but can be forced
g.SystemCallbacks.Screen.RefreshMap()
```

#### Callback Usage Patterns

**Pattern 1: Simple Success/Failure**
```go
func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
    // Validation logic...
    if validationFails {
        g.SystemCallbacks.Message.AddRowStr("Won't budge!")
        return false
    }
    
    // Core logic...
    if success {
        g.SystemCallbacks.Message.AddRowStr("Pushed!")
        g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
        g.SystemCallbacks.Flow.AdvanceTime(1)
        return true
    }
    
    g.SystemCallbacks.Message.AddRowStr("Won't budge!")
    return false
}
```

**Pattern 2: Multiple Outcome States**
```go
func (g *GameState) ActionGetSmallMap(direction references.Direction) bool {
    // Logic determines outcome...
    switch outcome {
    case GetSuccess:
        g.SystemCallbacks.Message.AddRowStr("Taken!")
        g.SystemCallbacks.Audio.PlaySoundEffect(SoundGetItem)
        g.SystemCallbacks.Flow.AdvanceTime(1)
        return true
    case GetTooHeavy:
        g.SystemCallbacks.Message.AddRowStr("Too heavy!")
        return false
    case GetNothing:
        g.SystemCallbacks.Message.AddRowStr("Nothing there!")
        return false
    }
}
```

**Pattern 3: Complex State Changes**
```go
func (g *GameState) ActionOpenSmallMap(direction references.Direction) bool {
    switch doorResult := g.MapState.OpenDoor(direction); doorResult {
    case map_state.OpenDoorOpened:
        g.SystemCallbacks.Message.AddRowStr("Opened!")
        g.SystemCallbacks.Audio.PlaySoundEffect(SoundOpenDoor)
        g.SystemCallbacks.Flow.AdvanceTime(1)
        return true
    case map_state.OpenDoorLocked:
        g.SystemCallbacks.Message.AddRowStr("Locked!")
        return false
    case map_state.OpenDoorLockedMagical:
        g.SystemCallbacks.Message.AddRowStr("Magically Locked!")
        return false
    default:
        g.SystemCallbacks.Message.AddRowStr("Bang to open!")
        return false
    }
}
```

### Command Flow Pattern

1. **Input Gathering**: UI captures command key and direction
2. **Early Validation**: UI performs simple checks (if any)
3. **Core Logic**: GameState `Action*` method performs game logic
4. **Callback Execution**: GameState uses SystemCallbacks for user feedback
5. **Result Handling**: UI interprets return value if additional logic needed

### When to Use Direct UI vs SystemCallbacks

**Use Direct UI** (`g.output.AddRowStr()`) when:
- The UI layer handles the logic entirely (early validation)
- GameState method isn't called
- Simple error cases that don't involve game state changes

**Use SystemCallbacks** (`g.SystemCallbacks.Message.AddRowStr()`) when:
- Inside GameState Action methods
- Game state is being modified
- Core game logic determines the message
- Sound effects or time advancement are involved

### Example: Complete Command Implementation

```go
// UI Layer - minimal, delegates to GameState
func (g *GameScene) smallMapPushSecondary(direction references.Direction) {
    pushThingPos := direction.GetNewPositionInDirection(&g.gameState.MapState.PlayerLocation.Position)
    pushThingTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(pushThingPos)

    // Early validation - avoid GameState call if obviously invalid
    if !g.gameState.IsPushable(pushThingTile) {
        g.output.AddRowStrWithTrim("Won't budge!") // Direct UI - no game logic
        return
    }

    // Delegate everything else to GameState
    g.gameState.ActionPushSmallMap(direction) // GameState handles all feedback
}

// GameState - complete logic with proper callback usage
func (g *GameState) ActionPushSmallMap(direction references.Direction) bool {
    // Auto-shut doors if needed (game logic)
    g.MapState.SmallMapProcessTurnDoors()

    pushableThingPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
    smallMap := g.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.MapState.PlayerLocation.Floor)
    pushableThingTile := smallMap.GetTopTile(pushableThingPosition)

    // Detailed validation with game state
    if g.ObjectPresentAt(pushableThingPosition) || !g.IsPushable(pushableThingTile) {
        g.SystemCallbacks.Message.AddRowStr("Won't budge!") // SystemCallback - in game logic
        return false
    }

    // Determine legal floor based on object type
    legalFloorIndex := indexes.SpriteIndex(indexes.BrickFloor)
    if pushableThingTile.IsCannon() {
        legalFloorIndex = indexes.SpriteIndex(indexes.HexMetalGridFloor)
    }

    farSidePosition := direction.GetNewPositionInDirection(pushableThingPosition)
    farSideTile := smallMap.GetTopTile(farSidePosition)
    playerTile := smallMap.GetTopTile(&g.MapState.PlayerLocation.Position)

    // Try to push object forward
    if !g.IsOutOfBounds(*farSidePosition) && 
       !g.ObjectPresentAt(farSidePosition) && 
       farSideTile.Index == legalFloorIndex {
        
        // Execute push
        g.pushIt(smallMap, pushableThingTile, farSideTile, pushableThingPosition, farSidePosition, direction)
        g.MapState.PlayerLocation.Position = *pushableThingPosition
        
        // User feedback via callbacks
        g.SystemCallbacks.Message.AddRowStr("Pushed!")
        g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
        g.SystemCallbacks.Flow.AdvanceTime(1)
        return true
        
    } else if playerTile.Index == legalFloorIndex {
        // Try swapping (pulling)
        g.swapIt(smallMap, playerTile, pushableThingTile, &g.MapState.PlayerLocation.Position, pushableThingPosition, direction)
        g.MapState.PlayerLocation.Position = *pushableThingPosition
        
        // User feedback via callbacks
        g.SystemCallbacks.Message.AddRowStr("Pulled!")
        g.SystemCallbacks.Audio.PlaySoundEffect(SoundPushObject)
        g.SystemCallbacks.Flow.AdvanceTime(1)
        return true
    }

    // Failed to push or pull
    g.SystemCallbacks.Message.AddRowStr("Won't budge!")
    return false
}
```

### GameState Helper Methods

GameState can provide helper methods that don't follow the Action pattern:
- `IsPushable(tile)` - validation helpers (no callbacks)
- `SelectCharacterForJimmy()` - selection logic (no callbacks)
- `JimmyDoor(direction, character)` - detailed operations returning enums (no callbacks)
- `OpenDoor(direction)` - specific operations with detailed results (may use callbacks)

### SystemCallbacks Interface Guidelines

- **Always use callbacks** for messages that result from game logic changes
- **Consistent naming**: Use present tense for actions ("Pushed!", "Opened!", "Taken!")
- **Audio timing**: Play sounds immediately when action succeeds
- **Time advancement**: Call at the end of successful actions
- **Error consistency**: Use similar phrasing for similar error types across commands

### Migration Guidelines

When refactoring existing code:
1. **Keep** `Action*Map` methods in GameState for core logic
2. **Use SystemCallbacks** for all messages, sounds, and time advancement in GameState
3. **Move early validation** to GameScene methods when appropriate
4. **Use return values** for GameScene to interpret when additional logic needed
5. **Update tests** to work with separated concerns

### Conventions:
- **Naming**: Always prefix with `Action`, use PascalCase command name, suffix with map type
- **Parameters**: Always include `direction references.Direction` parameter
- **Returns**: Return `bool` indicating success/failure
- **Imports**: Use direct import `"github.com/bradhannah/Ultima5ReduxGo/internal/references"`, avoid aliases
- **File organization**: Group related actions in `action_[command].go` files

### Examples:
- `ActionGetSmallMap(direction references.Direction) bool`
- `ActionLookLargeMap(direction references.Direction) bool`
- `ActionPushSmallMap(direction references.Direction) bool`

## Examples

Range vs index-based:
