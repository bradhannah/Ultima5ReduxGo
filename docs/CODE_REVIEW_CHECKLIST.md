# Code Review Checklist

**Purpose**: Prevent regression of issues addressed in the 2025-08-25 codebase remediation and maintain code quality standards.

---

## Critical Logic Checks ⚠️

### Deterministic Systems
- [ ] **No `time.Now()` usage in core game logic**
  - ✅ Use `GameState.DateTime` or central `GameClock` for time-dependent logic
  - ✅ Use `GetSpriteIndexWithAnimationBySpriteIndexTick()` for animations
  - ❌ Avoid `time.Now()` in `/internal/game_state/`, `/internal/ai/`, `/internal/sprites/`
  - **Rationale**: Tasks 2 & 3 - Ensures reproducible behavior for testing and consistency

- [ ] **All random number generation uses GameState RNG methods**
  - ✅ Use `gs.RollDice()`, `gs.OneInXOdds()`, `gs.GetRandomInt()`
  - ✅ Use `gs.SetRandomSeed()` for deterministic testing
  - ❌ Avoid raw `rand` package or `time.Now()` seeding
  - **Rationale**: Task 3 - Eliminates non-deterministic AI behavior, enables consistent testing

### Game Logic Correctness
- [ ] **Key consumption logic follows correct pattern**
  - ✅ Keys consumed BEFORE success check: `DecrementByOne()` then `if isSuccessful()`
  - ❌ Never consume keys only on failure (backwards logic)
  - **Rationale**: Task 1 - Prevents unlimited lock picking exploit

- [ ] **Object collision detection validates both layers**
  - ✅ Check terrain passability: `tile.IsPassable()`
  - ✅ Check object presence: `gs.ObjectPresentAt(position)`
  - ❌ Don't check only one layer (allows walking through objects)
  - **Rationale**: Task 4 - Prevents movement through NPCs, vehicles, objects

---

## Architecture & Package Standards 🏗️

### Package Boundaries
- [ ] **Core game logic packages don't import rendering libraries**
  - ✅ `/internal/game_state/`, `/internal/ai/`, `/internal/map_state/` stay rendering-independent
  - ❌ No `github.com/hajimehoshi/ebiten/v2` imports in core logic
  - **Rationale**: Task 8 - Maintains clean architecture, enables testing without graphics

- [ ] **Use DisplayManager for screen operations**
  - ✅ Use `DisplayManager.GetScreenSize()`, `DisplayManager.Update()`
  - ❌ No direct `ebiten.WindowSize()` calls in non-UI packages
  - **Rationale**: Task 8 - Centralized screen management, multi-resolution support

### Dependency Injection
- [ ] **SystemCallbacks used for UI interactions in GameState**
  - ✅ Use `gs.SystemCallbacks.Message.AddRowStr()` for user messages
  - ✅ Use `gs.SystemCallbacks.Audio.PlaySoundEffect()` for audio
  - ✅ Use `gs.SystemCallbacks.Flow.AdvanceTime()` for time progression
  - ❌ No direct UI calls from GameState methods
  - **Rationale**: Clean separation, testability with mock callbacks

### Import Organization
- [ ] **Clean import organization follows Go conventions**
  - ✅ Standard library first, external packages second, internal packages third
  - ✅ Blank lines between groups
  - ✅ Use `goimports` for consistent formatting
  - ✅ Minimal necessary aliases only (`etext` for ebiten/text conflicts)
  - ❌ No unnecessary aliases like `references2`, `ucolor`, `mainscreen2`
  - **Rationale**: Task 9 - Go standards compliance, code readability

---

## Error Handling Standards 🚨

### log.Fatal Usage
- [ ] **All `log.Fatal` calls have appropriate comments**
  - ✅ **Legitimate fatal cases**: Explanatory comment describing the unrecoverable condition
    ```go
    // System corruption: embedded game data files missing or corrupted
    log.Fatal("Failed to load embedded game data")
    ```
  - ✅ **Cases needing conversion**: TODO comment indicating future soft error handling
    ```go
    // TODO: soften to recoverable error - UI state issue shouldn't crash game
    log.Fatal("Failed to initialize dialog")
    ```
  - ❌ No uncommented `log.Fatal` calls
  - **Rationale**: Task 7 - Proper error categorization, planned conversion to soft errors

- [ ] **Error handling follows ERROR_HANDLING.md guidelines**
  - ✅ Return errors from functions instead of panicking in normal control flow
  - ✅ Wrap underlying errors with context: `fmt.Errorf("loading sprite: %w", err)`
  - ✅ Use `log.Fatal` only for unrecoverable system conditions
  - **Rationale**: Consistent error handling strategy, better user experience

---

## Testing Requirements 🧪

### Integration Testing
- [ ] **Integration tests use real game data**
  - ✅ Use `NewIntegrationTestBuilder()` with actual SAVED.GAM, LOOK2.DAT, TLK files
  - ❌ No stub/mock data for integration tests
  - **Rationale**: Task 10 - Validates against actual Ultima V behavior

- [ ] **Fixed seeds for deterministic testing**
  - ✅ Use `gs.SetRandomSeed(12345)` for reproducible test results
  - ✅ Use mock clocks instead of real time
  - **Rationale**: Reliable test results, regression detection

- [ ] **SystemCallbacks properly mocked for state validation**
  - ✅ Use `MockSystemCallbacks` with tracking and assertions
  - ✅ Verify message contents, time advancement, audio effects
  - **Rationale**: Validates UI integration without actual UI dependencies

### Command Workflow Testing
- [ ] **Command workflows tested end-to-end**
  - ✅ Test complete user action flows: input → logic → feedback
  - ✅ Validate SystemCallbacks integration
  - ✅ Test both success and failure scenarios
  - **Rationale**: Ensures user-facing functionality works correctly

---

## Code Quality Standards 📝

### Documentation
- [ ] **Exported identifiers have doc comments**
  - ✅ Explain the "why" not the obvious "what"
  - ✅ Keep comments up-to-date with behavior changes
  - **Package-level docs** for domains with non-trivial rules

### Performance Considerations
- [ ] **Avoid per-frame allocations in hot paths**
  - ✅ Reuse buffers where possible
  - ✅ Use small structs and pass pointers only when mutation required
  - **Note**: Performance optimization not a current priority, but good practice

---

## Security & Best Practices 🔒

### Secrets & Configuration
- [ ] **Never commit secrets or keys**
  - ❌ No API keys, passwords, or sensitive data in code
  - ✅ Use environment variables or secure configuration

### Data Integrity  
- [ ] **Validate external data inputs**
  - ✅ Check bounds on user inputs and file data
  - ✅ Handle malformed save files gracefully

---

## Review Process Guidelines

### Before Submitting PR
1. **Run automated checks**: `goimports`, `go test`, `go vet`
2. **Self-review against this checklist**
3. **Test with real game data** if touching game logic
4. **Verify no regressions** in deterministic behavior

### During Code Review
1. **Focus on architectural concerns** over style preferences  
2. **Validate against remediation lessons learned**
3. **Check for proper error handling patterns**
4. **Ensure test coverage** for new functionality

### Post-Review Actions
1. **Update documentation** if interfaces changed
2. **Add integration tests** for new commands/workflows
3. **Verify CI passes** before merging

---

## Quick Reference: Common Violations

### ❌ Bad Patterns
```go
// Non-deterministic time usage
animation := GetSpriteIndexWithAnimationBySpriteIndex(spriteIndex)

// Backwards key consumption
if successful {
    keys.DecrementByOne() // Wrong!
}

// Direct UI calls from game logic
fmt.Println("Debug message") // Should use SystemCallbacks

// Core logic importing rendering
import "github.com/hajimehoshi/ebiten/v2" // In game_state package
```

### ✅ Good Patterns  
```go
// Deterministic animations
animation := gs.GetSpriteIndexWithAnimationBySpriteIndexTick(spriteIndex, elapsedMs)

// Correct key consumption  
keys.DecrementByOne() // Consume first
if isSuccessful {
    return Success
}

// Proper UI interaction
gs.SystemCallbacks.Message.AddRowStr("Action completed")

// Clean dependency injection
type GameState struct {
    SystemCallbacks *SystemCallbacks
}
```

---

**Remember**: These standards prevent regression of critical fixes. When in doubt, refer to the original remediation tasks for context and rationale.