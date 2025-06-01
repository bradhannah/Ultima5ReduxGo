package party_state

import "unsafe"

const NPlayers = 6

type PartyState struct {
	Characters [NPlayers]PlayerCharacter
	Inventory  Inventory
	Karma      Karma
}

func (p *PartyState) LoadFromRaw(raw [4192]byte) {
	const lCharacters = 0x02
	characterPtr := (*[NPlayers]PlayerCharacter)(unsafe.Pointer(&raw[lCharacters]))
	p.Characters = *characterPtr
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

func (p *PartyState) SetMet(npcId int) {
	// noop
}
