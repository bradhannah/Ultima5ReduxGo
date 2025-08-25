package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/environment"
)

// messageCallbacksAdapter adapts GameState SystemCallbacks to environment.SystemCallbacks interface
type messageCallbacksAdapter struct {
	systemCallbacks *SystemCallbacks
}

func (m *messageCallbacksAdapter) AddRowStr(message string) {
	if m.systemCallbacks != nil && m.systemCallbacks.Message.AddRowStr != nil {
		m.systemCallbacks.Message.AddRowStr(message)
	}
}

// processEnvironmentalHazards checks and applies environmental hazards based on current player position
func (gs *GameState) processEnvironmentalHazards() {
	currentTile := gs.MapState.GetLayeredMapByCurrentLocation().GetTopTile(&gs.MapState.PlayerLocation.Position)

	if gs.environmentalHazards == nil {
		gs.environmentalHazards = environment.NewEnvironmentalHazards(gs.rng, &messageCallbacksAdapter{gs.SystemCallbacks})
	}

	hazardResult := gs.environmentalHazards.CheckTileHazards(
		currentTile,
		&gs.PartyState,
		gs.MapState.PlayerLocation.Location.GetMapType(),
	)

	// Log hazard for testing/debugging
	if hazardResult.Type != environment.NoHazard {
		gs.logEnvironmentalHazard(hazardResult)
	}
}

// logEnvironmentalHazard provides debug logging for environmental hazard events
func (gs *GameState) logEnvironmentalHazard(hazardResult environment.HazardResult) {
	// This could be expanded for more detailed logging if needed
	// For now, the SystemCallbacks already handle user-facing messages
}

// ProcessEnvironmentalHazardsAfterMovement should be called after player movement
func (gs *GameState) ProcessEnvironmentalHazardsAfterMovement() {
	gs.processEnvironmentalHazards()
}

// ProcessEnvironmentalHazardsOnTurnAdvancement should be called during turn processing for standing-on-tile effects
func (gs *GameState) ProcessEnvironmentalHazardsOnTurnAdvancement() {
	gs.processEnvironmentalHazards()
}
