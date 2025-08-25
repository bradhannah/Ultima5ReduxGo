package map_units

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

type (
	MapUnits      []MapUnit
	XyOccupiedMap map[int]map[int]bool
)

const MaximumNpcsPerMap = 32

func (m *MapUnits) getNextAvailableNPCIndexNumber() int {
	if len(*m) < MaximumNpcsPerMap {
		return len(*m)
	}

	return -1
}

func (m *MapUnits) RemoveNPCAtPosition(position references.Position) bool {
	// NOTE: the deletion could cause indexing issues - but shouldn't
	for i, mu := range *m {
		if mu.Pos() == position {
			mu.SetVisible(false)

			*m = helpers.Delete(*m, i)

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

func (m *MapUnits) AddVehicle(vehicle NPCFriendly) bool {
	index := m.getNextAvailableNPCIndexNumber()
	if index == -1 {
		return false
	}

	vehicle.mapUnitDetails.NPCNum = index
	vehicle.SetPos(references.Position{
		X: references.Coordinate(vehicle.NPCReference.Schedule.X[0]),
		Y: references.Coordinate(vehicle.NPCReference.Schedule.Y[0]),
	})

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
