package game_state

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"os"
	"unsafe"
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

func (g *GameState) LoadLegacySaveGame(savedGamFilePath string) error {
	// Open the file in read-only mode and as binary
	rawSaveGameBytesFromDisk, err := g.getLegacySavedGamRaw(savedGamFilePath)
	if err != nil {
		return err
	}

	//var saveGame = GameState{}
	g.RawSave = [savedGamFileSize]byte(rawSaveGameBytesFromDisk)

	// Overlay player characters over memory rawSaveGameBytesFromDisk to easily consume data
	const lCharacters = 0x02
	characterPtr := (*[NPlayers]PlayerCharacter)(unsafe.Pointer(&g.RawSave[lCharacters]))
	g.Characters = *characterPtr

	// world and position
	const lbLocation = 0x2ED
	const lbX = 0x2F0
	const lbY = 0x2F1
	const lbFloor = 0x2EF
	g.Location = references.Location(rawSaveGameBytesFromDisk[lbLocation])
	g.Position = references.Position{X: int16(rawSaveGameBytesFromDisk[lbX]), Y: int16(rawSaveGameBytesFromDisk[lbY])}
	//g.Position = ultimav.Position{X: 81, Y: 106}
	g.Floor = int8(rawSaveGameBytesFromDisk[lbFloor])

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
	g.Karma = rawSaveGameBytesFromDisk[lbKarma]
	g.QtyGold = getUint16(&rawSaveGameBytesFromDisk, lsGold)

	// Provisions
	const lsFood = 0x202
	const lbKeys = 0x206
	const lbGems = 0x207
	const lbTorches = 0x208
	const lbSkullKeys = 0x20B
	g.Provisions.QtyFood = getUint16(&rawSaveGameBytesFromDisk, lsFood)
	g.Provisions.QtyGems = rawSaveGameBytesFromDisk[lbGems]
	g.Provisions.QtyTorches = rawSaveGameBytesFromDisk[lbTorches]
	g.Provisions.QtyKeys = rawSaveGameBytesFromDisk[lbKeys]
	g.Provisions.QtySkullKeys = rawSaveGameBytesFromDisk[lbSkullKeys]

	g.LayeredMaps = *NewLayeredMaps()

	return nil
}

func getBytesAsUint16(data0 byte, data1 byte) uint16 {
	res := uint16(data0) | (uint16(data1) << 8)
	return res
}

func getUint16(bytes *[]byte, address uint16) uint16 {
	return getBytesAsUint16((*bytes)[address], (*bytes)[address+1])
}
