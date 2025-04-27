package references

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
)

type EnemyReferences []EnemyReference

const nFirstEnemyTileReferenceIndex = 384 // Sea horse
const nFramesPerEnemy = 4

func NewAllEnemyReferences(gameConfig *config.UltimaVConfiguration,
	tiles *Tiles) *EnemyReferences {
	rawEnemyReferences := newRawEnemyReferences(gameConfig)

	enemyRefs := EnemyReferences{}

	for nEnemy := 0; nEnemy < len(rawEnemyReferences.enemies); nEnemy++ {
		rawEnemyRef := rawEnemyReferences.enemies[nEnemy]
		enemyRef := EnemyReference{
			KeyFrameTile:         tiles.GetTile(indexes.SpriteIndex(nFirstEnemyTileReferenceIndex + (nEnemy * nFramesPerEnemy))),
			Armour:               rawEnemyRef.EnemyStats[EnemyStatArmour],
			Damage:               rawEnemyRef.EnemyStats[EnemyStatDamage],
			Dexterity:            rawEnemyRef.EnemyStats[EnemyStatDexterity],
			HitPoints:            rawEnemyRef.EnemyStats[EnemyStatHitPoints],
			Intelligence:         rawEnemyRef.EnemyStats[EnemyStatIntelligence],
			MaxPerMap:            rawEnemyRef.EnemyStats[EnemyStatMaxPerMap],
			Strength:             rawEnemyRef.EnemyStats[EnemyStatStrength],
			TreasureNumber:       rawEnemyRef.EnemyStats[EnemyStatTreasureNumber],
			enemyAbilities:       rawEnemyRef.EnemyAbilities,
			AttackRange:          int(rawEnemyRef.AttackRange),
			Friend:               nil,
			additionalEnemyFlags: rawEnemyRef.AdditionalEnemyFlags,
		}

		enemyRefs = append(enemyRefs, enemyRef)
	}

	// do the friend reference population aftewards
	for nEnemy := 0; nEnemy < len(enemyRefs); nEnemy++ {
		rawEnemyRef := rawEnemyReferences.enemies[nEnemy]
		enemyRef := &enemyRefs[nEnemy]
		enemyRef.Friend = &enemyRefs[rawEnemyRef.Friend]
	}

	return &enemyRefs
}
