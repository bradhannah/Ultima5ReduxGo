package ai

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/exp/rand"

	"github.com/bradhannah/Ultima5ReduxGo/internal/astar"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCAIControllerSmallMap struct {
	tileRefs *references.Tiles
	slr      *references.SmallLocationReference
	dateTime *datetime.UltimaDate
	// gameState *game_state.GameState
	mapState *map_state.MapState
	mapUnits map_units.MapUnits

	positionOccupiedChance *map_units.XyOccupiedMap
}

func NewNPCAIControllerSmallMap(
	slr *references.SmallLocationReference,
	tileRefs *references.Tiles,
	mapState *map_state.MapState,
	dateTime *datetime.UltimaDate,
	// gameState *game_state.GameState,
) *NPCAIControllerSmallMap {
	npcsAiCont := &NPCAIControllerSmallMap{}
	npcsAiCont.dateTime = dateTime
	npcsAiCont.tileRefs = tileRefs
	npcsAiCont.slr = slr
	npcsAiCont.mapState = mapState
	// npcsAiCont.gameState = gameState

	xy := make(map_units.XyOccupiedMap)
	npcsAiCont.positionOccupiedChance = &xy

	npcsAiCont.mapUnits = make(map_units.MapUnits, 0, map_units.MAXIMUM_NPCS_PER_MAP)

	return npcsAiCont
}

func (n *NPCAIControllerSmallMap) GetNpcs() *map_units.MapUnits {
	return &n.mapUnits
}

func (n *NPCAIControllerSmallMap) PopulateMapFirstLoad() {
	n.generateNPCs()

	for i, npc := range n.mapUnits {
		_ = i
		if npc.IsEmptyMapUnit() || !npc.IsVisible() {
			continue
		}

		switch mapUnit := npc.(type) {
		case *map_units.NPCFriendly:
			if n.mapState.PlayerLocation.Floor == npc.Floor() {
				indiv := mapUnit.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(*n.dateTime)
				mapUnit.SetIndividualNPCBehaviour(indiv)
			}
		case *map_units.NPCEnemy:
			// do not support NPC Enemy on small map
		}

	}
	n.placeNPCsOnLayeredMap()
}

func (n *NPCAIControllerSmallMap) AdvanceNextTurnCalcAndMoveNPCs() {
	n.mapState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
	n.updateAllNPCAiTypes()
	n.positionOccupiedChance = n.mapUnits.CreateFreshXyOccupiedMap()

	for _, mu := range n.mapUnits {
		// very lazy approach - but making sure every NPC is in correct spot on map
		// for every iteration makes sure next NPC doesn't assign the same tile space
		n.FreshenExistingNPCsOnMap()
		switch npc := mu.(type) {
		case *map_units.NPCFriendly:
			n.calculateNextNPCPosition(npc)
		}
	}
	n.FreshenExistingNPCsOnMap()
}

func (n *NPCAIControllerSmallMap) FreshenExistingNPCsOnMap() {
	n.mapState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
	n.placeNPCsOnLayeredMap()
}

func (n *NPCAIControllerSmallMap) generateNPCs() {
	npcs := make([]map_units.MapUnit, 0)
	// get the correct schedule
	npcsRefs := n.slr.GetNPCReferences()
	for nNpc, npcRef := range *npcsRefs {
		if npcRef.IsEmptyNPC() {
			continue
		}

		npcType := npcRef.GetNPCType()

		_ = npcType
		if npcRef.GetNPCType() == references.Vehicle {
			vehicle := map_units.NewNPCFriendlyVehicle(
				npcRef.GetVehicleType(), npcRef)
			npcs = append(npcs, vehicle)
		} else {
			friendly := map_units.NewNPCFriendly(npcRef, nNpc)
			if !friendly.IsEmptyMapUnit() {
				npcs = append(npcs, friendly)
			}
		}
	}
	n.mapUnits = npcs
}

func (n *NPCAIControllerSmallMap) updateAllNPCAiTypes() {
	for _, mu := range n.mapUnits {
		var indiv references.IndividualNPCBehaviour
		switch npc := mu.(type) {
		case *map_units.NPCFriendly:
			indiv = npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(*n.dateTime)

		}

		mu.MapUnitDetails().AiType = indiv.Ai
	}
}

func (n *NPCAIControllerSmallMap) placeNPCsOnLayeredMap() {
	lm := n.mapState.GetLayeredMapByCurrentLocation()

	for _, mu := range n.mapUnits {
		switch npc := mu.(type) {
		case *map_units.NPCFriendly:
			if !npc.IsVisible() {
				continue
			}
			if n.mapState.PlayerLocation.Floor == mu.Floor() {
				lm.SetTileByLayer(map_state.MapUnitLayer, mu.PosPtr(), npc.NPCReference.GetSpriteIndex())
			}
		}
	}
}

func (n *NPCAIControllerSmallMap) calculateNextNPCPosition(friendly *map_units.NPCFriendly) {
	refBehaviour := friendly.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(*n.dateTime)

	if friendly.GetVehicleDetails().VehicleType == references.CarpetVehicle {
		fmt.Sprint("oof")
	}

	// TEST: let's always finish what they are doing first before considering the next logic
	if n.moveNPCOnCalculatedPath(friendly) {
		return
	}

	if friendly.PosPtr().Equals(&refBehaviour.Position) && friendly.Floor() == refBehaviour.Floor {
		if n.performAiMovementOnAssignedPosition(friendly) {
			return
		}
	} else if friendly.Floor() != refBehaviour.Floor { // the NPC is on the wrong floor according to their schedule
		if friendly.Floor() == n.mapState.PlayerLocation.Floor { // the NPC is on the Avatar's current floor
			n.performAiMovementFromCurrentFloorToDifferentFloor(friendly)
			return
		}
		// the NPC is on another floor and needs to come to ours
		n.performAiMovementFromDifferentFloorToOurFloor(friendly)
	} else {
		if n.performAiMovementNotOnAssignedPosition(friendly) {
			return
		}
	}
}

// performAiMovementFromCurrentFloorToDifferentFloor From DIFFERENT floor to OUR floor
func (n *NPCAIControllerSmallMap) performAiMovementFromDifferentFloorToOurFloor(friendly *map_units.NPCFriendly) bool {
	// called if the NPC is currently on a different floor then the current floor
	refBehaviour := friendly.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(*n.dateTime)

	// current floor matters - if they are coming to your floor - then teleport them
	closestLadderPos := n.slr.GetClosestLadder(refBehaviour.Position, friendly.Floor(), n.mapState.PlayerLocation.Floor)

	// check if something or someone else is on the ladder, if so then we skip it for this turn
	// and try again next turn
	tile := n.mapState.GetLayeredMapByCurrentLocation().GetTopTile(&closestLadderPos)
	if !tile.IsWalkingPassable {
		return false
	}

	friendly.SetPos(closestLadderPos)
	friendly.SetFloor(refBehaviour.Floor)
	return true
}

// performAiMovementFromCurrentFloorToDifferentFloor From OUR floor to DIFFERENT floor
func (n *NPCAIControllerSmallMap) performAiMovementFromCurrentFloorToDifferentFloor(friendly *map_units.NPCFriendly) bool {
	// called if the NPC is currently on a different floor then the current floor
	refBehaviour := friendly.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(*n.dateTime)

	currentNpcMapTile := n.mapState.GetLayeredMapByCurrentLocation().GetTileTopMapOnlyTile(friendly.PosPtr())
	if references.IsSpecificLadderOrStairs(currentNpcMapTile.Index,
		references.GetLadderOfStairsType(friendly.Floor(), refBehaviour.Floor)) {
		// we have arrived at the ladder, so we will change their position as well
		// to make sure they "come down from" the correct spot as well
		friendly.SetFloor(refBehaviour.Floor)
		friendly.SetPos(refBehaviour.Position)
		return true
	}

	// // current floor matters - if they are coming to your floor - then teleport them
	closestLadderPos := n.slr.GetClosestLadder(refBehaviour.Position, friendly.Floor(), refBehaviour.Floor) // n.gameState.Floor)
	tile := n.mapState.GetLayeredMapByCurrentLocation().GetTopTile(&closestLadderPos)
	if !tile.IsWalkingPassable {
		return false
	}

	// the ladder is not used, so let's build a path
	if n.createFreshPathToScheduledLocation(friendly) {
		return true
	}

	return false
}

func (n *NPCAIControllerSmallMap) performAiMovementOnAssignedPosition(friendly *map_units.NPCFriendly) bool {
	npcSched := friendly.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(*n.dateTime)
	muDetails := friendly.MapUnitDetails()
	nWanderDistance := n.getWanderDistanceByAiType(muDetails.AiType)

	switch friendly.MapUnitDetails().AiType {
	case references.BlackthornGuardFixed, references.Fixed:
	case references.MerchantBuyingSellingCustom, references.MerchantBuyingSellingWander, references.Wander:
		n.wanderOneTileWithinN(muDetails, npcSched.Position, nWanderDistance)
		return true
	case references.BigWander, references.BlackthornGuardWander:
		n.wanderOneTileWithinN(muDetails, npcSched.Position, nWanderDistance)
		return true
	case references.ChildRunAway:
		return true
	case references.CustomAi, references.MerchantBuyingSelling:
		// don't think they move....?
		return true
	case references.DrudgeWorthThing:
		// try to approach avatar
		return true
	case references.ExtortOrAttackOrFollow:
		// set location of Avatar as way point, but only set the first movement from the list if within N of Avatar
		return true
	case references.HorseWander:
		if helpers.OneInXOdds(4) {
			return n.wanderOneTileWithinN(muDetails, npcSched.Position, nWanderDistance)
		}
	case references.StoneGargoyleTrigger:
		// if they are within 4 then change their AI to Drudgeworth (follow)
	case references.FixedExceptAttackWhenIsWantedByThePoPo:
		// if avatar is a wanted man/woman - then follow and get close
	case references.Begging, references.GenericExtortingGuard, references.HalfYourGoldExtortingGuard, references.SmallWanderWantsToChat:
		// let's have them try to hang out with the avatar most of the time, but not everytime
		// for a little randomness
		return true
	case references.FollowAroundAndBeAnnoyingThenNeverSeeAgain:
		// let's have them try to hang out with the avatar most of the time, but not everytime
		// for a little randomness
		return true
	default:
		log.Fatal("Unknown AiType")
	}
	return false
}

func (n *NPCAIControllerSmallMap) performAiMovementNotOnAssignedPosition(friendly *map_units.NPCFriendly) bool {
	npcSched := friendly.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(*n.dateTime)
	muDetails := friendly.MapUnitDetails()
	nWanderDistance := n.getWanderDistanceByAiType(friendly.MapUnitDetails().AiType)

	if n.moveNPCOnCalculatedPath(friendly) {
		return true
	}

	switch friendly.MapUnitDetails().AiType {
	case references.BlackthornGuardFixed, references.Fixed, references.CustomAi, references.MerchantBuyingSelling:
		if n.createFreshPathToScheduledLocation(friendly) {
			friendly.SetPos(friendly.MapUnitDetails().DequeueNextPosition())
			return true
		}
		return false
	case references.BigWander, references.BlackthornGuardWander, references.MerchantBuyingSellingCustom, references.MerchantBuyingSellingWander, references.Wander, references.HorseWander:
		if helpers.OneInXOdds(2) {
			if !npcSched.Position.IsWithinN(friendly.PosPtr(), nWanderDistance) {
				if n.createFreshPathToScheduledLocation(friendly) {
					friendly.SetPos(friendly.MapUnitDetails().DequeueNextPosition())
					return true
				}
				return false
			}
			return n.wanderOneTileWithinN(muDetails, npcSched.Position, nWanderDistance)
		}
	case references.ChildRunAway:
		// run away
		return true
	case references.DrudgeWorthThing:
		// try to approach avatar
		return true
	case references.ExtortOrAttackOrFollow:
		// set location of Avatar as way point, but only set the first movement from the list if within N of Avatar
		return true
	case references.FixedExceptAttackWhenIsWantedByThePoPo:
		// if avatar is a wanted man/woman - then follow and get close
		return true
	case references.StoneGargoyleTrigger:
		return true
	case references.FollowAroundAndBeAnnoyingThenNeverSeeAgain:
		return true
	case references.Begging,
		references.GenericExtortingGuard,
		references.HalfYourGoldExtortingGuard,
		references.SmallWanderWantsToChat:
		if !npcSched.Position.IsWithinN(friendly.PosPtr(), nWanderDistance) {
			if n.createFreshPathToScheduledLocation(friendly) {
				friendly.SetPos(friendly.MapUnitDetails().DequeueNextPosition())
				return true
			}
			return false
		}
		if helpers.OneInXOdds(3) {
			return n.wanderOneTileWithinN(muDetails, npcSched.Position, nWanderDistance)
		}
		return false
	default:
		log.Fatal("Unknown AiType")
	}
	return false
}

func (n *NPCAIControllerSmallMap) moveNPCOnCalculatedPath(friendly *map_units.NPCFriendly) bool {
	if !friendly.MapUnitDetails().HasAPathAlreadyCalculated() {
		return false
	}

	newPos := friendly.MapUnitDetails().DequeueNextPosition()
	newPosTile := n.mapState.GetLayeredMapByCurrentLocation().GetTopTile(&newPos)
	passable := newPosTile.IsWalkingPassable || newPosTile.Index.IsUnlockedDoor()
	if passable && n.mapState.PlayerLocation.Position != newPos {
		friendly.SetPos(newPos)
		return true
	}
	return false
}

func (n *NPCAIControllerSmallMap) createFreshPathToScheduledLocation(friendly *map_units.NPCFriendly) bool {
	// set up all the walkable and non walkable tiles plus the weights
	muDetails := friendly.MapUnitDetails()
	aStarMap := astar.NewAStarMap()
	// muDetails.AStarMap.InitializeByLayeredMap(
	aStarMap.InitializeByLayeredMap(
		friendly,
		n.mapState.GetLayeredMapByCurrentLocation(),
		[]references.Position{n.mapState.PlayerLocation.Position},
	)

	npcBehaviour := friendly.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(*n.dateTime)

	var path []references.Position
	if npcBehaviour.Floor != friendly.Floor() {
		// we prefer to find the best ladder or stairs
		closestFloorChangePosition := n.slr.GetClosestLadder(friendly.Pos(), friendly.Floor(), npcBehaviour.Floor)
		path = aStarMap.AStar(closestFloorChangePosition)
		// path = muDetails.AStarMap.AStar(closestFloorChangePosition)
	} else {
		path = aStarMap.AStar(npcBehaviour.Position)
	}

	muDetails.CurrentPath = path
	if len(path) == 0 {
		return false
	}
	// always pop the first because it is the current tile
	muDetails.DequeueNextPosition()
	return muDetails.HasAPathAlreadyCalculated()
}

func (n *NPCAIControllerSmallMap) wanderOneTileWithinN(muDetails *map_units.MapUnitDetails, anchorPos references.Position, withinN int) bool {

	rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator

	// Define possible moves: up, down, left, right
	directions := []references.Position{
		{X: 0, Y: -1}, // Up
		{X: 0, Y: 1},  // Down
		{X: -1, Y: 0}, // Left
		{X: 1, Y: 0},  // Right
	}

	// Shuffle the directions for randomness
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	// Try each direction to find a valid move
	for _, move := range directions {

		newPos := references.Position{
			X: muDetails.Position.X + move.X,
			Y: muDetails.Position.Y + move.Y,
		}
		if newPos.X < 0 || newPos.Y < 0 || newPos.X >= references.XSmallMapTiles || newPos.Y >= references.YSmallMapTiles {
			// we don't look outside of boundaries
			continue
		}

		// can't occupy same space as Avatar
		if n.mapState.PlayerLocation.Position.Equals(&newPos) {
			continue
		}

		if !n.mapState.GetLayeredMapByCurrentLocation().GetTopTile(&newPos).IsWalkableDuringWander() {
			continue
		}

		// Check if the new position is within N tiles of the anchorPos
		if helpers.AbsInt(int(newPos.X-anchorPos.X)) <= withinN && helpers.AbsInt(int(newPos.Y-anchorPos.Y)) <= withinN && n.mapState.IsNPCPassable(&newPos) {
			muDetails.Position.X = newPos.X
			muDetails.Position.Y = newPos.Y
			return true
		}
	}
	// If no valid moves are found, stay in the same position
	return false
}

func (n *NPCAIControllerSmallMap) getWanderDistanceByAiType(aiType references.AiType) int {
	switch aiType {
	case references.HorseWander:
		return 4
	case references.Wander:
		return 2
	case references.BigWander, references.BlackthornGuardWander, references.MerchantBuyingSellingCustom, references.MerchantBuyingSellingWander:
		return 4
	}
	return 0
}

func (n *NPCAIControllerSmallMap) RemoveAllEnemies() {
	// noop
}
