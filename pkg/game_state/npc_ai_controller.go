package game_state

import (
	"log"
	"time"

	"golang.org/x/exp/rand"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type XyOccupiedMap map[int]map[int]bool

type NPCAIController struct {
	tileRefs  *references.Tiles
	slr       *references.SmallLocationReference
	gameState *GameState

	npcs []*NPC

	positionOccupiedChance *XyOccupiedMap
}

func NewNPCAIController(
	slr *references.SmallLocationReference,
	tileRefs *references.Tiles,
	gameState *GameState,
) *NPCAIController {
	npcsAiCont := &NPCAIController{}

	npcsAiCont.tileRefs = tileRefs
	npcsAiCont.slr = slr
	npcsAiCont.gameState = gameState

	xy := make(XyOccupiedMap)
	npcsAiCont.positionOccupiedChance = &xy

	return npcsAiCont
}

func (n *NPCAIController) createFreshXyOccupiedMap() *XyOccupiedMap {
	xy := make(XyOccupiedMap)
	for _, npc := range n.npcs {
		if npc.IsEmptyNPC() {
			continue
		}
		// indiv := npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(n.gameState.DateTime)
		_, exists := xy[int(npc.Position.X)]
		if !exists {
			xy[int(npc.Position.X)] = make(map[int]bool)
		}
		xy[int(npc.Position.X)][int(npc.Position.Y)] = true
	}
	return &xy
}

func (n *NPCAIController) generateNPCs() {
	npcs := make([]*NPC, 0)
	// get the correct schedule
	npcsRefs := n.slr.GetNPCReferences()
	for _, npcRef := range *npcsRefs {
		npc := NewNPC(npcRef)
		npcs = append(npcs, &npc)
	}
	n.npcs = npcs
}

func (n *NPCAIController) PopulateMapFirstLoad() {
	// lm *LayeredMap,
	// ud datetime.UltimaDate) {

	n.generateNPCs()

	for i, npc := range n.npcs {
		_ = i
		if npc.IsEmptyNPC() {
			continue
		}
		indiv := npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(n.gameState.DateTime)

		npc.Position = indiv.Position
		npc.Floor = indiv.Floor
		npc.AiType = indiv.Ai
	}
	n.setAllNPCTiles()
}

func (n *NPCAIController) updateAllNPCAiTypes() {
	for i, npc := range n.npcs {
		_ = i
		if npc.IsEmptyNPC() {
			continue
		}

		indiv := npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(n.gameState.DateTime)

		npc.AiType = indiv.Ai
	}
}

func (n *NPCAIController) setAllNPCTiles() {
	lm := n.gameState.GetLayeredMapByCurrentLocation()

	for _, npc := range n.npcs {
		if npc.IsEmptyNPC() {
			continue
		}
		if n.gameState.Floor == npc.Floor {
			lm.SetTileByLayer(MapUnitLayer, &npc.Position, npc.NPCReference.GetTileIndex())
		}
	}
}

func (n *NPCAIController) clearMapUnitsFromMap() {
	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
}

func (n *NPCAIController) CalculateNextNPCPositions() {
	n.clearMapUnitsFromMap()
	n.updateAllNPCAiTypes()
	n.positionOccupiedChance = n.createFreshXyOccupiedMap()

	for _, npc := range n.npcs {
		if npc.IsEmptyNPC() {
			continue
		}
		n.calculateNextNPCPosition(npc)
	}
	n.setAllNPCTiles()
}

func (n *NPCAIController) calculateNextNPCPosition(npc *NPC) {
	refBehaviour := npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(n.gameState.DateTime)

	if npc.Position.Equals(refBehaviour.Position) {
		if n.performAiMovementOnAssignedPosition(npc) {
			return
		}
	} else {
		if n.performAiMovementNotOnAssignedPosition(npc) {
			return
		}
	}
}

func (n *NPCAIController) performAiMovementOnAssignedPosition(npc *NPC) bool {
	npcSched := npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(n.gameState.DateTime)

	switch npc.AiType {
	case references.BlackthornGuardFixed, references.Fixed:
	case references.MerchantBuyingSellingCustom, references.MerchantBuyingSellingWander, references.Wander:
		// wander 2
		n.WanderOneTileWithinN(npc, npcSched.Position, 2)
		return true
	case references.BigWander, references.BlackthornGuardWander:
		// wander 5
		n.WanderOneTileWithinN(npc, npcSched.Position, 5)
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
		// wander 4
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

func (n *NPCAIController) getWanderDistanceByAiType(aiType references.AiType) int {
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

func (n *NPCAIController) performAiMovementNotOnAssignedPosition(npc *NPC) bool {
	npcSched := npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(n.gameState.DateTime)
	nWanderDistance := n.getWanderDistanceByAiType(npc.AiType)

	if n.MoveNPCOnCalculatedPath(npc) {
		return true
	}

	switch npc.AiType {
	case references.BlackthornGuardFixed, references.Fixed, references.CustomAi, references.MerchantBuyingSelling:
		if n.CreateFreshPathToScheduledLocation(npc) {
			npc.Position = npc.DequeueNextPosition()
			return true
		}
		return false
	case references.BigWander, references.BlackthornGuardWander, references.MerchantBuyingSellingCustom, references.MerchantBuyingSellingWander, references.Wander, references.HorseWander:
		if !npcSched.Position.IsWithinN(&npc.Position, nWanderDistance) {
			if n.CreateFreshPathToScheduledLocation(npc) {
				npc.Position = npc.DequeueNextPosition()
				return true
			}
			return false
		}
		return n.WanderOneTileWithinN(npc, npcSched.Position, nWanderDistance)
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
		if !npcSched.Position.IsWithinN(&npc.Position, nWanderDistance) {
			if n.CreateFreshPathToScheduledLocation(npc) {
				npc.Position = npc.DequeueNextPosition()
				return true
			}
			return false
		}
		if helpers.OneInXOdds(3) {
			return n.WanderOneTileWithinN(npc, npcSched.Position, nWanderDistance)
		}
		return false
	default:
		log.Fatal("Unknown AiType")
	}
	return false
}

func (n *NPCAIController) MoveNPCOnCalculatedPath(npc *NPC) bool {
	if !npc.HasAPathAlreadyCalculated() {
		return false
	}

	newPos := npc.DequeueNextPosition()
	if n.gameState.GetLayeredMapByCurrentLocation().GetTopTile(&newPos).IsWalkingPassable && n.gameState.Position != newPos {
		npc.Position = newPos
		return true
	}
	return false
}

func (n *NPCAIController) CreateFreshPathToScheduledLocation(npc *NPC) bool {
	// set up all the walkable and non walkable tiles plus the weights
	npc.AStarMap.InitializeByLayeredMap(
		n.gameState.GetLayeredMapByCurrentLocation(),
		[]references.Position{n.gameState.Position},
	)

	path := npc.AStarMap.AStar(npc.Position, references.Position{X: 15, Y: 30})

	npc.CurrentPath = &path
	if len(path) == 0 {
		return false
	}
	// always pop the first because it is the current tile
	npc.DequeueNextPosition()
	return npc.HasAPathAlreadyCalculated()
}

func (n *NPCAIController) WanderOneTileWithinN(npc *NPC, anchorPos references.Position, withinN int) bool {
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
			X: npc.Position.X + move.X,
			Y: npc.Position.Y + move.Y,
		}

		// can't occupy same space as Avatar
		if n.gameState.Position.Equals(newPos) {
			continue
		}

		// Check if the new position is within N tiles of the anchorPos
		if helpers.AbsInt(int(newPos.X-anchorPos.X)) <= withinN && helpers.AbsInt(int(newPos.Y-anchorPos.Y)) <= withinN && n.gameState.IsNPCPassable(&newPos) {
			npc.Position.X = newPos.X
			npc.Position.Y = newPos.Y
			return true
		}
	}
	// If no valid moves are found, stay in the same position
	return false
}
