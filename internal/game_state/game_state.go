package game_state

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/ai"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const (
	DefaultSmallMapMinutesPerTurn              = 1
	DefaultLargeMapMinutesPerTurn              = 2
	DefaultNumberOfTurnsUntilTorchExtinguishes = 100
)

const (
	nightTowneCloseTime = 8 + 12
	nightTowneOpenTime  = 5
)

type GameState struct {
	PartyState   party_state.PartyState
	MapState     map_state.MapState
	AiController ai.NPCAIController

	RawSave         [savedGamFileSize]byte
	MoonstoneStatus MoonstoneStatus

	DebugOptions references.DebugOptions

	GameReferences *references.GameReferences

	PartyVehicle map_units.NPCFriendly

	CurrentNPCAIController  ai.NPCAIController
	LargeMapNPCAIController map[references.World]*ai.NPCAIControllerLargeMap

	LastLargeMapPosition references.Position
	LastLargeMapFloor    references.FloorNumber

	TheOdds references.TheOdds

	DateTime datetime.UltimaDate

	ItemStacksMap references.ItemStacksMap
}

func (g *GameState) IsAvatarAtPosition(pos *references.Position) bool {
	return g.MapState.PlayerLocation.Position.Equals(pos)
}

func (g *GameState) GetCurrentSmallLocationReference() *references.SmallLocationReference {
	return g.GameReferences.LocationReferences.GetLocationReference(g.MapState.PlayerLocation.Location)
}

func (g *GameState) GetLayeredMapByCurrentLocation() *map_state.LayeredMap {
	return g.MapState.LayeredMaps.GetLayeredMap(g.MapState.PlayerLocation.Location.GetMapType(), g.MapState.PlayerLocation.Floor)
}

func (g *GameState) IsOutOfBounds(position references.Position) bool {
	if g.MapState.PlayerLocation.Location.GetMapType() == references.LargeMapType {
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

func (g *GameState) GetCurrentLayeredMap() *map_state.LayeredMap {
	return g.MapState.LayeredMaps.GetLayeredMap(g.MapState.PlayerLocation.Location.GetMapType(), g.MapState.PlayerLocation.Floor)
}

func (g *GameState) GetCurrentLayeredMapAvatarTopTile() *references.Tile {
	theMap := g.GetCurrentLayeredMap()
	topTile := theMap.GetTopTile(&g.MapState.PlayerLocation.Position)
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

	return topTile.IsPassable(g.PartyVehicle.GetVehicleDetails().VehicleType)
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

func (g *GameState) BoardVehicle(vehicle map_units.NPCFriendly) bool {
	g.PartyVehicle = vehicle

	if !g.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(g.MapState.PlayerLocation.Position) {
		log.Fatalf("Unexpected - tried to remove NPC at position X=%d,Y=%d but failed", g.MapState.PlayerLocation.Position.X, g.MapState.PlayerLocation.Position.Y)
	}

	return true
}

func (g *GameState) GetTilesVisibleOnScreen() (int, int) {
	return g.GetCurrentLayeredMap().GetTilesVisibleOnScreen()
}

func (g *GameState) GetTopLeftExtent() references.Position {
	return g.GetCurrentLayeredMap().TopLeft
}

func (g *GameState) GetBottomRightExtent() references.Position {
	return g.GetCurrentLayeredMap().BottomRight
}

func (g *GameState) IsWrappedMap() bool {
	return g.GetCurrentLayeredMap().IsWrappedMap()
}

func (g *GameState) GetBottomRightWithoutOverflow() references.Position {
	return g.GetCurrentLayeredMap().GetBottomRightWithoutOverflow()
}

func (g *GameState) GetCurrentLargeMapNPCAIController() *ai.NPCAIControllerLargeMap {
	if g.MapState.PlayerLocation.Location.GetMapType() != references.LargeMapType {
		log.Fatalf("Expected large map type, got %d", g.MapState.PlayerLocation.Location.GetMapType())
	}
	var npcAiCon *ai.NPCAIControllerLargeMap
	if g.MapState.IsOverworld() {
		npcAiCon = g.LargeMapNPCAIController[references.OVERWORLD]
	} else {
		npcAiCon = g.LargeMapNPCAIController[references.UNDERWORLD]
	}
	return npcAiCon
}
