package game_state

type NPCAIController interface {
	PopulateMapFirstLoad()
	CalculateNextNPCPositions()
	FreshenExistingNPCsOnMap()
	GetNpcs() *NPCS
}
