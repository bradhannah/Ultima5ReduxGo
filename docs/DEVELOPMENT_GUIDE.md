# Development Guide

**Quick start guide for contributors to the Ultima 5 Redux Go project**

---

## Getting Started

### Prerequisites
- **Go 1.21+** (check `go.mod` for exact version)
- **Git** with hooks support
- **Real Ultima V game data** in `/Users/[username]/games/Ultima_5/Gold/` for integration testing

### Initial Setup
```bash
# Clone and setup
git clone https://github.com/bradhannah/Ultima5ReduxGo.git
cd Ultima5ReduxGo

# Install git hooks for validation
git config core.hooksPath .githooks

# Verify setup
./scripts/pre-commit-checks.sh
go test ./...
```

---

## Architecture Overview

### Core Principles 🏗️
1. **Core Logic Independence**: Game logic stays separate from rendering
2. **Deterministic Systems**: All time and randomness controlled centrally  
3. **Dependency Injection**: UI interactions via SystemCallbacks
4. **Real Data Testing**: Integration tests use actual Ultima V files

### Package Structure
```
internal/
├── game_state/      # Central game state and player actions
├── ai/             # NPC AI and pathfinding  
├── map_state/      # World maps and locations
├── display/        # Screen management (DisplayManager)
├── references/     # Game constants and enums
└── ui/             # UI widgets and components

cmd/ultimav/        # Main application and rendering
```

---

## Common Development Tasks

### Adding a New Player Command

1. **Follow the Action pattern**:
```go
// In internal/game_state/action_[command].go
func (gs *GameState) Action[Command]SmallMap(direction references.Direction) bool {
    // Validation
    if !gs.canPerformAction() {
        gs.SystemCallbacks.Message.AddRowStr("Can't do that!")
        return false
    }
    
    // Core logic
    success := gs.performAction(direction)
    
    // User feedback via SystemCallbacks
    if success {
        gs.SystemCallbacks.Message.AddRowStr("Success!")
        gs.SystemCallbacks.Audio.PlaySoundEffect(SoundSuccess)
        gs.SystemCallbacks.Flow.AdvanceTime(1)
        return true
    }
    
    gs.SystemCallbacks.Message.AddRowStr("Failed!")
    return false
}
```

2. **Write integration tests**:
```go
func TestNewCommandWorkflow_Integration(t *testing.T) {
    gs, mockCallbacks := NewIntegrationTestBuilder(t).
        WithLocation(references.Britain).
        WithPlayerAt(15, 15).
        WithSystemCallbacks().
        Build()
        
    result := gs.ActionNewCommandSmallMap(references.Up)
    
    // Validate SystemCallbacks integration
    assert.True(t, len(mockCallbacks.Messages) > 0)
    assert.Equal(t, "Expected message", mockCallbacks.Messages[0])
}
```

### Adding Animation or Time-Based Logic

```go
// ✅ GOOD: Use central game time
func (gs *GameState) GetAnimationFrame(sprite SpriteIndex) SpriteIndex {
    return GetSpriteIndexWithAnimationBySpriteIndexTick(sprite, gs.ElapsedMs)
}

// ❌ NEVER: Use time.Now() 
func GetAnimationFrame(sprite SpriteIndex) SpriteIndex {
    return GetSpriteIndexWithAnimationBySpriteIndex(sprite) // Uses time.Now()
}
```

### Adding Random Events or AI Logic

```go
// ✅ GOOD: Use GameState RNG
func (gs *GameState) RandomEvent() bool {
    return gs.OneInXOdds(10) // 10% chance
}

func (gs *GameState) RollAttackDamage() int {
    return gs.RollDice(6) + gs.GetPlayerStats().Strength
}

// ❌ NEVER: Direct rand package
func RandomEvent() bool {
    return rand.Intn(10) == 0 // Non-deterministic
}
```

---

## Testing Guide

### Test Types

1. **Unit Tests**: Fast, no external dependencies
2. **Integration Tests**: Use real Ultima V game data
3. **End-to-End Tests**: Complete user workflows

### Integration Testing Pattern
```go
// Always use IntegrationTestBuilder for complex scenarios
gs, mockCallbacks := NewIntegrationTestBuilder(t).
    WithLocation(references.Britain).        // Real map data
    WithPlayerAt(15, 15).                   // Starting position
    WithSystemCallbacks().                  // Mock UI callbacks
    Build()                                 // Loads real SAVED.GAM

// Use fixed seeds for deterministic tests  
gs.SetRandomSeed(12345)

// Test your logic
result := gs.ActionSomething(references.Up)

// Validate SystemCallbacks were called properly
assert.Equal(t, 1, len(mockCallbacks.Messages))
assert.Equal(t, "Expected output", mockCallbacks.Messages[0])
```

### Running Tests
```bash
# Unit tests
go test ./internal/game_state

# Integration tests (requires real game data)
go test ./internal/game_state -run "Integration"

# All tests
go test ./...
```

---

## Common Pitfalls & Solutions

### ❌ Problem: Non-Deterministic Behavior
```go
// This breaks testing
func AIDecision() bool {
    return time.Now().UnixNano()%2 == 0
}
```
**✅ Solution**: Use GameState RNG
```go
func (gs *GameState) AIDecision() bool {
    return gs.OneInXOdds(2)
}
```

### ❌ Problem: Package Boundary Violations
```go
// Core logic importing rendering
package game_state
import "github.com/hajimehoshi/ebiten/v2"
```
**✅ Solution**: Use DisplayManager
```go
package game_state
import "github.com/bradhannah/Ultima5ReduxGo/internal/display"

func (gs *GameState) GetScreenSize() (int, int) {
    return display.GetManager().GetScreenSize()
}
```

### ❌ Problem: Direct UI Calls from Game Logic
```go
func (gs *GameState) ActionPush() bool {
    fmt.Println("Pushed!") // Wrong layer
    return true
}
```
**✅ Solution**: Use SystemCallbacks
```go
func (gs *GameState) ActionPush() bool {
    gs.SystemCallbacks.Message.AddRowStr("Pushed!")
    gs.SystemCallbacks.Audio.PlaySoundEffect(SoundPush)
    return true
}
```

### ❌ Problem: Backwards Key Consumption
```go
// This allows unlimited attempts
if successful {
    return Success
} else {
    keys.DecrementByOne() // Wrong!
    return Failure
}
```
**✅ Solution**: Consume Before Attempt
```go
if keys.Get() <= 0 {
    return NoKeys
}

keys.DecrementByOne() // Always consume
if successful {
    return Success
} else {
    return Failure
}
```

---

## Debugging Tips

### Deterministic Debugging
```go
// For reproducible debugging, use fixed seeds
gs.SetRandomSeed(12345)

// Log state at key points
t.Logf("Game state: Time=%d:%02d, Turn=%d", 
    gs.DateTime.Hour, gs.DateTime.Minute, gs.DateTime.Turn)
```

### Integration Test Debugging
```go
// MockSystemCallbacks provides detailed logging
mockCallbacks.Reset() // Clear previous state
result := gs.ActionSomething(direction)

// Check what callbacks were triggered
t.Logf("Messages: %v", mockCallbacks.Messages)
t.Logf("Audio: %v", mockCallbacks.SoundEffectsPlayed)  
t.Logf("Time advanced: %v", mockCallbacks.TimeAdvanced)
```

### Real Game Data Issues
If integration tests fail due to missing game data:
```go
if gs == nil {
    t.Skip("Skipping test - real game data not available")
    return
}
```

---

## Code Review Process

### Before Submitting PR

1. **Run pre-commit checks**: `./scripts/pre-commit-checks.sh`
2. **Self-review**: Check against `docs/CODE_REVIEW_CHECKLIST.md`
3. **Test with real data**: Verify integration tests pass
4. **Document changes**: Update relevant docs if needed

### Pre-commit Hook Setup
```bash
# One-time setup to enable automatic validation
git config core.hooksPath .githooks
```

This will automatically run validation before each commit.

---

## Project-Specific Patterns

### SystemCallbacks Categories
- **Message**: User feedback text (`AddRowStr`, `ClearScreen`)
- **Audio**: Sound effects (`PlaySoundEffect`)
- **Visual**: Screen effects (`UpdateDisplay`, `RefreshMap`)  
- **Flow**: Time and state (`AdvanceTime`, `PushDialog`)
- **Talk**: Dialog system (`CreateTalkDialog`, `PushDialog`)

### Deterministic Systems
- **Time**: Use `gs.DateTime`, `gs.ElapsedMs`
- **Randomness**: Use `gs.RollDice()`, `gs.OneInXOdds()`
- **Animation**: Use `*Tick()` variants with `gs.ElapsedMs`

### Error Handling Strategy
- **System corruption**: `log.Fatal` with explanatory comment
- **User errors**: Return `error` with context
- **Config issues**: Return `error` or TODO comment for `log.Fatal`

---

## Quick Reference Commands

```bash
# Development workflow
./scripts/pre-commit-checks.sh    # Validate before commit
go test ./...                     # Run all tests
go test -run Integration ./...    # Integration tests only
goimports -w .                    # Fix import formatting
go vet ./...                      # Static analysis

# Integration testing
go test ./internal/game_state -v -run "TestJimmyWorkflow"
go test ./internal/game_state -v -run "Integration"

# Build and run
go build -o bin/ultimav cmd/ultimav/main.go
./bin/ultimav
```

---

## Getting Help

### Documentation
- `docs/CODE_REVIEW_CHECKLIST.md` - Code standards checklist
- `docs/ARCHITECTURE_GUIDELINES.md` - Architecture patterns
- `docs/CODING_CONVENTIONS.md` - General coding standards
- `docs/ERROR_HANDLING.md` - Error handling guidelines

### Common Issues
- **"Non-deterministic test results"** → Use `gs.SetRandomSeed()` and avoid `time.Now()`
- **"Package boundary violation"** → Use DisplayManager and SystemCallbacks
- **"Integration test fails"** → Check if real Ultima V game data is available
- **"Import formatting errors"** → Run `goimports -w .`

### Architecture Questions
When in doubt about architectural decisions, refer to the remediation lessons in `docs/CODING_CONVENTIONS.md` or the patterns in `docs/ARCHITECTURE_GUIDELINES.md`.

---

**Remember**: This codebase prioritizes deterministic behavior, clean architecture, and compatibility with the original Ultima V. When adding new features, always consider how they fit into these principles.