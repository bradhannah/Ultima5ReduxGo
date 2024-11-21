package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

type NPC struct {
	NPCReference references.NPCReference

	Position references.Position
	Floor    references.FloorNumber
	AiType   references.AiType
}

func NewNPC(npcReference references.NPCReference) NPC {
	npc := NPC{}
	npc.NPCReference = npcReference
	return npc
}

func (npc *NPC) IsEmptyNPC() bool {

	return npc.NPCReference.Type == 0 || (npc.NPCReference.Schedule.X[0] == 0 && npc.NPCReference.Schedule.Y[0] == 0)
}
