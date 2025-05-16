package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCVehicle struct {
	mapUnitDetails MapUnitDetails

	VehicleType   references.VehicleType
	SkiffQuantity int

	currentDirection  references.Direction
	previousDirection references.Direction
}

func NewNPCVehicle(vehicleType references.VehicleType, pos references.Position, floor references.FloorNumber) NPCVehicle {
	return NPCVehicle{
		mapUnitDetails: MapUnitDetails{
			NPCNum:      0,
			AStarMap:    NewAStarMap(),
			CurrentPath: nil,
			Position:    pos,
			Floor:       floor,
		},
		VehicleType:   vehicleType,
		SkiffQuantity: 0,
	}
}

func (n *NPCVehicle) GetMapUnitType() MapUnitType {
	return Vehicle
}

func (n *NPCVehicle) Pos() (_ references.Position) {
	return n.mapUnitDetails.Position
}

func (n *NPCVehicle) PosPtr() (_ *references.Position) {
	return &n.mapUnitDetails.Position
}

func (n *NPCVehicle) Floor() (_ references.FloorNumber) {
	return n.mapUnitDetails.Floor
}

func (n *NPCVehicle) SetPos(pos references.Position) {
	n.mapUnitDetails.Position = pos
}

func (n *NPCVehicle) SetFloor(floor references.FloorNumber) {
	n.mapUnitDetails.Floor = floor
}

func (n *NPCVehicle) MapUnitDetails() *MapUnitDetails {
	return &n.mapUnitDetails
}

func (n *NPCVehicle) SetVisible(visible bool) {
	n.mapUnitDetails.Visible = visible
}

func (n *NPCVehicle) IsVisible() (_ bool) {
	return n.mapUnitDetails.Visible
}

func (n *NPCVehicle) IsEmptyMapUnit() (_ bool) {
	return false
}

// func (n *NPCVehicle) TryToSwitchDirection(direction references.Direction) {
// 	if n.VehicleType.RequiresNewSprite(n.PreviousDirection, n.CurrentDirection) {

// 	}
// }

func (n *NPCVehicle) SetPartyVehicleDirection(direction references.Direction) {
	switch n.VehicleType {
	case references.HorseVehicle, references.CarpetVehicle:
		if direction == references.Up || direction == references.Down {
			return
		}
	case references.FrigateVehicle, references.SkiffVehicle, references.NoPartyVehicle:
	}
	n.previousDirection = n.currentDirection
	n.currentDirection = direction
}

func (n *NPCVehicle) GetSprite() indexes.SpriteIndex {
	return n.VehicleType.GetSpriteByDirection(n.previousDirection, n.currentDirection)
}

func (g *NPCVehicle) DoesMoveResultInMovement(newDirection references.Direction) bool {
	if g.VehicleType != references.FrigateVehicle {
		return true
	}
	if g.currentDirection == newDirection {
		return true
	}
	return false
}
