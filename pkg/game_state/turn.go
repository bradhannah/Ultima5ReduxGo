package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) FinishTurn() {

	switch g.Location.GetMapType() {
	case references.SmallMapType:
		g.smallMapProcessEndOfTurn()
	case references.LargeMapType:
		g.largeMapProcessEndOfTurn()
	default:
		panic("unhandled default case")
	}

	// process any residual damage such as lava or posion
	g.processDamageOnAdvanceTimeNonCombat()
	g.moveNonCombatMapMapUnitsToNextMove()
	g.GenerateAndCleanupEnemies()
}

func (g *GameState) largeMapProcessEndOfTurn() {
	topTile := g.GetCurrentLayeredMapAvatarTopTile()

	g.CurrentNPCAIController.AdvanceNextTurnCalcAndMoveNPCs()

	// we care about speed factor only for large maps
	g.DateTime.Advance(topTile.SpeedFactor)
}

func (g *GameState) smallMapProcessEndOfTurn() {
	g.smallMapProcessTurnDoors()
	// small maps are always 1 minute per turn
	g.DateTime.Advance(DefaultSmallMapMinutesPerTurn)

	g.smallMapProcessNPCs()
}

func (g *GameState) smallMapProcessNPCs() {
	g.CurrentNPCAIController.AdvanceNextTurnCalcAndMoveNPCs()
}

func (g *GameState) smallMapProcessTurnDoors() {
	if g.openDoorPos != nil {
		if g.openDoorTurns == 0 {
			tile := g.LayeredMaps.GetTileRefByPosition(references.SmallMapType, MapLayer, g.openDoorPos, g.Floor)
			var doorTileIndex indexes.SpriteIndex
			if tile.Index.IsWindowedDoor() {
				doorTileIndex = indexes.RegularDoorView
			} else {
				doorTileIndex = indexes.RegularDoor
			}
			g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).SetTileByLayer(MapOverrideLayer, g.openDoorPos, doorTileIndex)
			g.openDoorPos = nil
		} else {
			g.openDoorTurns--
		}
	}
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
	switch g.Location.GetMapType() {
	case references.SmallMapType:
		return
	case references.LargeMapType:
		//g.GetCurrentLayeredMap().
		//g.largeMapGenerateAndCleanupEnemies()
	case references.DungeonMapType:
	case references.CombatMapType:
	default:
		panic("unhandled default case")
	}
}

func (g *GameState) GetEra() references.Era {
	return references.GetEraByTurn(int(g.DateTime.Turn))
}
