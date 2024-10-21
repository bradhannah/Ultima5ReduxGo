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
	JimmyDoorDirectionInput
)

type DebugOptions struct {
	FreeMove bool
}

type GameState struct {
	RawSave         [savedGamFileSize]byte
	Characters      [NPlayers]PlayerCharacter
	MoonstoneStatus MoonstoneStatus

	DebugOptions DebugOptions

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
			tile := g.LayeredMaps.GetTileRefByPosition(SmallMap, MapLayer, g.openDoorPos)
			if tile.Index.IsWindowedDoor() {
				g.LayeredMaps.LayeredMaps[SmallMap].SetTile(MapOverrideLayer, g.openDoorPos, indexes.RegularDoorView)
			} else {
				g.LayeredMaps.LayeredMaps[SmallMap].SetTile(MapOverrideLayer, g.openDoorPos, indexes.RegularDoor)
			}
			g.openDoorPos = nil
		} else {
			g.openDoorTurns--
		}
	}
}

func (g *GameState) IsAvatarAtPosition(pos *references.Position) bool {
	return g.Position.Equals(*pos)
}
