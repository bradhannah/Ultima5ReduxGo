package game_state

import (
	"fmt"
	"os"

	"github.com/bradhannah/Ultima5ReduxGo/internal/ai"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const savedGamFileSize = 4192

type StartingMemoryAddressUb uint16
type StartingMemoryAddressU16 uint16

func (g *GameState) getLegacySavedGamRaw(savedGamFilePath string) ([]byte, error) {
	// Open the file in read-only mode and as binary
	file, err := os.OpenFile(savedGamFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	buffer := make([]byte, savedGamFileSize)
	n, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	if n != savedGamFileSize {
		return nil, fmt.Errorf("expected file of size 4192 but was %d", n)
	}

	return buffer, nil
}

func (g *GameState) LoadLegacySaveGame(savedGamFilePath string, gameRefs *references.GameReferences) error {
	// Open the file in read-only mode and as binary
	rawSaveGameBytesFromDisk, err := g.getLegacySavedGamRaw(savedGamFilePath)
	if err != nil {
		return err
	}

	// var saveGame = GameState{}
	g.RawSave = [savedGamFileSize]byte(rawSaveGameBytesFromDisk)

	// Overlay player characters over memory rawSaveGameBytesFromDisk to easily consume data
	g.PartyState.LoadFromRaw(g.RawSave)

	// world and position
	const lbLocation = 0x2ED
	const lbX = 0x2F0
	const lbY = 0x2F1
	const lbFloor = 0x2EF
	g.MapState.PlayerLocation.Location = references.Location(rawSaveGameBytesFromDisk[lbLocation])
	g.MapState.PlayerLocation.Position = references.Position{X: references.Coordinate(rawSaveGameBytesFromDisk[lbX]), Y: references.Coordinate(rawSaveGameBytesFromDisk[lbY])}
	g.MapState.PlayerLocation.Floor = references.FloorNumber(rawSaveGameBytesFromDisk[lbFloor])

	// Date/Time
	const lsYear = 0x2CE
	const lbMonth = 0x2D7
	const lbDay = 0x2D8
	const lbHour = 0x2D9
	const lbMinute = 0x2DB
	g.DateTime.Year = getUint16(&rawSaveGameBytesFromDisk, lsYear)
	g.DateTime.Month = rawSaveGameBytesFromDisk[lbMonth]
	g.DateTime.Day = rawSaveGameBytesFromDisk[lbDay]
	g.DateTime.Hour = rawSaveGameBytesFromDisk[lbHour]
	g.DateTime.Minute = rawSaveGameBytesFromDisk[lbMinute]

	// Various Things
	const lbKarma = 0x2e2
	const lsGold = 0x204
	g.PartyState.Karma = party_state.Karma(rawSaveGameBytesFromDisk[lbKarma])
	g.PartyState.Inventory.Gold = getUint16(&rawSaveGameBytesFromDisk, lsGold)

	// ProvisionsQuantity
	const lsFood = 0x202
	const lbKeys = 0x206
	const lbGems = 0x207
	const lbTorches = 0x208
	const lbSkullKeys = 0x20B
	g.PartyState.Inventory.Provisions.Food = getUint16(&rawSaveGameBytesFromDisk, lsFood)
	g.PartyState.Inventory.Provisions.Gems = rawSaveGameBytesFromDisk[lbGems]
	g.PartyState.Inventory.Provisions.Torches = rawSaveGameBytesFromDisk[lbTorches]
	g.PartyState.Inventory.Provisions.Keys = rawSaveGameBytesFromDisk[lbKeys]
	g.PartyState.Inventory.Provisions.SkullKeys = rawSaveGameBytesFromDisk[lbSkullKeys]

	g.MapState.LayeredMaps = *map_state.NewLayeredMaps(gameRefs.TileReferences,
		gameRefs.OverworldLargeMapReference,
		gameRefs.UnderworldLargeMapReference,
		g.MapState.XTilesVisibleOnGameScreen,
		g.MapState.YTilesVisibleOnGameScreen)

	g.LargeMapNPCAIController = make(map[references.World]*ai.NPCAIControllerLargeMap)
	overworldNPCAIInput := ai.NewNPCAIControllerLargeMapInput{
		World:           references.OVERWORLD,
		TileRefs:        gameRefs.TileReferences,
		MapState:        &g.MapState,
		DebugOptions:    &g.DebugOptions,
		TheOdds:         &g.TheOdds,
		EnemyReferences: g.GameReferences.EnemyReferences,
		DateTime:        &g.DateTime,
	}
	g.LargeMapNPCAIController[references.OVERWORLD] = ai.NewNPCAIControllerLargeMap(overworldNPCAIInput)
	// TODO: load the npcs from the save game
	underworldNPCAIInput := ai.NewNPCAIControllerLargeMapInput{
		World:           references.OVERWORLD,
		TileRefs:        gameRefs.TileReferences,
		MapState:        &g.MapState,
		DebugOptions:    &g.DebugOptions,
		TheOdds:         &g.TheOdds,
		EnemyReferences: g.GameReferences.EnemyReferences,
		DateTime:        &g.DateTime,
	}
	g.LargeMapNPCAIController[references.UNDERWORLD] = ai.NewNPCAIControllerLargeMap(underworldNPCAIInput)
	// TODO: load the npcs from the save game

	// if we are on a large map, we set this right away.
	if g.MapState.PlayerLocation.Location.GetMapType() == references.LargeMapType {
		g.CurrentNPCAIController = g.GetCurrentLargeMapNPCAIController()
	}

	g.ItemStacksMap = *references.NewItemStacksMap()

	g.TheOdds = references.NewDefaultTheOdds()
	g.DebugOptions = references.DebugOptions{
		FreeMove:   false,
		MonsterGen: true,
	}

	// g.Lighting = map_state.NewLighting(g)

	return nil
}

func getBytesAsUint16(data0 byte, data1 byte) uint16 {
	res := uint16(data0) | (uint16(data1) << 8)
	return res
}

func getUint16(bytes *[]byte, address uint16) uint16 {
	return getBytesAsUint16((*bytes)[address], (*bytes)[address+1])
}
