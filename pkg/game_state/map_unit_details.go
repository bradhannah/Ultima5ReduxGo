package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type MapUnitType int

const (
	NonPlayerCharacter MapUnitType = iota
	Enemy
)

type MapUnit interface {
	GetMapUnitType() MapUnitType

	Pos() references.Position
	PosPtr() *references.Position
	Floor() references.FloorNumber

	SetPos(position references.Position)
	SetFloor(floor references.FloorNumber)

	MapUnitDetails() *MapUnitDetails
	SetVisible(visible bool)
	IsVisible() bool
	IsEmptyMapUnit() bool
}

type MapUnitDetails struct {
	Position references.Position
	Floor    references.FloorNumber
	AiType   references.AiType

	Visible bool

	NPCNum int

	AStarMap *AStarMap

	CurrentPath *[]references.Position
}

func (npc *MapUnitDetails) DequeueNextPosition() references.Position {
	if !npc.HasAPathAlreadyCalculated() {
		log.Fatal("NPC has no path calculated")
	}
	pos := (*npc.CurrentPath)[0]
	*npc.CurrentPath = (*npc.CurrentPath)[1:]
	return pos
}

func (npc *MapUnitDetails) HasAPathAlreadyCalculated() bool {
	return npc.CurrentPath != nil && len(*npc.CurrentPath) > 0
}
