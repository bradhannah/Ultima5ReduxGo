# ADR-002: DisplayManager Pattern

**Status**: Accepted  
**Date**: 2025-08-25  
**Context**: Remediation Task 8 - Package boundary violations

---

## Context

The original codebase had architectural issues with package boundaries:

1. **Core Logic Dependencies**: `internal/config/` was importing `github.com/hajimehoshi/ebiten/v2` directly
2. **Tight Coupling**: Game logic was tightly coupled to specific rendering implementation
3. **Testing Difficulties**: Core logic tests required graphics initialization
4. **Architecture Violation**: Business logic mixed with presentation concerns

### Specific Problems
- `internal/config/config.go` calling `ebiten.WindowSize()` directly
- Core game logic packages importing rendering libraries
- Impossible to test game logic without graphics context
- Future multi-resolution or different renderer support blocked

---

## Decision

Create a **DisplayManager** as a centralized management layer (not abstraction) for screen operations:

### 1. **DisplayManager Architecture**
```go
// internal/display/manager.go
type Manager struct {
    currentWidth  int
    currentHeight int
    callbacks     []ResolutionCallback
}

func GetManager() *Manager {
    // Singleton pattern for centralized access
}

func (m *Manager) GetScreenSize() (width, height int)
func (m *Manager) Update() // Called each frame to detect changes  
func (m *Manager) OnResolutionChange(callback ResolutionCallback)
```

### 2. **Package Boundary Pattern**
```go
// BEFORE: Core logic importing rendering
package config
import "github.com/hajimehoshi/ebiten/v2"

func GetDisplayArea() (int, int) {
    return ebiten.WindowSize() // Violates boundaries
}

// AFTER: Core logic using DisplayManager
package config
import "github.com/bradhannah/Ultima5ReduxGo/internal/display"

func GetDisplayArea() (int, int) {
    return display.GetManager().GetScreenSize()
}
```

### 3. **Game Loop Integration**
```go
// cmd/ultimav/scene.go - Game loop calls Update()
func (g *GameScene) Update() error {
    display.GetManager().Update() // Detect resolution changes
    // ... rest of game logic
}
```

### 4. **Multi-Resolution Preparation**
```go
// Resolution change detection and callbacks
func (m *Manager) Update() {
    newWidth, newHeight := ebiten.WindowSize()
    if newWidth != m.currentWidth || newHeight != m.currentHeight {
        m.currentWidth, m.currentHeight = newWidth, newHeight
        m.notifyCallbacks()
    }
}
```

---

## Consequences

### ✅ **Positive Consequences**

1. **Clean Package Boundaries**
   - Core logic packages no longer import rendering libraries
   - Clear separation between business logic and presentation
   - Easier to reason about dependencies

2. **Improved Testability** 
   - Game logic can be tested without graphics initialization
   - Unit tests run faster without rendering overhead
   - Integration tests can mock display operations

3. **Future Extensibility**
   - Multi-resolution support infrastructure in place
   - Different rendering backends possible
   - Screen management centralized for easier changes

4. **Centralized Screen Management**
   - All screen size queries go through one place
   - Resolution change detection built-in
   - Consistent screen state across the application

### ⚠️ **Trade-offs**

1. **Additional Layer**
   - Introduces one more layer in the call chain
   - Slightly more complex than direct calls
   - Developers need to know about DisplayManager

2. **Singleton Pattern**
   - Uses singleton for global access convenience
   - Could be dependency injected instead for purity
   - Global state, though read-only for most consumers

3. **Not True Abstraction**
   - Still tied to Ebitengine underneath
   - Management layer, not abstraction layer
   - Wouldn't support completely different renderers without changes

### 🔄 **Neutral Consequences**

1. **Performance**: Negligible overhead from indirection
2. **Memory**: Minimal additional memory usage
3. **Learning Curve**: Simple concept, easy to adopt

---

## Implementation Details

### Files Created
- `/internal/display/manager.go` - DisplayManager implementation
- `/internal/display/resolution.go` - Resolution change detection

### Files Modified  
- `/internal/config/config.go` - Removed direct Ebitengine imports
- `/cmd/ultimav/scene.go` - Added DisplayManager.Update() call
- `/internal/text/scale.go` - Uses DisplayManager + legacy compatibility
- Multiple UI files updated to use DisplayManager

### Interface Design
```go
// Callback pattern for resolution changes
type ResolutionCallback func(oldWidth, oldHeight, newWidth, newHeight int)

// Simple management interface
type Manager interface {
    GetScreenSize() (width, height int)
    Update()
    OnResolutionChange(callback ResolutionCallback)
}
```

### Integration Pattern
```go
// In game loop - detect changes
display.GetManager().Update()

// In components - get current size
width, height := display.GetManager().GetScreenSize()

// For dynamic layouts - register callbacks
display.GetManager().OnResolutionChange(func(_, _, newW, newH int) {
    // Adjust layout for new size
})
```

---

## Alternatives Considered

### Alternative 1: Full Abstraction Layer
**Rejected**: Would require abstracting away Ebitengine completely, which is not a current goal. Management layer provides benefits without over-engineering.

### Alternative 2: Dependency Injection Throughout
**Rejected**: Would require passing display interface through many layers. Singleton provides pragmatic global access for read-only operations.

### Alternative 3: Keep Direct Ebitengine Calls
**Rejected**: Violates package boundaries and makes testing difficult. Benefits don't justify the architectural compromise.

### Alternative 4: Static Helper Functions
**Considered**: Could have used static functions instead of singleton. Chose singleton for state management and callback support.

---

## References

- **Remediation Task 8**: Fix Package Boundary Violations
- **Files**: `internal/display/manager.go`, `internal/display/resolution.go`
- **Related**: Enables clean testing in integration test framework
- **Future**: Prepares for multi-resolution support requests

---

## Validation

This decision was validated through:
1. **Zero package boundary violations**: No core logic imports rendering libraries
2. **Successful testing**: Game logic tests run without graphics context
3. **Clean architecture**: Clear separation between business and presentation logic
4. **Multi-resolution ready**: Infrastructure in place for future enhancements

The DisplayManager successfully resolved package boundary issues while preparing for future extensibility needs.

---

## Future Considerations

1. **Multi-Resolution Support**: DisplayManager provides foundation for responsive layouts
2. **Different Renderers**: Could potentially support different backends with DisplayManager as coordination point
3. **Performance Monitoring**: Could add frame rate monitoring and display performance metrics
4. **Mobile Support**: Resolution change detection supports orientation changes