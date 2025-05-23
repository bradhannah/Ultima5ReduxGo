package testdata

import (
	"bytes"
	"encoding/gob"
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
	baseState   *game_state.GameState
	loadOnce    sync.Once
	baseInitErr error
)

// getBaseState loads the save only once.
func getBaseState(t *testing.T) *game_state.GameState {
	loadOnce.Do(func() {
		cfg := config.NewUltimaVConfiguration()
		refs, err := references.NewGameReferences(cfg)
		if err != nil {
			baseInitErr = err
			return
		}
		baseState = game_state.NewGameStateFromLegacySaveBytes(
			saveFile, cfg, refs,
			xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen,
		)
	})
	if baseInitErr != nil {
		t.Fatalf("cannot init base GameState: %v", baseInitErr)
	}
	return baseState
}

func cloneGameState(src *game_state.GameState) *game_state.GameState {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	_ = enc.Encode(src)

	var dst game_state.GameState
	_ = dec.Decode(&dst)
	return &dst
}

func Test_BasicLoad(t *testing.T) {
	orig := cloneGameState(getBaseState(t))

	assert.Equal(t,
		references.FloorNumber(0),
		orig.MapState.PlayerLocation.Floor,
	)
}
