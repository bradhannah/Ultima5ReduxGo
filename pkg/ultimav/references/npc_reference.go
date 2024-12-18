package references

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
)

type NPCType byte

const (
	Blacksmith  NPCType = 0x81
	Barkeeper           = 0x82
	HorseSeller         = 0x83
	Shipwright          = 0x84
	Healer              = 0x87
	InnKeeper           = 0x88
	MagicSeller         = 0x85
	GuildMaster         = 0x86
	NoStatedNpc         = 0xFF
	Guard               = 0xFE
	WishingWell         = 0xFD
	// unknowns may be crown and sandlewood box
)

type NPCReference struct {
	Position     Position
	Location     Location
	DialogNumber byte
	Schedule     NPCSchedule
	Type         NPCType
	// script TalkScript
}

func (n *NPCReference) GetTileIndex() indexes.SpriteIndex {
	return indexes.SpriteIndex(int(n.Type) + 0x100)
}
