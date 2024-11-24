package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"

type FloorNumber int8

type LadderOrStairType int

const (
	NotLadderOrStair LadderOrStairType = iota
	LadderOrStairUp
	LadderOrStairDown
)

func GetLadderOfStairsType(currentFloor FloorNumber, targetFloor FloorNumber) LadderOrStairType {
	if currentFloor == targetFloor {
		return NotLadderOrStair
	}
	if targetFloor > currentFloor {
		return LadderOrStairUp
	}
	return LadderOrStairDown
}

func IsSpecificLadderOrStairs(s indexes.SpriteIndex, ladderOrStairType LadderOrStairType) bool {
	switch ladderOrStairType {
	case LadderOrStairUp:
		return s == indexes.LadderUp || s == indexes.StairsUp1 || s == indexes.StairsUp2
	case LadderOrStairDown:
		return s == indexes.LadderDown || s == indexes.StairsDown1 || s == indexes.StairsDown2
	}
	return false
}
