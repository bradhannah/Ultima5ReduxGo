package map_units

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type MapUnits []MapUnit

const MaximumNpcsPerMap = 32

func (m *MapUnits) getNextAvailableNPCIndexNumber() int {
	if len(*m) < MaximumNpcsPerMap {
		return len(*m)
	}

	return -1
}

func (m *MapUnits) RemoveNPCAtPosition(position references2.Position) bool {
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

func (m *MapUnits) GetMapUnitAtPositionOrNil(pos references2.Position) *MapUnit {
	for _, mu := range *m {
		if mu.Pos() == pos {
			return &mu
		}
	}
	return nil
}

func (m *MapUnits) GetVehicleAtPositionOrNil(pos references2.Position) *NPCFriendly {
	mu := m.GetMapUnitAtPositionOrNil(pos)

	if mu == nil {
		return nil
	}

	if friendly, ok := (*mu).(*NPCFriendly); ok {
		if friendly.NPCReference.GetNPCType() == references2.Vehicle {
			return friendly
		}

	}
	return nil
}

func (m *MapUnits) AddVehicle(vehicle NPCFriendly) bool {

	index := m.getNextAvailableNPCIndexNumber()
	if index == -1 {
		return false
	}

	vehicle.mapUnitDetails.NPCNum = index
	vehicle.SetPos(references2.Position{X: references2.Coordinate(vehicle.NPCReference.Schedule.X[0]), Y: references2.Coordinate(vehicle.NPCReference.Schedule.Y[0])})

	vehicle.SetVisible(true)

	*m = append(*m, &vehicle)

	return true
}

func (m *MapUnits) CreateFreshXyOccupiedMap() *XyOccupiedMap {
	xy := make(XyOccupiedMap)
	for _, mu := range *m {

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

// func As[T interface{ IsEmptyMapUnit() bool }](mu MapUnit) (T, bool) {
// 	if v, ok := mu.(T); ok && !v.IsEmptyMapUnit() {
// 		return v, true
// 	}
// 	var zero T
// 	return zero, false
// }

type XyOccupiedMap map[int]map[int]bool
