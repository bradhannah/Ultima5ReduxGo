package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const NPlayers = 6

type GameState struct {
	RawSave         [savedGamFileSize]byte
	Characters      [NPlayers]PlayerCharacter
	MoonstoneStatus MoonstoneStatus

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
