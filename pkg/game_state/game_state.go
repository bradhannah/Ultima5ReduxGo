package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const NPlayers = 6

type InputState int

const (
	PrimaryInput InputState = iota
	OpenDirectionInput
)

type GameState struct {
	RawSave         [savedGamFileSize]byte
	Characters      [NPlayers]PlayerCharacter
	MoonstoneStatus MoonstoneStatus

	SecondaryKeyState InputState

	Location references.Location
	Position references.Position
	Floor    int8

	LayeredMaps  LayeredMaps
	PartyVehicle references.PartyVehicle

	LastLargeMapPosition references.Position

	DateTime UltimaDate

	Provisions Provisions
	Karma      byte
	QtyGold    uint16

	// open door
	openDoorPos   *references.Position
	openDoorTurns int
}

type Provisions struct {
	QtyFood      uint16
	QtyGems      byte
	QtyTorches   byte
	QtyKeys      byte
	QtySkullKeys byte
}

func (g *GameState) GetMapType() GeneralMapType {
	if g.Location == references.Britannia_Underworld {
		return LargeMap
	}
	return SmallMap
}

func (g *GameState) ProcessEndOfTurn() {
	if g.openDoorPos != nil {
		if g.openDoorTurns == 0 {
			g.LayeredMaps.LayeredMaps[SmallMap].UnSetTile(MapOverrideLayer, g.openDoorPos)
			g.openDoorPos = nil
		} else {
			g.openDoorTurns--
		}
	}
}

func (g *GameState) OpenDoor(direction Direction) {
	const defaultTurnsForDoorOpen = 2
	mapType := GetMapTypeByLocation(g.Location)
	newPosition := direction.GetNewPositionInDirection(&g.Position)
	g.LayeredMaps.LayeredMaps[mapType].SetTile(MapOverrideLayer, newPosition, indexes.BrickFloor)

	if g.openDoorPos != nil {
		g.LayeredMaps.LayeredMaps[mapType].UnSetTile(MapOverrideLayer, g.openDoorPos)
	}
	g.openDoorPos = newPosition
	g.openDoorTurns = defaultTurnsForDoorOpen
}
