package game_state

import (
	"fmt"
	"os"

	"github.com/bradhannah/Ultima5ReduxGo/internal/ai"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	references "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

const savedGamFileSize = 4192

type (
	StartingMemoryAddressUb  uint16
	StartingMemoryAddressU16 uint16
)

func getLegacySavedGamRaw(savedGamFilePath string) ([]byte, error) {
	// Open the file in read-only mode and as binary
	file, err := os.OpenFile(savedGamFilePath, os.O_RDONLY, 0o666)
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

func (g *GameState) LoadLegacySaveGameFromBytes(rawSaveData []byte) error {

	// var saveGame = GameState{}
	g.RawSave = [savedGamFileSize]byte(rawSaveData)

	// Overlay player characters over memory rawSaveData to easily consume data
	g.PartyState.LoadFromRaw(g.RawSave)

	// world and position
	const lbLocation = 0x2ED
	const lbX = 0x2F0
	const lbY = 0x2F1
	const lbFloor = 0x2EF
	g.MapState.PlayerLocation.Location = references.Location(rawSaveData[lbLocation])
	g.MapState.PlayerLocation.Position = references.Position{X: references.Coordinate(rawSaveData[lbX]), Y: references.Coordinate(rawSaveData[lbY])}
	g.MapState.PlayerLocation.Floor = references.FloorNumber(rawSaveData[lbFloor])

	// Date/Time
	const lsYear = 0x2CE
	const lbMonth = 0x2D7
	const lbDay = 0x2D8
	const lbHour = 0x2D9
	const lbMinute = 0x2DB
	g.DateTime.Year = getUint16(&rawSaveData, lsYear)
	g.DateTime.Month = rawSaveData[lbMonth]
	g.DateTime.Day = rawSaveData[lbDay]
	g.DateTime.Hour = rawSaveData[lbHour]
	g.DateTime.Minute = rawSaveData[lbMinute]

	// Various Things
	const lbKarma = 0x2e2
	const lsGold = 0x204
	g.PartyState.Karma = party_state.Karma(rawSaveData[lbKarma])
	g.PartyState.Inventory.Gold = getUint16(&rawSaveData, lsGold)

	// ProvisionsQuantity
	const lsFood = 0x202
	const lbKeys = 0x206
	const lbGems = 0x207
	const lbTorches = 0x208
	const lbSkullKeys = 0x20B
	g.PartyState.Inventory.Provisions.Food = getUint16(&rawSaveData, lsFood)
	g.PartyState.Inventory.Provisions.Gems = rawSaveData[lbGems]
	g.PartyState.Inventory.Provisions.Torches = rawSaveData[lbTorches]
	g.PartyState.Inventory.Provisions.Keys = rawSaveData[lbKeys]
	g.PartyState.Inventory.Provisions.SkullKeys = rawSaveData[lbSkullKeys]

	g.MapState.LayeredMaps = *map_state.NewLayeredMaps(g.GameReferences.TileReferences,
		g.GameReferences.OverworldLargeMapReference,
		g.GameReferences.UnderworldLargeMapReference,
		g.MapState.XTilesVisibleOnGameScreen,
		g.MapState.YTilesVisibleOnGameScreen)

	g.LargeMapNPCAIController = make(map[references.World]*ai.NPCAIControllerLargeMap)
	overworldNPCAIInput := ai.NewNPCAIControllerLargeMapInput{
		World:           references.OVERWORLD,
		TileRefs:        g.GameReferences.TileReferences,
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
		TileRefs:        g.GameReferences.TileReferences,
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

//func (g *GameState) LoadLegacySaveGame(savedGamFilePath string) error {
//	// Open the file in read-only mode and as binary
//	rawSaveGameBytesFromDisk, err := g.getLegacySavedGamRaw(savedGamFilePath)
//	if err != nil {
//		return err
//	}
//	return g.LoadLegacySaveGameFromBytes(rawSaveGameBytesFromDisk)
//}

func getBytesAsUint16(data0, data1 byte) uint16 {
	res := uint16(data0) | (uint16(data1) << 8)
	return res
}

func getUint16(bytes *[]byte, address uint16) uint16 {
	return getBytesAsUint16((*bytes)[address], (*bytes)[address+1])
}
