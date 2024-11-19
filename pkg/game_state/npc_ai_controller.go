package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCAIController struct {
	tileRefs  *references.Tiles
	slr       *references.SmallLocationReference
	gameState *GameState
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

func (n *NPCAIController) PopulateMapFirstLoad(
	lm *LayeredMap,
	ud datetime.UltimaDate) {

	// get the correct schedule
	for _, npc := range *n.slr.GetNPCs() {
		npc.Schedule.GetIndividualNPCBehaviourByUltimaDate(ud)
	}

	// place all characters within that schedule

	// lm.SetTileByLayer(MapUnitLayer,
	// 	)
}
