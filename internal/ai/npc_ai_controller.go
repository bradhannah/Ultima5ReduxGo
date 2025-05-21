package ai

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
)

type NPCAIController interface {
	PopulateMapFirstLoad()
	AdvanceNextTurnCalcAndMoveNPCs()
	FreshenExistingNPCsOnMap()
	GetNpcs() *map_units.MapUnits
	RemoveAllEnemies()
}
