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

	if len(m.mapUnits) < map_units.MaximumNpcsPerMap && m.shouldGenerateTileBasedMonster() {
		m.generateTileBasedMonster()
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

func (m *NPCAIControllerLargeMap) generateTileBasedMonster() {
	const nYDistanceAway = 6
	const nXDistanceAway = 9
	const nTriesToGetValidEnemy = 10

	var dX, dY references2.Coordinate

	if !m.debugOptions.MonsterGen {
		return
	}

	// Probability checking now handled in shouldGenerateTileBasedMonster()
	// This method only executes when a monster should actually spawn

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

		// Use environment-based monster selection instead of era-based
		environment := m.determineEnvironmentType(tile)
		enemy := m.pickMonsterByEnvironment(environment, tile)
		if enemy == nil {
			continue
		}

		// Monster selection now handled by environment-based logic above

		// Simplified error handling for new system
		if enemy == nil {
			continue // Try another position
		}

		npc := map_units.NewEnemyNPC(*enemy, len(m.mapUnits))

		npc.SetPos(pos)
		npc.SetFloor(m.mapState.PlayerLocation.Floor)
		npc.SetVisible(true)
		m.mapUnits = append(m.mapUnits, &npc)
		return
	}
}

func (m *NPCAIControllerLargeMap) shouldGenerateTileBasedMonster() bool {
	// Calculate tile-based probability first
	probability := m.calculateTileBasedProbability()
	if probability == 0 {
		return false // Never spawn on roads
	}

	// Apply night bonus
	if m.dateTime.IsNight() {
		probability += 3
	}

	// Use the combined probability with base odds
	// This gives: Roads=never, Grass=1in32, Swamp=2in32, SwampNight=5in32
	combinedOdds := m.theOdds.GetOneInXLargeMapMonsterGeneration() / probability
	return helpers.OneInXOdds(combinedOdds)
}

func (m *NPCAIControllerLargeMap) ShouldEnemyMove() bool {
	return helpers.HappenedByPercentLikely(m.theOdds.GetPercentLikeyLargeMapMonsterMoves())
}

func (m *NPCAIControllerLargeMap) RemoveAllEnemies() {
	m.mapUnits = make(map_units.MapUnits, 0, map_units.MaximumNpcsPerMap)
}

// calculateTileBasedProbability implements the original MONSTER.C genprob() logic
func (m *NPCAIControllerLargeMap) calculateTileBasedProbability() int {
	playerTile := m.mapState.GetLayeredMapByCurrentLocation().GetTopTile(&m.mapState.PlayerLocation.Position)

	// Roads = 0 probability (no monsters spawn)
	if playerTile.IsRoad() {
		return 0
	}

	// Swamp, forest, mountains = 2 probability
	if playerTile.IsSwamp() || playerTile.IsForest() || playerTile.IsMountain() {
		return 2
	}

	// All other tiles = 1 probability
	return 1
}

// MonsterEnvironment represents the environment type for monster selection
type MonsterEnvironment int

const (
	WaterEnvironment MonsterEnvironment = iota
	DesertEnvironment
	LandEnvironment
	UnderworldEnvironment
)

// determineEnvironmentType categorizes a tile for monster selection
func (m *NPCAIControllerLargeMap) determineEnvironmentType(tile *references2.Tile) MonsterEnvironment {
	if tile.IsWater() {
		return WaterEnvironment
	}
	if tile.IsDesert() {
		return DesertEnvironment
	}
	// TODO: Add underworld detection based on floor/map type if needed

	return LandEnvironment
}

// pickMonsterByEnvironment selects a monster based on environment and tile compatibility using weighted selection
func (m *NPCAIControllerLargeMap) pickMonsterByEnvironment(environment MonsterEnvironment, tile *references2.Tile) *references2.EnemyReference {
	// Get monsters valid for this environment and tile
	validMonsters := make([]*references2.EnemyReference, 0)
	weights := make([]int, 0)

	for _, enemy := range *m.enemyReferences {
		if enemy.CanSpawnToTile(tile) {
			weight := m.getMonsterEnvironmentWeight(&enemy, environment)
			if weight > 0 {
				validMonsters = append(validMonsters, &enemy)
				weights = append(weights, weight)
			}
		}
	}

	if len(validMonsters) == 0 {
		return nil
	}

	// Weighted random selection based on environment
	return m.weightedRandomSelection(validMonsters, weights)
}

// getMonsterEnvironmentWeight returns the weight for a monster in a specific environment
func (m *NPCAIControllerLargeMap) getMonsterEnvironmentWeight(enemy *references2.EnemyReference, environment MonsterEnvironment) int {
	switch environment {
	case WaterEnvironment:
		return enemy.AdditionalEnemyFlags.WaterWeight
	case DesertEnvironment:
		return enemy.AdditionalEnemyFlags.DesertWeight
	case LandEnvironment:
		return enemy.AdditionalEnemyFlags.LandWeight
	case UnderworldEnvironment:
		return enemy.AdditionalEnemyFlags.UnderworldWeight
	default:
		return 1 // Default weight for unknown environments
	}
}

// weightedRandomSelection selects a monster based on weights
func (m *NPCAIControllerLargeMap) weightedRandomSelection(monsters []*references2.EnemyReference, weights []int) *references2.EnemyReference {
	// Calculate total weight
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	if totalWeight == 0 {
		return nil
	}

	// Generate random number in range [0, totalWeight)
	randomValue := helpers.RandomIntInRange(0, totalWeight-1)

	// Find the selected monster
	currentWeight := 0
	for i, weight := range weights {
		currentWeight += weight
		if randomValue < currentWeight {
			return monsters[i]
		}
	}

	// Fallback (should not happen)
	return monsters[len(monsters)-1]
}
