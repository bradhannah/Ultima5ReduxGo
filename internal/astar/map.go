package astar

import (
	"container/heap"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type Map struct {
	walkableMap map[references.Position]int
	bWrap       bool
	maxX, maxY  references.Coordinate
	mapUnit     map_units.MapUnit
}

func NewAStarMap() *Map {
	a := &Map{}
	return a
}

func (m *Map) AStar(goal references.Position) []references.Position {
	openSet := &priorityQueue{}
	heap.Init(openSet)

	startNode := &Node{
		Position: m.mapUnit.Pos(),
		GScore:   0,
		FScore:   m.mapUnit.PosPtr().HeuristicTileDistance(goal),
	}
	heap.Push(openSet, startNode)

	closedSet := make(map[references.Position]bool)
	gScore := map[references.Position]int{
		m.mapUnit.Pos(): 0,
	}

	for openSet.Len() > 0 {
		// Get the node with the lowest FScore
		current := heap.Pop(openSet).(*Node)

		// If we've reached the goal, reconstruct the path
		if current.Position == goal {
			return m.reconstructPath(current)
		}

		closedSet[current.Position] = true

		// Explore neighbors
		for _, neighbor := range current.Position.Neighbors() {
			// Check if the neighbor is walkable
			weight, exists := m.walkableMap[neighbor]
			if !exists || weight < 0 { // Impassable tile
				continue
			}

			// Skip if already in the closed set
			if closedSet[neighbor] {
				continue
			}

			// Calculate tentative GScore
			tentativeGScore := gScore[current.Position] + weight

			// Update scores and add to the open set
			if _, seen := gScore[neighbor]; !seen || tentativeGScore < gScore[neighbor] {
				gScore[neighbor] = tentativeGScore
				heap.Push(openSet, &Node{
					Position: neighbor,
					GScore:   tentativeGScore,
					FScore:   tentativeGScore + neighbor.HeuristicTileDistance(goal),
					Parent:   current,
				})
			}
		}
	}

	// No path found
	return nil
}

func (m *Map) reconstructPath(node *Node) []references.Position {
	var path []references.Position
	for node != nil {
		if m.bWrap {
			// if it is wrapped, then we check for wrapping at this last possible moment
			node.Position = *node.Position.GetWrapped(m.maxX, m.maxY)
		}

		path = append([]references.Position{node.Position}, path...)
		node = node.Parent
	}
	return path
}

func (m *Map) InitializeByLayeredMap(mapUnit map_units.MapUnit, lMap *map_state.LayeredMap, extraBlockTiles []references.Position) {
	m.mapUnit = mapUnit
	m.bWrap = false
	m.walkableMap = make(map[references.Position]int)
	for x := references.Coordinate(0); x < lMap.XMaxTilesPerMap; x++ {
		for y := references.Coordinate(0); y < lMap.YMaxTilesPerMap; y++ {
			pos := references.Position{
				X: x,
				Y: y,
			}
			topTile := lMap.GetTopTile(&pos)
			m.walkableMap[pos] = topTile.GetWalkableWeight()
		}
	}
	for i := 0; i < len(extraBlockTiles); i++ {
		m.walkableMap[extraBlockTiles[i]] = -1
	}
}

func (m *Map) InitializeByLayeredMapWithLimit(
	mapUnit map_units.MapUnit,
	lMap *map_state.LayeredMap,
	extraBlockTiles []references.Position,
	bWrap bool,
	centerPos references.Position,
	nMaxRadius int,
	maxX, maxY references.Coordinate,
) {
	m.mapUnit = mapUnit
	m.bWrap = bWrap
	m.maxX, m.maxY = maxX, maxY
	m.walkableMap = make(map[references.Position]int)

	for x := references.Coordinate(int(centerPos.X) - nMaxRadius); x < centerPos.X+references.Coordinate(nMaxRadius); x++ {
		for y := references.Coordinate(int(centerPos.Y) - nMaxRadius); y < centerPos.Y+references.Coordinate(nMaxRadius); y++ {
			pos := references.Position{
				X: x % maxX,
				Y: y % maxY,
			}
			if m.bWrap {
				pos = *pos.GetWrapped(m.maxX, m.maxY)
			}
			topTile := lMap.GetTopTile(&pos)
			if enemy, ok := m.mapUnit.(*map_units.NPCEnemy); ok {
				if enemy.EnemyReference.CanMoveToTile(topTile) {
					m.walkableMap[pos] = 1
				} else {
					m.walkableMap[pos] = -1
				}
			} else {
				m.walkableMap[pos] = topTile.GetWalkableWeight()
			}
		}
	}
	for i := 0; i < len(extraBlockTiles); i++ {
		m.walkableMap[extraBlockTiles[i]] = -1
	}
}
