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

	for i := 0; i < int(references.Serpents_Hold); i++ {
		ps.metNpcs[references.Location(i)] = make([]bool, 0)
		ps.deadNpcs[references.Location(i)] = make([]bool, 0)
	}
	return ps
}

func LoadFromRaw(raw [4192]byte) (p *PartyState) {
	ps := newPartyState()
	const lCharacters = 0x02
	characterPtr := (*[NPlayers]PlayerCharacter)(unsafe.Pointer(&raw[lCharacters]))
	ps.Characters = *characterPtr
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
