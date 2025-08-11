package main

import "github.com/bradhannah/Ultima5ReduxGo/internal/clock"

// AdvanceClockOnly advances the game clock and runs fixed steps.
// Left available as a helper; the main Update in gamescene_input.go calls similar logic.
func (g *GameScene) AdvanceClockOnly() {
	if g.clk == nil {
		// Safety: create a clock if not initialized for any reason
		g.clk = clock.NewGameClock(16)
	}
	steps := g.clk.Advance()

	for i := 0; i < steps; i++ {
		g.updateFixed()
	}
}

// updateFixed is a placeholder for core logic to run at a fixed rate.
func (g *GameScene) updateFixed() {
	// Intentionally left minimal; hook up AI, schedules, etc., later.
}
