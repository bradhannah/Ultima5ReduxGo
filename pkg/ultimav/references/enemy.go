package references

//"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"

type rawEnemyReference struct {
	EnemyStats           map[enemyStat]int
	AdditionalEnemyFlags AdditionalEnemyFlags
	EnemyAbilities       map[enemyAbility]bool
	AttackRange          byte
	AttackThing          byte
	Friend               byte
	Thing                byte
}

type enemyAbility int

const (
	Bludgeons         enemyAbility = 0
	PossessCharm      enemyAbility = 1
	Undead            enemyAbility = 2
	DivideOnHit       enemyAbility = 3
	Immortal          enemyAbility = 4
	PoisonAtRange     enemyAbility = 5
	StealsFood        enemyAbility = 6
	NoCorpse          enemyAbility = 7
	RangedMagic       enemyAbility = 8
	Teleport          enemyAbility = 9
	DisappearsOnDeath enemyAbility = 10
	Invisibility      enemyAbility = 11
	GatesInDaemon     enemyAbility = 12
	Poison            enemyAbility = 13
	InfectWithPlague  enemyAbility = 14
)

type enemyStat int

const (
	EnemyStatArmour         enemyStat = 0
	EnemyStatDamage         enemyStat = 1
	EnemyStatDexterity      enemyStat = 2
	EnemyStatHitPoints      enemyStat = 3
	EnemyStatIntelligence   enemyStat = 4
	EnemyStatMaxPerMap      enemyStat = 5
	EnemyStatStrength       enemyStat = 6
	EnemyStatTreasureNumber enemyStat = 7
)

func createEmptyEnemyReference() rawEnemyReference {
	var e rawEnemyReference
	e.AdditionalEnemyFlags = AdditionalEnemyFlags{}
	e.EnemyAbilities = make(map[enemyAbility]bool)
	e.EnemyStats = make(map[enemyStat]int)
	return e
}

func (e *rawEnemyReference) CanMoveToTile(tile *Tile) bool {
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

func (e *rawEnemyReference) isMonsterSpawnableOnTile(tile *Tile) bool {
	return tile.IsBoatPassable || tile.IsCarpetPassable || tile.IsHorsePassable ||
		tile.IsWalkingPassable || tile.IsWaterEnemyPassable || tile.IsLandEnemyPassable
}
