# Tile System and Conventions

## Overview

The tile system in Ultima5ReduxGo provides a structured way to identify and interact with different tile types in the game world. Tiles represent terrain, objects, and structures that make up the game maps.

## File Structure

- **Tile Definition**: `internal/references/tile.go` - Core tile struct and identification methods
- **Sprite Indexes**: `internal/sprites/indexes/sprites.go` - Tile sprite constants and index methods
- **Usage Examples**: Various files including `internal/references/enemy.go`

## Tile Identification Conventions

### Is* Methods Pattern

All tile type identification follows the `Is*()` method pattern on the `Tile` struct:

```go
func (t *Tile) IsTileType() bool {
    return t.Index == indexes.TileTypeConstant
}
```

### Terrain Type Methods

**Water Tiles**:
```go
func (t *Tile) IsWater() bool {
    return t.Index == indexes.Water1 || t.Index == indexes.Water2 || t.Index == indexes.WaterShallow
}
```

**Desert/Sand Tiles**:
```go
func (t *Tile) IsDesert() bool {
    return t.Index == indexes.Desert || t.Index == indexes.LeftDesert2 || t.Index == indexes.RightDesert2
}
```

**Swamp Tiles**:
```go
func (t *Tile) IsSwamp() bool {
    return t.Index == indexes.Swamp
}
```

**Mountain Tiles**:
```go
func (t *Tile) IsMountain() bool {
    return t.Index == indexes.SmallMountains
}
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
    return t.Index == indexes.LargeRockWall || 
        t.Index == indexes.StoneBrickWall || 
        t.Index == indexes.StoneBrickWallSecret
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

## Usage Patterns

### Monster Spawn Logic

Replace string-based tile identification:

```go
// ❌ AVOID: Hard-coded string comparison
if strings.HasPrefix(strings.ToLower(tile.Name), "sand") {
    // logic
}

// ✅ PREFERRED: Type-safe method
if tile.IsDesert() {
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
```

## Tile Properties

Beyond identification, tiles have various boolean properties:

### Movement Properties
- `IsWalkingPassable` - Can be walked on by player
- `IsBoatPassable` - Can be traversed by boat
- `IsHorsePassable` - Can be traversed by horse
- `IsCarpetPassable` - Can be traversed by magic carpet
- `IsLandEnemyPassable` - Can be occupied by land enemies
- `IsWaterEnemyPassable` - Can be occupied by water enemies

### Interaction Properties
- `IsOpenable` - Can be opened (doors, chests)
- `IsPushable` - Can be pushed/moved
- `IsTalkOverable` - Can talk over this tile
- `IsBoardable` - Can board vehicles on this tile
- `IsKlimable` - Can be climbed

### Visual Properties
- `BlocksLight` - Blocks line of sight
- `LightEmission` - Emits light (value indicates intensity)
- `IsUpright` - Rendered as upright sprite
- `DontDraw` - Should not be rendered

## Best Practices

1. **Always use Is* methods instead of string comparisons**
2. **Group related tile types in single methods when appropriate**
3. **Leverage existing sprite index constants**
4. **Follow the established naming convention: IsXxx()**
5. **Handle multiple sprite indexes for the same logical tile type**
6. **Use tile properties for behavior logic, Is* methods for type identification**

## Adding New Tile Types

When adding new tile identification methods:

1. Check if sprite indexes exist in `sprites.go`
2. **If missing tile constants are found, simply add them and proceed implementing your feature**
3. Add the Is* method to `tile.go` following the pattern
4. Handle all related sprite variations (facing directions, animation frames)
5. Update this documentation
6. Replace any existing string-based identification with the new method

## Integration with Game Systems

The tile system integrates with:
- **Monster Generation**: Environment-based spawning
- **Pathfinding**: Movement cost and validity
- **Combat**: Line of sight and positioning
- **Interaction**: Object manipulation and use
- **Rendering**: Visual representation and effects