package references

type PartyVehicle int

const (
	NoPartyVehicle PartyVehicle = iota
	CarpetVehicle
	HorseVehicle
)

type Tile struct {
	Index                     int
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
