# Architecture Guidelines

**Purpose**: Define architectural patterns and boundaries established during the 2025-08-25 codebase remediation.

---

## Core Architectural Principles

### 1. **Core Logic Independence** 🎯
**Principle**: Game logic must remain independent of rendering and input systems.

#### Package Boundaries
```
✅ ALLOWED DEPENDENCIES:
internal/game_state/ → internal/references/, internal/party_state/, internal/ai/
internal/ai/        → internal/references/, internal/game_state/ (interfaces)
internal/map_state/ → internal/references/

❌ FORBIDDEN DEPENDENCIES:  
internal/game_state/ → github.com/hajimehoshi/ebiten/v2
internal/ai/        → cmd/ultimav/ or internal/ui/
Core packages       → Rendering/UI packages
```

#### Benefits
- **Testability**: Core logic testable without graphics initialization
- **Portability**: Could support different rendering backends
- **Performance**: No rendering overhead in game logic calculations
- **Determinism**: Pure game state calculations unaffected by rendering timing

---

### 2. **Dependency Injection Pattern** 💉
**Principle**: External dependencies injected via interfaces, not imported directly.

#### SystemCallbacks Pattern
```go
// ✅ GOOD: Dependency injection
type GameState struct {
    SystemCallbacks *SystemCallbacks
}

func (gs *GameState) ActionJimmy() bool {
    // Game logic...
    gs.SystemCallbacks.Message.AddRowStr("Unlocked!")
    gs.SystemCallbacks.Audio.PlaySoundEffect(SoundUnlock)
    gs.SystemCallbacks.Flow.AdvanceTime(1)
    return true
}

// ❌ BAD: Direct coupling
func (gs *GameState) ActionJimmy() bool {
    // Game logic...
    fmt.Println("Unlocked!")              // Direct output
    PlaySound("unlock.wav")                // Direct audio call
    advanceGameTime(1)                     // Direct time manipulation
}
```

#### SystemCallbacks Interface Design
```go
type SystemCallbacks struct {
    Message  MessageCallbacks   // User feedback
    Audio    AudioCallbacks     // Sound effects
    Visual   VisualCallbacks    // Screen effects
    Screen   ScreenCallbacks    // Display updates  
    Flow     FlowCallbacks      // Time, state flow
    Talk     TalkCallbacks      // Dialog system
}
```

---

### 3. **DisplayManager Pattern** 🖥️
**Principle**: Centralize screen management to maintain package boundaries.

#### Usage Pattern
```go
// ✅ GOOD: Via DisplayManager
import "github.com/bradhannah/Ultima5ReduxGo/internal/display"

func GetGameArea() (width, height int) {
    return display.GetManager().GetScreenSize()
}

// ❌ BAD: Direct Ebitengine import
import "github.com/hajimehoshi/ebiten/v2"

func GetGameArea() (width, height int) {
    return ebiten.WindowSize() // Violates package boundaries
}
```

#### DisplayManager Responsibilities
- **Screen size management** across different resolutions
- **Resolution change detection** and callbacks
- **Centralized display state** (fullscreen, windowed modes)
- **Multi-monitor support** (future enhancement)

---

### 4. **Deterministic Systems** 🎲
**Principle**: All game systems use centralized, controllable sources of randomness and time.

#### Central GameClock Pattern
```go
// ✅ GOOD: Central game clock
type GameState struct {
    DateTime    UltimaDate
    ElapsedMs   uint64
}

func (gs *GameState) GetAnimation(sprite SpriteIndex) SpriteIndex {
    return GetSpriteIndexWithAnimationBySpriteIndexTick(sprite, gs.ElapsedMs)
}

// ❌ BAD: Non-deterministic time
func GetAnimation(sprite SpriteIndex) SpriteIndex {
    return GetSpriteIndexWithAnimationBySpriteIndex(sprite) // Uses time.Now()
}
```

#### Central RNG Pattern  
```go
// ✅ GOOD: Controlled randomness
type GameState struct {
    rng *rand.Rand
}

func (gs *GameState) RollDice(sides int) int {
    return gs.rng.Intn(sides) + 1
}

func (gs *GameState) SetRandomSeed(seed int64) {
    gs.rng = rand.New(rand.NewSource(seed))
}

// ❌ BAD: Uncontrolled randomness
func RollDice(sides int) int {
    rand.Seed(time.Now().UnixNano()) // Non-deterministic
    return rand.Intn(sides) + 1
}
```

---

### 5. **Error Handling Strategy** ⚠️
**Principle**: Distinguish between recoverable and unrecoverable errors.

#### Error Classification
```go
// ✅ LEGITIMATE FATAL ERRORS (System corruption, programming bugs)
if gameData == nil {
    // System corruption: embedded game data missing or corrupted  
    log.Fatal("Failed to load embedded game data - system integrity compromised")
}

// ✅ RECOVERABLE ERRORS (Return error, don't crash)
func LoadUserSaveFile(path string) (*SaveData, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("loading save file: %w", err)
    }
    // ... process data
}

// ❌ BAD: Crashing on user/config issues  
func LoadConfig(path string) Config {
    data, err := os.ReadFile(path)
    if err != nil {
        log.Fatal("Config file not found") // Should be recoverable
    }
}
```

#### Error Handling Guidelines
- **System corruption, embedded data issues**: `log.Fatal` with explanatory comment
- **User data, configuration, validation**: Return `error` with context
- **Programming bugs, impossible states**: `log.Fatal` with explanatory comment  
- **UI state, dialog issues**: `TODO: soften to recoverable error` comment

---

## Package Organization Patterns

### 1. **Core Game Logic Layer**
```
internal/
├── game_state/     # Central game state and actions
├── party_state/    # Player party data and inventory
├── map_state/      # World maps and location data
├── ai/            # NPC AI and pathfinding
└── references/    # Game data constants and enums
```
**Constraints**: No rendering imports, use SystemCallbacks for UI interaction

### 2. **Data and References Layer**  
```
internal/
├── references/    # Sprites, locations, directions, constants
├── conversation/  # Dialog system and TalkScript processing
└── text/         # Text processing and display utilities
```
**Constraints**: Pure data processing, minimal external dependencies

### 3. **System Integration Layer**
```
internal/
├── display/       # DisplayManager for screen operations
├── config/        # Configuration management  
└── sprites/       # Sprite processing and animations
```
**Constraints**: Bridge between core logic and presentation

### 4. **Presentation Layer**
```  
cmd/ultimav/       # Main game UI and rendering
internal/ui/       # UI widgets and components
```
**Constraints**: Can import anything needed, but shouldn't be imported by core logic

---

## Testing Architecture

### 1. **Unit Tests** - Package Level
- Test individual functions and methods
- No external dependencies (files, network, UI)
- Fast execution, deterministic results

### 2. **Integration Tests** - System Level  
- Test multiple packages working together
- Use real Ultima V game data (SAVED.GAM, LOOK2.DAT, TLK files)
- MockSystemCallbacks for UI validation
- IntegrationTestBuilder for complex scenarios

### 3. **End-to-End Tests** - Workflow Level
- Test complete user workflows  
- Validate cross-system integration
- Performance and regression detection

#### Testing Pattern
```go
// Integration test with real data
func TestJimmyWorkflow_Integration(t *testing.T) {
    gs, mockCallbacks := NewIntegrationTestBuilder(t).
        WithLocation(references.Britain).
        WithPlayerAt(15, 15).
        WithSystemCallbacks().
        Build()
        
    gs.SetRandomSeed(12345) // Deterministic
    result := gs.ActionJimmySmallMap(references.Up)
    
    // Validate SystemCallbacks integration
    assert.True(t, len(mockCallbacks.Messages) > 0)
    assert.Equal(t, "Not lock!", mockCallbacks.Messages[0])
}
```

---

## Command Implementation Pattern

### 1. **UI Layer** - Input Handling
```go
// GameScene methods - handle input, call GameState
func (g *GameScene) smallMapJimmySecondary(direction references.Direction) {
    // Early validation if needed
    if simpleCheck() {
        g.output.AddRowStr("Simple error")
        return  
    }
    
    // Delegate to GameState
    success := g.gameState.ActionJimmySmallMap(direction)
    
    // Additional UI logic if needed based on success
}
```

### 2. **GameState Layer** - Core Logic
```go
// GameState methods - core logic with SystemCallbacks
func (gs *GameState) ActionJimmySmallMap(direction references.Direction) bool {
    // Validation logic
    if !gs.hasKeys() {
        gs.SystemCallbacks.Message.AddRowStr("No lock picks!")
        return false
    }
    
    // Core game logic
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

---

## Migration Guidelines

### When Refactoring Existing Code

1. **Identify Dependencies**: What external systems does this code use?
2. **Extract Interfaces**: Create interfaces for external dependencies  
3. **Inject Dependencies**: Pass interfaces instead of importing directly
4. **Update Tests**: Use mocks/stubs for injected dependencies
5. **Validate Boundaries**: Ensure package constraints maintained

### When Adding New Features

1. **Follow Existing Patterns**: Use established SystemCallbacks, DisplayManager patterns
2. **Maintain Determinism**: Use GameState RNG and time sources
3. **Write Integration Tests**: Use real game data, test end-to-end workflows
4. **Document Decisions**: Update this file for new architectural patterns

---

## Common Anti-Patterns to Avoid

### ❌ **Circular Dependencies**
```
game_state → ui → game_state (BAD)
```
**Solution**: Use dependency injection, interfaces

### ❌ **God Objects**  
```go
type GameManager struct {
    // 50+ fields managing everything
}
```
**Solution**: Split responsibilities, focused packages

### ❌ **Hard-coded Dependencies**
```go
func ProcessAction() {
    PlaySound("effect.wav")     // Hard-coded audio
    UpdateScreen()              // Hard-coded UI
}
```
**Solution**: Inject callbacks/interfaces

### ❌ **Non-Deterministic Systems**
```go
func AIDecision() bool {
    return time.Now().UnixNano()%2 == 0  // Non-deterministic
}
```
**Solution**: Use central GameState RNG

---

## Architecture Decision Records (ADRs)

### ADR-001: Central GameClock over time.Now()
**Context**: Animation and AI systems were using `time.Now()` making behavior non-deterministic  
**Decision**: Use central GameClock with controllable time source  
**Consequences**: Reproducible behavior, easier testing, slight complexity increase

### ADR-002: DisplayManager Pattern  
**Context**: Core logic was importing Ebitengine directly, violating package boundaries  
**Decision**: Centralize screen operations in DisplayManager  
**Consequences**: Clean boundaries, multi-resolution ready, centralized screen management

### ADR-003: SystemCallbacks for UI Decoupling
**Context**: GameState methods needed UI interaction without direct coupling  
**Decision**: Dependency injection via SystemCallbacks interface  
**Consequences**: Testable game logic, clean separation of concerns, slight complexity increase

---

**Remember**: Architecture serves the codebase, not vice versa. These patterns emerged from solving real problems during remediation. Adapt them as the codebase evolves, but maintain the core principles of independence, determinism, and testability.