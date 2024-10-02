package references

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"log"
	"os"
	"path"
)

type Location byte

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

type SingleSmallMapReference struct {
	rawData [XSmallMapTiles][YSmallMapTiles]byte
}

type SingleMapReferences struct {
	maps   map[Location]map[int]SingleSmallMapReference
	config *config.UltimaVConfiguration
}

func (s *SingleMapReferences) AddLocation(location Location, bHasBasement bool, nFloors int, nOffset int) int {
	maps := make(map[int]SingleSmallMapReference)
	// get the file
	mapFileAndPath := path.Join(s.config.DataFilePath, getSmallMapFile(getMapMasterFromLocation(location)))

	// a bit wasteful since it's open a few times, but... fast computers...
	theChunksSerial, err := os.ReadFile(mapFileAndPath)
	if err != nil {
		log.Fatal("OOf")
	}

	floorModifier := 0
	if bHasBasement {
		floorModifier = -1
	}

	for i := 0; i < nFloors; i++ {
		smr := SingleSmallMapReference{}
		var x, y int
		for x = 0; x < int(XSmallMapTiles); x++ {
			for y = 0; y < int(YSmallMapTiles); y++ {
				byteIndex := nOffset + (i * smallMapSizeInBytes) + x + (y * int(YSmallMapTiles))
				smr.rawData[x][y] = theChunksSerial[byteIndex]
			}
		}
		maps[i+floorModifier] = smr
	}

	// returns the next offset - a handy way of keeping count
	s.maps[location] = maps
	return nFloors*smallMapSizeInBytes + nOffset
}
