package map_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

// var nChanceToGenerateEnemy = 32

const MaxTileDistanceBeforeCleanup = 22

func (g *MapState) UpdateLargeMap() {
	if g.PlayerLocation.Location.GetMapType() != references.LargeMapType {
		log.Fatalf("Expected large map type, got %d", g.PlayerLocation.Location.GetMapType())
	}
}

func (g *MapState) IsOverworld() bool {
	return g.PlayerLocation.Floor == references.FloorNumber(references.OVERWORLD) && g.PlayerLocation.Location.GetMapType() == references.LargeMapType
}

func (g *MapState) IsUnderworld() bool {
	return g.PlayerLocation.Floor == references.FloorNumber(references.UNDERWORLD) && g.PlayerLocation.Location.GetMapType() == references.LargeMapType
}

// func (g *MapState) GetCurrentLargeMapNPCAIController() *ai.NPCAIControllerLargeMap {
// 	if g.PlayerLocation.Location.GetMapType() != references.LargeMapType {
// 		log.Fatalf("Expected large map type, got %d", g.PlayerLocation.Location.GetMapType())
// 	}
// 	var npcAiCon *ai.NPCAIControllerLargeMap
// 	if g.IsOverworld() {
// 		npcAiCon = g.LargeMapNPCAIController[references.OVERWORLD]
// 	} else {
// 		npcAiCon = g.LargeMapNPCAIController[references.UNDERWORLD]
// 	}
// 	return npcAiCon
// }
