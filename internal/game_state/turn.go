package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) FinishTurn() {
	switch g.MapState.PlayerLocation.Location.GetMapType() {
	case references.SmallMapType:
		g.smallMapProcessEndOfTurn()
	case references.LargeMapType:
		g.largeMapProcessEndOfTurn()
	default:
		panic("unhandled default case")
	}

	// process any residual damage such as lava or position
	g.processDamageOnAdvanceTimeNonCombat()
	g.moveNonCombatMapMapUnitsToNextMove()
	g.GenerateAndCleanupEnemies()

	g.MapState.Lighting.AdvanceTurn()
}

func (g *GameState) largeMapProcessEndOfTurn() {
	topTile := g.GetCurrentLayeredMapAvatarTopTile()

	// g.LargeMapNPCAIController.AdvanceNextTurnCalcAndMoveNPCs()

	g.GetCurrentLargeMapNPCAIController().AdvanceNextTurnCalcAndMoveNPCs()

	// we care about the speed factor only for large maps
	g.DateTime.Advance(topTile.SpeedFactor)
}

func (g *GameState) smallMapProcessEndOfTurn() {
	g.MapState.SmallMapProcessTurnDoors()
	// small maps are always 1 minute per turn
	g.DateTime.Advance(DefaultSmallMapMinutesPerTurn)

	g.smallMapProcessNPCs()
}

func (g *GameState) smallMapProcessNPCs() {
	g.CurrentNPCAIController.AdvanceNextTurnCalcAndMoveNPCs()
}

// processDamageOnAdvanceTimeNonCombat
// Processes damage from lava, poison, etc. on the non-combat map
func (g *GameState) processDamageOnAdvanceTimeNonCombat() {
	// TODO: implement this
}

// moveNonCombatMapMapUnitsToNextMove
// Moves non-combat map units to their next move
func (g *GameState) moveNonCombatMapMapUnitsToNextMove() {
	// TODO: implement this
}

// GenerateAndCleanupEnemies
// Generates and cleans up enemies on the map
func (g *GameState) GenerateAndCleanupEnemies() {
	// TODO: implement this
	switch g.MapState.PlayerLocation.Location.GetMapType() {
	case references.SmallMapType:
		return
	case references.LargeMapType:
		// g.GetCurrentLayeredMap().
		// g.largeMapGenerateAndCleanupEnemies()
	case references.DungeonMapType:
	case references.CombatMapType:
	default:
		panic("unhandled default case")
	}
}

// func (g *GameState) GetEra() datetime.Era {
// 	return GetEraByTurn(int(g.DateTime.Turn))
// }
