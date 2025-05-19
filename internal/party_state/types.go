package party_state

const NMaxPlayerNameSize = 9

type PartyStatus byte

const (
	InTheParty     PartyStatus = 0x00
	HasntJoinedYet PartyStatus = 0xFF
	AtTheInn       PartyStatus = 0x01
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
