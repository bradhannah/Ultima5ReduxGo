package references

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
)

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

func (e *EnemyReference) GetEraWeight(era datetime.Era) int {
	switch era {
	case datetime.EarlyEra:
		return e.AdditionalEnemyFlags.Era1Weight
	case datetime.MiddleEra:
		return e.AdditionalEnemyFlags.Era2Weight
	case datetime.LateEra:
		return e.AdditionalEnemyFlags.Era3Weight
	default:
		log.Fatal("Unexpected Era")
		return 0
	}
}

func (e *EnemyReference) CanSpawnToTile(tile *Tile) bool {
	if !e.isMonsterSpawnableOnTile(tile) {
		return false
	}

	var bCanSpawnOnTile bool

	if e.AdditionalEnemyFlags.IsSandEnemy {
		bCanSpawnOnTile = strings.HasPrefix(strings.ToLower(tile.Name), "sand")
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
		bCanMoveToTile = strings.HasPrefix(strings.ToLower(tile.Name), "sand")
	} else if e.AdditionalEnemyFlags.IsWaterEnemy {
		bCanMoveToTile = tile.IsWaterEnemyPassable
	} else {
		bCanMoveToTile = tile.IsLandEnemyPassable
	}

	if e.AdditionalEnemyFlags.CanFlyOverWater {
		bCanMoveToTile = bCanMoveToTile || strings.Contains(strings.ToLower(tile.Name), "water")
	}

	if e.AdditionalEnemyFlags.CanPassThroughWalls {
		bCanMoveToTile = bCanMoveToTile || strings.Contains(strings.ToLower(tile.Name), "wall")
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

// GetRandomEnemyReferenceByEraAndTile returns a randomly selected enemy that is appropriate
// for the given nTurn era and is able to spawn on the provided tile.
// It returns an error if no enemies exist for the era.
// If none of the possible enemies can move onto the tile, it returns (nil, nil).
func (e *EnemyReferences) GetRandomEnemyReferenceByEraAndTile(era datetime.Era, tile *Tile) (*EnemyReference, error) {
	// Filter enemy references based on era weight.
	possibleEnemies := make([]*EnemyReference, 0)
	for _, v := range *e {
		// if (*e)[i].GetEraWeight(era) > 0 {
		if v.GetEraWeight(era) > 0 {
			possibleEnemies = append(possibleEnemies, &v) // &(*e)[i])
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

type EnemyReferenceSafe struct {
	KeyFrameTile         *Tile                `json:"key_frame_tile" yaml:"key_frame_tile"`
	Armour               int                  `json:"armour" yaml:"armour"`
	Damage               int                  `json:"damage" yaml:"damage"`
	Dexterity            int                  `json:"dexterity" yaml:"dexterity"`
	HitPoints            int                  `json:"hit_points" yaml:"hit_points"`
	Intelligence         int                  `json:"intelligence" yaml:"intelligence"`
	MaxPerMap            int                  `json:"max_per_map" yaml:"max_per_map"`
	Strength             int                  `json:"strength" yaml:"strength"`
	TreasureNumber       int                  `json:"treasure_number" yaml:"treasure_number"`
	EnemyAbilities       map[string]bool      `json:"enemy_abilities" yaml:"enemy_abilities"`
	AdditionalEnemyFlags AdditionalEnemyFlags `json:"additional_enemy_flags" yaml:"additional_enemy_flags"`
	AttackRange          int                  `json:"attack_range" yaml:"attack_range"`
	// Friend omitted to avoid cyclic reference
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

func toFriendlyAbilityMap(abilities map[EnemyAbility]bool) map[string]bool {
	friendly := make(map[string]bool)
	for k, v := range abilities {
		key := fmt.Sprintf("%d_%s", k, EnemyAbilityToString(k))
		friendly[key] = v
	}
	return friendly
}

func ToSafeEnemyReferences(refs []EnemyReference) []EnemyReferenceSafe {
	safe := make([]EnemyReferenceSafe, len(refs))
	for i, e := range refs {
		safe[i] = EnemyReferenceSafe{
			KeyFrameTile:         e.KeyFrameTile,
			Armour:               e.Armour,
			Damage:               e.Damage,
			Dexterity:            e.Dexterity,
			HitPoints:            e.HitPoints,
			Intelligence:         e.Intelligence,
			MaxPerMap:            e.MaxPerMap,
			Strength:             e.Strength,
			TreasureNumber:       e.TreasureNumber,
			EnemyAbilities:       toFriendlyAbilityMap(e.EnemyAbilities),
			AdditionalEnemyFlags: e.AdditionalEnemyFlags,
			AttackRange:          e.AttackRange,
		}
	}
	return safe
}
