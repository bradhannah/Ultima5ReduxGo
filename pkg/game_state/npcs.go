package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

type NPCS []*NPC

func (n *NPCS) getNextAvailableNPCIndexNumber() int {
	for i, npc := range *n {
		if npc.IsEmptyNPC() {
			return i
		}
	}
	return -1
}

func (n *NPCS) RemoveNPCAtPosition(position references.Position) bool {
	for _, npc := range *n {
		if npc.Position == position {
			npc.Visible = false
			return true
		}
	}
	return false
}

func (n *NPCS) AddVehicle(
	vehicle references.PartyVehicle,
	position references.Position,
	floorNumber references.FloorNumber) bool {

	index := n.getNextAvailableNPCIndexNumber()
	if index == -1 {
		return false
	}

	npcRef := references.NewNPCReferenceForVehicle(vehicle, position, floorNumber)

	npc := NewNPC(*npcRef, index)
	(*n)[index] = &npc
	return true
}
