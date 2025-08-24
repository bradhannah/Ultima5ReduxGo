# Time of Day System

This document describes the comprehensive time management system used in the Ultima 5 Redux Go project, centered around the `UltimaDate` struct and `TimeOfDay` enum.

## Architecture Overview

The time system is located in `internal/datetime/` and provides authentic Ultima V time mechanics including turn-based progression, day/night cycles, visibility calculations, and era tracking.

### Core Components

- **`UltimaDate`** - Complete date/time state with turn counter
- **`TimeOfDay` enum** - Discrete time periods (Morning, Noon, Evening, Midnight, Dusk)
- **`Era` system** - Long-term game progression tracking
- **Turn counter** - Increments with every time advancement

## UltimaDate Structure

```go
type UltimaDate struct {
    Year   uint16  // Game year
    Month  byte    // 1-12 (MonthsPerYear = 12)
    Day    byte    // 1-28 (DaysInMonth = 28)
    Hour   byte    // 0-23 (HoursPerDay = 24)
    Minute byte    // 0-59 (MinutesPerHour = 60)
    Turn   uint32  // Increments every time advancement
}
```

### Calendar System
- **12 months per year** - Standard calendar structure
- **28 days per month** - Consistent month length 
- **24 hours per day** - Standard day structure
- **60 minutes per hour** - Standard time granularity

## TimeOfDay Enum

```go
type TimeOfDay int

const (
    Morning  TimeOfDay = iota  // 5:00 AM
    Noon                       // 12:00 PM  
    Evening                    // 5:00 PM
    Midnight                   // 12:00 AM
    Dusk                       // 8:00 PM
)
```

### Time Period Mapping
- **Morning**: 5:00 AM - Dawn break, NPC schedule transitions
- **Noon**: 12:00 PM - Full daylight, maximum visibility
- **Evening**: 5:00 PM - Afternoon activities, pre-sunset
- **Dusk**: 8:00 PM - Sunset transition, reduced visibility
- **Midnight**: 12:00 AM - Deep night, minimum visibility

## Core Methods

### Time Advancement

```go
func (d *UltimaDate) Advance(nMinutes int)
```

**Features:**
- **Turn counter increment** - Every advancement increases turn count
- **Cascading updates** - Minutes → Hours → Days → Months → Years
- **Safety limits** - Maximum 9 hours advancement per call
- **Proper wrapping** - Handles month/year boundaries correctly

**Usage Examples:**
```go
gameState.DateTime.Advance(15)  // Advance 15 minutes
gameState.DateTime.Advance(120) // Advance 2 hours
```

### Time Setting

```go
func (d *UltimaDate) SetTimeOfDay(timeOfDay TimeOfDay)
```

**Purpose:** Jump to specific time periods (resets minutes to 0)

**Usage Examples:**
```go
gameState.DateTime.SetTimeOfDay(datetime.Morning)  // Jump to 5:00 AM
gameState.DateTime.SetTimeOfDay(datetime.Dusk)     // Jump to 8:00 PM
```

### Visibility System

```go
func (d *UltimaDate) GetVisibilityFactorWithoutTorch(baselineMin float32) float32
```

**Features:**
- **Dynamic dawn/dusk transitions** - Gradual visibility changes
- **Time-sensitive calculations** - Based on hour and minute
- **Baseline minimum** - Prevents complete darkness
- **Smooth interpolation** - Realistic light transitions

**Visibility Periods:**
- **Hour 5 (Dawn)**: Gradual 0.1 → 1.0 visibility increase
- **Hours 6-19**: Full daylight (1.0 visibility)  
- **Hour 20 (Dusk)**: Gradual 1.0 → 0.1 visibility decrease
- **Hours 21-4**: Night minimum visibility (baseline)

### Daylight Detection

```go
func (d *UltimaDate) IsDayLight() bool
```

**Purpose:** Simple day/night detection for game mechanics
**Logic:** Returns true between hours 6-20

## Era System Integration

### Era Progression
```go
type Era int
const (
    EarlyEra  = iota  // Turns 0-9,999
    MiddleEra         // Turns 10,000-29,999  
    LateEra           // Turns 30,000+
)
```

**Purpose:** Long-term game progression tracking
- **Early Era**: Game beginning, basic mechanics
- **Middle Era**: Advanced gameplay, more complex systems
- **Late Era**: End-game content, maximum complexity

## Integration with Game Systems

### GameState Integration

```go
type GameState struct {
    // ... other fields
    DateTime datetime.UltimaDate
}
```

**Usage in Game Logic:**
```go
// Check time for NPC schedules
if g.DateTime.Hour >= 8 && g.DateTime.Hour <= 17 {
    // Shopkeeper is open
}

// Visibility-dependent actions
visibilityFactor := g.DateTime.GetVisibilityFactorWithoutTorch(0.1)
if visibilityFactor < 0.3 {
    g.SystemCallbacks.Message.AddRowStr("It's too dark to see!")
    return false
}

// Era-specific content
if g.DateTime.GetEra() == datetime.LateEra {
    // Enable late-game mechanics
}
```

### Turn-Based Mechanics

Every game action that consumes time should:
1. Call `DateTime.Advance(minutes)` 
2. Trigger dependent systems (NPC schedules, visibility updates)
3. Update turn-sensitive mechanics

**Common Time Costs:**
```go
// Movement
g.DateTime.Advance(1)    // 1 minute per tile

// Combat actions  
g.DateTime.Advance(1)    // 1 minute per combat turn

// Resting
g.DateTime.Advance(480)  // 8 hours sleep

// Searching
g.DateTime.Advance(5)    // 5 minutes per search
```

## Display Formatting

### Date Display
```go
func (d *UltimaDate) GetDateAsString() string
// Returns: "3-15-142" (Month-Day-Year)
```

### Time Display  
```go
func (d *UltimaDate) GetTimeAsString() string
// Returns: " 9:45AM" or "11:30PM" (12-hour format with AM/PM)
```

**Format Features:**
- **12-hour format** - User-friendly time display
- **AM/PM indicators** - Clear morning/evening distinction  
- **Consistent spacing** - Aligned display formatting
- **Zero-padded minutes** - Professional time formatting

## System Dependencies

### Current Integration Points

**Game State:** 
- `GameState.DateTime` - Primary time state storage
- Turn advancement integrated with game actions

**Visibility System:**
- Light level calculations for gameplay
- Torch mechanics and visibility modifiers
- Day/night dependent actions

**NPC Systems:**
- Schedule-based NPC behavior
- Time-dependent availability
- Era-based NPC reactions

### Future Integration Opportunities

**Audio System:**
- Time-of-day ambient sounds
- Era-specific music tracks
- Clock chime effects

**Visual System:**  
- Day/night lighting effects
- Time-based color palettes
- Era-specific visual themes

**Game Flow:**
- Automatic time advancement
- Sleep/rest mechanics
- Time-sensitive events

## Usage Patterns

### Time-Dependent Actions

```go
func (g *GameState) ActionSearchSmallMap(direction references.Direction) bool {
    // Check visibility first
    if !g.DateTime.IsDayLight() {
        visibilityFactor := g.DateTime.GetVisibilityFactorWithoutTorch(0.1)
        if visibilityFactor < 0.5 {
            g.SystemCallbacks.Message.AddRowStr("Too dark to search effectively!")
            return false
        }
    }
    
    // Perform search
    // ... search logic ...
    
    // Advance time
    g.DateTime.Advance(5) // 5 minutes to search
    return true
}
```

### NPC Schedule Integration

```go
func (g *GameState) IsShopkeeperAvailable(shopId int) bool {
    hour := g.DateTime.Hour
    
    // Most shops open 8 AM to 6 PM
    if hour >= 8 && hour <= 18 {
        return true
    }
    
    // Taverns open later
    if shopId == TAVERN && hour >= 6 && hour <= 23 {
        return true  
    }
    
    return false
}
```

### Era-Based Content

```go
func (g *GameState) GetAvailableSpells() []Spell {
    era := g.DateTime.GetEra()
    spells := getBasicSpells()
    
    if era >= datetime.MiddleEra {
        spells = append(spells, getAdvancedSpells()...)
    }
    
    if era >= datetime.LateEra {
        spells = append(spells, getMasterSpells()...)
    }
    
    return spells
}
```

## Constants and Configuration

```go
const (
    DaysInMonth    = 28
    MonthsPerYear  = 12  
    MinutesPerHour = 60
    HoursPerDay    = 24
    
    hourOfSunrise = 5   // Dawn visibility transition
    hourOfSunset  = 20  // Dusk visibility transition
    
    beginningOfEra1 = 0      // Early Era start
    beginningOfEra2 = 10000  // Middle Era start  
    beginningOfEra3 = 30000  // Late Era start
)
```

## Testing Considerations

### Unit Tests
```go
func TestTimeAdvancement(t *testing.T) {
    date := datetime.UltimaDate{Hour: 23, Minute: 45}
    date.Advance(30) // Should wrap to next day
    assert.Equal(t, 0, date.Hour)
    assert.Equal(t, 15, date.Minute)
}

func TestVisibilityTransitions(t *testing.T) {
    date := datetime.UltimaDate{Hour: 5, Minute: 30} // Dawn
    visibility := date.GetVisibilityFactorWithoutTorch(0.1)
    assert.Greater(t, visibility, 0.5) // Partial daylight
}
```

### Integration Tests
- **Time-dependent actions** - Verify correct time advancement
- **Visibility mechanics** - Test day/night action restrictions  
- **Era transitions** - Validate content availability by era
- **NPC schedules** - Confirm time-based availability

## File Locations

- **Core Implementation**: `internal/datetime/ultima_date.go`
- **Era System**: `internal/datetime/era.go`  
- **Game State Integration**: `internal/game_state/game_state.go`
- **Documentation**: `docs/TIME_OF_DAY.md` (this file)