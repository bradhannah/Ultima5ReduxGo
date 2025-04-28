package references

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

// const (
//
//	enemyAttackRangeOffset = 0x15AC
//	enemyRangeThingOffset  = 0x15DC
//	enemyFriendsOffset     = 0x16E4
//	enemyThingOffset       = 0x1714
//
// )
const (
	nTotalEnemies = 48
)

var (
	//go:embed data/AdditionalEnemyFlags.json
	additionalEnemyFlagsRaw []byte
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

type rawEnemyReferences struct {
	enemies []rawEnemyReference
	//additionalEnemyFlags []AdditionalEnemyFlags
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

// func NewEnemyReferences() *EnemyReferences {
// 	var enemies EnemyReferences

// 	return &enemies
// }

func newRawEnemyReferences(gameConfig *config.UltimaVConfiguration) *rawEnemyReferences {
	var enemies rawEnemyReferences
	enemies.enemies = make([]rawEnemyReference, 0, nTotalEnemies)

	var allAdditionalEnemyFlags []AdditionalEnemyFlags
	err := json.Unmarshal(additionalEnemyFlagsRaw, &allAdditionalEnemyFlags)
	if err != nil {
		log.Fatalf("error unmarshaling JSON: %v", err)
	}

	// now persuse each monster individually
	for nEnemy := 0; nEnemy < nTotalEnemies; nEnemy++ {
		enemyRef := createEmptyEnemyReference()

		// Enemy Flags
		// 	_dataChunks.AddDataChunk(DataChunk.DataFormatType.ByteList, "Enemy Stats", 0x13CC, 0x30 * 8, 0x00, DataChunkName.ENEMY_STATS);
		const enemyFlagsOffset = 0x154C
		const nRawBytesPerEnemyFlag = 2
		enemyFlagsRaw := gameConfig.RawDataOvl[enemyFlagsOffset+(nEnemy*nRawBytesPerEnemyFlag) : enemyFlagsOffset+((nEnemy+1)*nRawBytesPerEnemyFlag)]
		flagsBools := helpers.GetAsBitmapBoolList(enemyFlagsRaw, 0, nRawBytesPerEnemyFlag)

		for nBit := 0; nBit < len(flagsBools); nBit++ {
			enemyRef.EnemyAbilities[EnemyAbility(nBit)] = flagsBools[nBit]
		}

		// Enemy Stats
		// _dataChunks.AddDataChunk(DataChunk.DataFormatType.ByteList, "Enemy Ability Flags", 0x154C, 0x30 * 2, 0x00, DataChunkName.ENEMY_FLAGS);
		const nEnemyStats = 8
		const enemyStatsOffset = 0x13CC
		enemyStatsRaw := gameConfig.RawDataOvl[enemyStatsOffset+(nEnemy*nEnemyStats) : enemyStatsOffset+((nEnemy+1)*nEnemyStats)]
		for nStat := 0; nStat < len(enemyStatsRaw); nStat++ {
			enemyRef.EnemyStats[rawEnemyStat(nStat)] = int(enemyStatsRaw[nStat])
		}

		// Enemy Attack range
		// _dataChunks.AddDataChunk(DataChunk.DataFormatType.ByteList, "Enemy Attack Range (1-9)", 0x15AC, 0x30, 0x00, DataChunkName.ENEMY_ATTACK_RANGE);
		const nRawBytesPerEnemyAttackRange = 1
		const enemyAttackRangeOffset = 0x15AC
		enemyRef.AttackRange = byte(gameConfig.RawDataOvl[enemyAttackRangeOffset+(nEnemy*nRawBytesPerEnemyAttackRange)])

		// Enemy Range THING
		// _dataChunks.AddDataChunk(DataChunk.DataFormatType.ByteList, "Enemy Range THING", 0x15DC, 0x30, 0x00, DataChunkName.ENEMY_RANGE_THING);
		// TODO: unsure what this even is.... ;)
		const nRawBytesPerEnemyRangeThing = 1
		const enemyRangeThingOffset = 0x15DC
		enemyRef.AttackThing = byte(gameConfig.RawDataOvl[enemyRangeThingOffset+(nEnemy*nRawBytesPerEnemyRangeThing)])

		// Enemy Friends
		// _dataChunks.AddDataChunk(DataChunk.DataFormatType.ByteList, "Enemy Friends", 0x16E4, 0x30, 0x00,	DataChunkName.ENEMY_FRIENDS);
		const nRawBytesPerEnemyFriends = 1
		const enemyFriendsOffset = 0x16E4
		enemyRef.Friend = byte(gameConfig.RawDataOvl[enemyFriendsOffset+(nEnemy*nRawBytesPerEnemyFriends)])

		// Enemy THING
		// _dataChunks.AddDataChunk(DataChunk.DataFormatType.ByteList, "Enemy THING", 0x1714, 0x30, 0x00, DataChunkName.ENEMY_THING);
		// TODO: unsure what this even is.... ;)
		const nRawBytesPerEnemyThing = 1
		const enemyThingOffset = 0x1714
		enemyRef.Thing = byte(gameConfig.RawDataOvl[enemyThingOffset+(nEnemy*nRawBytesPerEnemyThing)])

		//fmt.Printf("Monster %d: %v\n", nEnemy, flagsBools)

		enemyRef.AdditionalEnemyFlags = allAdditionalEnemyFlags[nEnemy]

		enemies.enemies = append(enemies.enemies, enemyRef)
	}

	return &enemies
}
