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

// Tile represents a single tile and it's properties.
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
	IsPushable                bool   `json:"IsPushable"`
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
