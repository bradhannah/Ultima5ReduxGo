package references

import (
	"log"
	"os"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

const totalLooks = 512

type LookReferences struct {
	GameConfig       *config.UltimaVConfiguration `json:"game_config" yaml:"game_config"`
	LookOffsets      [totalLooks]int              `json:"look_offsets" yaml:"look_offsets"`
	LookData         []byte                       `json:"look_data" yaml:"look_data"`
	TileDescriptions map[int]string               `json:"tile_descriptions" yaml:"tile_descriptions"`
}

func NewLookReferences(gameConfig *config.UltimaVConfiguration) *LookReferences {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	lookRefs := &LookReferences{}
	lookRefs.GameConfig = gameConfig

	lookDataPath := gameConfig.GetLookDataFilePath()
	log.Printf("Look data file path: %s", lookDataPath)
	if _, err := os.Stat(lookDataPath); os.IsNotExist(err) {
		log.Printf("ERROR: Look data file does not exist at: %s", lookDataPath)
		return lookRefs
	}

	var err error
	lookRefs.LookData, err = lookRefs.getLookFileAsBytes()
	if err != nil {
		log.Printf("Error reading look file: %v", err)
		log.Fatal("can't read look file") // Original game data corruption
	}

	count := 0
	for i := 0; i < totalLooks*2; i += 2 {
		lookRefs.LookOffsets[count] = int((lookRefs.LookData)[i]) | int((lookRefs.LookData)[i+1])<<8
		count++
	}

	// Populate tile descriptions
	lookRefs.TileDescriptions = make(map[int]string)
	for i := 0; i < totalLooks; i++ {
		lookRefs.TileDescriptions[i] = lookRefs.GetTileLookDescription(indexes.SpriteIndex(i))
	}

	return lookRefs
}

func (l *LookReferences) getLookFileAsBytes() ([]byte, error) {
	lookData, err := os.ReadFile(l.GameConfig.GetLookDataFilePath())
	if err != nil {
		return nil, err
	}
	return lookData, nil
}

func (l *LookReferences) GetTileLookDescription(tileIndex indexes.SpriteIndex) string {
	lookBytes := make([]byte, 0)

	for i := l.LookOffsets[tileIndex]; i < len(l.LookData)-1; i++ {
		if l.LookData[i] != 0 {
			lookBytes = append(lookBytes, l.LookData[i])
			continue
		}
		break
	}

	return string(lookBytes)
}
