package testdata

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

// shared across tests
var (
	//	baseState   *game_state.GameState
	baseGameReferences *references.GameReferences
	loadOnce           sync.Once
	baseInitErr        error
)

func getBaseGameReferences() *references.GameReferences {
	loadOnce.Do(func() {
		cfg := config.NewUltimaVConfiguration()

		var err error
		baseGameReferences, err = references.NewGameReferences(cfg)

		if err != nil {
			baseInitErr = err
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
