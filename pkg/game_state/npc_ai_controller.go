package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCAIController struct {
	tileRefs  *references.Tiles
	slr       *references.SmallLocationReference
	gameState *GameState

	npcs *[]NPC
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

	return npcsAiCont
}

func (n *NPCAIController) generateNPCs() {
	npcs := make([]NPC, 0)
	// get the correct schedule
	npcsRefs := n.slr.GetNPCReferences()
	for _, npcRef := range *npcsRefs {
		npc := NewNPC(npcRef)
		npcs = append(npcs, npc)
	}
	n.npcs = &npcs
}

func (n *NPCAIController) PopulateMapFirstLoad(
	lm *LayeredMap,
	ud datetime.UltimaDate) {

	n.generateNPCs()

	for _, npc := range *n.npcs {
		if npc.IsEmptyNPC() {
			continue
		}
		indiv := npc.NPCReference.Schedule.GetIndividualNPCBehaviourByUltimaDate(ud)

		if n.gameState.Floor == indiv.Floor {
			lm.SetTileByLayer(MapUnitLayer, &indiv.Position, npc.NPCReference.GetTileIndex())
		}
	}
}

func (n *NPCAIController) clearMapUnitsFromMap() {
	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
	// .GetTileByLayer(MapUnitLayer, &n.gameState.Position)

}

func (n *NPCAIController) CalculateNextNPCPositions() {
	// n.clearMapUnitsFromMap()
	for _, npc := range *n.npcs {
		if npc.IsEmptyNPC() {
			continue
		}
	}
}
