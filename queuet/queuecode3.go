// This example demonstrates a priority queue built using the heap interface.
package myqueue

import (
	"container/heap"
	"fmt"
	"math"
)

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
	// return pq[i].priority < pq[j].priority
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
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func QueueT3() {
	// Some items and their priorities.
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	fmt.Println("len items  = ", len(items))
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}

	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	fmt.Println("end push......")
	// pq.update(item, item.value, 3)
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
		// fmt.Printf("\nindex = %d\n", item.index)
	}
	fmt.Println("\n------------------------------")

	item2 := &Item{
		value:    "putao",
		priority: 3,
	}
	heap.Push(&pq, item2)
	// pq.update(item2, item2.value, 1)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
		// fmt.Printf("\nindex = %d\n", item.index)
	}
	fmt.Println("")
}
func QueueT4() {
	datas := map[string]int{
		"Frank": 3, "Tom": 2, "Alice": 4, "Bob": 5,
	}
	pq := new(PriorityQueue)
	// var t []byte
	// t = append(t, 1)
	maxn := math.MaxInt64
	i := 0
	for value, priority := range datas {
		heap.Push(pq, &Item{
			value:    value,
			priority: priority,
			index:    i,
		})
		i++
	}

	len := pq.Len()
	fmt.Println("len = ", len)
	for pq.Len() > 0 {
		if pq.Len() <= 0 {
			break
		}
		item := heap.Pop(pq).(*Item)
		if item.priority > maxn {
			msg := fmt.Sprintf("build min heap fail,[item.Priority - maxn]: ", item.priority, maxn)
			panic(msg)
		}
		maxn = item.priority
		fmt.Printf("value: %s, priority: %d \n", item.value, item.priority)
		fmt.Println("pq len = ", pq.Len())
	}

}
