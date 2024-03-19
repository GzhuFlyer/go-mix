package myqueue

import (
	"container/heap"
	"fmt"
	"math"
)

type SkyFendPQItem[T any] struct {
	Value    T
	Priority int64
	Index    int32
}
type SkyFendPQ[T any] struct {
	Items []*SkyFendPQItem[T]
	//SafeLocker sync.Mutex
}

func (spq *SkyFendPQ[T]) Len() int {
	//spq.SafeLocker.Lock()
	//defer spq.SafeLocker.Unlock()

	return len(spq.Items)
}
func (spq *SkyFendPQ[T]) Swap(i, j int) {
	//spq.SafeLocker.Lock()
	//defer spq.SafeLocker.Unlock()

	spq.Items[i], spq.Items[j] = spq.Items[j], spq.Items[i]
	spq.Items[i].Index = int32(i)
	spq.Items[j].Index = int32(j)
}

func (spq *SkyFendPQ[T]) Less(i, j int) bool {
	//spq.SafeLocker.Lock()
	//defer spq.SafeLocker.Unlock()

	return spq.Items[i].Priority < spq.Items[j].Priority
}
func (spq *SkyFendPQ[T]) Push(dataGen any) {
	//spq.SafeLocker.Lock()
	//defer spq.SafeLocker.Unlock()

	data := dataGen.(*SkyFendPQItem[T])
	n := len(spq.Items)
	data.Index = int32(n)
	spq.Items = append(spq.Items, data)
}
func (spq *SkyFendPQ[T]) Pop() any {
	//spq.SafeLocker.Lock()
	//defer spq.SafeLocker.Unlock()

	o := spq.Items
	n := len(o)

	item := o[n-1]
	item.Index = -1
	spq.Items = spq.Items[0 : n-1]
	return item
}

func QueueT5() {
	type skfString = string
	datas := map[skfString]int{
		"Frank": 3, "Tom": 2, "Alice": 4, "Bob": 5,
	}
	pq := new(SkyFendPQ[skfString])
	i := int32(0)
	for value, priority := range datas {
		heap.Push(pq, &SkyFendPQItem[skfString]{
			Value:    value,
			Priority: int64(priority),
			Index:    i,
		})
		i++
	}
	min := int64(math.MinInt64)
	for pq.Len() > 0 {
		item := heap.Pop(pq).(*SkyFendPQItem[skfString])
		if item.Priority < min {
			msg := fmt.Sprintf("build min heap fail,[item.Priority=%s - maxn=%d]: ", item.Priority, min)
			panic(msg)
		}
		min = item.Priority
		fmt.Printf("Value - priority : %s - %d \n", item.Value, item.Priority)
	}
}

func QueueT6() {
	type skfString = string
	datas := map[skfString]int{
		"Frank": 3, "Tom": 2, "Alice": 4, "Bob": 5,
	}
	pq := new(SkyFendPQ[skfString])
	i := int32(0)
	for value, priority := range datas {
		heap.Push(pq, &SkyFendPQItem[skfString]{
			Value:    value,
			Priority: int64(priority),
			Index:    i,
		})
		i++
	}
	min := int64(math.MinInt64)
	for pq.Len() > 0 {
		item := heap.Pop(pq).(*SkyFendPQItem[skfString])
		if item.Priority < min {
			msg := fmt.Sprintf("build min heap fail,[item.Priority=%s - maxn=%d]: ", item.Priority, min)
			panic(msg)
		}
		min = item.Priority
		fmt.Printf("Value - priority : %s - %d \n", item.Value, item.Priority)
	}
}
