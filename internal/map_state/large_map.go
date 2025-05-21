package map_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

// var nChanceToGenerateEnemy = 32

const MaxTileDistanceBeforeCleanup = 22

func (m *MapState) UpdateLargeMap() {
	if m.PlayerLocation.Location.GetMapType() != references.LargeMapType {
		log.Fatalf("Expected large map type, got %d", m.PlayerLocation.Location.GetMapType())
	}
}

func (m *MapState) IsOverworld() bool {
	return m.PlayerLocation.Floor == references.FloorNumber(references.OVERWORLD) && m.PlayerLocation.Location.GetMapType() == references.LargeMapType
}

func (m *MapState) IsUnderworld() bool {
	return m.PlayerLocation.Floor == references.FloorNumber(references.UNDERWORLD) && m.PlayerLocation.Location.GetMapType() == references.LargeMapType
}
