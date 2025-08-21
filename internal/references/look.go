package references

import (
	"log"
	"os"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

const totalLooks = 512

type LookReferences struct {
	gameConfig  *config.UltimaVConfiguration `json:"game_config" yaml:"game_config"`
	lookOffsets [totalLooks]int              `json:"look_offsets" yaml:"look_offsets"`
	lookData    []byte                       `json:"look_data" yaml:"look_data"`
}

func NewLookReferences(gameConfig *config.UltimaVConfiguration) *LookReferences {
	lookRefs := &LookReferences{}
	lookRefs.gameConfig = gameConfig

	var err error
	lookRefs.lookData, err = lookRefs.getLookFileAsBytes()
	if err != nil {
		log.Fatal("can't read look file")
	}

	count := 0
	for i := 0; i < totalLooks*2; i += 2 {
		lookRefs.lookOffsets[count] = int((lookRefs.lookData)[i]) | int((lookRefs.lookData)[i+1])<<8
		count++
	}

	return lookRefs
}

func (l *LookReferences) getLookFileAsBytes() ([]byte, error) {
	lookData, err := os.ReadFile(l.gameConfig.GetLookDataFilePath())
	if err != nil {
		return nil, err
	}
	return lookData, nil
}

func (l *LookReferences) GetTileLookDescription(tileIndex indexes.SpriteIndex) string {
	lookBytes := make([]byte, 0)

	for i := l.lookOffsets[tileIndex]; i < len(l.lookData)-1; i++ {
		if l.lookData[i] != 0 {
			lookBytes = append(lookBytes, l.lookData[i])
			continue
		}
		break
	}

	return string(lookBytes)
}
