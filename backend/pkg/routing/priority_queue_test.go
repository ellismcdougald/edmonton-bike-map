package routing

import (
	"container/heap"
	"testing"
)

func TestPriorityQueueBasic(t *testing.T) {
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Add some items
	items := []*Item{
		{NodeID: 1, Distance: 5},
		{NodeID: 2, Distance: 3},
		{NodeID: 3, Distance: 4},
	}

	for _, item := range items {
		heap.Push(pq, item)
	}

	// After pushing, Len should be 3
	if pq.Len() != 3 {
		t.Errorf("Len() = %d; want 3", pq.Len())
	}

	// Pop items and check order: expect NodeID 2 (3), 3 (4), 1 (5)
	expectedOrder := []struct {
		nodeID   int64
		distance float64
	}{
		{2, 3},
		{3, 4},
		{1, 5},
	}

	for i, exp := range expectedOrder {
		item := heap.Pop(pq).(*Item)
		if item.NodeID != exp.nodeID || item.Distance != exp.distance {
			t.Errorf("Pop #%d = (NodeID=%d, Distance=%f); want (NodeID=%d, Distance=%f)",
				i, item.NodeID, item.Distance, exp.nodeID, exp.distance)
		}
	}
}

func TestPriorityQueueUpdate(t *testing.T) {
	pq := &PriorityQueue{}
	heap.Init(pq)

	item1 := &Item{NodeID: 1, Distance: 5}
	item2 := &Item{NodeID: 2, Distance: 3}
	item3 := &Item{NodeID: 3, Distance: 4}

	heap.Push(pq, item1)
	heap.Push(pq, item2)
	heap.Push(pq, item3)

	// Update item1's distance to 2 (should bubble up to front)
	pq.Update(item1, 2)

	// Pop should now return item1 first
	item := heap.Pop(pq).(*Item)
	if item.NodeID != 1 {
		t.Errorf("After update, first Pop NodeID = %d; want 1", item.NodeID)
	}

	// Next pops in order
	item = heap.Pop(pq).(*Item)
	if item.NodeID != 2 {
		t.Errorf("Second Pop NodeID = %d; want 2", item.NodeID)
	}

	item = heap.Pop(pq).(*Item)
	if item.NodeID != 3 {
		t.Errorf("Third Pop NodeID = %d; want 3", item.NodeID)
	}
}

func TestPriorityQueueIndices(t *testing.T) {
	pq := &PriorityQueue{}
	heap.Init(pq)

	item1 := &Item{NodeID: 1, Distance: 5}
	item2 := &Item{NodeID: 2, Distance: 3}

	heap.Push(pq, item1)
	heap.Push(pq, item2)

	// Indices should be set correctly
	if item1.index == -1 || item2.index == -1 {
		t.Errorf("Expected valid indices, got item1.index=%d, item2.index=%d", item1.index, item2.index)
	}

	// Pop an item and check index reset for popped item
	popped := heap.Pop(pq).(*Item)
	if popped.index != -1 {
		t.Errorf("Popped item's index = %d; want -1", popped.index)
	}
}
