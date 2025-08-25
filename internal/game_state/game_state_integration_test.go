package game_state

import (
	"sync"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/stretchr/testify/assert"

	_ "embed"
)

const (
	xTilesVisibleOnGameScreen = 19
	yTilesVisibleOnGameScreen = 13
)

//go:embed testdata/britain2_SAVED.GAM
var saveFile []byte

//nolint:gochecknoglobals
var (
	baseGameReferences *references.GameReferences
	loadOnce           sync.Once
)

func getBaseGameReferences() *references.GameReferences {
	loadOnce.Do(func() {
		cfg := config.NewUltimaVConfiguration()

		var err error
		baseGameReferences, err = references.NewGameReferences(cfg)

		if err != nil {
			return
		}
	})

	return baseGameReferences
}

func getBritain2GameState() *GameState {
	getBaseGameReferences()

	cfg := config.NewUltimaVConfiguration()
	baseState := NewGameStateFromLegacySaveBytes(
		saveFile, cfg, baseGameReferences,
		xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen,
	)

	return baseState
}

func Test_LoadMetAndDeadNpcs(t *testing.T) {
	baseState := getBritain2GameState()
	ps := &baseState.PartyState

	castleLoc := references.Location(0)
	assert.False(t, ps.MetNpcs()[castleLoc][0], "NPC 0 in Castle should not be met")
	assert.False(t, ps.DeadNpcs()[castleLoc][0], "NPC 0 in Castle should not be dead")
	assert.False(t, ps.MetNpcs()[castleLoc][1], "NPC 1 in Castle should not be met")
	assert.False(t, ps.DeadNpcs()[castleLoc][1], "NPC 1 in Castle should not be dead")
	assert.False(t, ps.MetNpcs()[castleLoc][2], "NPC 2 in Castle should not be met")
	assert.False(t, ps.DeadNpcs()[castleLoc][2], "NPC 2 in Castle should not be dead")
	assert.True(t, ps.MetNpcs()[castleLoc][3], "NPC 3 in Castle should be met")
	assert.False(t, ps.DeadNpcs()[castleLoc][3], "NPC 3 in Castle should not be dead")
	assert.False(t, ps.MetNpcs()[castleLoc][4], "NPC 4 in Castle should not be met")
	assert.False(t, ps.DeadNpcs()[castleLoc][4], "NPC 4 in Castle should not be dead")

	britainLoc := references.Location(1)
	assert.False(t, ps.MetNpcs()[britainLoc][0], "NPC 0 in Britain should not be met")
	assert.False(t, ps.DeadNpcs()[britainLoc][0], "NPC 0 in Britain should not be dead")

	// Debug print actual values for the first few NPCs in Castle and Britain
	for i := 0; i < 5; i++ {
		metCastle := ps.MetNpcs()[castleLoc][i]
		deadCastle := ps.DeadNpcs()[castleLoc][i]
		t.Logf("Castle NPC %d: met=%v, dead=%v", i, metCastle, deadCastle)
	}
	for i := 0; i < 5; i++ {
		metBritain := ps.MetNpcs()[britainLoc][i]
		deadBritain := ps.DeadNpcs()[britainLoc][i]
		t.Logf("Britain NPC %d: met=%v, dead=%v", i, metBritain, deadBritain)
	}
}
