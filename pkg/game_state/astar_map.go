package game_state

import (
	"container/heap"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type AStarNode struct {
	Position references.Position
	GScore   int // Cost from start to current node
	FScore   int // GScore + heuristic cost to the goal
	Parent   *AStarNode
}

type AStarMap struct {
	walkableMap map[references.Position]int
	bWrap       bool
	maxX, maxY  references.Coordinate
	mapUnit     MapUnit
}

func NewAStarMap() *AStarMap {
	a := &AStarMap{}
	return a
}

func (m *AStarMap) AStar(goal references.Position) []references.Position {
	openSet := &aStarPriorityQueue{}
	heap.Init(openSet)

	startNode := &AStarNode{
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
		current := heap.Pop(openSet).(*AStarNode)

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
				heap.Push(openSet, &AStarNode{
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

func (m *AStarMap) reconstructPath(node *AStarNode) []references.Position {
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

func (m *AStarMap) InitializeByLayeredMap(mapUnit MapUnit, lMap *LayeredMap, extraBlockTiles []references.Position) {
	m.mapUnit = mapUnit
	m.bWrap = false
	m.walkableMap = make(map[references.Position]int)
	for x := references.Coordinate(0); x < lMap.xMax; x++ {
		for y := references.Coordinate(0); y < lMap.yMax; y++ {
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

func (m *AStarMap) InitializeByLayeredMapWithLimit(
	mapUnit MapUnit,
	lMap *LayeredMap,
	extraBlockTiles []references.Position,
	bWrap bool,
	centerPos references.Position,
	nMaxRadius int,
	maxX, maxY references.Coordinate) {

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
			if enemy, ok := m.mapUnit.(*NPCEnemy); ok {
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

type aStarPriorityQueue []*AStarNode

func (pq aStarPriorityQueue) Len() int { return len(pq) }

func (pq aStarPriorityQueue) Less(i, j int) bool {
	return pq[i].FScore < pq[j].FScore // Lower FScore has higher priority
}

func (pq aStarPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *aStarPriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*AStarNode))
}

func (pq *aStarPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
