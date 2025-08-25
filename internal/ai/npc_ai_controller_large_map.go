package ai

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/astar"
	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

type NPCAIControllerLargeMap struct {
	World           references.World
	tileRefs        *references.Tiles
	mapState        *map_state.MapState
	debugOptions    *references.DebugOptions
	enemyReferences *references.EnemyReferences
	theOdds         *references.TheOdds
	dateTime        *datetime.UltimaDate

	mapUnits               map_units.MapUnits
	positionOccupiedChance *map_units.XyOccupiedMap
}

type NewNPCAIControllerLargeMapInput struct {
	World           references.World
	TileRefs        *references.Tiles
	MapState        *map_state.MapState
	DebugOptions    *references.DebugOptions
	TheOdds         *references.TheOdds
	EnemyReferences *references.EnemyReferences
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

	// Respect per-enemy immobility flags (e.g., Reapers, Mimics, fields)
	if enemy, ok := mapUnit.(*map_units.NPCEnemy); ok {
		if enemy.EnemyReference.AdditionalEnemyFlags.DoNotMove {
			return
		}
	}

	if mapUnit.PosPtr().IsNextTo(m.mapState.PlayerLocation.Position) {
		// if the NPC is next to the player, we don't want to move them
		return
	}

	if !m.ShouldEnemyMoveThisTick(mapUnit) {
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
		[]references.Position{},
		true,
		m.mapState.PlayerLocation.Position,
		15,
		references.XLargeMapTiles,
		references.YLargeMapTiles)

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
	allDirections := mapUnit.PosPtr().GetFourDirectionsWrapped(references.XLargeMapTiles, references.YLargeMapTiles)
	// getting the current distance to the player will make sure they never move further away
	fCurrentShortestDistance := mapUnit.PosPtr().GetWrappedDistanceBetweenWrapped(&m.mapState.PlayerLocation.Position, references.XLargeMapTiles, references.YLargeMapTiles)
	bestPos := *mapUnit.PosPtr()
	bFound := false
	for _, newPos := range allDirections {
		fNewDistance := newPos.GetWrappedDistanceBetweenWrapped(&m.mapState.PlayerLocation.Position, references.XLargeMapTiles, references.YLargeMapTiles)

		if fNewDistance < fCurrentShortestDistance {
			topTile := m.mapState.GetLayeredMapByCurrentLocation().GetTopTile(&newPos)
			if enemy, ok := mapUnit.(*map_units.NPCEnemy); ok {
				if !enemy.EnemyReference.CanMoveToTile(topTile) {
					continue
				}
			} else {
				if !topTile.IsLandEnemyPassable() {
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

	var dX, dY references.Coordinate

	if !m.debugOptions.MonsterGen {
		return
	}

	// Probability checking now handled in shouldGenerateTileBasedMonster()
	// This method only executes when a monster should actually spawn

	for range nTriesToGetValidEnemy {
		if helpers.OneInXOdds(2) { // do dY
			dY = references.Coordinate(helpers.PickOneOf(nYDistanceAway, -nYDistanceAway))
			dX = references.Coordinate(helpers.RandomIntInRange(-nXDistanceAway, nXDistanceAway))

		} else { // do dX
			dY = references.Coordinate(helpers.RandomIntInRange(-nYDistanceAway, nYDistanceAway))
			dX = references.Coordinate(helpers.PickOneOf(nXDistanceAway, -nXDistanceAway))
		}

		pos := references.Position{X: m.mapState.PlayerLocation.Position.X + dX, Y: m.mapState.PlayerLocation.Position.Y + dY}
		pos = *pos.GetWrapped(references.XLargeMapTiles, references.YLargeMapTiles)
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
	combinedOdds := m.theOdds.GetOneInXMonsterGeneration() / probability
	return helpers.OneInXOdds(combinedOdds)
}

// ShouldEnemyMoveThisTick implements terrain-based movement throttling from Movement_Overworld.md pseudocode
func (m *NPCAIControllerLargeMap) ShouldEnemyMoveThisTick(mapUnit map_units.MapUnit) bool {
	tile := m.mapState.GetLayeredMapByCurrentLocation().GetTopTile(mapUnit.PosPtr())

	// Implement should_move_this_tick(tile) logic from pseudocode
	if m.isHeavyTerrain(tile) {
		// Heavy terrain (mountains, deep forest): 1-in-3 chance
		return helpers.RandomIntInRange(0, 2) == 2
	}

	if m.isDifficultTerrain(tile) {
		// Difficult terrain (swamps, light forest): 1-in-2 chance
		return helpers.RandomIntInRange(0, 1) == 0
	}

	// Open terrain (grass, roads): Always move
	return true
}

func (m *NPCAIControllerLargeMap) RemoveAllEnemies() {
	m.mapUnits = make(map_units.MapUnits, 0, map_units.MaximumNpcsPerMap)
}

// isHeavyTerrain classifies terrain that slows movement to 1-in-3 chance
func (m *NPCAIControllerLargeMap) isHeavyTerrain(tile *references.Tile) bool {
	// Heavy terrain: mountains only (game only has one forest type)
	return tile.IsMountain()
}

// isDifficultTerrain classifies terrain that slows movement to 1-in-2 chance
func (m *NPCAIControllerLargeMap) isDifficultTerrain(tile *references.Tile) bool {
	// Difficult terrain: swamps, forests (game only has one forest type)
	return tile.IsSwamp() || tile.IsForest()
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
func (m *NPCAIControllerLargeMap) determineEnvironmentType(tile *references.Tile) MonsterEnvironment {
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
func (m *NPCAIControllerLargeMap) pickMonsterByEnvironment(environment MonsterEnvironment, tile *references.Tile) *references.EnemyReference {
	// Get monsters valid for this environment and tile
	validMonsters := make([]*references.EnemyReference, 0)
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
func (m *NPCAIControllerLargeMap) getMonsterEnvironmentWeight(enemy *references.EnemyReference, environment MonsterEnvironment) int {
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
func (m *NPCAIControllerLargeMap) weightedRandomSelection(monsters []*references.EnemyReference, weights []int) *references.EnemyReference {
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
