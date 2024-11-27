package references

import (
	"log"
	"unsafe"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
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

type IndividualNPCBehaviour struct {
	Ai       AiType
	Position Position
	Floor    FloorNumber
}

func CreateNPCSchedule(rawData []byte) NPCSchedule {
	if len(rawData) < sizeOfNPCSchedule {
		log.Fatal("Not enough data to create NPCReference schedule")
	}
	npcSchedule := (*NPCSchedule)(unsafe.Pointer(&rawData[0]))
	return *npcSchedule
}

func (n *NPCSchedule) GetPreviousIndividualNPCBehaviourByUltimaDate(ud datetime.UltimaDate) IndividualNPCBehaviour {
	index := n.getScheduleIndex(ud)
	if index == 0 {
		index = len(n.Ai) - 1
	} else {
		index--
	}
	return IndividualNPCBehaviour{
		Ai:       AiType(n.Ai[index]),
		Position: Position{X: Coordinate(n.X[index]), Y: Coordinate(n.Y[index])},
		Floor:    FloorNumber(n.Floor[index]),
	}
}

func (n *NPCSchedule) GetIndividualNPCBehaviourByUltimaDate(ud datetime.UltimaDate) IndividualNPCBehaviour {
	index := n.getScheduleIndex(ud)
	return IndividualNPCBehaviour{
		Ai:       AiType(n.Ai[index]),
		Position: Position{X: Coordinate(n.X[index]), Y: Coordinate(n.Y[index])},
		Floor:    FloorNumber(n.Floor[index]),
	}
}

func (n *NPCSchedule) getScheduleIndex(date datetime.UltimaDate) int {
	const totalSchedules = totalScheduleItemsPerNpc // Alias for readability
	nHour := int(date.Hour)                         // Extract the hour from UltimaDate

	// Inline function to handle specific index logic
	getIndex := func(nOrigIndex int) int {
		if nOrigIndex == totalSchedules {
			return 1
		}
		return nOrigIndex
	}

	// If all times are zero, return 0
	if n.Time[0] == 0 && n.Time[1] == 0 && n.Time[2] == 0 && n.Time[3] == 0 {
		return 0
	}

	// Check if the hour matches any of the times
	for i := 0; i < totalSchedules+1; i++ {
		if n.Time[i] == byte(nHour) {
			return getIndex(i)
		}
	}

	// Check ranges
	if nHour > int(n.Time[3]) && nHour < int(n.Time[0]) {
		return 1
	}
	if nHour > int(n.Time[0]) && nHour < int(n.Time[1]) {
		return 0
	}
	if nHour > int(n.Time[1]) && nHour < int(n.Time[2]) {
		return 1
	}
	if nHour > int(n.Time[2]) && nHour < int(n.Time[3]) {
		return 2
	}

	// Find the index of the earliest and latest times
	nEarliestTimeIndex := n.getEarliestTimeIndex()
	nIndexPreviousToEarliest := 0
	if nEarliestTimeIndex == 0 {
		nIndexPreviousToEarliest = 1
	} else {
		nIndexPreviousToEarliest = nEarliestTimeIndex - 1
	}
	nLatestTimeIndex := n.getLatestTimeIndex()

	// Handle times outside the range
	if nHour < int(n.Time[nEarliestTimeIndex]) {
		return nIndexPreviousToEarliest
	}
	if nHour > int(n.Time[nLatestTimeIndex]) {
		return getIndex(nLatestTimeIndex)
	}

	// Fallback case
	panic("getScheduleIndex fell all the way through, which doesn't make sense.")
}

func (n *NPCSchedule) getEarliestTimeIndex() int {
	nEarliest := n.Time[0]
	nEarliestIndex := 0
	for i := 1; i < len(n.Time); i++ {
		if n.Time[i] >= nEarliest {
			continue
		}

		nEarliestIndex = i
		nEarliest = n.Time[i]
	}

	return nEarliestIndex
}

func (n *NPCSchedule) getLatestTimeIndex() int {
	nLargest := n.Time[0]
	nLargestIndex := 0

	for i := 1; i < len(n.Time); i++ {
		if n.Time[i] <= nLargest {
			continue
		}

		nLargestIndex = i
		nLargest = n.Time[i]
	}

	return nLargestIndex
}
