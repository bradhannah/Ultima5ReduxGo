package references

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/util"
)

type Location byte

//go:generate stringer -type=Location
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

type SmallMapReference struct {
	rawData              map[int]*[XSmallMapTiles][YSmallMapTiles]byte
	Location             Location
	FriendlyLocationName string
	EnteringText         string
	SmallMapType         SmallMapMasterTypes
	//config   *config.UltimaVConfiguration
}

func NewSingleSmallMapReference(location Location, dataOvl *DataOvl) *SmallMapReference {
	smr := &SmallMapReference{}
	smr.Location = location
	smr.rawData = make(map[int]*[XSmallMapTiles][YSmallMapTiles]byte)
	// NOTE: this needs to be moved to a higher level
	smr.FriendlyLocationName = dataOvl.LocationNames[location]
	smr.SmallMapType = getMapMasterFromLocation(location)
	smr.EnteringText = smr.getEnteringText()
	return smr
}

func (s *SmallMapReference) AddBlankFloor(index int) {
	// Initialize the array
	tileData := &[XSmallMapTiles][YSmallMapTiles]byte{}

	// Add the initialized array to the map at the given index
	s.rawData[index] = tileData
}

func (s *SmallMapReference) GetTileNumber(nFloor int, position *Position) int {
	return int(s.rawData[nFloor][position.X][position.Y])
}

func (s *SmallMapReference) GetEnteringText() string {
	return s.Location.String()
}

func (s *SmallMapReference) getEnteringText() string {
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
