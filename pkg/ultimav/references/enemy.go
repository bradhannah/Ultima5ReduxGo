package references

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
