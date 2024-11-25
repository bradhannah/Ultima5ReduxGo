package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"

type PartyVehicle int

const (
	NoPartyVehicle PartyVehicle = iota
	CarpetVehicle
	HorseVehicle
)

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

func (t *Tile) IsPassable(vehicle PartyVehicle) bool {
	switch vehicle {
	case CarpetVehicle:
		return t.IsCarpetPassable
	case HorseVehicle:
		return t.IsHorsePassable
	case NoPartyVehicle:
		return t.IsWalkingPassable
	}
	return false
}

func (t *Tile) IsChair() bool {
	return t.Index == indexes.ChairFacingDown || t.Index == indexes.ChairFacingUp || t.Index == indexes.ChairFacingRight || t.Index == indexes.ChairFacingLeft
}

func (t *Tile) IsCannon() bool {
	return t.Index == indexes.CannonFacingLeft || t.Index == indexes.CannonFacingRight || t.Index == indexes.CannonFacingUp || t.Index == indexes.CannonFacingDown
}

func (t *Tile) IsPath() bool {
	return t.Index >= indexes.PathUpDown && t.Index <= indexes.PathAllWays
}

func (t *Tile) isNPCNoPenaltyWalkable() bool {
	return t.Index == indexes.BrickFloor || t.Index == indexes.HexMetalGridFloor || t.Index == indexes.WoodenPlankVert1Floor || t.Index == indexes.WoodenPlankVert2Floor || t.Index == indexes.WoodenPlankHorizFloor
}

func (t *Tile) GetWalkableWeight() int {
	if t.Index.IsUnlockedDoor() {
		return 1
	}
	if !t.IsWalkingPassable {
		return -1
	}
	if t.isNPCNoPenaltyWalkable() {
		return 1
	}
	if t.IsPath() {
		return 2
	}
	if t.Index == indexes.Grass {
		return 3
	}
	return 10
}

func (t *Tile) IsWalkableDuringWander() bool {
	return t.IsWalkingPassable && !t.Index.IsBed() && !t.Index.IsDoor()
}
