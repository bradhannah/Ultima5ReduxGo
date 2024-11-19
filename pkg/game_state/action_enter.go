package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameState) EnterBuilding(
	slr *references.SmallLocationReference,
	tileRefs *references.Tiles,
	// npcRefs *references.NPCReferences,
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
		g.LayeredMaps.ResetAndCreateSmallMap(
			slr,
			tileRefs,
			// npcRefs,
			g.XTilesInMap,
			g.YTilesInMap)
		g.NPCAIController = *NewNPCAIController(slr, tileRefs, g)
		g.NPCAIController.PopulateMapFirstLoad(
			g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor),
			g.DateTime,
		)
	}
}

// func (g *GameState) initializeNPCs(
// 	npcReferences *references.NPCReferences,
// 	tileRefs *references.Tiles,
// ) {
// 	npcs := g.GameReferences.NPCReferences.GetNPCReferencesByLocation(g.Location)
//
// 	g.LayeredMaps.ResetAndCreateSmallMap()
// 	fmt.Sprintf("%d", npcs)
// }
