package references

type rawEnemyReference struct {
	EnemyStats           map[rawEnemyStat]int
	AdditionalEnemyFlags AdditionalEnemyFlags
	EnemyAbilities       map[EnemyAbility]bool
	AttackRange          byte
	AttackThing          byte
	Friend               byte
	Thing                byte
}

type rawEnemyStat int

const (
	EnemyStatArmour         rawEnemyStat = 0
	EnemyStatDamage         rawEnemyStat = 1
	EnemyStatDexterity      rawEnemyStat = 2
	EnemyStatHitPoints      rawEnemyStat = 3
	EnemyStatIntelligence   rawEnemyStat = 4
	EnemyStatMaxPerMap      rawEnemyStat = 5
	EnemyStatStrength       rawEnemyStat = 6
	EnemyStatTreasureNumber rawEnemyStat = 7
)

func createEmptyEnemyReference() rawEnemyReference {
	var e rawEnemyReference
	e.AdditionalEnemyFlags = AdditionalEnemyFlags{}
	e.EnemyAbilities = make(map[EnemyAbility]bool)
	e.EnemyStats = make(map[rawEnemyStat]int)
	return e
}
