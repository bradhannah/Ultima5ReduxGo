package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

type NPC struct {
	NPCReference references.NPCReference
}

func NewNPC(npcReference references.NPCReference) NPC {
	npc := NPC{}
	npc.NPCReference = npcReference
	return npc
}
