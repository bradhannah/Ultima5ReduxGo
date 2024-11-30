package game_state

import (
	"fmt"
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const (
	NPlayers                      = 6
	DefaultSmallMapMinutesPerTurn = 1
	DefaultLargeMapMinutesPerTurn = 2
)

const (
	nightTowneCloseTime = 8 + 12
	nightTowneOpenTime  = 5
)

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

	LayeredMaps     LayeredMaps
	NPCAIController NPCAIController
	PartyVehicle    references.PartyVehicle

	LastLargeMapPosition references.Position
	LastLargeMapFloor    references.FloorNumber

	DateTime datetime.UltimaDate

	Provisions Provisions
	Karma      Karma
	QtyGold    uint16

	// open door
	openDoorPos   *references.Position
	openDoorTurns int

	ItemStacks map[references.Position]references.ItemStack

	XTilesInMap int
	YTilesInMap int
}

type Provisions struct {
	QtyFood      uint16
	QtyGems      byte
	QtyTorches   byte
	QtyKeys      byte
	QtySkullKeys byte
}

func (g *GameState) IsAvatarAtPosition(pos *references.Position) bool {
	return g.Position.Equals(pos)
}

func (g *GameState) GetCurrentSmallLocationReference() *references.SmallLocationReference {
	return g.GameReferences.LocationReferences.GetLocationReference(g.Location)
}

func (g *GameState) GetLayeredMapByCurrentLocation() *LayeredMap {
	return g.LayeredMaps.GetLayeredMap(g.Location.GetMapType(), g.Floor)
}

func (g *GameState) IsOutOfBounds(position references.Position) bool {
	if g.Location.GetMapType() == references.LargeMapType {
		if position.X > references.XLargeMapTiles-1 || position.Y >= references.YLargeMapTiles {
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

func (g *GameState) IsPassable(pos *references.Position) bool {
	theMap := g.LayeredMaps.GetLayeredMap(g.Location.GetMapType(), g.Floor)
	topTile := theMap.GetTopTile(pos)
	if topTile == nil {
		return false
	}
	return topTile.IsPassable(g.PartyVehicle)
}

func (g *GameState) IsNPCPassable(pos *references.Position) bool {
	theMap := g.LayeredMaps.GetLayeredMap(g.Location.GetMapType(), g.Floor)
	topTile := theMap.GetTopTile(pos)
	if topTile == nil {
		return false
	}
	return topTile.IsPassable(references.NoPartyVehicle) || topTile.Index.IsUnlockedDoor()
}

func (g *GameState) GetArchwayPortcullisSpriteByTime() indexes.SpriteIndex {
	if g.DateTime.Hour >= nightTowneCloseTime || g.DateTime.Hour < nightTowneOpenTime {
		return indexes.Portcullis
	}
	return indexes.BrickWallArchway
}

func (g *GameState) GetDrawBridgeWaterByTime(origIndex indexes.SpriteIndex) indexes.SpriteIndex {
	if g.DateTime.Hour >= nightTowneCloseTime || g.DateTime.Hour < nightTowneOpenTime {
		return indexes.WaterShallow
	}
	return origIndex
}
