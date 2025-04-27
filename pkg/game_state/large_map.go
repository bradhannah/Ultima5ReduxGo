package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

var nChanceToGenerateEnemy = 10

func (g *GameState) largeMapGenerateAndCleanupEnemies() {
	if g.Location.GetMapType() != references.LargeMapType {
		log.Fatalf("Expected large map type, got %s", g.Location.GetMapType())
	}
	// get all the enemies

	if helpers.RandomIntInRange(0, 100) < nChanceToGenerateEnemy {
		g.generateNewEnemy()
	}
}

func (g *GameState) generateNewEnemy() {

}
