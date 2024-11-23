package game_state

import (
	"container/heap"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
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
	// topLeft     references.Position
}

func NewAStarMap() *AStarMap {
	a := &AStarMap{}
	// a.topLeft = topLeft
	return a
}

func (m *AStarMap) AStar(start, goal references.Position) []references.Position {
	openSet := &aStarPriorityQueue{}
	heap.Init(openSet)

	startNode := &AStarNode{
		Position: start,
		GScore:   0,
		FScore:   Heuristic(start, goal),
	}
	heap.Push(openSet, startNode)

	closedSet := make(map[references.Position]bool)
	gScore := map[references.Position]int{
		start: 0,
	}

	for openSet.Len() > 0 {
		// Get the node with the lowest FScore
		current := heap.Pop(openSet).(*AStarNode)

		// If we've reached the goal, reconstruct the path
		if current.Position == goal {
			return reconstructPath(current)
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
					FScore:   tentativeGScore + Heuristic(neighbor, goal),
					Parent:   current,
				})
			}
		}
	}

	// No path found
	return nil
}

func reconstructPath(node *AStarNode) []references.Position {
	var path []references.Position
	for node != nil {
		path = append([]references.Position{node.Position}, path...)
		node = node.Parent
	}
	return path
}

func (m *AStarMap) InitializeByLayeredMap(lMap *LayeredMap, extraBlockTiles []references.Position) {
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

func Heuristic(a, b references.Position) int {
	return helpers.AbsInt(int(a.X-b.X)) + helpers.AbsInt(int(a.Y-b.Y))
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
