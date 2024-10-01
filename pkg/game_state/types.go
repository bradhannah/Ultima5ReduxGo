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

type Location byte

//goland:noinspection ALL
const (
	Britannia_Underworld  Location = 0x00
	Moonglow              Location = 1
	Britain               Location = 2
	Jhelom                Location = 3
	Yew                   Location = 4
	Minoc                 Location = 5
	Trinsic               Location = 6
	Skara_Brae            Location = 7
	New_Magincia          Location = 8 // Town
	Fogsbane              Location = 9
	Stormcrow             Location = 10
	Greyhaven             Location = 11
	Waveguide             Location = 12
	Iolos_Hut             Location = 13
	Suteks_Hut            Location = 14
	SinVraals_Hut         Location = 15
	Grendels_Hut          Location = 16 // Dwelling
	Lord_Britishs_Castle  Location = 17
	Palace_of_Blackthorn  Location = 18
	West_Britanny         Location = 19
	North_Britanny        Location = 20
	East_Britanny         Location = 21
	Paws                  Location = 22
	Cove                  Location = 23
	Buccaneers_Den        Location = 24 // Castle
	Ararat                Location = 25
	Bordermarch           Location = 26
	Farthing              Location = 27
	Windemere             Location = 28
	Stonegate             Location = 29
	Lycaeum               Location = 30 // Keep
	Empath_Abbey          Location = 31
	Serpents_Hold         Location = 32
	Deceit                Location = 33 // Dungeons
	Despise               Location = 34
	Destard               Location = 35
	Wrong                 Location = 36
	Covetous              Location = 37
	Shame                 Location = 38
	Hythloth              Location = 39
	Doom                  Location = 40
	Combat_resting_shrine Location = 41
)

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
