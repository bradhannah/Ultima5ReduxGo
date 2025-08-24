package ai

import references "github.com/bradhannah/Ultima5ReduxGo/internal/references"

// ShouldMonsterMoveInCombat is a placeholder for combat AI movement gating.
// Legacy reference: OLD/COMSUBS1.C prevents REAPER and MIMIC from moving.
// This version generalizes to any enemy with AdditionalEnemyFlags.DoNotMove.
func ShouldMonsterMoveInCombat(enemy references.EnemyReference) bool {
	if enemy.AdditionalEnemyFlags.DoNotMove {
		return false
	}
	// TODO: When combat AI is implemented, incorporate other gating
	// conditions and special cases per legacy behavior.
	return true
}
