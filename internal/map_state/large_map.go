package map_state

import (
	"log"

	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// var nChanceToGenerateEnemy = 32

const MaxTileDistanceBeforeCleanup = 22

func (m *MapState) UpdateLargeMap() {
	if m.PlayerLocation.Location.GetMapType() != references2.LargeMapType {
		log.Fatalf("Expected large map type, got %d", m.PlayerLocation.Location.GetMapType())
	}
}

func (m *MapState) IsOverworld() bool {
	return m.PlayerLocation.Floor == references2.FloorNumber(references2.OVERWORLD) && m.PlayerLocation.Location.GetMapType() == references2.LargeMapType
}

func (m *MapState) IsUnderworld() bool {
	return m.PlayerLocation.Floor == references2.FloorNumber(references2.UNDERWORLD) && m.PlayerLocation.Location.GetMapType() == references2.LargeMapType
}
