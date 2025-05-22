package astar

type priorityQueue []*Node

func (pq *priorityQueue) Len() int {
	return len(*pq)
}

func (pq *priorityQueue) Less(i, j int) bool {
	return (*pq)[i].FScore < (*pq)[j].FScore // lower FScore = higher priority
}

func (pq *priorityQueue) Swap(i, j int) {
	q := *pq
	q[i], q[j] = q[j], q[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
