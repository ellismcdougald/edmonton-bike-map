package routing

import (
	"container/heap"
)

// Item represents a node in the priority queue with an associated distance value.
// The index field is used internally by the heap to track the item's position.
type Item struct {
	NodeID   int64   // Unique identifier of the node
	Distance float64 // Distance value used for priority ordering
	index    int     // Index of the item in the heap (maintained by the heap.Interface methods)
}

// PriorityQueue implements a min-heap priority queue of Items based on their Distance.
type PriorityQueue []*Item

// Len returns the number of items in the priority queue.
func (pq PriorityQueue) Len() int {
	return len(pq)
}

// Less reports whether the item with index i should sort before the item with index j.
// Here, Items with smaller Distance have higher priority.
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Distance < pq[j].Distance
}

// Swap swaps the items with indexes i and j and updates their indices accordingly.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Push adds a new item to the priority queue.
// This method is called by heap.Push.
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the item with the highest priority (lowest Distance).
// This method is called by heap.Pop.
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // Mark as removed
	*pq = old[0 : n-1]
	return item
}

// Update modifies the Distance of an item in the priority queue and re-establishes heap order.
// Use this when the priority of an existing item changes.
func (pq *PriorityQueue) Update(item *Item, distance float64) {
	item.Distance = distance
	heap.Fix(pq, item.index)
}
