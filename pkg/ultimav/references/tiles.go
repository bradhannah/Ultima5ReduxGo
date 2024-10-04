package references

import (
	_ "embed"
	"encoding/json"
	"log"
	"strconv"
)

var (
	//go:embed data/TileData.json
	tileDataRaw []byte
)

//type Tiles struct {
//	Tiles []Tile
//}

type Tiles map[int]Tile

type Tile struct {
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

func (t *Tiles) UnmarshalJSON(data []byte) error {
	// Temporary map to unmarshal JSON with string keys
	var tempMap map[string]Tile
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	// Convert the keys from strings to integers
	tilesMap := make(Tiles)
	for key, value := range tempMap {
		intKey, err := strconv.Atoi(key) // Convert string key to int
		if err != nil {
			return err
		}
		tilesMap[intKey] = value
	}

	// Set the result to the original Tiles map
	*t = tilesMap
	return nil
}

func NewTileReferences() *Tiles {
	var tiles Tiles
	err := json.Unmarshal(tileDataRaw, &tiles)
	if err != nil {
		log.Fatalf("error unmarshaling JSON: %v", err)
	}
	return &tiles
}
