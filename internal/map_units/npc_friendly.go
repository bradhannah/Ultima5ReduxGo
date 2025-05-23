package map_units

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// NPCFriendly is a friendly NPC that is not a player character
// It is used to represent NPCS who are either inanimate objects (ex. carpet) or non-combative NPCS
type NPCFriendly struct {
	NPCReference   references.NPCReference
	mapUnitDetails MapUnitDetails

	vehicleDetails VehicleDetails
}

func NewNPCFriendly(npcReference references.NPCReference, npcNum int) *NPCFriendly {
	var friendly NPCFriendly
	// friendly := NPCFriendly{}
	friendly.NPCReference = npcReference
	friendly.mapUnitDetails.NPCNum = npcNum

	// friendly.mapUnitDetails.AStarMap = map_state.NewAStarMap()
	friendly.mapUnitDetails.CurrentPath = nil

	if !friendly.IsEmptyMapUnit() {
		friendly.mapUnitDetails.Visible = true
	}

	return &friendly
}

func NewNPCFriendlyVehicle(vehicleType references.VehicleType, npcRef references.NPCReference) *NPCFriendly {
	friendly := NewNPCFriendly(npcRef, int(npcRef.DialogNumber))

	friendly.vehicleDetails = NewVehicleDetails(vehicleType)
	friendly.SetFloor(references.FloorNumber(npcRef.Schedule.Floor[0]))
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

func (friendly *NPCFriendly) SetDirectionBasedOnNewPos(position references.Position) {
	if friendly.GetVehicleDetails().VehicleType != references.NoPartyVehicle {
		// we change the direction of the vehicle based on the dominate direction
		friendly.GetVehicleDetails().SetPartyVehicleDirection(
			friendly.MapUnitDetails().Position.GetDominateDirection(position))
	}
}

func (friendly *NPCFriendly) SetPos(position references.Position) {
	friendly.SetDirectionBasedOnNewPos(position)
	friendly.mapUnitDetails.Position = position
}

func (friendly *NPCFriendly) SetFloor(floor references.FloorNumber) {
	friendly.mapUnitDetails.Floor = floor
}

func (friendly *NPCFriendly) GetVehicleDetails() *VehicleDetails {
	if friendly == nil {
		_ = "a"
	}
	if friendly.NPCReference.GetNPCType() == references.Vehicle {
		return &friendly.vehicleDetails
	}
	// log.Fatal("Wrong type, not a vehicle")
	return &VehicleDetails{}
}
