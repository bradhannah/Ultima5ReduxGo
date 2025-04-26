package game_state

type XyOccupiedMap map[int]map[int]bool

type NPCAIController interface {
	PopulateMapFirstLoad()
	CalculateNextNPCPositions()
	FreshenExistingNPCsOnMap()
	GetNpcs() *NPCS
}
