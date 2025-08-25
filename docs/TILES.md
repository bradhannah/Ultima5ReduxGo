# Tile System and Conventions

## Overview

The tile system in Ultima5ReduxGo provides a structured way to identify and interact with different tile types in the game world. Tiles represent terrain, objects, and structures that make up the game maps.

## File Structure

- **Tile Definition**: `internal/references/tile.go` - Core tile struct and identification methods
- **Sprite Indexes**: `internal/sprites/indexes/sprites.go` - Tile sprite constants and index methods
- **Usage Examples**: Various files including `internal/references/enemy.go`

## Tile Identification Conventions

### Is* Methods Pattern

Tile type identification uses two patterns on the `Tile` struct:

**Single Tile Check (Generic Pattern)**:
```go
func (t *Tile) Is(spriteIndex indexes.SpriteIndex) bool {
    return t.Index == spriteIndex
}

// Usage: tile.Is(indexes.Grass), tile.Is(indexes.Barrel)
```

**Logical Groupings (Specific Methods)**:
```go
func (t *Tile) IsTileGroup() bool {
    return t.Is(indexes.Type1) || t.Is(indexes.Type2) || /* other variants */
}
```

### Terrain Type Methods

**Water Tiles (Multiple Variants)**:
```go
func (t *Tile) IsWater() bool {
    return t.Is(indexes.Water1) || t.Is(indexes.Water2) || t.Is(indexes.WaterShallow)
}
```

**Desert/Sand Tiles (Multiple Variants)**:
```go
func (t *Tile) IsDesert() bool {
    return t.Is(indexes.Desert) || t.Is(indexes.LeftDesert2) || t.Is(indexes.RightDesert2)
}
```

**Swamp Tiles (Single Tile - Use Generic Pattern)**:
```go
func (t *Tile) IsSwamp() bool {
    return t.Is(indexes.Swamp)
}

// Or directly: tile.Is(indexes.Swamp)
```

**Mountain Tiles (Single Tile - Use Generic Pattern)**:
```go
func (t *Tile) IsMountain() bool {
    return t.Is(indexes.SmallMountains)
}

// Or directly: tile.Is(indexes.SmallMountains)
```

**Path/Road Tiles**:
```go
func (t *Tile) IsRoad() bool {
    return t.IsPath() // Leverages existing IsPath() method
}

func (t *Tile) IsPath() bool {
    return t.Index >= indexes.PathUpDown && t.Index <= indexes.PathAllWays
}
```

### Object Type Methods

**Furniture**:
```go
func (t *Tile) IsChair() bool {
    return t.Index == indexes.ChairFacingDown ||
        t.Index == indexes.ChairFacingUp ||
        t.Index == indexes.ChairFacingRight ||
        t.Index == indexes.ChairFacingLeft
}

func (t *Tile) IsCannon() bool {
    return t.Index == indexes.CannonFacingLeft ||
        t.Index == indexes.CannonFacingRight ||
        t.Index == indexes.CannonFacingUp ||
        t.Index == indexes.CannonFacingDown
}
```

**Structural Elements**:
```go
func (t *Tile) IsWall() bool {
    return t.Is(indexes.LargeRockWall) || 
        t.Is(indexes.StoneBrickWall) || 
        t.Is(indexes.StoneBrickWallSecret)
}
```

### Sprite Index Methods

The `SpriteIndex` type also provides identification methods:

```go
func (s SpriteIndex) IsDoor() bool
func (s SpriteIndex) IsUnlockedDoor() bool
func (s SpriteIndex) IsPushableFloor() bool
func (s SpriteIndex) IsBed() bool
func (s SpriteIndex) IsStairs() bool
```

## Coding Conventions

### When to Use Generic `Is()` vs Specific Methods

**❌ AVOID: Creating individual functions for single tile checks**
```go
// Don't create these:
func (t *Tile) IsCactus() bool { return t.Index == indexes.Cactus }
func (t *Tile) IsTree() bool { return t.Index == indexes.Tree }
func (t *Tile) IsRock() bool { return t.Index == indexes.Rock }
```

**✅ PREFERRED: Use generic Is() for single tiles**
```go
// Use the generic pattern:
if tile.Is(indexes.Cactus) { /* handle cactus */ }
if tile.Is(indexes.Tree) { /* handle tree */ }
if tile.Is(indexes.Rock) { /* handle rock */ }
```

**✅ GOOD: Create specific functions for logical groupings**
```go
// For multiple variants of the same concept:
func (t *Tile) IsDoor() bool {
    return t.Index.IsDoor() // Uses sprite index method
}

func (t *Tile) IsChair() bool {
    return t.Index == indexes.ChairFacingDown ||
        t.Index == indexes.ChairFacingUp ||
        t.Index == indexes.ChairFacingRight ||
        t.Index == indexes.ChairFacingLeft
}

func (t *Tile) IsPushable() bool {
    return t.IsChair() || t.IsCannon() || // Logical groupings
           t.Is(indexes.Barrel) || t.Is(indexes.Chest) // Single tiles
}
```

### Usage Patterns

**Monster Spawn Logic**:

```go
// ❌ AVOID: Hard-coded string comparison
if strings.HasPrefix(strings.ToLower(tile.Name), "sand") {
    // logic
}

// ✅ PREFERRED: Type-safe method
if tile.IsDesert() {
    // logic
}

// ✅ ALSO GOOD: Direct generic check for single tiles
if tile.Is(indexes.Grass) {
    // logic
}
```

### Environment Detection

```go
func determineEnvironmentType(tile *Tile) MonsterEnvironment {
    if tile.IsWater() {
        return WaterEnvironment
    }
    if tile.IsDesert() {
        return DesertEnvironment
    }
    return LandEnvironment
}
```

### Movement and Pathfinding

```go
func calculateTileBasedProbability(tile *Tile) int {
    if tile.IsRoad() {
        return 0 // No monsters on roads
    }
    
    if tile.IsSwamp() || tile.IsForest() || tile.IsMountain() {
        return 2 // Higher spawn probability
    }
    
    return 1 // Default probability
}

// Modern function-based approach:
func (t *Tile) IsWalkingPassable() bool {
    return t.Is(indexes.Grass) ||
        t.IsPath() ||
        t.Is(indexes.BrickFloor) ||
        t.IsDesert() ||
        t.IsSwamp() ||
        (!t.IsMountain() && !t.IsWater() && !t.IsWall())
}
```

## Tile Properties

Tile properties are implemented as function-based logic rather than data fields:

### Movement Properties (Functions)
- `IsWalkingPassable()` - Can be walked on by player (logic-based)
- `IsBoatPassable()` - Can be traversed by boat (logic-based)
- `IsHorsePassable()` - Can be traversed by horse (logic-based)
- `IsCarpetPassable()` - Can be traversed by magic carpet (logic-based)
- `IsLandEnemyPassable()` - Can be occupied by land enemies (logic-based)
- `IsWaterEnemyPassable()` - Can be occupied by water enemies (logic-based)
- `IsRangeWeaponPassable()` - Can ranged weapons pass through (logic-based)

### Interaction Properties (Functions)
- `IsOpenable()` - Can be opened (doors, chests) (logic-based)
- `IsPushable()` - Can be pushed/moved (logic-based)
- `IsKlimable()` - Can be climbed (logic-based)

### Visual Properties (Data Fields)
- `BlocksLight` - Blocks line of sight
- `LightEmission` - Emits light (value indicates intensity)
- `DontDraw` - Should not be rendered
- `IsPartOfAnimation` - Part of animated sequence

### Example Function-Based Properties

```go
func (t *Tile) IsWalkingPassable() bool {
    // Logic determines walkability based on tile type
    return t.Is(indexes.Grass) ||
           t.IsPath() ||
           t.IsDesert() ||
           (!t.IsMountain() && !t.IsWater() && !t.IsWall())
}

func (t *Tile) IsOpenable() bool {
    // Doors and chests can be opened
    return t.Index.IsDoor() || t.Is(indexes.Chest)
}

func (t *Tile) IsPushable() bool {
    // Combines logical groupings with single tile checks
    return t.IsChair() || t.IsCannon() || // Logical groupings
           t.Is(indexes.Barrel) || t.Is(indexes.Chest) // Single tiles
}
```

## Best Practices

1. **Use generic `Is()` for single tile checks** - `tile.Is(indexes.Grass)` not `tile.IsGrass()`
2. **Create specific methods only for logical groupings** - `IsChair()` for multiple chair orientations
3. **Prefer existing helper functions** - `tile.IsDesert()` over `tile.Is(indexes.Desert)` when available
4. **Use function-based properties over data fields** - Implement logic in methods for better maintainability
5. **Always use Is* methods instead of string comparisons**
6. **Leverage existing sprite index constants and methods**
7. **Follow the established naming convention: IsXxx()**
8. **Handle multiple sprite indexes for the same logical tile type**

## Adding New Tile Types

When adding new tile identification:

**For Single Tiles:**
1. Check if sprite indexes exist in `sprites.go`
2. **If missing tile constants are found, simply add them and proceed implementing your feature**
3. **Use the generic `Is()` pattern directly** - no need to create new methods
4. Example: `tile.Is(indexes.NewTileType)`

**For Logical Groupings Only:**
1. Create a specific `IsXxx()` method when you have multiple variants of the same concept
2. Handle all related sprite variations (facing directions, animation frames)
3. Use the `Is()` pattern within the method for consistency
4. Example:
   ```go
   func (t *Tile) IsNewFurniture() bool {
       return t.Is(indexes.NewFurnitureType1) ||
              t.Is(indexes.NewFurnitureType2) ||
              t.Is(indexes.NewFurnitureType3)
   }
   ```
5. Update this documentation
6. Replace any existing string-based identification with the new pattern

## Pattern Summary

**Use Generic `Is()` Pattern For:**
- Single tile checks: `tile.Is(indexes.Grass)`, `tile.Is(indexes.Barrel)`
- One-off comparisons in algorithms
- Simple boolean checks

**Create Specific Methods For:**
- Multiple variants of same concept: `IsChair()` (4 orientations)
- Complex terrain logic: `IsWalkingPassable()` (multiple conditions)
- Logical groupings: `IsPushable()` (chairs + cannons + individual items)
- Commonly used combinations: `IsWater()` (3 water types)

**Modern Function-Based Approach:**
- Replace data fields with logic functions where possible
- Use `t.Index.IsDoor()` for sprite index methods
- Combine patterns: `t.IsChair() || t.Is(indexes.Barrel)`
- Maintain consistent coding style throughout

## Integration with Game Systems

The tile system integrates with:
- **Monster Generation**: Environment-based spawning using terrain methods
- **Pathfinding**: Movement cost and validity using passability functions
- **Combat**: Line of sight and positioning using wall/obstacle detection
- **Interaction**: Object manipulation using pushable/openable detection
- **Rendering**: Visual representation and effects using property fields