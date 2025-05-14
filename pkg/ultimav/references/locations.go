package references

type Location byte

type GeneralMapType int

const (
	LargeMapType GeneralMapType = iota
	SmallMapType
	CombatMapType
	DungeonMapType
)

const (
	EmptyLocation         Location = 0xFF
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

func (l Location) GetMapType() GeneralMapType {
	switch l {
	case Britannia_Underworld:
		return LargeMapType
	case Deceit, Despise, Destard, Wrong, Covetous, Shame, Hythloth, Doom:
		return DungeonMapType
	case Combat_resting_shrine:
		return CombatMapType
	}
	return SmallMapType
}

func GetListOfAllSmallMaps() []string {
	names := make([]string, 0)

	for loc := Moonglow; loc <= Serpents_Hold; loc++ {
		names = append(names, loc.String())
	}
	return names
}
