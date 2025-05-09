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

	n.positionOccupiedChance = n.mapUnits.createFreshXyOccupiedMap()

	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()

	// let's filter out all map units that are too far away, or empty before we even begin out
	// path computing
	n.mapUnits = helpers.FilterFromSlice(n.mapUnits,
		func(v MapUnit) bool {
			bRemove := v.IsEmptyMapUnit() || v.PosPtr().HeuristicTileDistance(n.gameState.Position) > maxTileDistanceBeforeCleanup
			return !bRemove
		})

	for _, npc := range n.mapUnits {
		// 	// very lazy approach - but making sure every NPC is in correct spot on map
		// 	// for every iteration makes sure next NPC doesn't assign the same tile space
		n.FreshenExistingNPCsOnMap()
		n.calculateNextNPCPosition(npc)
	}
	n.FreshenExistingNPCsOnMap()

	if len(n.mapUnits) < maxNPCS && n.ShouldGenerateLargeMapMonster() {
		n.generateEraBoundMonster()
	}
}

func (n *NPCAIControllerLargeMap) calculateNextNPCPosition(mapUnit MapUnit) {
	if mapUnit.PosPtr().IsNextTo(n.gameState.Position) {
		// if the NPC is next to the player, we don't want to move them
		return
	}

	if !n.ShouldEnemyMove() {
		return
	}

	n.setBestNextPositionToMoveTowardsWalkablePoint(mapUnit)

}

func (n *NPCAIControllerLargeMap) setBestNextPositionToMoveTowardsWalkablePoint(mapUnit MapUnit) {
	// this is an optimized a* pathfinding algorithm that limits the size of the map that it reads from
	mapUnit.MapUnitDetails().AStarMap.InitializeByLayeredMapWithLimit(
		mapUnit,
		n.gameState.GetLayeredMapByCurrentLocation(),
		[]references.Position{},
		true,
		n.gameState.Position,
		15,
		references.Coordinate(references.XLargeMapTiles),
		references.Coordinate(references.YLargeMapTiles))

	// make sure the correct enemies spawn on the correct tiles (water, sand, ground)

	path := mapUnit.MapUnitDetails().AStarMap.AStar(n.gameState.Position)
	if len(path) > 1 {
		// if the path is empty, we don't move
		mapUnit.SetPos(path[1])
	} else {
		// if we don't find a new position using AStar, then we at least try to get them closer to the avatar using
		// basic pathing
		n.setBestNextPositionToMoveTowardsWalkablePointDumb(mapUnit)
		return
	}

}

func (n *NPCAIControllerLargeMap) setBestNextPositionToMoveTowardsWalkablePointDumb(mapUnit MapUnit) {
	allDirections := mapUnit.PosPtr().GetFourDirectionsWrapped(references.XLargeMapTiles, references.YLargeMapTiles)
	// getting the current distance to the player will make sure they never move further away
	var fCurrentShortestDistance float64 = mapUnit.PosPtr().GetWrappedDistanceBetweenWrapped(&n.gameState.Position, references.XLargeMapTiles, references.YLargeMapTiles)
	var bestPos references.Position = *mapUnit.PosPtr()
	bFound := false
	for _, newPos := range allDirections {
		fNewDistance := newPos.GetWrappedDistanceBetweenWrapped(&n.gameState.Position, references.XLargeMapTiles, references.YLargeMapTiles)

		if fNewDistance < fCurrentShortestDistance {
			topTile := n.gameState.GetCurrentLayeredMap().GetTopTile(&newPos)
			if enemy, ok := mapUnit.(*NPCEnemy); ok {
				if !enemy.EnemyReference.CanMoveToTile(topTile) {
					continue
				}
			} else {
				if !topTile.IsLandEnemyPassable {
					continue
				}
			}

			bestPos = newPos
			fCurrentShortestDistance = fNewDistance
			bFound = true
		}
	}

	if !bFound {
		// if we don't find a new position, we don't try to move
		return
	}
	mapUnit.SetPos(bestPos)
}

func (n *NPCAIControllerLargeMap) FreshenExistingNPCsOnMap() {
	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
	n.placeNPCsOnLayeredMap()
}

func (n *NPCAIControllerLargeMap) generateEraBoundMonster() {
	const nYDistanceAway = 6
	const nXDistanceAway = 9
	const nTriesToGetValidEnemy = 10

	var dX, dY references.Coordinate

	if !n.gameState.DebugOptions.MonsterGen {
		return
	}

	for range nTriesToGetValidEnemy {
		if helpers.OneInXOdds(2) { // do dY
			dY = references.Coordinate(helpers.PickOneOf(nYDistanceAway, -nYDistanceAway))
			dX = references.Coordinate(helpers.RandomIntInRange(-nXDistanceAway, nXDistanceAway))

		} else { // do dX
			dY = references.Coordinate(helpers.RandomIntInRange(-nYDistanceAway, nYDistanceAway))
			dX = references.Coordinate(helpers.PickOneOf(nXDistanceAway, -nXDistanceAway))
		}

		pos := references.Position{X: n.gameState.Position.X + dX, Y: n.gameState.Position.Y + dY}
		pos = *pos.GetWrapped(references.XLargeMapTiles, references.YLargeMapTiles)
		if pos.X < 0 || pos.Y < 0 {
			log.Fatalf("Unexpected negative position X=%d Y=%d", pos.X, pos.Y)
		}

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

func (o *NPCAIControllerLargeMap) ShouldGenerateLargeMapMonster() bool {
	return helpers.OneInXOdds(o.gameState.TheOdds.GetOneInXLargeMapMonsterGeneration())
}

func (o *NPCAIControllerLargeMap) ShouldEnemyMove() bool {
	return helpers.HappenedByPercentLikely(o.gameState.TheOdds.GetPercentLikeyLargeMapMonsterMoves())
}

func (o *NPCAIControllerLargeMap) RemoveAllEnemies() {
	o.mapUnits = make(MapUnits, 0, maxNPCS)
}
