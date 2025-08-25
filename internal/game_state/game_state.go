package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"

	"golang.org/x/exp/rand"

	"github.com/bradhannah/Ultima5ReduxGo/internal/ai"
	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/environment"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	references "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

const (
	DefaultSmallMapMinutesPerTurn = 1
)

const (
	nightTowneCloseTime = 8 + 12
	nightTowneOpenTime  = 5
)

type GameState struct {
	PartyState   party_state.PartyState
	MapState     map_state.MapState
	AiController ai.NPCAIController

	RawSave         [savedGamFileSize]byte
	MoonstoneStatus MoonstoneStatus

	DebugOptions references.DebugOptions

	GameReferences *references.GameReferences

	PartyVehicle map_units.NPCFriendly

	CurrentNPCAIController  ai.NPCAIController
	LargeMapNPCAIController map[references.World]*ai.NPCAIControllerLargeMap

	LastLargeMapPosition references.Position
	LastLargeMapFloor    references.FloorNumber

	TheOdds references.TheOdds

	DateTime datetime.UltimaDate

	ItemStacksMap references.ItemStacksMap

	// Dependency injection callbacks for external systems
	SystemCallbacks *SystemCallbacks

	// Centralized deterministic random number generator
	rng *rand.Rand

	// Environmental hazards system
	environmentalHazards *environment.EnvironmentalHazards

	// Testing overrides
	jimmySuccessForTesting func(*party_state.PlayerCharacter) bool
}

func initBlankGameState(gameConfig *config.UltimaVConfiguration,
	gameReferences *references.GameReferences,
	xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen int,
) *GameState {
	gameState := &GameState{} //nolint:exhaustruct

	gameState.GameReferences = gameReferences

	// Initialize deterministic RNG with fixed seed for consistent behavior
	gameState.rng = rand.New(rand.NewSource(1))
	gameState.MapState = *map_state.NewMapState(
		map_state.NewMapStateInput{
			GameDimensions:            gameState,
			XTilesVisibleOnGameScreen: xTilesVisibleOnGameScreen,
			YTilesVisibleOnGameScreen: yTilesVisibleOnGameScreen,
		},
	)

	return gameState
}

func NewGameStateFromLegacySaveFile(savedGamFilePath string,
	gameConfig *config.UltimaVConfiguration,
	gameReferences *references.GameReferences,
	xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen int,
) *GameState {
	rawSaveData, err := getLegacySavedGamRaw(savedGamFilePath)
	if err != nil {
		log.Fatalf("Error loading saved gam raw data: %v", err)
	}

	gameState := NewGameStateFromLegacySaveBytes(rawSaveData,
		gameConfig,
		gameReferences,
		xTilesVisibleOnGameScreen,
		yTilesVisibleOnGameScreen)
	return gameState
}

func NewGameStateFromLegacySaveBytes(rawSaveData []byte,
	gameConfig *config.UltimaVConfiguration,
	gameReferences *references.GameReferences,
	xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen int,
) *GameState {
	gameState := initBlankGameState(gameConfig, gameReferences, xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen)
	err := gameState.LoadLegacySaveGameFromBytes(rawSaveData)
	if err != nil {
		log.Fatalf("Error loading legacy save game from bytes: %v", err)
	}
	return gameState
}

func (g *GameState) IsAvatarAtPosition(pos *references.Position) bool {
	return g.MapState.PlayerLocation.Position.Equals(pos)
}

func (g *GameState) GetCurrentSmallLocationReference() *references.SmallLocationReference {
	return g.GameReferences.LocationReferences.GetLocationReference(g.MapState.PlayerLocation.Location)
}

func (g *GameState) GetLayeredMapByCurrentLocation() *map_state.LayeredMap {
	return g.MapState.LayeredMaps.GetLayeredMap(g.MapState.PlayerLocation.Location.GetMapType(), g.MapState.PlayerLocation.Floor)
}

func (g *GameState) IsOutOfBounds(position references.Position) bool {
	if g.MapState.PlayerLocation.Location.GetMapType() == references.LargeMapType {
		if position.X > references.XLargeMapTiles-1 || position.Y >= references.YLargeMapTiles {
			log.Fatalf("Exceeded large map tiles: X=%d, Y=%d", position.X, position.Y)
		}
		return false
	}

	if position.X < 0 || position.Y < 0 {
		return true
	}

	if position.X > g.GetCurrentSmallLocationReference().GetMaxX() ||
		position.Y > g.GetCurrentSmallLocationReference().GetMaxY() {
		return true
	}

	return false
}

func (g *GameState) GetCurrentLayeredMap() *map_state.LayeredMap {
	return g.MapState.LayeredMaps.GetLayeredMap(
		g.MapState.PlayerLocation.Location.GetMapType(),
		g.MapState.PlayerLocation.Floor)
}

func (g *GameState) GetCurrentLayeredMapAvatarTopTile() *references.Tile {
	theMap := g.GetCurrentLayeredMap()
	topTile := theMap.GetTopTile(&g.MapState.PlayerLocation.Position)
	if topTile == nil {
		return nil
	}
	return topTile
}

func (g *GameState) IsPassable(pos *references.Position) bool {
	theMap := g.GetCurrentLayeredMap()
	topTile := theMap.GetTopTile(pos)

	if topTile == nil {
		return false
	}

	// Check terrain passability first (like legalmove() in original)
	if !topTile.IsPassable(g.PartyVehicle.GetVehicleDetails().VehicleType) {
		return false
	}

	// Check for object collision (like looklist() in original)
	// This is what was missing - objects were not blocking movement!
	if g.ObjectPresentAt(pos) {
		return false
	}

	return true
}

func (g *GameState) GetArchwayPortcullisSpriteByTime() indexes.SpriteIndex {
	if g.DateTime.Hour >= nightTowneCloseTime || g.DateTime.Hour < nightTowneOpenTime {
		return indexes.Portcullis
	}

	return indexes.BrickWallArchway
}

func (g *GameState) GetDrawBridgeWaterByTime(origIndex indexes.SpriteIndex) indexes.SpriteIndex {
	if g.DateTime.Hour >= nightTowneCloseTime || g.DateTime.Hour < nightTowneOpenTime {
		return indexes.WaterShallow
	}

	return origIndex
}

type BoardVehicleResult int

const (
	BoardVehicleResultSuccess BoardVehicleResult = iota
	BoardVehicleResultNoVehicle
	BoardVehicleResultNoSkiffs
)

func (g *GameState) BoardVehicle(vehicle map_units.NPCFriendly) BoardVehicleResult {
	result := BoardVehicleResultSuccess

	if vehicle.GetVehicleDetails().VehicleType == references.FrigateVehicle {
		if g.PartyVehicle.GetVehicleDetails().VehicleType == references.SkiffVehicle {
			vehicle.GetVehicleDetails().IncrementSkiffQuantity()
		} else {
			result = BoardVehicleResultNoSkiffs
		}
	}

	g.PartyVehicle = vehicle

	if !g.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(g.MapState.PlayerLocation.Position) {
		log.Fatalf("Unexpected - tried to remove NPC at position X=%d,Y=%d but failed",
			g.MapState.PlayerLocation.Position.X, g.MapState.PlayerLocation.Position.Y)
	}

	return result
}

func (g *GameState) GetTilesVisibleOnScreen() (int, int) {
	return g.GetCurrentLayeredMap().GetTilesVisibleOnScreen()
}

func (g *GameState) GetTopLeftExtent() references.Position {
	return g.GetCurrentLayeredMap().TopLeft
}

func (g *GameState) GetBottomRightExtent() references.Position {
	return g.GetCurrentLayeredMap().BottomRight
}

func (g *GameState) IsWrappedMap() bool {
	return g.GetCurrentLayeredMap().IsWrappedMap()
}

func (g *GameState) GetBottomRightWithoutOverflow() references.Position {
	return g.GetCurrentLayeredMap().GetBottomRightWithoutOverflow()
}

func (g *GameState) GetCurrentLargeMapNPCAIController() *ai.NPCAIControllerLargeMap {
	if g.MapState.PlayerLocation.Location.GetMapType() != references.LargeMapType {
		log.Fatalf("Expected large map type, got %d", g.MapState.PlayerLocation.Location.GetMapType())
	}
	var npcAiCon *ai.NPCAIControllerLargeMap

	if g.MapState.IsOverworld() {
		npcAiCon = g.LargeMapNPCAIController[references.OVERWORLD]
	} else {
		npcAiCon = g.LargeMapNPCAIController[references.UNDERWORLD]
	}
	return npcAiCon
}

// Centralized RNG methods - replace helpers package usage for deterministic behavior

// SetRandomSeed sets the seed for the centralized RNG (for testing and deterministic behavior)
func (g *GameState) SetRandomSeed(seed uint64) {
	g.rng.Seed(seed)
}

// OneInXOdds returns true with probability 1/odds (deterministic replacement for helpers.OneInXOdds)
func (g *GameState) OneInXOdds(odds int) bool {
	if odds <= 0 {
		return false
	}
	return g.rng.Intn(odds) == 0
}

// HappenedByPercentLikely returns true with given likelihood percentage (deterministic replacement for helpers.HappenedByPercentLikely)
func (g *GameState) HappenedByPercentLikely(likelihoodToSucceedPercent int) bool {
	if likelihoodToSucceedPercent >= 100 {
		return true
	}
	if likelihoodToSucceedPercent <= 0 {
		return false
	}
	return g.RandomIntInRange(0, 100) < likelihoodToSucceedPercent
}

// RandomIntInRange returns a random integer in [min, max] inclusive (deterministic replacement for helpers.RandomIntInRange)
func (g *GameState) RandomIntInRange(min, max int) int {
	if min > max {
		log.Fatalf("RandomIntInRange: min (%d) cannot be greater than max (%d)", min, max)
	}
	return g.rng.Intn(max-min+1) + min
}

// PickOneOfInt returns one of two int values randomly (deterministic replacement for helpers.PickOneOf)
func (g *GameState) PickOneOfInt(a, b int) int {
	if g.OneInXOdds(2) {
		return a
	}
	return b
}

// ShuffleDirections shuffles a slice of positions for deterministic randomness in AI movement
func (g *GameState) ShuffleDirections(directions []references.Position) {
	g.rng.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
}
