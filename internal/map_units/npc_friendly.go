package map_units

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type NPCFriendly struct {
	NPCReference   references2.NPCReference
	mapUnitDetails MapUnitDetails

	vehicleDetails VehicleDetails
}

func NewNPCFriendly(npcReference references2.NPCReference, npcNum int) *NPCFriendly {
	friendly := NPCFriendly{}
	friendly.NPCReference = npcReference
	friendly.mapUnitDetails.NPCNum = npcNum

	// friendly.mapUnitDetails.AStarMap = map_state.NewAStarMap()
	friendly.mapUnitDetails.CurrentPath = nil

	if !friendly.IsEmptyMapUnit() {
		friendly.mapUnitDetails.Visible = true
	}

	return &friendly
}

func NewNPCFriendlyVehicle(vehicleType references2.VehicleType, npcRef references2.NPCReference) *NPCFriendly {
	// npcReference := references.NewNPCReferenceForVehicle(vehicleType, references.Position{X: 15, Y: 15}, 0)
	friendly := NewNPCFriendly(npcRef, int(npcRef.DialogNumber))

	friendly.vehicleDetails = VehicleDetails{
		currentDirection:  references2.NoneDirection,
		previousDirection: references2.NoneDirection,
		VehicleType:       vehicleType,
		SkiffQuantity:     0,
	}
	friendly.SetFloor(references2.FloorNumber(npcRef.Schedule.Floor[0]))
	friendly.SetVisible(true)

	return friendly
}

func NewNPCFriendlyVehiceNewRef(vehicletype references2.VehicleType, pos references2.Position, floor references2.FloorNumber) *NPCFriendly {
	npcRef := references2.NewNPCReferenceForVehicle(vehicletype, pos, floor)
	return NewNPCFriendlyVehicle(vehicletype, *npcRef)
}

func NewNPCFriendlyVehiceNoVehicle() NPCFriendly {
	return *NewNPCFriendlyVehiceNewRef(references2.NoPartyVehicle, references2.Position{X: 0, Y: 0}, 0)
}

func (friendly *NPCFriendly) IsEmptyMapUnit() bool {
	return friendly.NPCReference.GetNPCType() == 0 || (friendly.NPCReference.Schedule.X[0] == 0 && friendly.NPCReference.Schedule.Y[0] == 0)
}

func (friendly *NPCFriendly) GetMapUnitType() MapUnitType {
	return NonPlayerCharacter
}

func (friendly *NPCFriendly) Pos() references2.Position {
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

func (friendly *NPCFriendly) Floor() references2.FloorNumber {
	return friendly.mapUnitDetails.Floor
}

func (friendly *NPCFriendly) PosPtr() *references2.Position {
	return &friendly.mapUnitDetails.Position
}

func (friendly *NPCFriendly) SetIndividualNPCBehaviour(indiv references2.IndividualNPCBehaviour) {
	friendly.mapUnitDetails.Position = indiv.Position
	friendly.mapUnitDetails.Floor = indiv.Floor
	friendly.mapUnitDetails.AiType = indiv.Ai
}

func (friendly *NPCFriendly) SetPos(position references2.Position) {
	friendly.mapUnitDetails.Position = position
}

func (friendly *NPCFriendly) SetFloor(floor references2.FloorNumber) {
	friendly.mapUnitDetails.Floor = floor
}

func (friendly *NPCFriendly) GetVehicleDetails() *VehicleDetails {
	if friendly == nil {
		_ = "a"
	}
	if friendly.NPCReference.GetNPCType() == references2.Vehicle {
		return &friendly.vehicleDetails
	}
	// log.Fatal("Wrong type, not a vehicle")
	return &VehicleDetails{}
}
