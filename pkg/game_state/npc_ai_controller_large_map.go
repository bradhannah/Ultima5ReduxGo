package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCAIControllerLargeMap struct {
	tileRefs  *references.Tiles
	World     references.World
	gameState *GameState

	npcs MapUnits

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

func (n *NPCAIControllerLargeMap) GetNpcs() *MapUnits {
	return &n.npcs
}

func (n *NPCAIControllerLargeMap) PopulateMapFirstLoad() {
}

func (n *NPCAIControllerLargeMap) AdvanceNextTurnCalcAndMoveNPCs() {
	//n.clearMapUnitsFromMap()

	// n.updateAllNPCAiTypes()
	n.positionOccupiedChance = n.npcs.createFreshXyOccupiedMap()

	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()

	for _, npc := range n.npcs {
		if npc.IsEmptyMapUnit() {
			continue
		}
		// 	// very lazy approach - but making sure every NPC is in correct spot on map
		// 	// for every iteration makes sure next NPC doesn't assign the same tile space
		n.FreshenExistingNPCsOnMap()
		// 	n.calculateNextNPCPosition(npc)
	}
	n.FreshenExistingNPCsOnMap()

	// should we spawn units after these ones have moved? probably
}

func (n *NPCAIControllerLargeMap) FreshenExistingNPCsOnMap() {
	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
	n.placeNPCsOnLayeredMap()
}

// func (n *NPCAIControllerLargeMap) clearMapUnitsFromMap() {
// 	n.gameState.GetLayeredMapByCurrentLocation().ClearMapUnitTiles()
// }

func (n *NPCAIControllerLargeMap) placeNPCsOnLayeredMap() {
	lm := n.gameState.GetLayeredMapByCurrentLocation()

	for _, npc := range n.npcs {
		if npc.IsEmptyMapUnit() || !npc.IsVisible() {
			continue
		}
		switch mapUnit := npc.(type) {
		case *NPCFriendly:
			if n.gameState.Floor == npc.Floor() {
				lm.SetTileByLayer(MapUnitLayer, npc.PosPtr(), mapUnit.NPCReference.GetTileIndex())
			}
		case *NPCEnemy:
			if n.gameState.Floor == npc.Floor() {
				//lm.SetTileByLayer(MapUnitLayer, npc.PosPtr(), mapUnit.GetTileIndex())
			}
		}
	}
}

func (n *NPCAIControllerLargeMap) generateEraBoundMonster() {
	//enemy := NewFriendlyNPC(references.NewNPCReferenceForMonster(references.MonsterType(0), references.Position{}, references.FloorNumber(0)), 0)
}
