package game_state

import (
	"fmt"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) EnterBuilding(
	slr *references.SmallLocationReference,
	tileRefs *references.Tiles,
	npcRefs *references.NPCReferences,
) {
	if slr.Location != references.EmptyLocation {
		g.LastLargeMapPosition = g.Position
		g.LastLargeMapFloor = g.Floor
		g.Position = references.Position{
			X: 15,
			Y: 30,
		}
		g.Location = slr.Location
		g.Floor = 0
		g.LayeredMaps.ResetAndCreateSmallMap(slr, tileRefs, g.XTilesInMap, g.YTilesInMap)
		g.initializeNPCs(npcRefs, tileRefs)
	}
}

func (g *GameState) initializeNPCs(
	npcReferences *references.NPCReferences,
	tileRefs *references.Tiles,
) {
	npcs := g.GameReferences.NPCReferences.GetNPCReferencesByLocation(g.Location)
	fmt.Sprintf("%d", npcs)
}
