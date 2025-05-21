package map_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type MapState struct {
	PlayerLocation references.PlayerLocation
	LayeredMaps    LayeredMaps

	XTilesVisibleOnGameScreen int
	YTilesVisibleOnGameScreen int

	Lighting       Lighting
	gameDimensions GameDimensions

	// open door
	openDoorPos   *references.Position
	openDoorTurns int
}

type NewMapStateInput struct {
	GameDimensions            GameDimensions
	XTilesVisibleOnGameScreen int
	YTilesVisibleOnGameScreen int
}

func NewMapState(input NewMapStateInput) *MapState {
	mapState := &MapState{}
	mapState.Lighting = NewLighting(input.GameDimensions)
	mapState.XTilesVisibleOnGameScreen = input.XTilesVisibleOnGameScreen
	mapState.YTilesVisibleOnGameScreen = input.YTilesVisibleOnGameScreen
	return mapState
}

func (m *MapState) GetLayeredMapByCurrentLocation() *LayeredMap {
	return m.LayeredMaps.GetLayeredMap(m.PlayerLocation.Location.GetMapType(), m.PlayerLocation.Floor)
}

func (g *MapState) IsNPCPassable(pos *references.Position) bool {
	theMap := g.GetLayeredMapByCurrentLocation()
	// theMap := g.LayeredMaps.GetLayeredMap(g.PlayerLocation.Location.GetMapType(), g.PlayerLocation.Floor)
	topTile := theMap.GetTopTile(pos)

	if topTile == nil {
		return false
	}
	return topTile.IsPassable(references.NPC) || topTile.Index.IsUnlockedDoor()
}
