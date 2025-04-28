package references

import (
	"fmt"
	"math/rand"
)

type EnemyAbility int

const (
	beginningOfEra1 = 0
	beginningOfEra2 = 10000
	beginningOfEra3 = 30000
)

const (
	Bludgeons         EnemyAbility = 0
	PossessCharm      EnemyAbility = 1
	Undead            EnemyAbility = 2
	DivideOnHit       EnemyAbility = 3
	Immortal          EnemyAbility = 4
	PoisonAtRange     EnemyAbility = 5
	StealsFood        EnemyAbility = 6
	NoCorpse          EnemyAbility = 7
	RangedMagic       EnemyAbility = 8
	Teleport          EnemyAbility = 9
	DisappearsOnDeath EnemyAbility = 10
	Invisibility      EnemyAbility = 11
	GatesInDaemon     EnemyAbility = 12
	Poison            EnemyAbility = 13
	InfectWithPlague  EnemyAbility = 14
)

type EnemyReference struct {
	KeyFrameTile   *Tile
	Armour         int
	Damage         int
	Dexterity      int
	HitPoints      int
	Intelligence   int
	MaxPerMap      int
	Strength       int
	TreasureNumber int

	enemyAbilities       map[EnemyAbility]bool
	additionalEnemyFlags AdditionalEnemyFlags

	AttackRange int
	Friend      *EnemyReference
}

func (e *EnemyReference) CanMoveToTile(tile *Tile) bool {
	if !e.isMonsterSpawnableOnTile(tile) {
		return false
	}
	switch tile.CombatMapIndex {
	case "":
		return false
	default:
		return true
	}
}

func (e *EnemyReference) isMonsterSpawnableOnTile(tile *Tile) bool {
	return tile.IsBoatPassable || tile.IsCarpetPassable || tile.IsHorsePassable ||
		tile.IsWalkingPassable || tile.IsWaterEnemyPassable || tile.IsLandEnemyPassable
}

func (e *EnemyReference) HasAbility(ability EnemyAbility) bool {
	return e.enemyAbilities[ability]
}

func (e *EnemyReference) GetEraWeightByTurn(nTurn int) int {
	if nTurn >= beginningOfEra3 {
		return e.additionalEnemyFlags.Era3Weight
	}
	if nTurn >= beginningOfEra2 {
		return e.additionalEnemyFlags.Era2Weight
	}
	return e.additionalEnemyFlags.Era1Weight
}

// GetRandomEnemyReferenceByEraAndTile returns a randomly selected enemy that is appropriate
// for the given nTurn era and is able to spawn on the provided tile.
// It returns an error if no enemies exist for the era.
// If none of the possible enemies can move onto the tile, it returns (nil, nil).
func (e *EnemyReferences) GetRandomEnemyReferenceByEraAndTile(nTurn int, tile *Tile) (*EnemyReference, error) {
	// Filter enemy references based on era weight.
	possibleEnemies := make([]*EnemyReference, 0)
	for i := range *e {
		if (*e)[i].GetEraWeightByTurn(nTurn) > 0 {
			possibleEnemies = append(possibleEnemies, &(*e)[i])
		}
	}

	// if 0, then no possible enemies based on era.
	if len(possibleEnemies) == 0 {
		return nil, fmt.Errorf("you should always have more than zero enemies to fight in each era")
	}

	// Filter the enemies that can move onto the given tile.
	enemiesThatCanGoOnTile := make([]*EnemyReference, 0)
	for _, enemy := range possibleEnemies {
		if enemy.CanMoveToTile(tile) {
			enemiesThatCanGoOnTile = append(enemiesThatCanGoOnTile, enemy)
		}
	}

	// if no enemy can go on that tile, return nil.
	if len(enemiesThatCanGoOnTile) == 0 {
		return nil, nil
	}

	// Choose a random enemy reference from the filtered list.
	idx := rand.Intn(len(enemiesThatCanGoOnTile))
	return enemiesThatCanGoOnTile[idx], nil
}
