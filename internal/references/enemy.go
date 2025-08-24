package references

import ()

type EnemyAbility int

//goland:noinspection GoUnusedConst
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
	KeyFrameTile   *Tile `json:"key_frame_tile" yaml:"key_frame_tile"`
	Armour         int   `json:"armour" yaml:"armour"`
	Damage         int   `json:"damage" yaml:"damage"`
	Dexterity      int   `json:"dexterity" yaml:"dexterity"`
	HitPoints      int   `json:"hit_points" yaml:"hit_points"`
	Intelligence   int   `json:"intelligence" yaml:"intelligence"`
	MaxPerMap      int   `json:"max_per_map" yaml:"max_per_map"`
	Strength       int   `json:"strength" yaml:"strength"`
	TreasureNumber int   `json:"treasure_number" yaml:"treasure_number"`

	EnemyAbilities       map[EnemyAbility]bool `json:"enemy_abilities" yaml:"enemy_abilities"`
	AdditionalEnemyFlags AdditionalEnemyFlags  `json:"additional_enemy_flags" yaml:"additional_enemy_flags"`

	AttackRange int             `json:"attack_range" yaml:"attack_range"`
	Friend      *EnemyReference `json:"friend" yaml:"friend"`
}

func (e *EnemyReference) CanSpawnToTile(tile *Tile) bool {
	if !e.isMonsterSpawnableOnTile(tile) {
		return false
	}

	var bCanSpawnOnTile bool

	if e.AdditionalEnemyFlags.IsSandEnemy {
		bCanSpawnOnTile = tile.IsDesert()
	} else if e.AdditionalEnemyFlags.IsWaterEnemy {
		bCanSpawnOnTile = tile.IsWaterEnemyPassable
	} else {
		bCanSpawnOnTile = tile.IsLandEnemyPassable
	}
	return bCanSpawnOnTile
}

func (e *EnemyReference) CanMoveToTile(tile *Tile) bool {
	if !e.isMonsterSpawnableOnTile(tile) {
		return false
	}

	var bCanMoveToTile bool

	if e.AdditionalEnemyFlags.IsSandEnemy {
		bCanMoveToTile = tile.IsDesert()
	} else if e.AdditionalEnemyFlags.IsWaterEnemy {
		bCanMoveToTile = tile.IsWaterEnemyPassable
	} else {
		bCanMoveToTile = tile.IsLandEnemyPassable
	}

	if e.AdditionalEnemyFlags.CanFlyOverWater {
		bCanMoveToTile = bCanMoveToTile || tile.IsWater()
	}

	if e.AdditionalEnemyFlags.CanPassThroughWalls {
		bCanMoveToTile = bCanMoveToTile || tile.IsWall()
	}

	return bCanMoveToTile
}

func (e *EnemyReference) isMonsterSpawnableOnTile(tile *Tile) bool {
	return tile.IsBoatPassable || tile.IsCarpetPassable || tile.IsHorsePassable ||
		tile.IsWalkingPassable || tile.IsWaterEnemyPassable || tile.IsLandEnemyPassable
}

func (e *EnemyReference) HasAbility(ability EnemyAbility) bool {
	return e.EnemyAbilities[ability]
}

func EnemyAbilityToString(ability EnemyAbility) string {
	switch ability {
	case Bludgeons:
		return "Bludgeons"
	case PossessCharm:
		return "PossessCharm"
	case Undead:
		return "Undead"
	case DivideOnHit:
		return "DivideOnHit"
	case Immortal:
		return "Immortal"
	case PoisonAtRange:
		return "PoisonAtRange"
	case StealsFood:
		return "StealsFood"
	case NoCorpse:
		return "NoCorpse"
	case RangedMagic:
		return "RangedMagic"
	case Teleport:
		return "Teleport"
	case DisappearsOnDeath:
		return "DisappearsOnDeath"
	case Invisibility:
		return "Invisibility"
	case GatesInDaemon:
		return "GatesInDaemon"
	case Poison:
		return "Poison"
	case InfectWithPlague:
		return "InfectWithPlague"
	default:
		return "Unknown"
	}
}
