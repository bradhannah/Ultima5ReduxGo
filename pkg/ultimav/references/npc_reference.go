package references

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
)

type NPCType byte

const (
	Blacksmith  NPCType = 0x81
	Barkeeper   NPCType = 0x82
	HorseSeller NPCType = 0x83
	Shipwright  NPCType = 0x84
	Healer      NPCType = 0x87
	InnKeeper   NPCType = 0x88
	MagicSeller NPCType = 0x85
	GuildMaster NPCType = 0x86
	NoStatedNpc NPCType = 0xFF
	Guard       NPCType = 0xFE
	WishingWell NPCType = 0xFD
	Vehicle     NPCType = 0x1F
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

func (n *NPCReference) SetKeyIndex(index indexes.SpriteIndex) {
	n.Type = NPCType(index - 0x100)
}

func NewNPCReferenceForVehicle(vehicle VehicleType, position Position, floorNumber FloorNumber) *NPCReference {
	npcRef := &NPCReference{}
	npcRef.Position = position
	npcRef.Type = Vehicle
	npcRef.SetKeyIndex(vehicle.GetSpriteByDirection(NoneDirection, NoneDirection))
	npcRef.Schedule = CreateNPCScheduledFixedOneLocation(position, floorNumber)
	return npcRef
}
