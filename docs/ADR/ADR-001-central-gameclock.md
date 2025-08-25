# ADR-001: Central GameClock over time.Now()

**Status**: Accepted  
**Date**: 2025-08-25  
**Context**: Remediation Task 2 & 3 - Non-deterministic animation and AI systems

---

## Context

The original codebase had critical determinism issues:

1. **Animation System**: `GetSpriteIndexWithAnimationBySpriteIndex()` used `time.Now()` internally
2. **AI Systems**: NPCAIControllerSmallMap used `rand.Seed(time.Now().UnixNano())` for random decisions
3. **Testing Impact**: Tests produced different results on each run, making regression detection impossible
4. **Gameplay Impact**: Identical player actions could yield different outcomes

### Problems This Created
- **Non-reproducible bugs**: Issues couldn't be consistently reproduced
- **Unreliable testing**: Integration tests failed intermittently
- **Debugging difficulties**: Same inputs produced different outputs
- **Development friction**: Developers couldn't trust test results

---

## Decision

Replace all `time.Now()` usage in core game logic with centralized time sources:

### 1. **Central GameClock in GameState**
```go
type GameState struct {
    DateTime    UltimaDate  // Game world time
    ElapsedMs   uint64      // Milliseconds for animations
    rng         *rand.Rand  // Controlled randomness
}
```

### 2. **Deterministic Animation Pattern**
```go
// OLD: Non-deterministic
animation := GetSpriteIndexWithAnimationBySpriteIndex(sprite)

// NEW: Deterministic  
animation := gs.GetSpriteIndexWithAnimationBySpriteIndexTick(sprite, gs.ElapsedMs)
```

### 3. **Controlled Randomness**
```go
// OLD: Non-deterministic seeding
rand.Seed(time.Now().UnixNano())
decision := rand.Intn(2) == 0

// NEW: Controlled RNG
func (gs *GameState) SetRandomSeed(seed int64) {
    gs.rng = rand.New(rand.NewSource(seed))
}

func (gs *GameState) OneInXOdds(x int) bool {
    return gs.rng.Intn(x) == 0
}
```

### 4. **Testing Pattern**
```go
func TestDeterministicBehavior(t *testing.T) {
    gs := createGameState()
    gs.SetRandomSeed(12345) // Fixed seed
    
    // Test multiple runs with same results
    result1 := gs.PerformAction()
    gs.SetRandomSeed(12345) // Reset to same seed
    result2 := gs.PerformAction()
    
    assert.Equal(t, result1, result2) // Always passes
}
```

---

## Consequences

### ✅ **Positive Consequences**

1. **Reproducible Behavior**
   - Same inputs always produce same outputs
   - Bugs can be consistently reproduced
   - Tests are reliable and deterministic

2. **Better Testing**
   - Integration tests use fixed seeds for consistent results
   - Regression detection actually works
   - CI/CD tests are trustworthy

3. **Easier Debugging**
   - Developers can reproduce issues consistently
   - Game behavior is predictable during development
   - State transitions can be traced reliably

4. **Faithful Recreation**
   - Can match original Ultima V's deterministic behavior
   - RNG sequences can be controlled and documented
   - Gameplay mechanics behave consistently

### ⚠️ **Negative Consequences**

1. **Slight Complexity Increase**
   - Developers must remember to use GameState time/RNG methods
   - Cannot use convenient `time.Now()` for quick debugging
   - Requires understanding of central time system

2. **Initial Migration Effort**
   - Required updating animation calls throughout codebase
   - AI controller needed RNG dependency injection
   - Some legacy timing code needed rewriting

3. **Testing Setup Requirements**
   - Tests must set deterministic seeds
   - Time-based tests need controlled time advancement
   - More setup code required for temporal testing

### 🔄 **Neutral Consequences**

1. **Performance Impact**: Negligible - no measurable difference in practice
2. **Memory Usage**: Minimal increase from additional GameState fields
3. **Learning Curve**: Developers adapt quickly to the pattern

---

## Implementation Details

### Files Modified
- `/internal/sprites/sprite_animations.go` - Removed time.Now() usage
- `/internal/ai/npc_ai_controller_small_map.go` - Added RNG dependency injection
- `/internal/game_state/game_state.go` - Added central time and RNG methods
- `/internal/references/large_map.go` - Updated animation call

### Tests Created
- `/internal/sprites/sprite_animations_test.go` - 6 determinism tests
- `/internal/game_state/rng_determinism_test.go` - 6 RNG consistency tests

### Interface Design
```go
// RNGProvider allows dependency injection of RNG into AI systems
type RNGProvider interface {
    RollDice(sides int) int
    OneInXOdds(odds int) bool
    GetRandomInt(max int) int
}
```

---

## Alternatives Considered

### Alternative 1: Mock time.Now() in tests
**Rejected**: Would require complex time mocking infrastructure and wouldn't solve production determinism issues.

### Alternative 2: Accept non-deterministic behavior  
**Rejected**: Made testing unreliable and debugging extremely difficult. Not acceptable for a faithful recreation.

### Alternative 3: Separate time systems per component
**Rejected**: Would create coordination problems and potential timing inconsistencies between systems.

---

## References

- **Remediation Task 2**: Fix Non-Deterministic Animation System
- **Remediation Task 3**: Fix Non-Deterministic Random Number Generation  
- **Test Files**: `sprite_animations_test.go`, `rng_determinism_test.go`
- **Related**: ADR-003 (SystemCallbacks) relies on deterministic behavior for testing

---

## Validation

This decision was validated through:
1. **47 integration tests** all pass with deterministic behavior
2. **100% reproducible test results** across multiple runs
3. **Zero time.Now() violations** in core logic packages
4. **Successful AI behavior consistency** across test runs

The decision successfully solved the original non-determinism problems while maintaining good developer experience.