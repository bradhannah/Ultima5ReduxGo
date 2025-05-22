package map_units

import (
	"log"

	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type MapUnitType int

const (
	NonPlayerCharacter MapUnitType = iota
	Enemy
	Vehicle
)

type MapUnit interface {
	GetMapUnitType() MapUnitType

	Pos() references2.Position
	PosPtr() *references2.Position
	Floor() references2.FloorNumber

	SetPos(position references2.Position)
	SetFloor(floor references2.FloorNumber)

	MapUnitDetails() *MapUnitDetails
	SetVisible(visible bool)
	IsVisible() bool
	IsEmptyMapUnit() bool
}

type MapUnitDetails struct {
	Position references2.Position
	Floor    references2.FloorNumber
	AiType   references2.AiType

	Visible bool

	NPCNum int

	// AStarMap *map_state.AStarMap

	CurrentPath []references2.Position
}

func (mu *MapUnitDetails) DequeueNextPosition() references2.Position {
	if !mu.HasAPathAlreadyCalculated() {
		log.Fatal("NPC has no path calculated")
	}
	pos := mu.CurrentPath[0]
	mu.CurrentPath = mu.CurrentPath[1:]
	return pos
}

func (mu *MapUnitDetails) HasAPathAlreadyCalculated() bool {
	return mu.CurrentPath != nil && len(mu.CurrentPath) > 0
}
func (mu *MapUnitDetails) SetCurrentPath(path []references2.Position) {
	mu.CurrentPath = path
}
