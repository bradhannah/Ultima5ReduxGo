package main

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameScene) debugMoveOnMap(position references2.Position) {
	g.gameState.MapState.PlayerLocation.Position = position
}

func (g *GameScene) debugFloorUp() bool {
	if g.gameState.MapState.PlayerLocation.Location == references2.Britannia_Underworld {
		if g.gameState.MapState.PlayerLocation.Floor == -1 {
			g.gameState.MapState.PlayerLocation.Floor = 0
			return true
		}
		return false
	}

	singleMap := g.GetCurrentLocationReference()

	_, maxFloor := singleMap.GetFloorMinMax()
	if g.gameState.MapState.PlayerLocation.Floor+1 > maxFloor {
		return false
	}
	g.gameState.MapState.PlayerLocation.Floor++
	g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)

	return true
}

func (g *GameScene) debugFloorDown() bool {
	if g.gameState.MapState.PlayerLocation.Location == references2.Britannia_Underworld {
		if g.gameState.MapState.PlayerLocation.Floor == 0 {
			g.gameState.MapState.PlayerLocation.Floor = -1
			return true
		}
		return false
	}

	singleMap := g.GetCurrentLocationReference()
	minFloor, _ := singleMap.GetFloorMinMax()
	if g.gameState.MapState.PlayerLocation.Floor-1 < minFloor {
		return false
	}
	g.gameState.MapState.PlayerLocation.Floor--
	g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
	return true
}

func (g *GameScene) debugFloorY(nFloor references2.FloorNumber) bool {
	if g.gameState.MapState.PlayerLocation.Location == references2.Britannia_Underworld {
		if nFloor < -1 {
			g.gameState.MapState.PlayerLocation.Floor = -1
		} else if nFloor > 0 {
			g.gameState.MapState.PlayerLocation.Floor = 0
		} else {
			g.gameState.MapState.PlayerLocation.Floor = nFloor
		}
		return true
	}
	singleMap := g.GetCurrentLocationReference()
	minFloor, maxFloor := singleMap.GetFloorMinMax()
	if nFloor < minFloor || nFloor > maxFloor {
		// bad
		return false
	}
	g.gameState.MapState.PlayerLocation.Floor = nFloor
	return true
}
