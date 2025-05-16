package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type MapUnits []MapUnit

const MAXIMUM_NPCS_PER_MAP = 32

func (n *MapUnits) getNextAvailableNPCIndexNumber() int {
	if len(*n) < MAXIMUM_NPCS_PER_MAP {
		return len(*n)
	}

	return -1
}

func (m *MapUnits) RemoveNPCAtPosition(position references.Position) bool {
	// TODO: this doesn't really remove them correctly... I think we don't have to care too much
	// about their position in the slice as long as we have the index recorded at the time of
	// reading in the NPC reference data
	for _, mu := range *m {
		if mu.Pos() == position {
			mu.SetVisible(false)
			return true
		}
	}
	return false
}

func (m *MapUnits) GetMapUnitAtPositionOrNil(pos references.Position) *MapUnit {
	for _, mu := range *m {
		if mu.Pos() == pos {
			return &mu
		}
	}
	return nil
}

func (m *MapUnits) GetVehicleAtPositionOrNil(pos references.Position) *NPCVehicle {
	mu := m.GetMapUnitAtPositionOrNil(pos)

	if mu == nil {
		return nil
	}

	if vehicle, ok := (*mu).(*NPCVehicle); ok {
		return vehicle

	}
	return nil
}

func (n *MapUnits) AddVehicle(vehicle NPCVehicle) bool {

	index := n.getNextAvailableNPCIndexNumber()
	if index == -1 {
		return false
	}

	vehicle.mapUnitDetails.NPCNum = index

	// npcRef := references.NewNPCReferenceForVehicle(vehicle, position, floorNumber)

	// npc := NewNPCFriendly(*npcRef, index)
	// npc.SetPos(position)
	// npc.SetFloor(0)
	//npc.SetVisible(true)
	vehicle.SetVisible(true)

	*n = append(*n, &vehicle)

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

func As[T interface{ IsEmptyMapUnit() bool }](mu MapUnit) (T, bool) {
	if v, ok := mu.(T); ok && !v.IsEmptyMapUnit() {
		return v, true
	}
	var zero T
	return zero, false
}

// func getMapUnitAsFriendlyOrNil(mu MapUnit) *NPCFriendly {
// 	if (*mu).IsEmptyMapUnit() {
// 		return nil
// 	}
// 	switch mapUnit := (*mu).(type) {
// 	case *NPCFriendly:
// 		return mapUnit
// 	case *NPCVehicle:
// 	case *NPCEnemy:
// 		return nil
// 	default:
// 		log.Fatal("Unknown Map Unit Type")
// 	}
// 	return nil
// }

// func getMapUnitAsEnemyOrNil(mu MapUnit) *NPCEnemy {
// 	if (*mu).IsEmptyMapUnit() {
// 		return nil
// 	}
// 	switch mapUnit := (*mu).(type) {
// 	case *NPCVehicle:
// 	case *NPCFriendly:
// 		return nil
// 	case *NPCEnemy:
// 		return mapUnit
// 	default:
// 		log.Fatal("Unknown Map Unit Type")
// 	}
// 	return nil
// }
