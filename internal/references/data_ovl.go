package references

import (
	"fmt"
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
)

//_compressedWords = dataRef.GetDataChunk(DataOvlReference.DataChunkName.TALK_COMPRESSED_WORDS).GetChunkAsStringList();
//_dataChunks.AddDataChunk(DataChunk.DataFormatType.StringList,
//	"Compressed words used in the conversation files", 0x104c, 0x24e, 0,
//	DataChunkName.TALK_COMPRESSED_WORDS);

const (
	nLocationNameOffset  = 0xa4d
	nLocationNameOffset2 = 0xac1
)
const (
	nTalkCompressedWordsOffset = 0x104c
	nTalkCompressedWordsLength = 0x24e
)

type DataOvl struct {
	LocationNames   []string
	CompressedWords []string
}

func readNullTerminatedStrings(data *[]byte, offset, n int) ([]string, error) {
	var result []string
	start := offset
	for i := 0; i < n; i++ {
		// Find the null terminator (0 byte) to end the current string
		end := start
		for end < len(*data) && (*data)[end] != 0 {
			end++
		}

		// Check if we have reached the end of data without finding enough null terminators
		if end == len(*data) {
			return nil, fmt.Errorf("unexpected end of data before reading %d strings", n)
		}

		// Extract the string
		result = append(result, string((*data)[start:end]))

		// Move the start position to the byte after the null terminator
		start = end + 1
	}

	return result, nil
}

func NewDataOvl(config *config.UltimaVConfiguration) *DataOvl {
	dataOvl := DataOvl{}

	var err error
	dataOvl.LocationNames, err = readNullTerminatedStrings(&config.RawDataOvl, nLocationNameOffset, int(Iolos_Hut))
	dataOvl.LocationNames = append([]string{""}, dataOvl.LocationNames...)
	dataOvl.LocationNames = append(dataOvl.LocationNames,
		[]string{"SUTEK'S HUT", "SIN VRAAL'S HUT", "GRENDAL'S HUT", "LORD BRITISH'S CASTLE", "PALACE OF BLACKTHORN"}...)
	secondHalf, _ := readNullTerminatedStrings(&config.RawDataOvl, nLocationNameOffset2, 27-int(Iolos_Hut))
	dataOvl.LocationNames = append(dataOvl.LocationNames, secondHalf...)

	dataOvl.CompressedWords, err = readNullTerminatedStrings(
		&config.RawDataOvl, nTalkCompressedWordsOffset, nTalkCompressedWordsLength)

	if err != nil {
		log.Fatalf("error reading compressed words: %v", err)
	}

	if err != nil {
		panic(err)
	}
	return &dataOvl
}
