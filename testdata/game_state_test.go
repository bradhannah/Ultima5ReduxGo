package testdata_test

import (
	"sync"
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/stretchr/testify/assert"

	_ "embed"
)

const (
	xTilesVisibleOnGameScreen = 19
	yTilesVisibleOnGameScreen = 13
)

//go:embed britain2_SAVED.GAM
var saveFile []byte

//nolint:gochecknoglobals
var (
	//	baseState   *game_state.GameState
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

func getBritain2GameState() *game_state.GameState {
	getBaseGameReferences()

	cfg := config.NewUltimaVConfiguration()
	baseState := game_state.NewGameStateFromLegacySaveBytes(
		saveFile, cfg, baseGameReferences,
		xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen,
	)

	return baseState
}

func Test_BasicLoad(t *testing.T) {
	t.Parallel()

	baseState := getBritain2GameState()

	assert.Equal(t,
		references.FloorNumber(0),
		baseState.MapState.PlayerLocation.Floor,
	)
	assert.Equal(t,
		references.Position{X: 84, Y: 108},
		baseState.MapState.PlayerLocation.Position,
	)
}

// Test_IgniteTorchAndExtinguish make sure torch extinguishes itself
func Test_IgniteTorchAndExtinguish(t *testing.T) {
	t.Parallel()

	baseState := getBritain2GameState()

	baseState.IgniteTorch()

	for i := range 100 {
		assert.True(t, baseState.MapState.Lighting.HasTorchLit(), "Iteration %d", i)
		baseState.FinishTurn()
	}

	assert.False(t, baseState.MapState.Lighting.HasTorchLit())
}
