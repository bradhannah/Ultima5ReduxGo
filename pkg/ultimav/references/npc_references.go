package references

import (
	"log"
	"os"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
)

const (
	npcsPerTown     = 32
	townsPerNPCFile = 8

	startingNpcTypeOffset         = sizeOfNPCSchedule * npcsPerTown
	startingNpcDialogNumberOffset = startingNpcTypeOffset + npcsPerTown

	singleTownSize = (sizeOfNPCSchedule * npcsPerTown) + (npcsPerTown * 2)
)

type NPCType byte

const (
	Blacksmith  NPCType = 0x81
	Barkeeper           = 0x82
	HorseSeller         = 0x83
	Shipwright          = 0x84
	Healer              = 0x87
	InnKeeper           = 0x88
	MagicSeller         = 0x85
	GuildMaster         = 0x86
	NoStatedNpc         = 0xFF
	Guard               = 0xFE
	WishingWell         = 0xFD
	// unknowns may be crown and sandlewood box
)

type NPCReferences struct {
	npcs []NPCReference
}

type NPCReference struct {
	Position     Position
	Location     Location
	DialogNumber byte
	Schedule     NPCSchedule
	Type         NPCType
	// script TalkScript
}

func NewNPCReferences(config *config.UltimaVConfiguration) *NPCReferences {
	allNpcs := &NPCReferences{}

	npcFiles := config.GetAllNpcFilePaths()

	for i, filePath := range npcFiles {
		npcs, err := getNPCsFromFile(filePath, i/townsPerNPCFile)
		if err != nil {
			log.Fatal(err)
		}
		allNpcs.npcs = append(allNpcs.npcs, npcs...)
	}

	return allNpcs
}

func getNPCsFromFile(path string, locationOffset int) ([]NPCReference, error) {
	npcRaw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	npcs := make([]NPCReference, 0) // npcsPerTown*townsPerNPCFile)
	for townIndex := 0; townIndex < townsPerNPCFile; townIndex++ {
		townOffset := singleTownSize * townIndex
		townRawData := npcRaw[townOffset : townOffset+singleTownSize]
		for npcIndex := 0; npcIndex < npcsPerTown; npcIndex++ {
			npc := NPCReference{}
			npc.Location = Location(locationOffset + townIndex + 1)

			npc.Schedule = CreateNPCSchedule(townRawData[npcIndex*sizeOfNPCSchedule : (npcIndex*sizeOfNPCSchedule)+sizeOfNPCSchedule])
			npc.DialogNumber = townRawData[startingNpcDialogNumberOffset+npcIndex]
			npc.Type = NPCType(townRawData[startingNpcTypeOffset+npcIndex])
			npcs = append(npcs, npc)
		}
	}
	return npcs, nil
}

func (n *NPCReferences) getNPCIndexesByLocation(location Location) (startIndex, endIndex int) {
	adjLocationIndex := int(location) - 1
	return adjLocationIndex * npcsPerTown, adjLocationIndex*npcsPerTown + npcsPerTown
}

func (n *NPCReferences) GetNPCReferencesByLocation(location Location) *[]NPCReference {
	startIndex, endIndex := n.getNPCIndexesByLocation(location)
	npcs := n.npcs[startIndex:endIndex]
	return &npcs
}
