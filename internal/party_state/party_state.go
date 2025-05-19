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
