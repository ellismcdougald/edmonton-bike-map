package routing

import (
	"container/heap"
)

type Item struct {
	NodeID int64
	Distance float64
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Distance < pq[j].Distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n - 1]
	item.index = -1
	*pq = old[0 : n - 1]
	return item
}

func (pq *PriorityQueue) Update(item *Item, distance float64) {
	item.Distance = distance
	heap.Fix(pq, item.index)
}