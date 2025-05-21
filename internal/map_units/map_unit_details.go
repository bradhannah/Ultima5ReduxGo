package map_units

import (
	"log"

	// "github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type MapUnitType int

const (
	NonPlayerCharacter MapUnitType = iota
	Enemy
	Vehicle
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

	// AStarMap *map_state.AStarMap

	CurrentPath []references.Position
}

func (mu *MapUnitDetails) DequeueNextPosition() references.Position {
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
func (mu *MapUnitDetails) SetCurrentPath(path []references.Position) {
	mu.CurrentPath = path
}
