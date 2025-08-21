package references

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
)

// LocationReferences a collection of SingleSmallMapReferences. Provides an easier way to keep
// the raw map data organized and accessible
type LocationReferences struct {
	maps           map[Location]*SmallLocationReference
	mapsByStr      map[string]*SmallLocationReference
	config         *config.UltimaVConfiguration
	dataOvl        *DataOvl
	WorldLocations *WorldLocations `json:"world_locations" yaml:"world_locations"`

	npcRefs *NPCReferences
}

const MapsPerLocationFile = 8

func (s *LocationReferences) GetLocationReference(location Location) *SmallLocationReference {
	return s.maps[location]
}

func (s *LocationReferences) GetSmallLocationReference(name string) *SmallLocationReference {
	return s.mapsByStr[strings.ToLower(name)]
}

func newSingleMapReferences(config *config.UltimaVConfiguration, dataOvl *DataOvl) *LocationReferences {
	smr := &LocationReferences{}
	smr.config = config
	smr.dataOvl = dataOvl
	smr.WorldLocations = NewWorldLocations(smr.config)

	smr.maps = make(map[Location]*SmallLocationReference)
	smr.mapsByStr = make(map[string]*SmallLocationReference)

	smr.npcRefs = NewNPCReferences(smr.config)

	return smr
}

func (s *LocationReferences) addLocation(location Location, bHasBasement bool, nFloors, nOffset int) int {
	maps := make(map[int]*SmallLocationReference)
	// get the file
	mapFileAndPath := path.Join(s.config.SavedConfigData.DataFilePath, GetSmallMapFile(getMapMasterFromLocation(location)))

	// a bit wasteful since it's open a few times, but... fast computers...
	theChunksSerial, err := os.ReadFile(mapFileAndPath)
	if err != nil {
		log.Fatal("OOf")
	}

	floorModifier := 0
	if bHasBasement {
		floorModifier = -1
	}

	smr := NewSingleSmallMapReference(
		location,
		s.npcRefs.GetNPCReferencesByLocation(location),
		s.dataOvl)

	for floor := range nFloors {
		actualFloor := floor + floorModifier
		smr.addBlankFloor(actualFloor)

		for x := range int(XSmallMapTiles) {
			for y := range int(YSmallMapTiles) {
				byteIndex := nOffset + (floor * smallMapSizeInBytes) + x + (y * int(YSmallMapTiles))
				smr.rawData[floor+floorModifier][x][y] = theChunksSerial[byteIndex]
			}
		}

		maps[actualFloor] = smr
	}

	// returns the next offset - a handy way of keeping count
	s.maps[location] = smr
	s.mapsByStr[strings.ToLower(location.String())] = smr

	return nFloors*smallMapSizeInBytes + nOffset
}

func GetLocationBySmallMapAndIndex(smallMap SmallMapMasterTypes, index int) Location {
	var nMultiplier int

	switch smallMap {
	case Towne:
		nMultiplier = 0
	case Dwelling:
		nMultiplier = 1
	case Castle:
		nMultiplier = 2
	case Keep:
		nMultiplier = 3
	default:
		log.Fatalf("Unhandled small map type: %v", smallMap)
	}

	return Location((MapsPerLocationFile * nMultiplier) + index)

	//{
	//switch (location)
	//{
	//case Location.Lord_Britishs_Castle:
	//case Location.Palace_of_Blackthorn:
	//case Location.East_Britanny:
	//case Location.West_Britanny:
	//case Location.North_Britanny:
	//case Location.Paws:
	//case Location.Cove:
	//case Location.Buccaneers_Den:
	//return SmallMapMasterFiles.Castle;
	//case Location.Moonglow:
	//case Location.Britain:
	//case Location.Jhelom:
	//case Location.Yew:
	//case Location.Minoc:
	//case Location.Trinsic:
	//case Location.Skara_Brae:
	//case Location.New_Magincia:
	//return SmallMapMasterFiles.Towne;
	//case Location.Fogsbane:
	//case Location.Stormcrow:
	//case Location.Waveguide:
	//case Location.Greyhaven:
	//case Location.Iolos_Hut:
	////case _location.spektran
	//case Location.Suteks_Hut:
	//case Location.SinVraals_Hut:
	//case Location.Grendels_Hut:
	//return SmallMapMasterFiles.Dwelling;
	//case Location.Ararat:
	//case Location.Bordermarch:
	//case Location.Farthing:
	//case Location.Windemere:
	//case Location.Stonegate:
	//case Location.Lycaeum:
	//case Location.Empath_Abbey:
	//case Location.Serpents_Hold:
	//return SmallMapMasterFiles.Keep;
	//case Location.Deceit:
	//case Location.Despise:
	//case Location.Destard:
	//case Location.Wrong:
	//case Location.Covetous:
	//case Location.Shame:
	//case Location.Hythloth:
	//case Location.Doom:
	//return SmallMapMasterFiles.Dungeon;
	//case Location.Combat_resting_shrine:
	//return SmallMapMasterFiles.CutOrIntroScene;
	//case Location.Britannia_Underworld:
	//return SmallMapMasterFiles.None;
	//default:
	//throw new Ultima5ReduxException("EH?");
	//}
	//}
}
