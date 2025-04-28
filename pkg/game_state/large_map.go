package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

var nChanceToGenerateEnemy = 10

func (g *GameState) UpdateLargeMap() {
	if g.Location.GetMapType() != references.LargeMapType {
		log.Fatalf("Expected large map type, got %d", g.Location.GetMapType())
	}
}

func (g *GameState) largeMapGenerateAndCleanupEnemies() {
	if g.Location.GetMapType() != references.LargeMapType {
		log.Fatalf("Expected large map type, got %d", g.Location.GetMapType())
	}
	// get all the enemies

	// ClearEnemiesIfFarAway()
	_ = g.LargeMapNPCAIController[references.OVERWORLD]
	g.Location.GetMapType()
	if helpers.RandomIntInRange(0, 100) < nChanceToGenerateEnemy {
		g.generateNewEnemy()
	}
}

func (g *GameState) generateNewEnemy() {

}

func (g *GameState) IsOverworld() bool {
	return g.Floor == references.FloorNumber(references.OVERWORLD) && g.Location.GetMapType() == references.LargeMapType
}

func (g *GameState) IsUnderworld() bool {
	return g.Floor == references.FloorNumber(references.UNDERWORLD) && g.Location.GetMapType() == references.LargeMapType
}

func (g *GameState) GetCurrentLargeMapNPCAIController() *NPCAIControllerLargeMap {
	if g.Location.GetMapType() != references.LargeMapType {
		log.Fatalf("Expected large map type, got %d", g.Location.GetMapType())
	}
	var npcAiCon *NPCAIControllerLargeMap
	if g.IsOverworld() {
		npcAiCon = g.LargeMapNPCAIController[references.OVERWORLD]
	} else {
		npcAiCon = g.LargeMapNPCAIController[references.UNDERWORLD]
	}
	return npcAiCon
}
