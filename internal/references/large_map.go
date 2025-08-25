package references

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/files"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

const (
	// YLargeMapTiles Total map height in tiles
	YLargeMapTiles = Coordinate(256)
	// XLargeMapTiles Total map width in tiles
	XLargeMapTiles = Coordinate(256)
)

const (
	totalChunks           = 256    // Number of chunks in the map
	tilesPerChunkX        = 16     // Tiles per chunk on the x-axis
	tilesPerChunkY        = 16     // Tiles per chunk on the y-axis
	datOverlayBritMap int = 0x3886 // Offset for overlay data (adjust as needed)
)

// LargeMapReference A reference to a large map. This is a 256x256 map with 16x16 tiles.
// The map is stored in a 256x256 array of bytes, with each byte representing a tile.
type LargeMapReference struct {
	rawData [XLargeMapTiles][YLargeMapTiles]byte
}

// World represents the world of the map. Overworld or Underworld.
type World int

const (
	OVERWORLD World = iota
	UNDERWORLD
)

func NewLargeMapReference(gameConfig *config.UltimaVConfiguration, world World) (*LargeMapReference, error) {
	if world == OVERWORLD {
		return loadLargeMapFromFile(OVERWORLD, gameConfig)
	} else {
		return loadLargeMapFromFile(UNDERWORLD, gameConfig)
	}
}

func (m *LargeMapReference) GetSpriteIndex(x, y Coordinate) indexes.SpriteIndex {
	pos := Position{X: x, Y: y}
	wrappedPos := pos.GetWrapped(XLargeMapTiles, YLargeMapTiles)

	// Use deterministic animation function with fixed time (0) for map initialization
	// This ensures consistent map setup regardless of when initialization occurs
	return sprites.GetSpriteIndexWithAnimationBySpriteIndexTick(indexes.SpriteIndex(m.rawData[wrappedPos.X][wrappedPos.Y]), 0, 0)
}

func loadLargeMapFromFile(world World, gameConfig *config.UltimaVConfiguration) (*LargeMapReference, error) {
	mapFileAndPath, dataOvlFileAndPath := path.Join(gameConfig.SavedConfigData.DataFilePath, files.BRIT_DAT),
		path.Join(gameConfig.SavedConfigData.DataFilePath, files.DATA_OVL)
	if world == UNDERWORLD {
		mapFileAndPath = path.Join(gameConfig.SavedConfigData.DataFilePath, files.UNDER_DAT)
		dataOvlFileAndPath = ""
	}

	theChunksSerial, err := os.ReadFile(mapFileAndPath)
	if err != nil {
		return nil, err
	}
	// defer file.Close() // Ensure the file is closed when done

	// Create a buffer for data overlay chunks
	dataOvlChunks := make([]byte, totalChunks)

	var dataOvl *os.File

	ignoreOverlay := dataOvlFileAndPath == ""
	if !ignoreOverlay {
		// Open the overlay file for reading (case-insensitive path)
		var err error
		dataOvl, err = os.Open(dataOvlFileAndPath)
		if err != nil {
			panic(fmt.Sprintf("failed to open overlay file: %v", err))
		}
		defer dataOvl.Close()

		// Seek to the overlay data offset
		_, err = dataOvl.Seek(int64(datOverlayBritMap), os.SEEK_SET)
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
	for chunk := 0; chunk < totalChunks; chunk++ {
		// Calculate the chunk column and row
		col := chunk % tilesPerChunkX
		row := chunk / tilesPerChunkY

		// Read overlay chunk value (or set to 0 if ignoring overlay)
		if ignoreOverlay {
			dataOvlChunks[chunkCount] = 0x00
		} else {
			binary.Read(dataOvl, binary.LittleEndian, &dataOvlChunks[chunkCount])
		}

		// Process each row (horizon) in the chunk
		for curRow := row * tilesPerChunkY; curRow < (row*tilesPerChunkY)+tilesPerChunkY; curRow++ {
			for curCol := col * tilesPerChunkX; curCol < (col*tilesPerChunkX)+tilesPerChunkX; curCol++ {
				// Check if it's a water tile (0xFF in overlay means water)
				const waterTile = 0xFF
				if dataOvlChunks[chunkCount] == waterTile {
					mapRef.rawData[curCol][curRow] = 0x01 // Water1 tile
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
