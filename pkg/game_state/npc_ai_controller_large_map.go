package game_state

import (
	"fmt"

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

func (n *NPCAIControllerLargeMap) CalculateNextNPCPositions() {
	fmt.Print("hello")
}

func (n *NPCAIControllerLargeMap) FreshenExistingNPCsOnMap() {
}
