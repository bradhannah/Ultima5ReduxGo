package game_state

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"log"
)

const NPlayers = 6

type DebugOptions struct {
	FreeMove bool
}

type GameState struct {
	RawSave         [savedGamFileSize]byte
	Characters      [NPlayers]PlayerCharacter
	MoonstoneStatus MoonstoneStatus

	DebugOptions DebugOptions

	GameReferences *references.GameReferences

	Location       references.Location
	Position       references.Position
	Floor          references.FloorNumber
	avatarPosition references.Position

	LayeredMaps  LayeredMaps
	PartyVehicle references.PartyVehicle

	LastLargeMapPosition references.Position

	DateTime UltimaDate

	Provisions Provisions
	Karma      Karma
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

func (g *GameState) LargeMapProcessEndOfTurn() {
	return
}

func (g *GameState) SmallMapProcessEndOfTurn() {
	if g.openDoorPos != nil {
		if g.openDoorTurns == 0 {
			tile := g.LayeredMaps.GetTileRefByPosition(references.SmallMapType, MapLayer, g.openDoorPos, g.Floor)
			if tile.Index.IsWindowedDoor() {
				g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).SetTileByLayer(MapOverrideLayer, g.openDoorPos, indexes.RegularDoorView)
			} else {
				g.LayeredMaps.GetLayeredMap(references.SmallMapType, g.Floor).SetTileByLayer(MapOverrideLayer, g.openDoorPos, indexes.RegularDoor)
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

func (g *GameState) GetCurrentSmallLocationReference() *references.SmallLocationReference {
	return g.GameReferences.LocationReferences.GetLocationReference(g.Location)
}

func (g *GameState) IsOutOfBounds(position references.Position) bool {
	if g.Location.GetMapType() == references.LargeMapType {
		if position.X > references.XLargeMapTiles-1 || position.Y > references.YLargeMapTiles {
			//|| position.X < 0 || position.Y < 0 {
			log.Fatal(fmt.Sprintf("Exceeded large map tiles: X=%d, Y=%d", position.X, position.Y))
		}
		return false
	}

	if position.X < 0 || position.Y < 0 {
		return true
	}

	if position.X > g.GetCurrentSmallLocationReference().GetMaxX() || position.Y > g.GetCurrentSmallLocationReference().GetMaxY() {
		return true
	}

	return false
}
