package map_state

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type MapState struct {
	PlayerLocation references2.PlayerLocation
	LayeredMaps    LayeredMaps

	XTilesVisibleOnGameScreen int
	YTilesVisibleOnGameScreen int

	Lighting       Lighting
	gameDimensions GameDimensions

	// open door
	openDoorPos   *references2.Position
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

func (m *MapState) IsNPCPassable(pos *references2.Position) bool {
	theMap := m.GetLayeredMapByCurrentLocation()
	// theMap := m.LayeredMaps.GetLayeredMap(m.PlayerLocation.Location.GetMapType(), m.PlayerLocation.Floor)
	topTile := theMap.GetTopTile(pos)

	if topTile == nil {
		return false
	}
	return topTile.IsPassable(references2.NPC) || topTile.Index.IsUnlockedDoor()
}
