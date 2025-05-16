package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const (
	NPlayers                                   = 6
	DefaultSmallMapMinutesPerTurn              = 1
	DefaultLargeMapMinutesPerTurn              = 2
	DefaultNumberOfTurnsUntilTorchExtinguishes = 100
)

const (
	nightTowneCloseTime = 8 + 12
	nightTowneOpenTime  = 5
)

type DebugOptions struct {
	FreeMove   bool
	MonsterGen bool
}

type GameState struct {
	RawSave         [savedGamFileSize]byte
	Characters      [NPlayers]PlayerCharacter
	MoonstoneStatus MoonstoneStatus

	DebugOptions DebugOptions

	GameReferences *references.GameReferences

	Location references.Location
	Position references.Position
	Floor    references.FloorNumber
	//avatarPosition references.Position

	LayeredMaps             LayeredMaps
	CurrentNPCAIController  NPCAIController
	LargeMapNPCAIController map[references.World]*NPCAIControllerLargeMap

	PartyVehicle NPCVehicle

	LastLargeMapPosition references.Position
	LastLargeMapFloor    references.FloorNumber

	TheOdds TheOdds

	DateTime datetime.UltimaDate

	Inventory Inventory
	Karma     Karma

	// open door
	openDoorPos   *references.Position
	openDoorTurns int

	ItemStacksMap references.ItemStacksMap

	XTilesVisibleOnGameScreen int
	YTilesVisibleOnGameScreen int

	Lighting Lighting
}

const (
	MaxGold               = 9999
	MaxProvisionFood      = 9999
	MaxProvisionGems      = 99
	MaxProvisionTorches   = 99
	MaxProvisionKey       = 99
	MaxProvisionSkullKeys = 99
)

type ProvisionsQuantity struct {
	Food      uint16
	Gems      byte
	Torches   byte
	Keys      byte
	SkullKeys byte
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
			log.Fatalf("Exceeded large map tiles: X=%d, Y=%d", position.X, position.Y)
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

func (g *GameState) GetCurrentLayeredMap() *LayeredMap {
	return g.LayeredMaps.GetLayeredMap(g.Location.GetMapType(), g.Floor)
}

func (g *GameState) GetCurrentLayeredMapAvatarTopTile() *references.Tile {
	theMap := g.GetCurrentLayeredMap()
	topTile := theMap.GetTopTile(&g.Position)
	if topTile == nil {
		return nil
	}
	return topTile
}

func (g *GameState) IsPassable(pos *references.Position) bool {
	theMap := g.GetCurrentLayeredMap()
	topTile := theMap.GetTopTile(pos)
	if topTile == nil {
		return false
	}

	return topTile.IsPassable(g.PartyVehicle.VehicleType)
}

func (g *GameState) IsNPCPassable(pos *references.Position) bool {
	theMap := g.LayeredMaps.GetLayeredMap(g.Location.GetMapType(), g.Floor)
	topTile := theMap.GetTopTile(pos)

	if topTile == nil {
		return false
	}
	return topTile.IsPassable(references.NPC) || topTile.Index.IsUnlockedDoor()
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

func (g *GameState) BoardVehicle(vehicle NPCVehicle) bool {
	g.PartyVehicle = vehicle

	if !g.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(g.Position) {
		log.Fatalf("Unexpected - tried to remove NPC at position X=%d,Y=%d but failed", g.Position.X, g.Position.Y)
	}

	return true
}

func (g *GameState) GetTilesVisibleOnScreen() (int, int) {
	return g.GetCurrentLayeredMap().GetTilesVisibleOnScreen()
}

func (g *GameState) GetTopLeftExtent() references.Position {
	return g.GetCurrentLayeredMap().topLeft
}

func (g *GameState) GetBottomRightExtent() references.Position {
	return g.GetCurrentLayeredMap().bottomRight
}

func (g *GameState) IsWrappedMap() bool {
	return g.GetCurrentLayeredMap().IsWrappedMap()
}

func (g *GameState) GetBottomRightWithoutOverflow() references.Position {
	return g.GetCurrentLayeredMap().GetBottomRightWithoutOverflow()
}
