package references

import (
	_ "embed"
	"encoding/json"
	"log"
	"strconv"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
)

var (
	//go:embed data/TileData.json
	tileDataRaw []byte
)

type Tiles map[indexes.SpriteIndex]Tile

func (t *Tiles) GetTile(tileNum indexes.SpriteIndex) *Tile {
	tile, exists := (*t)[tileNum]
	if !exists {
		return nil // Return nil if the tileNum doesn't exist
	}
	return &tile
}

func (t *Tiles) UnmarshalJSON(data []byte) error {
	// Temporary map to unmarshal JSON with string keys
	var tempMap map[string]Tile
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	// Convert the keys from strings to integers
	tilesMap := make(Tiles)
	for key, value := range tempMap {
		intKey, err := strconv.Atoi(key) // Convert string key to int
		if err != nil {
			return err
		}
		value.Index = indexes.SpriteIndex(intKey)
		tilesMap[indexes.SpriteIndex(intKey)] = value
	}

	// Set the result to the original Tiles map
	*t = tilesMap
	return nil
}

func NewTileReferences() *Tiles {
	var tiles Tiles
	err := json.Unmarshal(tileDataRaw, &tiles)
	if err != nil {
		log.Fatalf("error unmarshaling JSON: %v", err)
	}
	return &tiles
}
