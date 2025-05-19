package game_state

// const (
// 	ubPartyMembers StartingMemoryAddressUb = 0x2B5
// 	ubActivePlayer StartingMemoryAddressUb = 0x2D5
// )

// const (
// 	u16CurrentYear StartingMemoryAddressU16 = 0x2CE
// )

//goland:noinspection ALL

type BritOrUnderworld byte

const (
	Britannia  BritOrUnderworld = 0x00
	Underworld BritOrUnderworld = 0xFF
)

type MoonstoneStatus byte

const (
	Buried      MoonstoneStatus = 0x00
	InInventory MoonstoneStatus = 0xFF
)
