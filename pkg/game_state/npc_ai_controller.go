package game_state

type XyOccupiedMap map[int]map[int]bool

func createFreshXyOccupiedMap() *XyOccupiedMap {
	xy := make(XyOccupiedMap)
	for _, npc := range n.npcs {
		if npc.IsEmptyNPC() || !npc.Visible {
			continue
		}
		_, exists := xy[int(npc.Position.X)]
		if !exists {
			xy[int(npc.Position.X)] = make(map[int]bool)
		}
		xy[int(npc.Position.X)][int(npc.Position.Y)] = true
	}
	return &xy
}

type NPCAIController interface {
	PopulateMapFirstLoad()
	AdvanceNextTurnCalcAndMoveNPCs()
	FreshenExistingNPCsOnMap()
	GetNpcs() *NPCS
}
