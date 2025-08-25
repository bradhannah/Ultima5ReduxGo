# Error Handling Guidelines

This document defines error handling standards for the Ultima 5 Redux Go codebase, specifically addressing when `log.Fatal()` should be used versus recoverable error handling.

## Core Principle

**Fatal errors should represent programmer errors or catastrophic system failures that indicate the program cannot continue safely. User-facing errors, transient failures, and recoverable conditions should use proper error returns.**

## ✅ KEEP AS FATAL - Legitimate Uses

### 1. Critical Game Data Loading (Embedded/Internal Files)
**Rule**: Essential game files that are shipped with the binary should always work.

```go
// KEEP FATAL: Internal embedded data loading
data, err := gameFiles.ReadFile("data/tiles.json") 
if err != nil {
    log.Fatal(err) // ✅ This should never fail - shipped with binary
}
```

**Examples**:
- `/internal/sprites/sprite_utils.go:14` - Loading embedded sprite data
- `/internal/text/fonts.go:28` - Loading embedded font data  
- `/internal/references/*.go` - Loading embedded JSON data files

### 2. Original Game Data Parsing (User Files)
**Rule**: User-provided Ultima V data files that are malformed should be fatal - indicates data corruption.

```go
// KEEP FATAL: Original game data that should be well-formed
if len(rawData) < expectedSize {
    log.Fatal("Not enough data to create NPCReference schedule") // ✅ Data corruption
}
```

**Examples**:
- `/internal/references/npc_schedule.go:32` - Invalid schedule data structure
- `/internal/config/config.go:85` - Cannot read DATA.OVL (core Ultima V file)
- `/internal/references/data_ovl.go:89,98` - Corrupted compressed data

### 3. Programmer Logic Errors (Contract Violations)
**Rule**: Programming mistakes that violate documented contracts should be fatal during development.

```go
// KEEP FATAL: Logic errors that should never happen
func (v *VehicleDetails) DecrementSkiffQuantity() {
    if v.skiffQuantity == 0 {
        log.Fatal("skiff quantity is 0, cannot decrement it.") // ✅ Programming bug
    }
}
```

**Examples**:
- `/internal/map_units/npc_vehicle.go:70` - Decrementing zero quantity (caller bug)
- `/internal/game_state/game_state.go:293` - Invalid range parameters (caller bug)
- `/internal/references/item_stack.go:88,98` - Pop/peek from empty stack (caller bug)

### 4. Critical System State Violations
**Rule**: System state that would make the game unplayable or unsafe should be fatal.

```go
// KEEP FATAL: System invariant violations  
func (g *GameState) getCurrentLargeMapNPCAIController() *ai.NPCAIControllerLargeMap {
    if g.MapState.PlayerLocation.Location.GetMapType() != references.LargeMapType {
        log.Fatalf("Expected large map type, got %d", g.MapState.PlayerLocation.Location.GetMapType()) // ✅ Game state corruption
    }
}
```

**Examples**:
- `/internal/game_state/game_state.go:131` - Map bounds overflow (memory safety)
- `/internal/game_state/game_state.go:252` - Wrong map type (state corruption)
- `/internal/map_state/layered_map.go:243` - Nil position in critical path (memory safety)

### 5. Main Function and Core Engine Failures
**Rule**: Application startup and core engine failures should be fatal.

```go
// KEEP FATAL: Core engine startup
if err := ebiten.RunGame(game); err != nil {
    log.Fatal(err) // ✅ Cannot run game engine
}
```

**Examples**:
- `/cmd/ultimav/main.go:37` - Game engine startup failure
- `/cmd/ultimav/gamescene.go:103,121,165` - Critical system initialization failures

## ❌ CONVERT TO SOFT ERRORS - Should Be Recoverable

### 1. UI State Management
**Rule**: UI inconsistencies should not crash the game - log warnings and recover gracefully.

```go
// CONVERT TO SOFT ERROR: UI state issues
dialogIndex := d.FindDialogIndex(inputDialogBox)
if dialogIndex == -1 {
    log.Printf("Warning: input dialog box not found, creating new one")
    // Recovery logic here
    return d.createNewInputDialog()
}
```

**Examples**:
- `/internal/ui/widgets/dialog_stack.go:64,72` - Dialog not found (recoverable UI state)
- `/cmd/ultimav/gamescene.go:239` - Debug dialog index (non-critical)

### 2. Gameplay Logic Validation  
**Rule**: Invalid gameplay actions should show error messages to player, not crash.

```go
// CONVERT TO SOFT ERROR: Gameplay validation
func AdvanceHours(hours int) error {
    if hours > 9 {
        return fmt.Errorf("cannot advance more than 9 hours at a time") // ❌ User input error
    }
}
```

**Examples**:
- `/internal/datetime/ultima_date.go:69` - Invalid time advancement (user input)
- `/internal/ai/npc_ai_controller_small_map.go:303,378` - Unknown AI type (data validation)

### 3. Missing NPC/Object Data
**Rule**: Missing game objects should degrade gracefully, not crash the experience.

```go
// CONVERT TO SOFT ERROR: Missing game data
func findNearbyLadder(pos Position) (*Ladder, error) {
    ladder := searchForLadder(pos)
    if ladder == nil {
        return nil, fmt.Errorf("no ladder found near position %v", pos) // ❌ Recoverable
    }
    return ladder, nil
}
```

**Examples**:
- `/internal/references/small_location_reference.go:183` - Missing ladder/stair (level design issue)
- `/internal/map_units/map_unit_details.go:55` - No path calculated (pathfinding failure)

### 4. Config File Operations
**Rule**: Configuration problems should use defaults and warn user, not prevent startup.

```go
// CONVERT TO SOFT ERROR: Config file issues
func writeConfigFile(path string, data []byte) error {
    if err := os.WriteFile(path, data, 0644); err != nil {
        log.Printf("Warning: could not write config file %s: %v", path, err)
        return err // ❌ Use defaults, don't crash
    }
    return nil
}
```

**Examples**:
- `/internal/config/config.go:63,66` - Config file read/write errors (use defaults)

### 5. Development/Debug Code  
**Rule**: Debug and development-only code should never crash production.

```go
// CONVERT TO SOFT ERROR: Debug/development code
func drawMapUnit(tile Tile) {
    if tile.Index < 0 || tile.Index > maxIndex {
        log.Printf("Debug: unexpected tile index %d, using default", tile.Index)
        tile.Index = defaultTileIndex // ❌ Fallback in production
    }
}
```

**Examples**:
- `/cmd/ultimav/gamescene_tiles.go:137,156,179,203` - Rendering edge cases (use fallbacks)

## Implementation Strategy for Task #7

### Phase 1: Add TODO Comments (Immediate)
All existing `log.Fatal` calls get TODO comments indicating their category:

```go
// Keep as fatal
log.Fatal("data corruption") // TODO: KEEP FATAL - critical data corruption

// Convert to soft error  
log.Fatal("missing dialog") // TODO: CONVERT TO SOFT ERROR - UI recovery needed
```

### Phase 2: Convert Soft Errors (Future Work)
Replace soft error fatals with proper error handling:

```go
// Before
func doSomething() {
    if condition {
        log.Fatal("bad thing happened") // TODO: CONVERT TO SOFT ERROR
    }
}

// After  
func doSomething() error {
    if condition {
        return fmt.Errorf("bad thing happened")
    }
    return nil
}
```

## Error Categories Summary

| Category | Action | Count (Est.) | Priority |
|----------|--------|--------------|----------|
| **Critical Data Loading** | Keep Fatal | ~15 | N/A |
| **Programming Logic Errors** | Keep Fatal | ~20 | N/A | 
| **UI State Issues** | Convert to Soft | ~8 | High |
| **Gameplay Validation** | Convert to Soft | ~5 | High |
| **Config Operations** | Convert to Soft | ~3 | Medium |
| **Debug/Development** | Convert to Soft | ~10 | Medium |

## Testing Strategy

- **Fatal Errors**: Should be tested with unit tests that expect panics
- **Soft Errors**: Should be tested with proper error return validation
- **Recovery Logic**: Should be tested to ensure graceful degradation

## Code Review Checklist

- [ ] Does this error indicate a programming bug? → Keep fatal
- [ ] Does this error indicate system/data corruption? → Keep fatal  
- [ ] Can the user/game recover from this condition? → Convert to soft error
- [ ] Is this a UI/presentation issue? → Convert to soft error
- [ ] Is this debug/development code? → Convert to soft error

---

**Note**: This document was created as part of Task #7 remediation planning. All `log.Fatal` usage should follow these guidelines to ensure appropriate error handling throughout the codebase.