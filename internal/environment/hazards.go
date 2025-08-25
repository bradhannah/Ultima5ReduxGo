package environment

import (
	"golang.org/x/exp/rand"

	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

// SystemCallbacks interface to avoid import cycle with game_state
type SystemCallbacks interface {
	AddRowStr(message string)
}

type HazardType int

const (
	NoHazard HazardType = iota
	PoisonSwamp
	LavaBurn
	FireplaceBurn
)

type EnvironmentalHazards struct {
	rng       *rand.Rand
	callbacks SystemCallbacks
}

type HazardResult struct {
	Type         HazardType
	Affected     []int // Party member indices affected
	MessageShown bool
	DamageDealt  int
}

func NewEnvironmentalHazards(rng *rand.Rand, callbacks SystemCallbacks) *EnvironmentalHazards {
	return &EnvironmentalHazards{
		rng:       rng,
		callbacks: callbacks,
	}
}

func (h *EnvironmentalHazards) CheckTileHazards(tile *references.Tile, party *party_state.PartyState, location references.GeneralMapType) HazardResult {
	hazardType := h.detectHazardType(tile, location)

	switch hazardType {
	case PoisonSwamp:
		return h.checkSwampPoison(party)
	case LavaBurn, FireplaceBurn:
		return h.checkBurningHazard(party, hazardType)
	default:
		return HazardResult{Type: NoHazard}
	}
}

func (h *EnvironmentalHazards) detectHazardType(tile *references.Tile, location references.GeneralMapType) HazardType {
	if tile == nil {
		return NoHazard
	}

	switch tile.Index {
	case indexes.Swamp:
		return PoisonSwamp
	case indexes.Lava:
		return LavaBurn
	case indexes.Fireplace:
		if location == references.SmallMapType {
			return FireplaceBurn // Only dangerous indoors
		}
	}

	return NoHazard
}

func (h *EnvironmentalHazards) checkSwampPoison(party *party_state.PartyState) HazardResult {
	result := HazardResult{Type: PoisonSwamp}

	for i, member := range party.Characters {
		if member.Status == party_state.Dead || member.Status == party_state.Poisoned || member.PartyStatus != party_state.InTheParty {
			continue // Skip dead, already poisoned, or not active party members
		}

		// Original: random(1, 30) > member.dexterity
		poisonRoll := h.rng.Intn(30) + 1
		if poisonRoll > int(member.Dexterity) {
			party.Characters[i].Status = party_state.Poisoned
			result.Affected = append(result.Affected, i)
		}
	}

	if len(result.Affected) > 0 {
		// Build message with character names
		message := "Poisoned: "
		names := make([]string, 0, len(result.Affected))
		for _, memberIndex := range result.Affected {
			if memberIndex < len(party.Characters) {
				// Get character name, trimming null bytes
				name := string(party.Characters[memberIndex].Name[:])
				// Remove null terminator if present
				if nullIndex := len(name); nullIndex > 0 {
					for i, b := range []byte(name) {
						if b == 0 {
							nullIndex = i
							break
						}
					}
					name = name[:nullIndex]
				}
				if name != "" {
					names = append(names, name)
				}
			}
		}

		if len(names) > 0 {
			message += joinNames(names)
			h.callbacks.AddRowStr(message)
			result.MessageShown = true
		} else {
			// Fallback if no names found
			h.callbacks.AddRowStr("Poisoned!")
			result.MessageShown = true
		}
	}

	return result
}

func (h *EnvironmentalHazards) checkBurningHazard(party *party_state.PartyState, hazardType HazardType) HazardResult {
	result := HazardResult{Type: hazardType}

	// Original: damage_party_on_land() - each living member takes 1..8
	totalDamage := 0
	for i, member := range party.Characters {
		if member.Status != party_state.Dead && member.PartyStatus == party_state.InTheParty {
			damage := h.rng.Intn(8) + 1 // 1..8 range
			newHP := helpers.Max(0, int(member.CurrentHp)-damage)
			party.Characters[i].CurrentHp = uint16(newHP)
			totalDamage += damage
			result.Affected = append(result.Affected, i)

			if party.Characters[i].CurrentHp <= 0 {
				party.Characters[i].Status = party_state.Dead
			}
		}
	}

	result.DamageDealt = totalDamage

	// Show burning message with character names
	if len(result.Affected) > 0 {
		message := "Burning: "
		names := make([]string, 0, len(result.Affected))
		for _, memberIndex := range result.Affected {
			if memberIndex < len(party.Characters) {
				// Get character name, trimming null bytes
				name := string(party.Characters[memberIndex].Name[:])
				// Remove null terminator if present
				if nullIndex := len(name); nullIndex > 0 {
					for i, b := range []byte(name) {
						if b == 0 {
							nullIndex = i
							break
						}
					}
					name = name[:nullIndex]
				}
				if name != "" {
					names = append(names, name)
				}
			}
		}

		if len(names) > 0 {
			message += joinNames(names)
			h.callbacks.AddRowStr(message)
			result.MessageShown = true
		} else {
			// Fallback if no names found
			h.callbacks.AddRowStr("Burning!")
			result.MessageShown = true
		}
	}

	return result
}

// joinNames creates a grammatically correct list of names
func joinNames(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0]
	case 2:
		return names[0] + " and " + names[1]
	default:
		// For 3+ names: "Alice, Bob, and Charlie"
		result := ""
		for i, name := range names {
			if i == 0 {
				result = name
			} else if i == len(names)-1 {
				result += ", and " + name
			} else {
				result += ", " + name
			}
		}
		return result
	}
}
