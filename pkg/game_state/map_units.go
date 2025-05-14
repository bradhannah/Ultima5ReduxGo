package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type MapUnits []MapUnit

const MAXIMUM_NPCS_PER_MAP = 32

func (n *MapUnits) getNextAvailableNPCIndexNumber() int {
	if len(*n) < MAXIMUM_NPCS_PER_MAP {
		return len(*n)
	}
	// for i, mu := range *n {
	// 	if mu.IsEmptyMapUnit() {
	// 		return i
	// 	}
	// }
	return -1
}

func (n *MapUnits) RemoveNPCAtPosition(position references.Position) bool {
	for _, mu := range *n {
		if mu.Pos() == position {
			mu.SetVisible(false)
			return true
		}
	}
	return false
}

func (n *MapUnits) AddVehicle(
	vehicle references.PartyVehicle,
	position references.Position,
	floorNumber references.FloorNumber) bool {

	index := n.getNextAvailableNPCIndexNumber()
	if index == -1 {
		return false
	}

	npcRef := references.NewNPCReferenceForVehicle(vehicle, position, floorNumber)

	npc := NewNPCFriendly(*npcRef, index)
	npc.SetPos(position)
	npc.SetFloor(0)
	npc.SetVisible(true)

	*n = append(*n, npc)
	//(*n)[index] = npc
	return true
}

func (n *MapUnits) createFreshXyOccupiedMap() *XyOccupiedMap {
	xy := make(XyOccupiedMap)
	for _, mu := range *n {

		if mu.IsEmptyMapUnit() || !mu.IsVisible() {
			continue
		}
		_, exists := xy[int(mu.Pos().X)]
		if !exists {
			xy[int(mu.Pos().X)] = make(map[int]bool)
		}
		xy[int(mu.Pos().X)][int(mu.Pos().Y)] = true
	}
	return &xy
}

func getMapUnitAsFriendlyOrNil(mu *MapUnit) *NPCFriendly {
	if (*mu).IsEmptyMapUnit() {
		return nil
	}
	switch mapUnit := (*mu).(type) {
	case *NPCFriendly:
		return mapUnit
	case *NPCEnemy:
		return nil
	default:
		log.Fatal("Unknown Map Unit Type")
	}
	return nil
}

func getMapUnitAsEnemyOrNil(mu *MapUnit) *NPCEnemy {
	if (*mu).IsEmptyMapUnit() {
		return nil
	}
	switch mapUnit := (*mu).(type) {
	case *NPCFriendly:
		return nil
	case *NPCEnemy:
		return mapUnit
	default:
		log.Fatal("Unknown Map Unit Type")
	}
	return nil
}
