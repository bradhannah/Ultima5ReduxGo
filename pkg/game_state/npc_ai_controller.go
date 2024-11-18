package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

type NPCAIController struct {
	tileRefs *references.Tiles
	slr      *references.SmallLocationReference
}

func NewNPCAIController(
	slr *references.SmallLocationReference,
	tileRefs *references.Tiles,
) *NPCAIController {
	npcsAiCont := &NPCAIController{}

	npcsAiCont.tileRefs = tileRefs
	npcsAiCont.slr = slr

	return npcsAiCont
}
