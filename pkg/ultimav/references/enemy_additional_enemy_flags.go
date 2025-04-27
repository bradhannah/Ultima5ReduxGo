package references

import (
	"encoding/json"
)

type AdditionalEnemyFlags struct {
	ActivelyAttacks      bool        `json:"ActivelyAttacks"`
	CanFlyOverWater      bool        `json:"CanFlyOverWater"`
	CanPassThroughWalls  bool        `json:"CanPassThroughWalls"`
	DoNotMove            bool        `json:"DoNotMove"`
	Era1Weight           int         `json:"Era1Weight"`
	Era2Weight           int         `json:"Era2Weight"`
	Era3Weight           int         `json:"Era3Weight"`
	Experience           int         `json:"Experience"`
	IsSandEnemy          bool        `json:"IsSandEnemy"`
	IsWaterEnemy         bool        `json:"IsWaterEnemy"`
	LargeMapMissileStr   string      `json:"LargeMapMissile"`
	LargeMapMissile      MissileType `json:"-"`
	LargeMapMissileRange int         `json:"LargeMapMissileRange"`
	Name                 string      `json:"Name"`
}

func (e *AdditionalEnemyFlags) UnmarshalJSON(data []byte) error {
	type rawAdditionalEnemyFlags AdditionalEnemyFlags
	var tempMap rawAdditionalEnemyFlags

	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	*e = AdditionalEnemyFlags(tempMap)
	e.LargeMapMissile = GetMissileTypeFromString(tempMap.LargeMapMissileStr)

	return nil
}
