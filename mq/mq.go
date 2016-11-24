package mq

import (
	"container/heap"
	"sync"
)

type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

var lock = sync.Locker{}

// A PriorityQueue implements heap.Interface and holds Items.
type Queue []*Item

func (pq *Queue) Len() int { return len(pq) }

func (pq *Queue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq *Queue) Swap(i, j int) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *Queue) Push(x interface{}) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *Queue) Pop() interface{} {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *Queue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
