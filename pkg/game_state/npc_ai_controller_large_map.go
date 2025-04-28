package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCAIControllerLargeMap struct {
	tileRefs  *references.Tiles
	World     references.World
	gameState *GameState

	npcs NPCS

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

	return npcsAiCont
}

func (n *NPCAIControllerLargeMap) GetNpcs() *NPCS {
	return &n.npcs
}

func (n *NPCAIControllerLargeMap) PopulateMapFirstLoad() {
}

func (n *NPCAIControllerLargeMap) AdvanceNextTurnCalcAndMoveNPCs() {
	n.clearMapUnitsFromMap()

	n.positionOccupiedChance = createFreshXyOccupiedMap()

	// n.updateAllNPCAiTypes()
	// n.positionOccupiedChance = n.createFreshXyOccupiedMap()

	// for _, npc := range n.npcs {
	// 	if npc.IsEmptyNPC() {
	// 		continue
	// 	}
	// 	// very lazy approach - but making sure every NPC is in correct spot on map
	// 	// for every iteration makes sure next NPC doesn't assign the same tile space
	// 	n.FreshenExistingNPCsOnMap()
	// 	n.calculateNextNPCPosition(npc)
	// }
	// n.FreshenExistingNPCsOnMap()
}

func (n *NPCAIControllerLargeMap) FreshenExistingNPCsOnMap() {
}

func (n *NPCAIControllerLargeMap) clearMapUnitsFromMap() {
	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
}
