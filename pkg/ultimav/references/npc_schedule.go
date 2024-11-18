package references

import (
	"log"
	"unsafe"
)

const (
	totalScheduleItemsPerNpc = 3

	sizeOfNPCSchedule = int(unsafe.Sizeof(NPCSchedule{}))
)

type NPCSchedule struct {
	Ai    [totalScheduleItemsPerNpc]byte
	X     [totalScheduleItemsPerNpc]byte
	Y     [totalScheduleItemsPerNpc]byte
	Floor [totalScheduleItemsPerNpc]byte
	Time  [totalScheduleItemsPerNpc + 1]byte
}

func CreateNPCSchedule(rawData []byte) NPCSchedule {
	if len(rawData) < sizeOfNPCSchedule {
		log.Fatal("Not enough data to create NPCReference schedule")
	}
	npcSchedule := (*NPCSchedule)(unsafe.Pointer(&rawData[0]))
	return *npcSchedule
}

// struct NPC_File {
//  NPC_Info info[8]; // each NPCReference file has information for 8 maps
// };
//
// For each city, we have an information entry for the Npcs of the map:
//
// struct NPC_Info {
//  NPC_Schedule schedule[32];
//  uint8 type[32]; // merchant, guard, etc.
//  uint8 dialog_number[32];
// };
// The dialog number gives the entry index to the *.TLK file.
//
// Finally, the schedule says how the Npc moves around in the city and especially when:
//
// struct NPC_Schedule {
//  uint8 AI_types[3];
//  uint8 x_coordinates[3];
//  uint8 y_coordinates[3];
//  sint8 z_coordinates[3];
//  uint8 times[4];
// };
