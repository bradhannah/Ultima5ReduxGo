# Testing Documentation

This document outlines testing strategies, conventions, and best practices for the Ultima 5 Redux Go project. The testing approach emphasizes practical test coverage, regression prevention, and efficient development workflows.

## Testing Philosophy

### Core Principles
1. **Go Conventions First**: Follow standard Go testing patterns (`*_test.go` files, `testing` package)
2. **Regression Prevention**: Every bug fix gets a test to prevent reoccurrence  
3. **Helper Function Reuse**: Avoid duplicative test code through well-designed test utilities
4. **Real Data Integration**: Use actual Ultima V data files for integration tests when possible
5. **Deterministic Testing**: Use fixed seeds, mock clocks, and controlled inputs for reproducible results

### Test Categories
- **Unit Tests**: Test individual functions and methods in isolation
- **Integration Tests**: Test component interactions with real game data
- **Regression Tests**: Prevent known bugs from returning (heavily documented)
- **Property Tests**: Test invariants and edge cases with generated data

## Project Setup

### Test Data Sources

#### Original Ultima V Data
Integration tests rely on a local copy of the original Ultima V. Configuration via Viper:
```yaml
# config.yaml
game:
  ultima_v_path: "/path/to/original/ultima5"
```

#### Fan Remake Reference (@OLD)
The `@OLD` directory contains a fan remake implementation for behavioral reference:
- **Files to Reference**: `*.C`, `*.H`, `*.ASM` files only
- **Use Case**: When pseudocode in docs is unclear, reference original behavior
- **Key Files**: `U5DEFS.H`, `COMSUBS*.C`, `TALKNPC.C`, `NPCTRAK.C`

### Test Structure
```
/
├── internal/
│   ├── package_name/
│   │   ├── file.go
│   │   └── file_test.go          # Unit tests
│   └── integration/
│       └── component_test.go      # Integration tests
├── testdata/
│   ├── fixtures/
│   └── mocks/
└── test/
    ├── helpers/                   # Shared test utilities
    └── regression/                # Regression test suite
```

## Test Conventions

### File Naming
- **Unit Tests**: `filename_test.go` alongside source files
- **Integration Tests**: `integration/*_test.go` 
- **Regression Tests**: `test/regression/bug_YYYYMMDD_description_test.go`

### Function Naming
```go
func TestFunctionName(t *testing.T) {}           // Basic test
func TestFunctionName_EdgeCase(t *testing.T) {}  // Specific scenario
func BenchmarkFunctionName(b *testing.B) {}     // Performance test
```

### Test Organization
```go
func TestConversationEngine_ProcessInput(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected ConversationResponse
        setup    func() *LinearConversationEngine
    }{
        {
            name:  "basic_name_request",
            input: "name",
            expected: ConversationResponse{
                Output: "I am Alistair",
                NeedsInput: false,
            },
            setup: func() *LinearConversationEngine {
                return setupEngineWithNPC("alistair")
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            engine := tt.setup()
            response := engine.ProcessInput(tt.input)
            assert.Equal(t, tt.expected, response)
        })
    }
}
```

## Test Helpers and Utilities

### Real Game Data Test Helpers

For complex integration tests that require actual game data, use this infrastructure:

```go
// Helper for tests that need a fully initialized game with real data
func createTestGameStateWithRealData(t *testing.T, location references.Location) (*GameState, *config.UltimaVConfiguration, *references.GameReferences) {
    t.Helper()
    
    config := loadTestConfiguration(t)
    gameRefs := loadTestGameReferences(t, config)
    
    gs := game_state.NewGameState(gameRefs, config)
    
    // Initialize specific map if needed
    if location != references.Britannia_Underworld {
        gs.LoadSmallMap(location, references.FloorNumber(0))
    }
    
    return gs, config, gameRefs
}

// Helper for loading test configuration (assumes Ultima V data is available)
func loadTestConfiguration(t *testing.T) *config.UltimaVConfiguration {
    t.Helper()
    
    // Try common test data locations
    testPaths := []string{
        "/Users/bradhannah/games/Ultima_5/Gold/",  // Developer setup
        "./testdata/",                              // Relative test data
        os.Getenv("ULTIMA5_DATA_PATH"),            // Environment variable
    }
    
    for _, path := range testPaths {
        if path != "" && pathExists(path) {
            config, err := config.LoadUltimaVConfiguration(path)
            if err == nil {
                return config
            }
        }
    }
    
    t.Skip("Ultima V game data not found - set ULTIMA5_DATA_PATH environment variable")
    return nil
}

func pathExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func loadTestGameReferences(t *testing.T, config *config.UltimaVConfiguration) *references.GameReferences {
    t.Helper()
    
    gameRefs, err := references.LoadGameReferences(config)
    if err != nil {
        t.Fatalf("Failed to load game references: %v", err)
    }
    
    return gameRefs
}
```

This approach ensures tests use actual game data and proper initialization. Tests requiring complex setup should be marked with `t.Skip("Converting to use real game data - see TESTING.md")` until converted to use these helpers.

### Test Categories for Real Data Approach

**Simple Unit Tests**: Test individual functions with minimal setup (e.g., `IsPushable()`)
- No skip needed
- Use direct struct initialization
- Fast execution

**Integration Tests**: Test complex interactions requiring game data
- Use `t.Skip()` with real data message for now
- Will be converted to use helpers above
- Require full game initialization

**Data-Driven Tests**: Test against actual game files
- Load real NPCs, maps, items from game data
- Verify behavior matches original implementation
- Most reliable for ensuring compatibility

### Common Test Helpers

#### Mock Game Clock
```go
// test/helpers/clock.go
type MockGameClock struct {
    currentTime time.Time
    fixedSeed   int64
}

func NewMockGameClock(seed int64) *MockGameClock {
    return &MockGameClock{
        currentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
        fixedSeed:   seed,
    }
}

func (m *MockGameClock) GetCurrentTime() time.Time {
    return m.currentTime
}

func (m *MockGameClock) AdvanceHours(hours int) {
    m.currentTime = m.currentTime.Add(time.Duration(hours) * time.Hour)
}
```

#### Test Game State Builder
```go
// test/helpers/gamestate.go
type GameStateBuilder struct {
    state *game_state.GameState
}

func NewGameStateBuilder() *GameStateBuilder {
    return &GameStateBuilder{
        state: game_state.NewGameState(),
    }
}

func (b *GameStateBuilder) WithPlayerAt(x, y int) *GameStateBuilder {
    b.state.MapState.PlayerLocation.Position = references.Position{
        X: references.Coordinate(x),
        Y: references.Coordinate(y),
    }
    return b
}

func (b *GameStateBuilder) WithTime(hour int) *GameStateBuilder {
    b.state.DateTime.Hour = byte(hour)
    return b
}

func (b *GameStateBuilder) Build() *game_state.GameState {
    return b.state
}
```

#### NPC Test Setup
```go
// test/helpers/npc.go
func CreateTestNPCWithDialog(npcID int, dialog string) *references.NPCReference {
    // Helper to create NPCs with specific dialog for testing
    return &references.NPCReference{
        NPCIndex: npcID,
        Name:     fmt.Sprintf("TestNPC_%d", npcID),
        // ... setup test dialog
    }
}

func LoadRealNPCFromTLK(location string, npcID int) (*references.NPCReference, error) {
    // Helper to load actual NPCs from TLK files for integration tests
    // Uses Viper config to find Ultima V data path
}
```

#### Assertion Helpers
```go
// test/helpers/assertions.go
func AssertPositionEquals(t *testing.T, expected, actual references.Position) {
    t.Helper()
    if expected.X != actual.X || expected.Y != actual.Y {
        t.Errorf("Position mismatch: expected (%d,%d), got (%d,%d)", 
                expected.X, expected.Y, actual.X, actual.Y)
    }
}

func AssertConversationOutput(t *testing.T, expected string, response ConversationResponse) {
    t.Helper()
    if !strings.Contains(response.Output, expected) {
        t.Errorf("Expected output to contain '%s', got: %s", expected, response.Output)
    }
}
```

## Integration Testing

### TLK File Testing
```go
func TestConversationEngine_RealTLKData(t *testing.T) {
    // Integration test using actual CASTLE.TLK data
    tlkPath := viper.GetString("game.ultima_v_path") + "/CASTLE.TLK"
    if _, err := os.Stat(tlkPath); os.IsNotExist(err) {
        t.Skip("Original Ultima V data not available")
    }
    
    script, err := references.LoadTalkScript(tlkPath)
    require.NoError(t, err)
    
    engine := conversation.NewLinearConversationEngine(script, &MockCallbacks{})
    
    // Test with real NPCs (Alistair, Treanna, etc.)
    response := engine.Start(1) // Alistair
    assert.Contains(t, response.Output, "greet")
}
```

### Pathfinding Integration
```go
func TestNPCPathfinding_RealMapData(t *testing.T) {
    // Test pathfinding with actual map layouts
    mapData := loadTestMapData("BRITAIN.DAT")
    pathfinder := astar.NewPathfinder(mapData)
    
    start := references.Position{X: 10, Y: 10}
    goal := references.Position{X: 20, Y: 20}
    
    path, found := pathfinder.FindPath(start, goal)
    assert.True(t, found, "Should find path between valid positions")
    assert.Greater(t, len(path), 0, "Path should not be empty")
}
```

## Regression Testing

### Documentation Requirements
Every regression test must include:
1. **Bug Report Reference**: Link to issue or description
2. **Root Cause**: What caused the bug
3. **Test Purpose**: Why this test prevents regression
4. **Setup Requirements**: Any special data or conditions needed

### Example Regression Test
```go
// test/regression/bug_20250115_goldprompt_parsing_test.go
package regression

import (
    "testing"
    "github.com/bradhannah/Ultima5ReduxGo/internal/conversation"
)

// REGRESSION TEST: Bug discovered 2025-01-15
// 
// BUG DESCRIPTION: 
// GoldPrompt commands in TLK files were incorrectly using the Num field (always 0)
// instead of parsing the gold amount from the numeric prefix in the following 
// PlainString command (e.g., "005We thank thee." means 5 gold).
//
// ROOT CAUSE: 
// LinearConversationEngine processed GoldPrompt.Num field directly without
// checking for embedded amount in subsequent PlainString.
//
// WHY THIS TEST MATTERS:
// - Prevents incorrect gold deduction amounts in temple donations
// - Ensures proper TLK file format compatibility with original Ultima V
// - Validates correct parsing of complex command sequences
//
// TEST VALIDATES:
// 1. GoldPrompt with Num=0 correctly extracts amount from following PlainString
// 2. PlainString output has numeric prefix stripped
// 3. ActionCallbacks receives correct gold amount
func TestGoldPromptParsingRegression(t *testing.T) {
    // Setup: Create TLK script with problematic GoldPrompt sequence
    mockCallbacks := &MockActionCallbacks{}
    script := createMockTLKWithGoldPrompt(0, "005We thank thee.")
    engine := conversation.NewLinearConversationEngine(script, mockCallbacks)
    
    // Execute: Process the GoldPrompt command
    engine.Start(1)
    response := engine.ProcessInput("yes")
    
    // Verify: Gold amount extracted correctly (5, not 0)
    assert.Equal(t, 5, mockCallbacks.LastGoldAmount)
    assert.Equal(t, "We thank thee.", response.Output)
    
    // Additional validation: Ensure no "005" prefix in output
    assert.NotContains(t, response.Output, "005")
}
```

### Regression Test Categories

#### Data Parsing Regressions
```go
// Tests for file format parsing issues
func TestSavedGameLoadingRegression_20250110(t *testing.T) {
    // Test for SAVED.GAM character record parsing bug
}
```

#### Logic Flow Regressions  
```go
// Tests for game logic bugs
func TestNPCScheduleTransitionRegression_20250108(t *testing.T) {
    // Test for NPC getting stuck during floor transitions
}
```

#### UI/Input Regressions
```go
// Tests for user interface bugs  
func TestDebugConsoleAutoCompleteRegression_20250112(t *testing.T) {
    // Test for debug console command parsing edge case
}
```

## Property-Based Testing

For complex game systems, consider property tests:

```go
func TestPathfinding_Properties(t *testing.T) {
    properties := []struct {
        name string
        test func(*testing.T, references.Position, references.Position)
    }{
        {
            name: "path_exists_if_reachable",
            test: func(t *testing.T, start, goal references.Position) {
                if isReachable(start, goal) {
                    path, found := findPath(start, goal)
                    assert.True(t, found)
                    assert.Greater(t, len(path), 0)
                }
            },
        },
        {
            name: "path_starts_and_ends_correctly",
            test: func(t *testing.T, start, goal references.Position) {
                path, found := findPath(start, goal)
                if found {
                    assert.Equal(t, start, path[0])
                    assert.Equal(t, goal, path[len(path)-1])
                }
            },
        },
    }
    
    // Generate test cases
    for _, prop := range properties {
        t.Run(prop.name, func(t *testing.T) {
            for i := 0; i < 100; i++ {
                start := randomValidPosition()
                goal := randomValidPosition()
                prop.test(t, start, goal)
            }
        })
    }
}
```

## Performance Testing

### Benchmark Tests
```go
func BenchmarkConversationEngine_ProcessInput(b *testing.B) {
    engine := setupBenchmarkEngine()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine.ProcessInput("name")
    }
}

func BenchmarkPathfinding_LargeMap(b *testing.B) {
    mapData := loadLargeTestMap()
    pathfinder := astar.NewPathfinder(mapData)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        start := randomPosition()
        goal := randomPosition()
        pathfinder.FindPath(start, goal)
    }
}
```

## Running Tests

### Standard Test Commands
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./internal/conversation

# Run regression tests only
go test ./test/regression

# Run benchmarks
go test -bench=. ./internal/astar

# Run tests with race detection
go test -race ./...

# Verbose output for debugging
go test -v ./internal/conversation
```

### CI/CD Integration (GitHub Actions)
```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run Tests
        run: |
          go test -v -cover ./...
          go test -race ./...
      
      - name: Run Regression Tests
        run: go test -v ./test/regression
```

### Test Configuration
Create `testing.conf` for test-specific settings:
```yaml
# testing.conf
test:
  data_path: "testdata/"
  mock_seed: 12345
  timeout: "30s"
  
integration:
  ultima_v_path: "/opt/ultima5"  # CI environment path
  skip_if_missing: true
```

## Debugging Tests

### Useful Test Debugging Techniques

#### Test-Specific Debug Output
```go
func TestComplexFeature(t *testing.T) {
    if testing.Verbose() {
        t.Log("Detailed debugging information")
        t.Logf("State: %+v", gameState)
    }
}
```

#### Test Data Dumps
```go
func dumpGameStateForTest(t *testing.T, state *game_state.GameState) {
    t.Helper()
    if testing.Verbose() {
        t.Logf("Player Position: %+v", state.MapState.PlayerLocation)
        t.Logf("Game Time: %+v", state.DateTime)
    }
}
```

#### Conditional Test Skipping
```go
func TestExpensiveOperation(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping expensive test in short mode")
    }
    // Long-running test code
}
```

## Best Practices Summary

1. **Write Tests for Bug Fixes**: Every bug fix should include a regression test
2. **Use Test Helpers**: Create reusable utilities to reduce boilerplate
3. **Test with Real Data**: Integration tests should use actual game data files
4. **Document Regression Tests**: Explain why each regression test exists
5. **Use Fixed Seeds**: Ensure deterministic behavior in tests with randomness
6. **Leverage @OLD**: Reference fan remake code when behavior is unclear
7. **Benchmark Critical Paths**: Performance test pathfinding, conversation processing
8. **Test Edge Cases**: Consider boundary conditions and invalid inputs
9. **Mock External Dependencies**: Use interfaces and mocks for testability
10. **Run Tests in CI**: Automate testing in GitHub Actions

This testing approach ensures reliable, maintainable code while supporting the rapid iteration needed for game development. The combination of unit tests, integration tests with real data, and comprehensive regression testing provides confidence in both current functionality and protection against future regressions.