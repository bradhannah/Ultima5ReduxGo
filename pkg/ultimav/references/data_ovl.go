package references

import (
	"fmt"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
)

const nLocationNameOffset = 0xa4d
const nSecondOffset = 0xac1

type DataOvl struct {
	LocationNames []string
}

func readNullTerminatedStrings(data *[]byte, offset int, n int) ([]string, error) {
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
	secondHalf, _ := readNullTerminatedStrings(&config.RawDataOvl, nSecondOffset, 27-int(Iolos_Hut))
	dataOvl.LocationNames = append(dataOvl.LocationNames, secondHalf...)
	if err != nil {
		panic(err)

	}
	return &dataOvl
}
