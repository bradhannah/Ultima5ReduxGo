package references

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"log"
	"os"
	"path"
)

// LocationReferences a collection of SingleSmallMapReferences. Provides an easier way to keep
// the raw map data organized and accessible
type LocationReferences struct {
	maps           map[Location]*SmallLocationReference
	config         *config.UltimaVConfiguration
	dataOvl        *DataOvl
	WorldLocations *WorldLocations
}

func (s *LocationReferences) GetLocationReference(location Location) *SmallLocationReference {
	return s.maps[location]
}

func newSingleMapReferences(config *config.UltimaVConfiguration, dataOvl *DataOvl) *LocationReferences {
	smr := &LocationReferences{}
	smr.config = config
	smr.dataOvl = dataOvl
	smr.WorldLocations = NewWorldLocations(smr.config)
	return smr
}

func (s *LocationReferences) addLocation(location Location, bHasBasement bool, nFloors int, nOffset int) int {
	maps := make(map[int]*SmallLocationReference)
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

	smr := NewSingleSmallMapReference(location, s.dataOvl) //SmallLocationReference{}
	for i := 0; i < nFloors; i++ {
		actualFloor := i + floorModifier
		smr.addBlankFloor(actualFloor)
		var x, y int
		for x = 0; x < int(XSmallMapTiles); x++ {
			for y = 0; y < int(YSmallMapTiles); y++ {
				byteIndex := nOffset + (i * smallMapSizeInBytes) + x + (y * int(YSmallMapTiles))
				smr.rawData[i+floorModifier][x][y] = theChunksSerial[byteIndex]
			}
		}
		maps[actualFloor] = smr
	}

	// returns the next offset - a handy way of keeping count
	s.maps[location] = smr
	return nFloors*smallMapSizeInBytes + nOffset
}
