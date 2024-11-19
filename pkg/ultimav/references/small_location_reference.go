package references

import (
	"fmt"
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/util"
)

type SmallLocationReference struct {
	rawData              map[int]*[XSmallMapTiles][YSmallMapTiles]byte
	Location             Location
	FriendlyLocationName string
	EnteringText         string
	SmallMapType         SmallMapMasterTypes
	ListOfFloors         []FloorNumber

	npcRefs *[]NPCReference

	// config   *config.UltimaVConfiguration
}

func NewSingleSmallMapReference(
	location Location,
	npcRefs *[]NPCReference,
	dataOvl *DataOvl) *SmallLocationReference {

	smr := &SmallLocationReference{}
	smr.Location = location
	smr.rawData = make(map[int]*[XSmallMapTiles][YSmallMapTiles]byte)

	// NOTE: this needs to be moved to a higher level
	smr.FriendlyLocationName = dataOvl.LocationNames[location]
	smr.SmallMapType = getMapMasterFromLocation(location)
	smr.EnteringText = smr.getEnteringText()

	smr.npcRefs = npcRefs
	return smr
}

func (s *SmallLocationReference) addBlankFloor(index int) {
	// Initialize the array
	tileData := &[XSmallMapTiles][YSmallMapTiles]byte{}

	// Add the initialized array to the map at the given index
	s.rawData[index] = tileData

	s.ListOfFloors = s.getListOfFloors()
}

func (s *SmallLocationReference) GetFloorMinMax() (FloorNumber, FloorNumber) {
	if s.HasBasement() {
		return -1, FloorNumber(len(s.rawData) - 2)
	}
	return 0, FloorNumber(len(s.rawData) - 1)
}

func (s *SmallLocationReference) HasBasement() bool {
	_, ok := s.rawData[-1]
	return ok
}

// func (s *SmallLocationReference) GetTileNumberWithAnimation(nFloor FloorNumber, position *Position) indexes.SpriteIndex {
//
//	mainTile := s.GetTileNumber(nFloor, position.X, position.Y)
//	//indexes.SpriteIndex(s.rawData[int(nFloor)][position.X][position.Y])
//
//	if (mainTile >= indexes.Waterfall_KeyIndex && mainTile <= indexes.Waterfall_KeyIndex+3) || mainTile == indexes.Fountain_KeyIndex || mainTile >= indexes.AvatarSittingAndEatingFacingDown {
//		return sprites.GetSpriteIndexWithAnimationBySpriteIndex(mainTile)
//	}
//	return mainTile
// }

func (s *SmallLocationReference) GetTileNumber(nFloor FloorNumber, x Coordinate, y Coordinate) indexes.SpriteIndex {
	return indexes.SpriteIndex(s.rawData[int(nFloor)][x][y])
}

func (s *SmallLocationReference) GetEnteringText() string {
	return s.Location.String()
}

func (s *SmallLocationReference) getEnteringText() string {
	switch s.Location {
	case Lord_Britishs_Castle:
		return "Enter the Castle of Lord British!"
	case Palace_of_Blackthorn:
		return "Enter the Palace of Lord Blackthorn"
	case Fogsbane, Stormcrow, Waveguide, Greyhaven:
		return fmt.Sprintf("Enter Lighthouse\n\n%s", util.GetCenteredText(s.FriendlyLocationName))
	case West_Britanny, East_Britanny, North_Britanny, Paws, Cove:
		return fmt.Sprintf("Enter Village\n\n%s", util.GetCenteredText(s.FriendlyLocationName))
	case Moonglow, Britain, Jhelom, Yew, Minoc, Trinsic, Skara_Brae, New_Magincia:
		return fmt.Sprintf("Enter Towne\n\n%s", util.GetCenteredText(s.FriendlyLocationName))
	case Iolos_Hut, Grendels_Hut, SinVraals_Hut, Suteks_Hut:
		return fmt.Sprintf("Enter Hut\n\n%s", util.GetCenteredText(s.FriendlyLocationName))
	case Ararat:
		return fmt.Sprintf("Enter Ruins\n\n%s", util.GetCenteredText(s.FriendlyLocationName))
	case Bordermarch, Farthing, Windemere, Stonegate, Lycaeum, Empath_Abbey, Serpents_Hold, Buccaneers_Den:
		return fmt.Sprintf("Enter Keep\n\n%s", util.GetCenteredText(s.FriendlyLocationName))
	case Deceit, Despise, Destard, Wrong, Covetous, Shame, Hythloth, Doom:
		return fmt.Sprintf("Enter Dungeon\n\n%s", util.GetCenteredText(s.FriendlyLocationName))
	}
	return "NOT IMPLEMENTED"
}

func (s *SmallLocationReference) GetOuterTile() indexes.SpriteIndex {
	switch s.Location {
	case SinVraals_Hut:
		return indexes.Desert
	case Grendels_Hut:
		return indexes.Swamp
	default:
		return indexes.Grass
	}
}

func (s *SmallLocationReference) GetNumberOfFloors() int {
	return len(s.rawData)
}

func (s *SmallLocationReference) GetMaxY() Coordinate {
	return s.GetMaxX()
}

func (s *SmallLocationReference) GetMaxX() Coordinate {
	if s.Location.GetMapType() == LargeMapType {
		return XLargeMapTiles - 1
	} else if s.Location.GetMapType() == SmallMapType {
		return XSmallMapTiles - 1
	}
	log.Fatal("missing max tiles")
	return 0
}

func (s *SmallLocationReference) getListOfFloors() []FloorNumber {
	numFloors := s.GetNumberOfFloors()
	startIndex := FloorNumber(0)

	if s.HasBasement() {
		startIndex = -1
	}

	// Initialize the floors slice with values directly
	floors := make([]FloorNumber, numFloors)
	for i := FloorNumber(0); i < FloorNumber(numFloors); i++ {
		floors[i] = startIndex + i
	}

	return floors
}

func (s *SmallLocationReference) CanGoUpOneFloor(currentFloor FloorNumber) bool {
	floorIndex := 0
	if s.HasBasement() {
		floorIndex = -1
	}
	nextFloor := currentFloor + 1
	return int(nextFloor) < s.GetNumberOfFloors()+floorIndex
}

func (s *SmallLocationReference) CanGoDownOneFloor(currentFloor FloorNumber) bool {
	if currentFloor < 0 {
		if s.HasBasement() {
			return true
		}
		return false
	}
	return true
}

func (s *SmallLocationReference) GetNPCs() *[]NPCReference {
	return s.npcRefs
}
