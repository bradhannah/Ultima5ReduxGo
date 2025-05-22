package references

import (
	"log"
	"os"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
)

const (
	npcsPerTown     = 32
	townsPerNPCFile = 8

	startingNpcTypeOffset         = sizeOfNPCSchedule * npcsPerTown
	startingNpcDialogNumberOffset = startingNpcTypeOffset + npcsPerTown

	singleTownSize = (sizeOfNPCSchedule * npcsPerTown) + (npcsPerTown * 2)
)

type NPCReferences struct {
	npcs []NPCReference
}

func NewNPCReferences(config *config.UltimaVConfiguration) *NPCReferences {
	allNpcs := &NPCReferences{}

	npcFiles := config.GetAllNpcFilePaths()

	for i, filePath := range npcFiles {
		npcs, err := getNPCsFromFile(filePath, i*townsPerNPCFile)
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

	npcs := make([]NPCReference, 0)
	for townIndex := 0; townIndex < townsPerNPCFile; townIndex++ {
		townOffset := singleTownSize * townIndex
		townRawData := npcRaw[townOffset : townOffset+singleTownSize]
		for npcIndex := 0; npcIndex < npcsPerTown; npcIndex++ {
			npc := NPCReference{}
			npc.Location = Location(locationOffset + townIndex + 1)
			npc.Schedule = CreateNPCSchedule(townRawData[npcIndex*sizeOfNPCSchedule : (npcIndex*sizeOfNPCSchedule)+sizeOfNPCSchedule])
			npc.npcType = NPCType(townRawData[startingNpcTypeOffset+npcIndex])
			npc.DialogNumber = townRawData[startingNpcDialogNumberOffset+npcIndex]

			sprite := npc.GetSpriteIndex()
			if sprite.IsHorseUnBoarded() || sprite == 274 || sprite == 275 {
				npc.Schedule.OverrideAllAI(HorseWander)
			}

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
