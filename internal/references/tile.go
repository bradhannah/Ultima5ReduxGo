package references

import (
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

const (
	weightImpassable = -1
	weightIdealPath  = 1
	weightPath       = 2
	weightGrass      = 3
	weightDefault    = 10
)

// Tile represents a single tile and its properties.
//
// The Tile struct contains both static properties (loaded from JSON) and methods
// for tile-specific behavior checks. Methods that return boolean values based solely
// on tile type/index (like IsPushable, IsOpenable, IsKlimable) should be placed here
// rather than in external structs like GameState, as these are intrinsic properties
// of the tile itself.
//
//nolint:tagliatelle
type Tile struct {
	Index                indexes.SpriteIndex
	Name                 string `json:"Name"`
	Description          string `json:"Description"`
	SpeedFactor          int    `json:"SpeedFactor"`
	LightEmission        int    `json:"LightEmission"`
	IsPartOfAnimation    bool   `json:"IsPartOfAnimation"`
	TotalAnimationFrames int    `json:"TotalAnimationFrames"`
	AnimationIndex       int    `json:"AnimationIndex"`
	IsEnemy              bool   `json:"IsEnemy"`
	IsNPC                bool   `json:"IsNPC"`
	IsBuilding           bool   `json:"IsBuilding"`
	DontDraw             bool   `json:"DontDraw"`
	IsGuessableFloor     bool   `json:"IsGuessableFloor"`
	BlocksLight          bool   `json:"BlocksLight"`
	IsWindow             bool   `json:"IsWindow"`
	CombatMapIndex       string `json:"CombatMapIndex"`
}

func (t *Tile) IsPassable(vehicle VehicleType) bool {
	switch vehicle {
	case CarpetVehicle:
		return t.IsCarpetPassable()
	case HorseVehicle:
		return t.IsHorsePassable()
	case SkiffVehicle:
		return t.IsSkiffPassable()
	case FrigateVehicle:
		return t.IsBoatPassable()
	case NoPartyVehicle:
		return t.IsWalkingPassable()
	case NPC:
		return t.IsLandEnemyPassable()
	}

	return false
}

// Is provides a generic method for checking if a tile matches a specific sprite index
func (t *Tile) Is(spriteIndex indexes.SpriteIndex) bool {
	return t.Index == spriteIndex
}

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

func (t *Tile) IsBarrel() bool {
	return t.Is(indexes.Barrel)
}

func (t *Tile) IsPath() bool {
	return t.Index >= indexes.PathUpDown && t.Index <= indexes.PathAllWays
}

func (t *Tile) GetStairsFloorDirection() LadderOrStairType {
	switch t.Index {
	case indexes.Stairs1, indexes.Stairs2:
		return LadderOrStairUp
	case indexes.Stair3, indexes.Stairs4:
		return LadderOrStairDown
	default:
		return NotLadderOrStair
	}
}

func (t *Tile) isNPCNoPenaltyWalkable() bool {
	return t.Is(indexes.BrickFloor) || t.Is(indexes.HexMetalGridFloor) || t.Is(indexes.WoodenPlankVert1Floor) || t.Is(indexes.WoodenPlankVert2Floor) || t.Is(indexes.WoodenPlankHorizFloor)
}

func (t *Tile) GetWalkableWeight() int {
	if t.Index.IsUnlockedDoor() {
		return weightIdealPath
	}

	if !t.IsWalkingPassable() {
		return weightImpassable
	}

	if t.isNPCNoPenaltyWalkable() {
		return weightIdealPath
	}

	if t.IsPath() {
		return weightPath
	}

	if t.Is(indexes.Grass) {
		return weightGrass
	}

	return weightDefault
}

func (t *Tile) IsWalkableDuringWander() bool {
	return t.IsWalkingPassable() && !t.Index.IsBed() && !t.Index.IsDoor()
}

func (t *Tile) GetExtraMovementString() string {
	switch t.SpeedFactor {
	case 4:
		return "Slow Progress!"
	case 6:
		return "Very Slow!"
	case 1, 2, -1:
		return ""
	default:
		return "Untrodden Combat Tile"
	}
}

func (t *Tile) IsWall() bool {
	return t.Is(indexes.LargeRockWall) || t.Is(indexes.StoneBrickWall) || t.Is(indexes.StoneBrickWallSecret)
}

func (t *Tile) IsRoad() bool {
	return t.IsPath() // Roads are paths in this context
}

func (t *Tile) IsSwamp() bool {
	return t.Is(indexes.Swamp)
}

func (t *Tile) IsWater() bool {
	return t.Is(indexes.Water1) || t.Is(indexes.Water2) || t.Is(indexes.WaterShallow)
}

func (t *Tile) IsDesert() bool {
	return t.Is(indexes.Desert) || t.Is(indexes.LeftDesert2) || t.Is(indexes.RightDesert2)
}

func (t *Tile) IsMountain() bool {
	return t.Is(indexes.SmallMountains)
}

func (t *Tile) IsForest() bool {
	// Forest tiles are identified as passable land tiles that are not other terrain types
	// This logic may need refinement based on actual forest tile indexes
	return t.IsLandEnemyPassable() &&
		t.Index != indexes.Grass &&
		t.Index != indexes.Desert &&
		t.Index != indexes.Swamp &&
		!t.IsPath() &&
		!t.IsMountain()
}

// IsPushable checks if this tile can be pushed/moved by the player.
// This is an intrinsic property of the tile type, not dependent on game state.
func (t *Tile) IsPushable() bool {
	// Chair variants (logical grouping - keep IsChair() for multiple chairs)
	if t.IsChair() {
		return true
	}

	// Cannon variants (logical grouping - keep IsCannon() for multiple cannons)
	if t.IsCannon() {
		return true
	}

	// Single tile checks using generic Is() pattern - based on original game data
	return t.Is(indexes.Barrel) ||
		t.Is(indexes.EndTable) ||
		t.Is(indexes.Vanity) ||
		t.Is(indexes.WaterJugTable) ||
		t.Is(indexes.Dresser) ||
		t.Is(indexes.Box) ||
		t.Is(indexes.Plant) ||
		// Additional pushable items (extended beyond original data)
		t.Is(indexes.TableMiddle) ||
		t.Is(indexes.TableFoodTop) ||
		t.Is(indexes.TableFoodBottom) ||
		t.Is(indexes.TableFoodBoth) ||
		t.Is(indexes.Mirror) ||
		t.Is(indexes.Well) ||
		t.Is(indexes.Brazier) ||
		t.Is(indexes.CookStove) ||
		t.Is(indexes.Chest) ||
		t.Is(indexes.WoodenBox)
}

// IsBoatPassable returns true if boats (frigates) can pass through this tile
func (t *Tile) IsBoatPassable() bool {
	// Boats need deep water
	return t.Is(indexes.Water1) || t.Is(indexes.Water2)
}

// IsSkiffPassable returns true if skiffs can pass through this tile
func (t *Tile) IsSkiffPassable() bool {
	// Skiffs can handle shallow water and some coastal areas
	return t.IsWater() || t.Is(indexes.Grass) || t.IsSwamp()
}

// IsCarpetPassable returns true if magic carpets can pass through this tile
func (t *Tile) IsCarpetPassable() bool {
	// Magic carpets can fly over most terrain, limited by impassable obstacles
	return !t.IsMountain() // Can't fly through solid mountains
}

// IsHorsePassable returns true if horses can pass through this tile
func (t *Tile) IsHorsePassable() bool {
	// Horses are land-based like walking but avoid difficult terrain
	return t.IsWalkingPassable() &&
		!t.IsWater() &&
		!t.IsSwamp() &&
		!t.IsMountain()
}

// IsKlimable returns true if this tile can be climbed by players
func (t *Tile) IsKlimable() bool {
	// Based on typical Ultima V mechanics: mountains, ladders, and some structures
	return t.IsMountain() ||
		t.Is(indexes.LadderUp) ||
		t.Is(indexes.LadderDown)
}

// IsLandEnemyPassable returns true if land-based enemies can move through this tile
func (t *Tile) IsLandEnemyPassable() bool {
	// Land enemies can move on walkable land tiles
	return t.IsWalkingPassable() && !t.IsWater()
}

// IsWaterEnemyPassable returns true if water-based enemies can move through this tile
func (t *Tile) IsWaterEnemyPassable() bool {
	// Water enemies move in water tiles
	return t.IsWater()
}

// IsOpenable returns true if this tile can be opened (doors, chests, etc.)
func (t *Tile) IsOpenable() bool {
	// Doors and chests can be opened
	return t.Index.IsDoor() || t.Is(indexes.Chest)
}

// IsWalkingPassable returns true if this tile can be walked through by players on foot
func (t *Tile) IsWalkingPassable() bool {
	// Use game data properties when available for more accurate logic
	// If SpeedFactor is -1, tile is explicitly impassable
	if t.SpeedFactor == -1 {
		return false
	}

	// Buildings are generally not walkable, but entrances should be passable
	if t.IsBuilding {
		// Allow entrances to be walkable
		if len(t.Name) > 0 && (strings.Contains(strings.ToLower(t.Name), "entrance") ||
			strings.Contains(strings.ToLower(t.Name), "entrace")) { // Handle typo in data
			return true
		}

		// Allow specific building tiles that represent location entrances
		if t.Is(indexes.Village) || t.Is(indexes.Keep) || t.Is(indexes.Hut) ||
			t.Is(indexes.Castle) || t.Is(indexes.Cave) || t.Is(indexes.Mine) ||
			t.Is(indexes.Shrine) || t.Is(indexes.RuinedShrine) || t.Is(indexes.Lighthouse) ||
			// Additional building entrances by index (from TileData.csv)
			t.Index == 20 || // SmallCastle
			t.Index == 21 || // LargeCastle
			t.Index == 24 || // DoomEntrance
			t.Index == 57 || // EvilCastleEntrance
			t.Index == 62 || // CastleBritianEntrace
			t.Index == 71 { // Dock
			return true
		}

		return false
	}

	// Castle tiles should not be passable even if isBuilding=false
	if len(t.Name) > 0 && strings.Contains(strings.ToLower(t.Name), "castle") {
		return false
	}

	// Basic walkable terrain: grass, paths, floors
	if t.Is(indexes.Grass) || t.IsPath() ||
		t.Is(indexes.BrickFloor) || t.Is(indexes.HexMetalGridFloor) ||
		t.Is(indexes.WoodenPlankVert1Floor) || t.Is(indexes.WoodenPlankVert2Floor) ||
		t.Is(indexes.WoodenPlankHorizFloor) {
		return true
	}

	// Terrain types that should be walkable
	if t.IsDesert() || t.IsSwamp() || t.Is(indexes.Beach) ||
		t.Is(indexes.Brush) || t.Is(indexes.ThickBrush) ||
		t.Is(indexes.Forest) || t.Is(indexes.Hills) ||
		t.Is(indexes.LeftHills) || t.Is(indexes.RightHills) {
		return true
	}

	// For tiles with names containing "Path", they should be walkable
	// This handles cases like "OutsidePath3" that aren't in the standard path index range
	if len(t.Name) > 0 &&
		(strings.Contains(strings.ToLower(t.Name), "path") ||
			strings.Contains(strings.ToLower(t.Name), "road")) {
		return true
	}

	// Use original catch-all but with more restrictions
	// Only allow for basic terrain tiles, not structures or unknown tiles
	return !t.IsMountain() && !t.IsWater() && !t.IsWall() &&
		t.SpeedFactor > 0 && t.SpeedFactor <= 6 && // Reasonable speed factors
		!t.IsBuilding // Not a building
}

// IsRangeWeaponPassable returns true if ranged weapons (arrows, etc.) can pass through this tile
func (t *Tile) IsRangeWeaponPassable() bool {
	// Range weapons can pass through most open spaces but not walls, mountains, etc.
	return !t.IsWall() && !t.IsMountain() && !t.Index.IsDoor()
}
