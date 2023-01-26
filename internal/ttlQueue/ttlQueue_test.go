package ttlQueue

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	val1 = int64(1)
	val2 = int64(2)
	val3 = int64(3)
	val4 = int64(4)
	val5 = int64(5)
)

var items = map[string]*int64{"0": nil, "5": &val5, "2": &val2, "4": &val4, "3": &val3, "1": &val1}

func TestPushToQueue(t *testing.T) {
	tq := TTLQueue{}
	heap.Init(&tq)
	for key, val := range items {
		item := Item{Value: key, ExpirationTimestamp: val}
		heap.Push(&tq, &item)
	}
	// Take the items out; they arrive in increasing ts order.
	expectedItems := []string{"1", "2", "3", "4", "5", "0"}
	var gotItems []string
	for tq.Len() > 0 {
		item := heap.Pop(&tq).(*Item)
		gotItems = append(gotItems, item.Value)
	}
	assert.Equal(t, gotItems, expectedItems)
}

func TestPopFromQueue(t *testing.T) {
	tq := TTLQueue{}
	heap.Init(&tq)
	for key, val := range items {
		item := Item{Value: key, ExpirationTimestamp: val}
		heap.Push(&tq, &item)
	}
	ts := int64(1)
	expectedItem := Item{Value: "1", ExpirationTimestamp: &ts, index: -1}
	item := heap.Pop(&tq).(*Item)
	assert.Equal(t, *item, expectedItem)

	expectedItems := []string{"2", "3", "4", "5", "0"}
	var gotItems []string
	for tq.Len() > 0 {
		item := heap.Pop(&tq).(*Item)
		gotItems = append(gotItems, item.Value)
	}
	assert.Equal(t, gotItems, expectedItems)
}

func TestDeleteFromQueue(t *testing.T) {
	tq := TTLQueue{}
	heap.Init(&tq)
	for key, val := range items {
		item := Item{Value: key, ExpirationTimestamp: val}
		heap.Push(&tq, &item)
	}

	err := tq.Delete("3")
	assert.Equal(t, err, nil)

	expectedItems := []string{"1", "2", "4", "5", "0"}
	var gotItems []string
	for tq.Len() > 0 {
		item := heap.Pop(&tq).(*Item)
		gotItems = append(gotItems, item.Value)
	}
	assert.Equal(t, gotItems, expectedItems)
}

func TestUpdateQueue(t *testing.T) {
	tq := TTLQueue{}
	heap.Init(&tq)
	for key, val := range items {
		item := Item{Value: key, ExpirationTimestamp: val}
		heap.Push(&tq, &item)
	}

	newTTL := int64(124)
	err := tq.Update("3", &newTTL)
	assert.Equal(t, err, nil)

	itemChanged, err := tq.find("3")
	assert.Equal(t, err, nil)
	assert.Equal(t, itemChanged.ExpirationTimestamp, &newTTL)

	expectedItems := []string{"1", "2", "4", "5", "3", "0"}
	var gotItems []string
	var item *Item
	for tq.Len() > 0 {
		item = heap.Pop(&tq).(*Item)
		gotItems = append(gotItems, item.Value)
	}
	assert.Equal(t, gotItems, expectedItems)
}
