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

func (m *MapUnits) GetVehicleAtPositionOrNil(pos references.Position) *NPCFriendly {
	mu := m.GetMapUnitAtPositionOrNil(pos)

	if mu == nil {
		return nil
	}

	if friendly, ok := (*mu).(*NPCFriendly); ok {
		if friendly.NPCReference.GetNPCType() == references.Vehicle {
			return friendly
		}

	}
	return nil
}

func (n *MapUnits) AddVehicle(vehicle NPCFriendly) bool {

	index := n.getNextAvailableNPCIndexNumber()
	if index == -1 {
		return false
	}

	vehicle.mapUnitDetails.NPCNum = index
	vehicle.SetPos(references.Position{X: references.Coordinate(vehicle.NPCReference.Schedule.X[0]), Y: references.Coordinate(vehicle.NPCReference.Schedule.Y[0])})

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
