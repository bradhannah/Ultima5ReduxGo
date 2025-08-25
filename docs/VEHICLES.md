# Vehicle System

## Overview

The vehicle system in Ultima5ReduxGo manages player transportation including horses, magic carpets, boats (skiffs and frigates), and NPC vehicles. The system handles boarding, exiting, movement sprites, and vehicle-specific mechanics.

## Core Types

### VehicleType Enum

```go
type VehicleType int

const (
    NoPartyVehicle VehicleType = iota
    CarpetVehicle
    HorseVehicle
    SkiffVehicle
    FrigateVehicle
    NPC
)
```

### Key Methods

- `GetUnBoardedSpriteByDirection(direction)` - Returns sprite for unboarded vehicle
- `GetBoardedSpriteByDirection(previousDirection, direction)` - Returns sprite for player on vehicle
- `RequiresNewSprite(currentDirection, newDirection)` - Determines if sprite update needed
- `GetMovementPrefix()` - Returns movement text prefix ("Fly", "Ride", "Row", etc.)
- `GetLowerCaseName()` - Returns vehicle name for messages
- `GetExitString()` / `GetBoardString()` - Returns boarding/exiting messages

## Vehicle Mechanics

### Movement Prefixes

| Vehicle | Prefix | Example |
|---------|--------|---------|
| Carpet  | "Fly " | "Fly North" |
| Horse   | "Ride " | "Ride South" |
| Skiff   | "Row " | "Row East" |
| Frigate | "Row " | "Row West" |
| On foot | "" | "North" |

### Sprite Management

#### Direction-Based Sprites
- **Horses/Carpets**: Only have left/right sprites, maintain previous direction for vertical movement
- **Ships**: Have sprites for all four directions (Up/Down/Left/Right)
- **Skiffs**: Have sprites for all four directions

#### Sprite States
- **Furled vs Unfurled**: Frigates have both furled (anchored) and unfurled (sailing) sprites
- **Boarded vs Unboarded**: Separate sprites for empty vehicles and vehicles with player

### Boarding Rules

#### Horse Boarding
- Must be on foot
- In towns, checks if NPC refuses mount (`npc_refuses_mount()`)
- Sets appropriate riding sprite based on direction

#### Carpet Boarding
- Must be on foot
- Cannot board on mountains (overworld restriction)
- Cannot board while on ship ("X-it ship first!")

#### Ship Boarding (Skiffs/Frigates)
- Must be on foot or have compatible vehicle (carpet can board ship)
- Ship condition warnings:
  - Hull < 10: "DANGER: SHIP BADLY DAMAGED!"
  - No skiffs: "WARNING: NO SKIFFS ON BOARD!"
- Transfers skiff/carpet counts to ship inventory

### Vehicle Interactions

#### Board Command Pattern
```go
// In GameState - handles all logic and messaging via SystemCallbacks
func (g *GameState) ActionBoard() bool {
    // Validation, vehicle detection, boarding logic
    // All messages sent via g.SystemCallbacks.Message
    return success
}

// In GameScene input handler - direct call to GameState
case ebiten.KeyB:
    g.gameState.ActionBoard()
```

#### Exit Command Pattern  
```go
// In GameState - handles all logic and messaging via SystemCallbacks
func (g *GameState) ActionExit() bool {
    // Exit logic, vehicle-specific messaging
    // All messages sent via g.SystemCallbacks.Message
    return success
}

// In GameScene input handler - direct call to GameState
case ebiten.KeyX:
    g.gameState.ActionExit()
```

## Vehicle Detection

### From Sprite Index
```go
func GetVehicleFromSpriteIndex(s indexes.SpriteIndex) VehicleType {
    switch s {
    case indexes.Carpet2_MagicCarpet:
        return CarpetVehicle
    case indexes.HorseRight, indexes.HorseLeft:
        return HorseVehicle
    case indexes.FrigateDownFurled, indexes.FrigateUpFurled, 
         indexes.FrigateLeftFurled, indexes.FrigateRightFurled:
        return FrigateVehicle
    case indexes.SkiffLeft, indexes.SkiffRight, 
         indexes.SkiffUp, indexes.SkiffDown:
        return SkiffVehicle
    default:
        return NoPartyVehicle
    }
}
```

## Integration Points

### Game State Integration
- `BoardVehicle(vehicle)` - Handles boarding logic and state changes
- `ExitVehicle()` - Handles exiting and returns exited vehicle
- Vehicle state affects movement mechanics and available actions

### NPC System Integration
- `GetVehicleAtPositionOrNil(position)` - Finds vehicles at given position
- Horses can be NPC-owned and may refuse mounting

### Map System Integration
- Vehicle passability checks for different terrain types
- Context-specific restrictions (e.g., no carpets on mountains)
- Map transitions preserve vehicle state

## Special Cases

### Ship Repair (Hole Up & Camp)
- Only available when aboard ship with furled sails
- Repairs hull points over time with monster movement
- Must remain aboard during repair process

### Combat Integration
- Vehicle state affects available actions in combat
- Some vehicles provide protection or movement advantages
- Vehicle damage affects functionality

### Sailing Mechanics
- Frigates can unfurl/furl sails for different movement modes
- Wind affects sailing movement (handled by movement system)
- Anchored vs sailing states affect available actions

## Error Handling

### Common Error Messages
- "Board what?" - No vehicle at position
- "X-it what?" - Not in/on a vehicle  
- "On foot!" - Action requires being on foot
- "X-it ship first!" - Cannot use carpet while on ship
- "Not here!" - Action not allowed in current context
- "Nay!" - NPC refuses to let player mount horse

### State Validation
- Vehicle system validates state transitions
- Prevents invalid combinations (e.g., boarding while already on vehicle)
- Ensures proper cleanup when switching vehicles

## Implementation Files

- `internal/references/vehicles.go` - Core vehicle type definitions and logic
- `internal/game_state/action_board.go` - ActionBoard() implementation
- `internal/game_state/action_exit.go` - ActionExit() and ExitVehicle() implementations  
- `cmd/ultimav/gamescene_input_*.go` - Input handlers that call GameState actions directly
- `internal/sprites/indexes/sprites.go` - Vehicle sprite constants

## Usage Patterns

### Adding New Vehicle Types
1. Add to `VehicleType` enum
2. Implement sprite mapping in `GetUnBoardedSpriteByDirection()` and `GetBoardedSpriteByDirection()`
3. Add movement prefix and name methods
4. Update `GetVehicleFromSpriteIndex()` detection
5. Add any special boarding/exiting rules

### Vehicle State Checks
```go
// Check if player can perform action based on vehicle state
if g.gameState.PartyState.CurrentVehicle != references.NoPartyVehicle {
    // Handle vehicle-specific logic
}

// Get current vehicle for context-sensitive behavior
currentVehicle := g.gameState.GetCurrentVehicleType()
```

This system provides a flexible foundation for all vehicle-related mechanics while maintaining clean separation between vehicle logic, game state, and UI presentation.