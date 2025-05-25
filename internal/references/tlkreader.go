package references

import (
	"encoding/binary"
	"fmt"
	"os"
)

type TalkBytesForSmallMapType map[int][]byte

// NpcTalkOffset is exactly four little‑endian bytes in the TLK header table.
type NpcTalkOffset struct {
	NpcIndex   uint16 // which NPC this chunk belongs to
	FileOffset uint16 // offset, in bytes, from the start of the file
}

// LoadFile reads an entire *.tlk file from disk and hands you a
// map[npcIndex] = rawDialogueBytes.
//
// Callers who already have the file in memory can jump straight to
// Load(data []byte).
func LoadFile(path string) (map[int][]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Load(data)
}

// Load parses a TLK buffer already resident in memory.
func Load(data []byte) (TalkBytesForSmallMapType, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("file too small: need at least 2 bytes")
	}

	entryCount := binary.LittleEndian.Uint16(data[:2])
	const entrySize = 4 // sizeof(NpcTalkOffset)
	tableBytes := int(entryCount) * entrySize

	if len(data) < 2+tableBytes {
		return nil, fmt.Errorf("truncated offset table: need %d bytes, have %d",
			2+tableBytes, len(data))
	}

	// --- read the offset table -------------------------------------------
	offsets := make([]NpcTalkOffset, entryCount)
	for i := 0; i < int(entryCount); i++ {
		base := 2 + i*entrySize

		offsets[i].NpcIndex = binary.LittleEndian.Uint16(data[base : base+2])
		offsets[i].FileOffset = binary.LittleEndian.Uint16(data[base+2 : base+4])
	}
	// ---------------------------------------------------------------------

	// --- slice the file into per‑NPC chunks ------------------------------
	result := make(TalkBytesForSmallMapType, entryCount)
	for i, off := range offsets {
		start := int(off.FileOffset)

		// final chunk runs to EOF, otherwise to next offset
		end := len(data)
		if i+1 < len(offsets) {
			end = int(offsets[i+1].FileOffset)
		}

		if start > end || start >= len(data) {
			return nil, fmt.Errorf("offset out of bounds for NPC %d: start=%d end=%d",
				off.NpcIndex, start, end)
		}
		// *share* underlying storage – no copy.
		result[int(off.NpcIndex)] = data[start:end]
	}
	// ---------------------------------------------------------------------
	return result, nil
}
