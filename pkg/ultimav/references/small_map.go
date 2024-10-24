package references

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/legacy"
)

const (
	XSmallMapTiles         = Coordinate(32)
	YSmallMapTiles         = Coordinate(32)
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
		return legacy.CASTLE_DAT
	case Towne:
		return legacy.TOWNE_DAT
	case Dwelling:
		return legacy.DWELLING_DAT
	case Keep:
		return legacy.KEEP_DAT
	case Dungeon:
		return legacy.DUNGEON_DAT
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

func NewSmallMapReferences(gameConfig *config.UltimaVConfiguration, dataOvl *DataOvl) (*LocationReferences, error) {
	smr := newSingleMapReferences(gameConfig, dataOvl)

	smr.maps = make(map[Location]*SmallLocationReference)

	// Keepe.dat
	nextOffset := smr.addLocation(Lord_Britishs_Castle, true, 5, 0)
	nextOffset = smr.addLocation(Palace_of_Blackthorn, true, 5, nextOffset)
	nextOffset = smr.addLocation(West_Britanny, false, 1, nextOffset)
	nextOffset = smr.addLocation(North_Britanny, false, 1, nextOffset)
	nextOffset = smr.addLocation(East_Britanny, false, 1, nextOffset)
	nextOffset = smr.addLocation(Paws, false, 1, nextOffset)
	nextOffset = smr.addLocation(Cove, false, 1, nextOffset)
	nextOffset = smr.addLocation(Buccaneers_Den, false, 1, nextOffset)

	// Towne.dat
	nextOffset = smr.addLocation(Moonglow, false, 2, 0)
	nextOffset = smr.addLocation(Britain, false, 2, nextOffset)
	nextOffset = smr.addLocation(Jhelom, false, 2, nextOffset)
	nextOffset = smr.addLocation(Yew, true, 2, nextOffset)
	nextOffset = smr.addLocation(Minoc, false, 2, nextOffset)
	nextOffset = smr.addLocation(Trinsic, false, 2, nextOffset)
	nextOffset = smr.addLocation(Skara_Brae, false, 2, nextOffset)
	nextOffset = smr.addLocation(New_Magincia, false, 2, nextOffset)

	// Dwelling.dat
	nextOffset = smr.addLocation(Fogsbane, false, 3, 0)
	nextOffset = smr.addLocation(Stormcrow, false, 3, nextOffset)
	nextOffset = smr.addLocation(Greyhaven, false, 3, nextOffset)
	nextOffset = smr.addLocation(Waveguide, false, 3, nextOffset)
	nextOffset = smr.addLocation(Iolos_Hut, false, 1, nextOffset)
	nextOffset = smr.addLocation(Suteks_Hut, false, 1, nextOffset)
	nextOffset = smr.addLocation(SinVraals_Hut, false, 1, nextOffset)
	nextOffset = smr.addLocation(Grendels_Hut, false, 1, nextOffset)

	// Keep.dat
	nextOffset = smr.addLocation(Ararat, false, 2, 0)
	nextOffset = smr.addLocation(Bordermarch, false, 2, nextOffset)
	nextOffset = smr.addLocation(Farthing, false, 1, nextOffset)
	nextOffset = smr.addLocation(Windemere, false, 1, nextOffset)
	nextOffset = smr.addLocation(Stonegate, false, 1, nextOffset)
	nextOffset = smr.addLocation(Lycaeum, false, 3, nextOffset)
	nextOffset = smr.addLocation(Empath_Abbey, false, 3, nextOffset)
	nextOffset = smr.addLocation(Serpents_Hold, true, 3, nextOffset)

	// DUNGEONS.DAT
	//nextOffset = smr.addLocation(Deceit, false, 8, 0)
	//nextOffset = smr.addLocation(Despise, false, 8, nextOffset)
	//nextOffset = smr.addLocation(Destard, false, 8, nextOffset)
	//nextOffset = smr.addLocation(Wrong, false, 8, nextOffset)
	//nextOffset = smr.addLocation(Covetous, false, 8, nextOffset)
	//nextOffset = smr.addLocation(Shame, false, 8, nextOffset)
	//nextOffset = smr.addLocation(Hythloth, false, 8, nextOffset)
	//nextOffset = smr.addLocation(Doom, false, 8, nextOffset)

	//addLocation(Britannia_Underworld, true, 2, 0)
	//addLocation(Combat_resting_shrine, false, 1, 0)

	return smr, nil
}
