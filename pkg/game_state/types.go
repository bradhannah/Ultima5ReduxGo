package game_state

const NMaxPlayerNameSize = 9

const (
	ubPartyMembers StartingMemoryAddressUb = 0x2B5
	ubActivePlayer                         = 0x2D5
)

const (
	u16CurrentYear StartingMemoryAddressU16 = 0x2CE
)

type PartyStatus byte

const (
	InTheParty     PartyStatus = 0x00
	HasntJoinedYet PartyStatus = 0xFF
	AtTheInn       PartyStatus = 0x01
)

//goland:noinspection ALL

type BritOrUnderworld byte

const (
	Britannia  BritOrUnderworld = 0x00
	Underworld                  = 0xFF
)

type MoonstoneStatus byte

const (
	Buried      MoonstoneStatus = 0x00
	InInventory                 = 0xFF
)

type PlayerCharacter struct {
	Name         [NMaxPlayerNameSize]byte
	Gender       CharacterGender
	Class        CharacterClass
	Status       CharacterStatus
	Strength     byte
	Dexterity    byte
	Intelligence byte
	CurrentMp    byte
	CurrentHp    uint16
	MaxHp        uint16
	Exp          uint16
	Level        byte
	MonthsAtInn  byte
	Unknown      byte
	Helmet       byte
	Armor        byte
	Weapon       byte
	Shield       byte
	Ring         byte
	Amulet       byte
	PartyStatus  PartyStatus
}
