package references

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

// NPCType general type of NPC.
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
	// unknowns may be crown and sandlewood box.
)

type NPCReference struct {
	Position     Position
	Location     Location
	DialogNumber byte
	Schedule     NPCSchedule
	npcType      NPCType
	// script TalkScript
}

const secondHalfSpriteTableIndex = 0x100

func (n *NPCReference) GetSpriteIndex() indexes.SpriteIndex {
	return indexes.SpriteIndex(int(n.npcType) + secondHalfSpriteTableIndex)
}

func (n *NPCReference) SetKeyIndex(index indexes.SpriteIndex) {
	n.npcType = NPCType(index - secondHalfSpriteTableIndex)
}

func (n *NPCReference) GetVehicleType() VehicleType {
	spriteIndex := n.GetSpriteIndex()
	if spriteIndex.IsSkiff() {
		return SkiffVehicle
	}
	if spriteIndex.IsHorseUnBoarded() {
		return HorseVehicle
	}
	if spriteIndex.IsFrigateFurled() {
		return FrigateVehicle
	}
	if spriteIndex.IsMagicCarpetUnboarded() {
		return CarpetVehicle
	}

	return NoPartyVehicle
}

func (n *NPCReference) IsEmptyNPC() bool {
	return false
}

func (n *NPCReference) GetNPCType() NPCType {
	tileIndex := n.GetSpriteIndex()

	if tileIndex.IsPartOfAnimation(indexes.Guard_KeyIndex) {
		return Guard
	}
	// TODO: requires a check on MapLocation == SmallMapReferences.SingleMapReference.Location.Windemere)
	if tileIndex.IsPartOfAnimation(indexes.Daemon1_KeyIndex) {
		_ = "a" // skip for now
	}
	const avatarTileIndexInNPCRefes = 256
	if tileIndex == avatarTileIndexInNPCRefes { // it's the Avatar
		return NoStatedNpc
	}

	if isValidType(int(n.DialogNumber)) {
		return NPCType(n.DialogNumber)
	}

	if isValidType(int(n.npcType)) {
		return n.npcType
	}

	if tileIndex.IsHorseUnBoarded() || tileIndex.IsMagicCarpetUnboarded() || tileIndex.IsFrigateFurled() || tileIndex.IsSkiff() {
		return Vehicle
	}

	return NoStatedNpc
}

func isValidType(nType int) bool {
	switch NPCType(nType) {
	case Blacksmith:
	case Barkeeper:
	case HorseSeller:
	case Shipwright:
	case Healer:
	case InnKeeper:
	case MagicSeller:
	case GuildMaster:
	case NoStatedNpc:
	case Guard:
	case WishingWell:
	case Vehicle:
		return true
	}
	return false
}

func NewNPCReferenceForVehicle(vehicle VehicleType, position Position, floorNumber FloorNumber) *NPCReference {
	npcRef := &NPCReference{}
	npcRef.Position = position
	npcRef.npcType = Vehicle
	npcRef.SetKeyIndex(vehicle.GetSpriteByDirection(NoneDirection, NoneDirection))
	npcRef.Schedule = CreateNPCScheduledFixedOneLocation(position, floorNumber)
	return npcRef
}
