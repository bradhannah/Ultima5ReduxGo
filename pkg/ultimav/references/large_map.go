package references

import (
	"encoding/binary"
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
	"os"
	"path"
)

const (
	YTiles                = 256    // Total map height in tiles
	XTiles                = 256    // Total map width in tiles
	TotalChunks           = 256    // Number of chunks in the map
	TilesPerChunkX        = 16     // Tiles per chunk on the x-axis
	TilesPerChunkY        = 16     // Tiles per chunk on the y-axis
	DatOverlayBritMap int = 0x3886 // Offset for overlay data (adjust as needed)
)

type LargeMapReference struct {
	rawData [XTiles][YTiles]byte
}

type World int

const (
	OVERWORLD World = iota
	UNDERWORLD
)

func NewLargeMapReference(gameConfig *config.UltimaVConfiguration, world World) (*LargeMapReference, error) {
	if world == OVERWORLD {
		return loadLargeMapFromFile(path.Join(gameConfig.DataFilePath, BRIT_DAT), path.Join(gameConfig.DataFilePath, DATA_OVL))
	} else {
		return loadLargeMapFromFile(path.Join(gameConfig.DataFilePath, UNDER_DAT), "")
	}

}

func loadLargeMapFromFile(mapFileAndPath string, dataOvlFileAndPath string) (*LargeMapReference, error) {
	ignoreOverlay := dataOvlFileAndPath == ""

	theChunksSerial, err := os.ReadFile(mapFileAndPath)
	if err != nil {
		return nil, err
	}
	//defer file.Close() // Ensure the file is closed when done

	// Create a buffer for data overlay chunks
	dataOvlChunks := make([]byte, TotalChunks)

	var dataOvl *os.File
	if !ignoreOverlay {
		// Open the overlay file for reading (case-insensitive path)
		//filePath := GetFirstFileAndPathCaseInsensitive(overlayFilename)
		var err error
		dataOvl, err = os.Open(dataOvlFileAndPath)
		if err != nil {
			panic(fmt.Sprintf("failed to open overlay file: %v", err))
		}
		defer dataOvl.Close()

		// Seek to the overlay data offset
		_, err = dataOvl.Seek(int64(DatOverlayBritMap), os.SEEK_SET)
		if err != nil {
			return nil, err
		}
	}

	mapRef := LargeMapReference{}

	// Counter for serial chunks from brit.dat
	britDatChunkCount := 0

	// Chunk count
	chunkCount := 0

	// Process each chunk
	for chunk := 0; chunk < TotalChunks; chunk++ {
		// Calculate the chunk column and row
		col := chunk % TilesPerChunkX
		row := chunk / TilesPerChunkY

		// Read overlay chunk value (or set to 0 if ignoring overlay)
		if ignoreOverlay {
			dataOvlChunks[chunkCount] = 0x00
		} else {
			binary.Read(dataOvl, binary.LittleEndian, &dataOvlChunks[chunkCount])
		}

		// Process each row (horizon) in the chunk
		for curRow := row * TilesPerChunkY; curRow < (row*TilesPerChunkY)+TilesPerChunkY; curRow++ {
			for curCol := col * TilesPerChunkX; curCol < (col*TilesPerChunkX)+TilesPerChunkX; curCol++ {
				// Check if it's a water tile (0xFF in overlay means water)
				if dataOvlChunks[chunkCount] == 0xFF {
					mapRef.rawData[curCol][curRow] = 0x01 // Water tile
				} else {
					// Land tile, read from brit.dat
					mapRef.rawData[curCol][curRow] = theChunksSerial[britDatChunkCount]
					britDatChunkCount++
				}
			}
		}

		// Increment the chunk count
		chunkCount++
	}
	return &mapRef, nil
}
