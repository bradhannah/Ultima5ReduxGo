package game_state

import (
	"fmt"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type NPCFriendly struct {
	NPCReference   references.NPCReference
	mapUnitDetails MapUnitDetails

	vehicleDetails VehicleDetails
}

func NewNPCFriendly(npcReference references.NPCReference, npcNum int) *NPCFriendly {
	friendly := NPCFriendly{}
	friendly.NPCReference = npcReference
	friendly.mapUnitDetails.NPCNum = npcNum

	friendly.mapUnitDetails.AStarMap = NewAStarMap()
	friendly.mapUnitDetails.CurrentPath = nil

	if !friendly.IsEmptyMapUnit() {
		friendly.mapUnitDetails.Visible = true
	}

	return &friendly
}

func NewNPCFriendlyVehicle(vehicleType references.VehicleType, npcRef references.NPCReference) *NPCFriendly {
	//npcReference := references.NewNPCReferenceForVehicle(vehicleType, references.Position{X: 15, Y: 15}, 0)
	friendly := NewNPCFriendly(npcRef, int(npcRef.DialogNumber))

	friendly.vehicleDetails = VehicleDetails{
		currentDirection:  references.NoneDirection,
		previousDirection: references.NoneDirection,
		VehicleType:       vehicleType,
		SkiffQuantity:     0,
	}

	friendly.SetVisible(true)

	return friendly
}

func NewNPCFriendlyVehiceNewRef(vehicletype references.VehicleType, pos references.Position, floor references.FloorNumber) *NPCFriendly {
	npcRef := references.NewNPCReferenceForVehicle(vehicletype, pos, floor)
	return NewNPCFriendlyVehicle(vehicletype, *npcRef)
}

func NewNPCFriendlyVehiceNoVehicle() NPCFriendly {
	return *NewNPCFriendlyVehiceNewRef(references.NoPartyVehicle, references.Position{X: 0, Y: 0}, 0)
}

func (friendly *NPCFriendly) IsEmptyMapUnit() bool {
	return friendly.NPCReference.GetNPCType() == 0 || (friendly.NPCReference.Schedule.X[0] == 0 && friendly.NPCReference.Schedule.Y[0] == 0)
}

func (friendly *NPCFriendly) GetMapUnitType() MapUnitType {
	return NonPlayerCharacter
}

func (friendly *NPCFriendly) Pos() references.Position {
	return friendly.mapUnitDetails.Position
}

func (friendly *NPCFriendly) MapUnitDetails() *MapUnitDetails {
	return &friendly.mapUnitDetails
}

func (friendly *NPCFriendly) SetVisible(visible bool) {
	friendly.mapUnitDetails.Visible = visible
}
func (friendly *NPCFriendly) IsVisible() bool {
	return friendly.mapUnitDetails.Visible
}

func (friendly *NPCFriendly) Floor() references.FloorNumber {
	return friendly.mapUnitDetails.Floor
}

func (friendly *NPCFriendly) PosPtr() *references.Position {
	return &friendly.mapUnitDetails.Position
}

func (friendly *NPCFriendly) SetIndividualNPCBehaviour(indiv references.IndividualNPCBehaviour) {
	friendly.mapUnitDetails.Position = indiv.Position
	friendly.mapUnitDetails.Floor = indiv.Floor
	friendly.mapUnitDetails.AiType = indiv.Ai
}

func (friendly *NPCFriendly) SetPos(position references.Position) {
	friendly.mapUnitDetails.Position = position
}

func (friendly *NPCFriendly) SetFloor(floor references.FloorNumber) {
	friendly.mapUnitDetails.Floor = floor
}

func (friendly *NPCFriendly) GetVehicleDetails() *VehicleDetails {
	if friendly == nil {
		fmt.Sprint("oof")
	}
	if friendly.NPCReference.GetNPCType() == references.Vehicle {
		return &friendly.vehicleDetails
	}
	//log.Fatal("Wrong type, not a vehicle")
	return &VehicleDetails{}
}
