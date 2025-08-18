package party_state

import (
	"unsafe"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

const NPlayers = 6

type PartyState struct {
	Characters [NPlayers]PlayerCharacter
	Inventory  Inventory
	Karma      Karma
	metNpcs    map[references.Location][]bool
	deadNpcs   map[references.Location][]bool
}

func newPartyState() *PartyState {
	ps := new(PartyState)
	ps.metNpcs = make(map[references.Location][]bool)
	ps.deadNpcs = make(map[references.Location][]bool)
	ps.Inventory = *NewInventory()

	for i := 0; i < int(references.Serpents_Hold); i++ {
		ps.metNpcs[references.Location(i)] = make([]bool, 32) // 32 NPCs per location
		ps.deadNpcs[references.Location(i)] = make([]bool, 32)
	}
	return ps
}

func LoadFromRaw(raw [4192]byte) (p *PartyState) {
	ps := newPartyState()
	const lCharacters = 0x02
	const metNpcsOffset = 0x634  // Corrected offset for met NPCs
	const deadNpcsOffset = 0x5B4 // Corrected offset for dead NPCs
	const npcsPerTown = 32
	const numNpcBits = 128
	const bitsPerByte = 8

	characterPtr := (*[NPlayers]PlayerCharacter)(unsafe.Pointer(&raw[lCharacters]))
	ps.Characters = *characterPtr

	// Load metNpcs and deadNpcs from bitfields
	for bitIdx := 0; bitIdx < numNpcBits; bitIdx++ {
		byteIdx := bitIdx / bitsPerByte
		bitInByte := bitIdx % bitsPerByte
		met := (raw[metNpcsOffset+byteIdx] & (1 << bitInByte)) != 0
		dead := (raw[deadNpcsOffset+byteIdx] & (1 << bitInByte)) != 0
		location := references.Location(bitIdx / npcsPerTown)
		npcIdx := bitIdx % npcsPerTown
		if slice, ok := ps.metNpcs[location]; ok && npcIdx < len(slice) {
			ps.metNpcs[location][npcIdx] = met
		}
		if slice, ok := ps.deadNpcs[location]; ok && npcIdx < len(slice) {
			ps.deadNpcs[location][npcIdx] = dead
		}
	}

	return ps
}

func (p *PartyState) HasRoom() bool {
	for n, c := range p.Characters {
		if c.GetNameAsString() == "" {
			if n < NPlayers-1 {
				return true
			}

		}
	}
	return false
	//return len(p.Characters) < NPlayers
}

func (p *PartyState) JoinNPC(_ interface{}) error {
	// noop
	return nil
}
func (p *PartyState) HasMet(npc interface{}) bool {
	return false
}

func (p *PartyState) AvatarName() string {
	return p.Characters[0].GetNameAsString()
}

func (p *PartyState) SetMet(location references.Location, npcId int) {
	p.metNpcs[location][npcId] = true
}

func (p *PartyState) SetDeadNpc(location references.Location, npcId int) {
	p.deadNpcs[location][npcId] = true
}

// MetNpcs returns the metNpcs map for testing purposes.
func (p *PartyState) MetNpcs() map[references.Location][]bool {
	return p.metNpcs
}

// DeadNpcs returns the deadNpcs map for testing purposes.
func (p *PartyState) DeadNpcs() map[references.Location][]bool {
	return p.deadNpcs
}
