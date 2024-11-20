package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameScene) DebugMoveOnMap(position references.Position) {
	g.gameState.Position = position
}

func (g *GameScene) DebugFloorUp() bool {
	if g.gameState.Location == references.Britannia_Underworld {
		if g.gameState.Floor == -1 {
			g.gameState.Floor = 0
			return true
		}
		return false
	}

	singleMap := g.GetCurrentLocationReference()

	_, maxFloor := singleMap.GetFloorMinMax()
	if g.gameState.Floor+1 > maxFloor {
		return false
	}
	g.gameState.Floor++
	g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)

	return true
}

func (g *GameScene) DebugFloorDown() bool {
	if g.gameState.Location == references.Britannia_Underworld {
		if g.gameState.Floor == 0 {
			g.gameState.Floor = -1
			return true
		}
		return false
	}

	singleMap := g.GetCurrentLocationReference()
	minFloor, _ := singleMap.GetFloorMinMax()
	if g.gameState.Floor-1 < minFloor {
		return false
	}
	g.gameState.Floor--
	g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
	return true
}

func (g *GameScene) DebugFloorY(nFloor references.FloorNumber) bool {
	if g.gameState.Location == references.Britannia_Underworld {
		if nFloor < -1 {
			g.gameState.Floor = -1
		} else if nFloor > 0 {
			g.gameState.Floor = 0
		} else {
			g.gameState.Floor = nFloor
		}
		return true
	}
	singleMap := g.GetCurrentLocationReference()
	minFloor, maxFloor := singleMap.GetFloorMinMax()
	if nFloor < minFloor || nFloor > maxFloor {
		// bad
		return false
	}
	g.gameState.Floor = nFloor
	return true
}
