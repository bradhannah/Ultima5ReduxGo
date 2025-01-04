package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPC struct {
	NPCReference references.NPCReference

	Position references.Position
	Floor    references.FloorNumber
	AiType   references.AiType

	Visible bool

	NPCNum int

	AStarMap *AStarMap

	CurrentPath *[]references.Position
}

func NewNPC(npcReference references.NPCReference, npcNum int) NPC {
	npc := NPC{}
	npc.NPCReference = npcReference
	npc.NPCNum = npcNum

	npc.AStarMap = NewAStarMap()
	npc.CurrentPath = nil

	// npc.Position = npcReference.Position
	// npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate()
	// npc.Floor = npcReference.

	if !npc.IsEmptyNPC() {
		npc.Visible = true
	}

	return npc
}

func (npc *NPC) IsEmptyNPC() bool {
	return npc.NPCReference.Type == 0 || (npc.NPCReference.Schedule.X[0] == 0 && npc.NPCReference.Schedule.Y[0] == 0)
}

func (npc *NPC) DequeueNextPosition() references.Position {
	if !npc.HasAPathAlreadyCalculated() {
		log.Fatal("NPC has no path calculated")
	}
	pos := (*npc.CurrentPath)[0]
	*npc.CurrentPath = (*npc.CurrentPath)[1:]
	return pos
}

func (npc *NPC) HasAPathAlreadyCalculated() bool {
	return npc.CurrentPath != nil && len(*npc.CurrentPath) > 0
}
