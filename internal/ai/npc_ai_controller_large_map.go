package ai

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/astar"
	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

type NPCAIControllerLargeMap struct {
	World           references2.World
	tileRefs        *references2.Tiles
	mapState        *map_state.MapState
	debugOptions    *references2.DebugOptions
	enemyReferences *references2.EnemyReferences
	theOdds         *references2.TheOdds
	dateTime        *datetime.UltimaDate

	mapUnits               map_units.MapUnits
	positionOccupiedChance *map_units.XyOccupiedMap
}

type NewNPCAIControllerLargeMapInput struct {
	World           references2.World
	TileRefs        *references2.Tiles
	MapState        *map_state.MapState
	DebugOptions    *references2.DebugOptions
	TheOdds         *references2.TheOdds
	EnemyReferences *references2.EnemyReferences
	DateTime        *datetime.UltimaDate
}

func NewNPCAIControllerLargeMap(input NewNPCAIControllerLargeMapInput) *NPCAIControllerLargeMap {
	npcsAiCont := &NPCAIControllerLargeMap{}
	npcsAiCont.mapState = input.MapState
	npcsAiCont.debugOptions = input.DebugOptions
	npcsAiCont.tileRefs = input.TileRefs
	npcsAiCont.World = input.World
	npcsAiCont.enemyReferences = input.EnemyReferences
	npcsAiCont.theOdds = input.TheOdds
	npcsAiCont.dateTime = input.DateTime

	xy := make(map_units.XyOccupiedMap)
	npcsAiCont.positionOccupiedChance = &xy

	npcsAiCont.mapUnits = make(map_units.MapUnits, 0, map_units.MaximumNpcsPerMap)

	return npcsAiCont
}

func (m *NPCAIControllerLargeMap) GetNpcs() *map_units.MapUnits {
	return &m.mapUnits
}

func (m *NPCAIControllerLargeMap) PopulateMapFirstLoad() {
}

func (m *NPCAIControllerLargeMap) placeNPCsOnLayeredMap() {
	lm := m.mapState.GetLayeredMapByCurrentLocation()

	for _, mu := range m.mapUnits {
		if !mu.IsVisible() || m.mapState.PlayerLocation.Floor != mu.Floor() {
			continue
		}
		switch npc := mu.(type) {
		case *map_units.NPCEnemy:
			lm.SetTileByLayer(map_state.MapUnitLayer, npc.PosPtr(), npc.EnemyReference.KeyFrameTile.Index)
		case *map_units.NPCFriendly:
			lm.SetTileByLayer(map_state.MapUnitLayer, npc.PosPtr(), npc.NPCReference.GetSpriteIndex())
		}
	}
}

func (m *NPCAIControllerLargeMap) AdvanceNextTurnCalcAndMoveNPCs() {
	m.positionOccupiedChance = m.mapUnits.CreateFreshXyOccupiedMap()

	m.mapState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()

	// let's filter out all map units that are too far away, or empty before we even begin out
	// path computing
	m.mapUnits = helpers.FilterFromSlice(m.mapUnits,
		func(v map_units.MapUnit) bool {
			bRemove := v.IsEmptyMapUnit() || v.PosPtr().HeuristicTileDistance(m.mapState.PlayerLocation.Position) > MaxTileDistanceBeforeCleanup
			return !bRemove
		})

	for _, npc := range m.mapUnits {
		// 	// very lazy approach - but making sure every NPC is in correct spot on map
		// 	// for every iteration makes sure next NPC doesn't assign the same tile space
		m.FreshenExistingNPCsOnMap()
		m.calculateNextNPCPosition(npc)
	}
	m.FreshenExistingNPCsOnMap()

	if len(m.mapUnits) < map_units.MaximumNpcsPerMap && m.ShouldGenerateLargeMapMonster() {
		m.generateEraBoundMonster()
	}
}

func (m *NPCAIControllerLargeMap) calculateNextNPCPosition(mapUnit map_units.MapUnit) {
	if _, ok := mapUnit.(*map_units.NPCEnemy); !ok {
		// Friendly units do not currently move in the large maps (ie. Frigates)
		return
	}

	if mapUnit.PosPtr().IsNextTo(m.mapState.PlayerLocation.Position) {
		// if the NPC is next to the player, we don't want to move them
		return
	}

	if !m.ShouldEnemyMove() {
		return
	}

	m.setBestNextPositionToMoveTowardsWalkablePoint(mapUnit)
}

func (m *NPCAIControllerLargeMap) setBestNextPositionToMoveTowardsWalkablePoint(mapUnit map_units.MapUnit) {
	// this is an optimized a* pathfinding algorithm that limits the size of the map that it reads from
	aStarMap := astar.NewAStarMap()
	aStarMap.InitializeByLayeredMapWithLimit(
		// mapUnit.MapUnitDetails().AStarMap.InitializeByLayeredMapWithLimit(
		mapUnit,
		m.mapState.GetLayeredMapByCurrentLocation(),
		[]references2.Position{},
		true,
		m.mapState.PlayerLocation.Position,
		15,
		references2.XLargeMapTiles,
		references2.YLargeMapTiles)

	// make sure the correct enemies spawn on the correct tiles (water, sand, ground)

	path := aStarMap.AStar(m.mapState.PlayerLocation.Position)
	// path := mapUnit.MapUnitDetails().AStarMap.AStar(m.mapState.PlayerLocation.Position)
	if len(path) > 1 {
		// if the path is empty, we don't move
		mapUnit.SetPos(path[1])
	} else {
		// if we don't find a new position using AStar, then we at least try to get them closer to the avatar using
		// basic pathing
		m.setBestNextPositionToMoveTowardsWalkablePointDumb(mapUnit)
		return
	}
}

func (m *NPCAIControllerLargeMap) setBestNextPositionToMoveTowardsWalkablePointDumb(mapUnit map_units.MapUnit) {
	allDirections := mapUnit.PosPtr().GetFourDirectionsWrapped(references2.XLargeMapTiles, references2.YLargeMapTiles)
	// getting the current distance to the player will make sure they never move further away
	fCurrentShortestDistance := mapUnit.PosPtr().GetWrappedDistanceBetweenWrapped(&m.mapState.PlayerLocation.Position, references2.XLargeMapTiles, references2.YLargeMapTiles)
	bestPos := *mapUnit.PosPtr()
	bFound := false
	for _, newPos := range allDirections {
		fNewDistance := newPos.GetWrappedDistanceBetweenWrapped(&m.mapState.PlayerLocation.Position, references2.XLargeMapTiles, references2.YLargeMapTiles)

		if fNewDistance < fCurrentShortestDistance {
			topTile := m.mapState.GetLayeredMapByCurrentLocation().GetTopTile(&newPos)
			if enemy, ok := mapUnit.(*map_units.NPCEnemy); ok {
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

func (m *NPCAIControllerLargeMap) FreshenExistingNPCsOnMap() {
	m.mapState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
	m.placeNPCsOnLayeredMap()
}

func (m *NPCAIControllerLargeMap) generateEraBoundMonster() {
	const nYDistanceAway = 6
	const nXDistanceAway = 9
	const nTriesToGetValidEnemy = 10

	var dX, dY references2.Coordinate

	if !m.debugOptions.MonsterGen {
		return
	}

	for range nTriesToGetValidEnemy {
		if helpers.OneInXOdds(2) { // do dY
			dY = references2.Coordinate(helpers.PickOneOf(nYDistanceAway, -nYDistanceAway))
			dX = references2.Coordinate(helpers.RandomIntInRange(-nXDistanceAway, nXDistanceAway))

		} else { // do dX
			dY = references2.Coordinate(helpers.RandomIntInRange(-nYDistanceAway, nYDistanceAway))
			dX = references2.Coordinate(helpers.PickOneOf(nXDistanceAway, -nXDistanceAway))
		}

		pos := references2.Position{X: m.mapState.PlayerLocation.Position.X + dX, Y: m.mapState.PlayerLocation.Position.Y + dY}
		pos = *pos.GetWrapped(references2.XLargeMapTiles, references2.YLargeMapTiles)
		if pos.X < 0 || pos.Y < 0 {
			log.Fatalf("Unexpected negative position X=%d Y=%d", pos.X, pos.Y)
		}

		tile := m.mapState.GetLayeredMapByCurrentLocation().GetTopTile(&pos)
		enemy, err := m.enemyReferences.GetRandomEnemyReferenceByEraAndTile(m.dateTime.GetEra(), tile)

		if err != nil {
			log.Printf("Error getting random enemy reference: %v", err)
			continue
			// return
		} else if enemy == nil {
			log.Fatal("Unexpected nil")
		}

		npc := map_units.NewEnemyNPC(*enemy, len(m.mapUnits))

		npc.SetPos(pos)
		npc.SetFloor(m.mapState.PlayerLocation.Floor)
		npc.SetVisible(true)
		m.mapUnits = append(m.mapUnits, &npc)
		return
	}
}

func (m *NPCAIControllerLargeMap) ShouldGenerateLargeMapMonster() bool {
	return helpers.OneInXOdds(m.theOdds.GetOneInXLargeMapMonsterGeneration())
}

func (m *NPCAIControllerLargeMap) ShouldEnemyMove() bool {
	return helpers.HappenedByPercentLikely(m.theOdds.GetPercentLikeyLargeMapMonsterMoves())
}

func (m *NPCAIControllerLargeMap) RemoveAllEnemies() {
	m.mapUnits = make(map_units.MapUnits, 0, map_units.MaximumNpcsPerMap)
}
