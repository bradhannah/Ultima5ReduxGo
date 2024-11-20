package game_state

import (
	// _ "github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCAIController struct {
	tileRefs  *references.Tiles
	slr       *references.SmallLocationReference
	gameState *GameState

	npcs []*NPC
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

	for _, npc := range n.npcs {
		if npc.IsEmptyNPC() {
			continue
		}
		n.calculateNextNPCPosition(npc)
	}
	n.setAllNPCTiles()
}

func (n *NPCAIController) calculateNextNPCPosition(npc *NPC) {
	// npc.Position = *npc.Position.GetPositionUp()
}
