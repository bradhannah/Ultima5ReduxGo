package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCAIControllerLargeMap struct {
	tileRefs  *references.Tiles
	World     references.World
	gameState *GameState

	mapUnits MapUnits

	positionOccupiedChance *XyOccupiedMap
}

func NewNPCAIControllerLargeMap(
	world references.World,
	tileRefs *references.Tiles,
	gameState *GameState,
) *NPCAIControllerLargeMap {
	npcsAiCont := &NPCAIControllerLargeMap{}

	npcsAiCont.tileRefs = tileRefs
	npcsAiCont.World = world
	npcsAiCont.gameState = gameState

	xy := make(XyOccupiedMap)
	npcsAiCont.positionOccupiedChance = &xy

	npcsAiCont.mapUnits = make(MapUnits, 0, maxNPCS)

	return npcsAiCont
}

func (n *NPCAIControllerLargeMap) GetNpcs() *MapUnits {
	return &n.mapUnits
}

func (n *NPCAIControllerLargeMap) PopulateMapFirstLoad() {
}

func (n *NPCAIControllerLargeMap) placeNPCsOnLayeredMap() {
	lm := n.gameState.GetLayeredMapByCurrentLocation()

	for _, npc := range n.mapUnits {
		enemy := getMapUnitAsEnemyOrNil(&npc)
		if enemy == nil || !enemy.IsVisible() {
			continue
		}
		if n.gameState.Floor == npc.Floor() {
			//_ = lm
			lm.SetTileByLayer(MapUnitLayer, npc.PosPtr(), enemy.EnemyReference.KeyFrameTile.Index)
		}
	}
}

func (n *NPCAIControllerLargeMap) AdvanceNextTurnCalcAndMoveNPCs() {
	//n.clearMapUnitsFromMap()
	n.generateEraBoundMonster()

	n.positionOccupiedChance = n.mapUnits.createFreshXyOccupiedMap()

	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()

	for _, npc := range n.mapUnits {
		if npc.IsEmptyMapUnit() {
			continue
		}
		// 	// very lazy approach - but making sure every NPC is in correct spot on map
		// 	// for every iteration makes sure next NPC doesn't assign the same tile space
		n.FreshenExistingNPCsOnMap()
		// 	n.calculateNextNPCPosition(npc)
	}
	n.FreshenExistingNPCsOnMap()

	// should we spawn units after these ones have moved? probably
}

func (n *NPCAIControllerLargeMap) FreshenExistingNPCsOnMap() {
	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
	n.placeNPCsOnLayeredMap()
}

func (n *NPCAIControllerLargeMap) generateEraBoundMonster() {
	const nYDistanceAway = 6
	const nXDistanceAway = 9

	var dX, dY references.Coordinate

	for i := 0; i < 10; i = i + 1 {
		if helpers.OneInXOdds(2) { // do dY
			dY = references.Coordinate(helpers.PickOneOf(nYDistanceAway, -nYDistanceAway))
			dX = references.Coordinate(helpers.RandomIntInRange(-nXDistanceAway, nXDistanceAway))

		} else { // do dX
			dY = references.Coordinate(helpers.RandomIntInRange(-nYDistanceAway, nYDistanceAway))
			dX = references.Coordinate(helpers.PickOneOf(nXDistanceAway, -nXDistanceAway))
		}

		pos := references.Position{X: n.gameState.Position.X + dX, Y: n.gameState.Position.Y + dY}

		tile := n.gameState.GetLayeredMapByCurrentLocation().GetTopTile(&pos)
		enemy, err := n.gameState.GameReferences.EnemyReferences.GetRandomEnemyReferenceByEraAndTile(n.gameState.GetEra(), tile)

		if err != nil {
			log.Printf("Error getting random enemy reference: %v", err)
			continue
			//return
		} else if enemy == nil {
			log.Fatal("Unexpected nil")
		}

		npc := NewEnemyNPC(*enemy, len(n.mapUnits))

		npc.SetPos(pos)
		npc.SetFloor(n.gameState.Floor)
		npc.SetVisible(true)
		n.mapUnits = append(n.mapUnits, &npc)
		return
	}
}
