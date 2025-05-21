package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/ai"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (m *GameState) UpdateSmallMap(tileRefs *references.Tiles, locs *references.LocationReferences) {
	slr := locs.GetLocationReference(m.MapState.PlayerLocation.Location)
	m.MapState.LayeredMaps.ResetAndCreateSmallMap(
		slr,
		tileRefs,
		m.MapState.XTilesVisibleOnGameScreen,
		m.MapState.YTilesVisibleOnGameScreen)
	m.CurrentNPCAIController = ai.NewNPCAIControllerSmallMap(slr, tileRefs, &m.MapState, &m.DateTime) // , m.MapState.PlayerLocation)
	m.CurrentNPCAIController.PopulateMapFirstLoad()
}
