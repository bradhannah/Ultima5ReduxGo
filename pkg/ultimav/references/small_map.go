package references

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
)

const (
	XSmallMapTiles         = int16(32)
	YSmallMapTiles         = int16(32)
	TotalSmallMapLocations = 32
	smallMapSizeInBytes    = int(XSmallMapTiles * YSmallMapTiles)
)

type SmallMapMasterTypes int

const (
	Castle SmallMapMasterTypes = iota
	Towne
	Dwelling
	Keep
	Dungeon
	CutOrIntroScene
	None
)

func getSmallMapFile(smallMap SmallMapMasterTypes) string {
	switch smallMap {
	case Castle:
		return CASTLE_DAT
	case Towne:
		return TOWNE_DAT
	case Dwelling:
		return DWELLING_DAT
	case Keep:
		return KEEP_DAT
	case Dungeon:
		return DUNGEON_DAT
	default:

		panic("unhandled default case")
	}
}

func getMapMasterFromLocation(location Location) SmallMapMasterTypes {
	switch location {
	case Lord_Britishs_Castle, Palace_of_Blackthorn, East_Britanny, West_Britanny, North_Britanny, Paws, Cove, Buccaneers_Den:
		return Castle
	case Moonglow, Britain, Jhelom, Yew, Minoc, Trinsic, Skara_Brae, New_Magincia:
		return Towne
	case Fogsbane, Stormcrow, Waveguide, Greyhaven, Iolos_Hut, Suteks_Hut, SinVraals_Hut, Grendels_Hut:
		return Dwelling
		//case _location.spektran
	case Ararat, Bordermarch, Farthing, Windemere, Stonegate, Lycaeum, Empath_Abbey, Serpents_Hold:
		return Keep
	case Deceit, Despise, Destard, Wrong, Covetous, Shame, Hythloth, Doom:
		return Dungeon
	case Combat_resting_shrine:
		return CutOrIntroScene
	case Britannia_Underworld:
		return None
	}
	return None
}

func NewSmallMapReferences(gameConfig *config.UltimaVConfiguration) (*SingleMapReferences, error) {
	smr := &SingleMapReferences{
		config: gameConfig,
	}

	smr.maps = make(map[Location]map[int]SingleSmallMapReference)

	// Keepe.dat
	nextOffset := smr.AddLocation(Lord_Britishs_Castle, true, 5, 0)
	nextOffset = smr.AddLocation(Palace_of_Blackthorn, true, 5, nextOffset)
	nextOffset = smr.AddLocation(West_Britanny, false, 1, nextOffset)
	nextOffset = smr.AddLocation(North_Britanny, false, 1, nextOffset)
	nextOffset = smr.AddLocation(East_Britanny, false, 1, nextOffset)
	nextOffset = smr.AddLocation(Paws, false, 1, nextOffset)
	nextOffset = smr.AddLocation(Cove, false, 1, nextOffset)
	nextOffset = smr.AddLocation(Buccaneers_Den, false, 1, nextOffset)

	// Towne.dat
	nextOffset = smr.AddLocation(Moonglow, false, 2, 0)
	nextOffset = smr.AddLocation(Britain, false, 2, nextOffset)
	nextOffset = smr.AddLocation(Jhelom, false, 2, nextOffset)
	nextOffset = smr.AddLocation(Yew, true, 2, nextOffset)
	nextOffset = smr.AddLocation(Minoc, false, 2, nextOffset)
	nextOffset = smr.AddLocation(Trinsic, false, 2, nextOffset)
	nextOffset = smr.AddLocation(Skara_Brae, false, 2, nextOffset)
	nextOffset = smr.AddLocation(New_Magincia, false, 2, nextOffset)

	// Dwelling.dat
	nextOffset = smr.AddLocation(Fogsbane, false, 3, 0)
	nextOffset = smr.AddLocation(Stormcrow, false, 3, nextOffset)
	nextOffset = smr.AddLocation(Greyhaven, false, 3, nextOffset)
	nextOffset = smr.AddLocation(Waveguide, false, 3, nextOffset)
	nextOffset = smr.AddLocation(Iolos_Hut, false, 1, nextOffset)
	nextOffset = smr.AddLocation(Suteks_Hut, false, 1, nextOffset)
	nextOffset = smr.AddLocation(SinVraals_Hut, false, 1, nextOffset)
	nextOffset = smr.AddLocation(Grendels_Hut, false, 1, nextOffset)

	// Keep.dat
	nextOffset = smr.AddLocation(Ararat, false, 2, 0)
	nextOffset = smr.AddLocation(Bordermarch, false, 2, nextOffset)
	nextOffset = smr.AddLocation(Farthing, false, 1, nextOffset)
	nextOffset = smr.AddLocation(Windemere, false, 1, nextOffset)
	nextOffset = smr.AddLocation(Stonegate, false, 1, nextOffset)
	nextOffset = smr.AddLocation(Lycaeum, false, 3, nextOffset)
	nextOffset = smr.AddLocation(Empath_Abbey, false, 3, nextOffset)
	nextOffset = smr.AddLocation(Serpents_Hold, true, 3, nextOffset)

	// DUNGEONS.DAT
	//nextOffset = smr.AddLocation(Deceit, false, 8, 0)
	//nextOffset = smr.AddLocation(Despise, false, 8, nextOffset)
	//nextOffset = smr.AddLocation(Destard, false, 8, nextOffset)
	//nextOffset = smr.AddLocation(Wrong, false, 8, nextOffset)
	//nextOffset = smr.AddLocation(Covetous, false, 8, nextOffset)
	//nextOffset = smr.AddLocation(Shame, false, 8, nextOffset)
	//nextOffset = smr.AddLocation(Hythloth, false, 8, nextOffset)
	//nextOffset = smr.AddLocation(Doom, false, 8, nextOffset)

	//AddLocation(Britannia_Underworld, true, 2, 0)
	//AddLocation(Combat_resting_shrine, false, 1, 0)

	return smr, nil
}
