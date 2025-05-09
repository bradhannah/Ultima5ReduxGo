package references

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
)

type EnemyAbility int

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

func (e *EnemyReference) GetEraWeight(era Era) int {
	switch era {
	case EarlyEra:
		return e.additionalEnemyFlags.Era1Weight
	case MiddleEra:
		return e.additionalEnemyFlags.Era2Weight
	case LateEra:
		return e.additionalEnemyFlags.Era3Weight
	default:
		log.Fatal("Unepxected Era")
		return 0
	}

	// if nTurn >= beginningOfEra3 {
	// 	return e.additionalEnemyFlags.Era3Weight
	// }
	// if nTurn >= beginningOfEra2 {
	// 	return e.additionalEnemyFlags.Era2Weight
	// }
	// return e.additionalEnemyFlags.Era1Weight
}

//  protected override bool CanMoveToDumb(Map map, Point2D mapUnitPosition)
// {
// 	if (EnemyReference.DoesNotMove) return false;

// 	bool bCanMove = false;
// 	TileReference tileReference = map.GetTileReference(mapUnitPosition);

// 	bool bIsMapUnitOnTile = map.IsMapUnitOccupiedTile(mapUnitPosition);
// 	if (bIsMapUnitOnTile) return false;

// 	if (EnemyReference.IsSandEnemy)
// 	{
// 		// if tile is sand
// 		bCanMove |= tileReference.Name.IndexOf("sand", 0, StringComparison.CurrentCultureIgnoreCase) >= 0;
// 	}
// 	else if (EnemyReference.IsWaterEnemy)
// 	{
// 		// if tile is water
// 		bCanMove |= tileReference.IsWaterEnemyPassable;
// 	}
// 	else
// 	{
// 		// the enemy is a land monster by process of elimination
// 		bCanMove |= tileReference.IsLandEnemyPassable;
// 	}

// 	if (EnemyReference.CanFlyOverWater)
// 	{
// 		// if tile is water
// 		bCanMove |= tileReference.IsWaterTile;
// 	}

// 	if (EnemyReference.CanPassThroughWalls)
// 	{
// 		// if tile is wall
// 		bCanMove |=
// 			tileReference.Name.IndexOf("wall", 0, StringComparison.CurrentCultureIgnoreCase) >= 0;
// 	}

// 	return bCanMove;
// }
// }

func (e *EnemyReference) CanSpawnToTile(tile *Tile) bool {

	if !e.isMonsterSpawnableOnTile(tile) {
		return false
	}

	bCanSpawnOnTile := false

	if e.additionalEnemyFlags.IsSandEnemy {
		bCanSpawnOnTile = bCanSpawnOnTile || strings.HasPrefix(strings.ToLower(tile.Name), "sand")
	} else if e.additionalEnemyFlags.IsWaterEnemy {
		bCanSpawnOnTile = bCanSpawnOnTile || tile.IsWaterEnemyPassable
	} else {
		bCanSpawnOnTile = bCanSpawnOnTile || tile.IsLandEnemyPassable
	}
	return bCanSpawnOnTile
}

func (e *EnemyReference) CanMoveToTile(tile *Tile) bool {

	if !e.isMonsterSpawnableOnTile(tile) {
		return false
	}

	bCanMoveToTile := false

	if e.additionalEnemyFlags.IsSandEnemy {
		bCanMoveToTile = bCanMoveToTile || strings.HasPrefix(strings.ToLower(tile.Name), "sand")
	} else if e.additionalEnemyFlags.IsWaterEnemy {
		bCanMoveToTile = bCanMoveToTile || tile.IsWaterEnemyPassable
	} else {
		bCanMoveToTile = bCanMoveToTile || tile.IsLandEnemyPassable
	}

	if e.additionalEnemyFlags.CanFlyOverWater {
		bCanMoveToTile = bCanMoveToTile || strings.Contains(strings.ToLower(tile.Name), "water")
	}

	if e.additionalEnemyFlags.CanPassThroughWalls {
		bCanMoveToTile = bCanMoveToTile || strings.Contains(strings.ToLower(tile.Name), "wall")
	}

	return bCanMoveToTile

	// switch tile.CombatMapIndex {
	// case "":
	// 	return false
	// default:
	// 	return true
	// }
}

func (e *EnemyReference) isMonsterSpawnableOnTile(tile *Tile) bool {
	return tile.IsBoatPassable || tile.IsCarpetPassable || tile.IsHorsePassable ||
		tile.IsWalkingPassable || tile.IsWaterEnemyPassable || tile.IsLandEnemyPassable
}

func (e *EnemyReference) HasAbility(ability EnemyAbility) bool {
	return e.enemyAbilities[ability]
}

// GetRandomEnemyReferenceByEraAndTile returns a randomly selected enemy that is appropriate
// for the given nTurn era and is able to spawn on the provided tile.
// It returns an error if no enemies exist for the era.
// If none of the possible enemies can move onto the tile, it returns (nil, nil).
func (e *EnemyReferences) GetRandomEnemyReferenceByEraAndTile(era Era, tile *Tile) (*EnemyReference, error) {
	// Filter enemy references based on era weight.
	possibleEnemies := make([]*EnemyReference, 0)
	for _, v := range *e {
		// if (*e)[i].GetEraWeight(era) > 0 {
		if v.GetEraWeight(era) > 0 {
			possibleEnemies = append(possibleEnemies, &v) //&(*e)[i])
		}
	}

	// if 0, then no possible enemies based on era.
	if len(possibleEnemies) == 0 {
		return nil, fmt.Errorf("you should always have more than zero enemies to fight in each era")
	}

	// Filter the enemies that can move onto the given tile.
	enemiesThatCanGoOnTile := make([]*EnemyReference, 0)
	for _, enemy := range possibleEnemies {
		if enemy.CanSpawnToTile(tile) {
			enemiesThatCanGoOnTile = append(enemiesThatCanGoOnTile, enemy)
		}
	}

	// if no enemy can go on that tile, return nil.
	if len(enemiesThatCanGoOnTile) == 0 {
		return nil, fmt.Errorf("no enemies can go on tile %s", tile.Name)
	}

	// Choose a random enemy reference from the filtered list.
	idx := rand.Intn(len(enemiesThatCanGoOnTile))
	return enemiesThatCanGoOnTile[idx], nil
}
