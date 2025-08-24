package references

import (
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
	Index                     indexes.SpriteIndex
	Name                      string `json:"Name"`
	Description               string `json:"Description"`
	IsWalkingPassable         bool   `json:"IsWalking_Passable"`
	RangeWeaponPassable       bool   `json:"RangeWeapon_Passable"`
	IsBoatPassable            bool   `json:"IsBoat_Passable"`
	IsSkiffPassable           bool   `json:"IsSkiff_Passable"`
	IsCarpetPassable          bool   `json:"IsCarpet_Passable"`
	IsHorsePassable           bool   `json:"IsHorse_Passable"`
	IsKlimable                bool   `json:"IsKlimable"`
	IsOpenable                bool   `json:"IsOpenable"`
	IsLandEnemyPassable       bool   `json:"IsLandEnemyPassable"`
	IsWaterEnemyPassable      bool   `json:"IsWaterEnemyPassable"`
	SpeedFactor               int    `json:"SpeedFactor"`
	LightEmission             int    `json:"LightEmission"`
	IsPartOfAnimation         bool   `json:"IsPartOfAnimation"`
	TotalAnimationFrames      int    `json:"TotalAnimationFrames"`
	AnimationIndex            int    `json:"AnimationIndex"`
	IsUpright                 bool   `json:"IsUpright"`
	FlatTileSubstitutionIndex int    `json:"FlatTileSubstitutionIndex"`
	FlatTileSubstitutionName  string `json:"FlatTileSubstitutionName"`
	IsEnemy                   bool   `json:"IsEnemy"`
	IsNPC                     bool   `json:"IsNPC"`
	IsBuilding                bool   `json:"IsBuilding"`
	DontDraw                  bool   `json:"DontDraw"`
	IsTalkOverable            bool   `json:"IsTalkOverable"`
	IsBoardable               bool   `json:"IsBoardable"`
	IsGuessableFloor          bool   `json:"IsGuessableFloor"`
	BlocksLight               bool   `json:"BlocksLight"`
	IsWindow                  bool   `json:"IsWindow"`
	CombatMapIndex            string `json:"CombatMapIndex"`
}

func (t *Tile) IsPassable(vehicle VehicleType) bool {
	switch vehicle {
	case CarpetVehicle:
		return t.IsCarpetPassable
	case HorseVehicle:
		return t.IsHorsePassable
	case SkiffVehicle:
		return t.IsSkiffPassable
	case FrigateVehicle:
		return t.IsWaterEnemyPassable
	case NoPartyVehicle:
		return t.IsWalkingPassable
	case NPC:
		return t.IsLandEnemyPassable
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
	return t.Index == indexes.Barrel
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
	return t.Index == indexes.BrickFloor || t.Index == indexes.HexMetalGridFloor || t.Index == indexes.WoodenPlankVert1Floor || t.Index == indexes.WoodenPlankVert2Floor || t.Index == indexes.WoodenPlankHorizFloor
}

func (t *Tile) GetWalkableWeight() int {
	if t.Index.IsUnlockedDoor() {
		return weightIdealPath
	}

	if !t.IsWalkingPassable {
		return weightImpassable
	}

	if t.isNPCNoPenaltyWalkable() {
		return weightIdealPath
	}

	if t.IsPath() {
		return weightPath
	}

	if t.Index == indexes.Grass {
		return weightGrass
	}

	return weightDefault
}

func (t *Tile) IsWalkableDuringWander() bool {
	return t.IsWalkingPassable && !t.Index.IsBed() && !t.Index.IsDoor()
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
	return t.Index == indexes.LargeRockWall || t.Index == indexes.StoneBrickWall || t.Index == indexes.StoneBrickWallSecret
}

func (t *Tile) IsRoad() bool {
	return t.IsPath() // Roads are paths in this context
}

func (t *Tile) IsSwamp() bool {
	return t.Index == indexes.Swamp
}

func (t *Tile) IsWater() bool {
	return t.Index == indexes.Water1 || t.Index == indexes.Water2 || t.Index == indexes.WaterShallow
}

func (t *Tile) IsDesert() bool {
	return t.Index == indexes.Desert || t.Index == indexes.LeftDesert2 || t.Index == indexes.RightDesert2
}

func (t *Tile) IsMountain() bool {
	return t.Index == indexes.SmallMountains
}

func (t *Tile) IsForest() bool {
	// Forest tiles are identified as passable land tiles that are not other terrain types
	// This logic may need refinement based on actual forest tile indexes
	return t.IsLandEnemyPassable &&
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
