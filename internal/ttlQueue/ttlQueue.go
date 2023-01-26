package ttlQueue

import (
	"container/heap"
	"errors"
)

type Item struct {
	Value               string
	ExpirationTimestamp *int64 // UTC timestamp, when item should be removed. If nil, store forever
	// The index is needed by update and is maintained by the ttlQueue.Interface methods.
	index int // The index of the item in the ttlQueue.
}

type TTLQueue []*Item

func (tq TTLQueue) Len() int {
	return len(tq)
}

func (tq TTLQueue) Peek() *Item { // get top element of ttlQueue
	if tq.Len() == 0 {
		return nil
	}
	return tq[0]
}

func (tq TTLQueue) Less(i, j int) bool {
	// on top of heap should be item with the lowest (closest to Now) exp timestamp
	if tq[i].ExpirationTimestamp == nil {
		return false
	}
	if tq[j].ExpirationTimestamp == nil {
		return true
	}
	return *tq[i].ExpirationTimestamp < *tq[j].ExpirationTimestamp
}

func (tq TTLQueue) Swap(i, j int) {
	tq[i], tq[j] = tq[j], tq[i]
	tq[i].index = i
	tq[j].index = j
}

func (tq *TTLQueue) Push(x any) {
	n := len(*tq)
	item := x.(*Item)
	item.index = n
	*tq = append(*tq, item)
}

func (tq *TTLQueue) Pop() any {
	old := *tq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*tq = old[0 : n-1]
	return item
}

func (tq *TTLQueue) Update(value string, expiration *int64) error {
	item, err := tq.find(value)
	if err != nil {
		return err
	}
	tq.update(item, value, expiration)
	return nil
}

func (tq *TTLQueue) Delete(value string) error {
	item, err := tq.find(value)
	if err != nil {
		return err
	}
	heap.Remove(tq, item.index)
	return nil
}

func (tq *TTLQueue) update(item *Item, value string, expiration *int64) {
	item.Value = value
	item.ExpirationTimestamp = expiration
	heap.Fix(tq, item.index)
}

func (tq *TTLQueue) find(value string) (itemFound *Item, err error) {
	for _, item := range *tq {
		if item.Value == value {
			itemFound = item
			break
		}
	}
	if itemFound == nil {
		return nil, errors.New("item not found")
	}
	return itemFound, nil
}
