package ttlQueue

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushToQueue(t *testing.T) {
	items := map[string]int{"5": 5, "2": 2, "4": 4, "3": 3, "1": 1}
	tq := TTLQueue{}
	heap.Init(&tq)
	for key, val := range items {
		ts := int64(val)
		item := Item{Value: key, ExpirationTimestamp: &ts}
		heap.Push(&tq, &item)
	}
	// Take the items out; they arrive in increasing ts order.
	expectedItems := []string{"1", "2", "3", "4", "5"}
	var gotItems []string
	for tq.Len() > 0 {
		item := heap.Pop(&tq).(*Item)
		gotItems = append(gotItems, item.Value)
	}
	assert.Equal(t, gotItems, expectedItems)
}

func TestPopFromQueue(t *testing.T) {
	items := map[string]int{"5": 5, "2": 2, "4": 4, "3": 3, "1": 1}
	tq := TTLQueue{}
	heap.Init(&tq)
	for key, val := range items {
		ts := int64(val)
		item := Item{Value: key, ExpirationTimestamp: &ts}
		heap.Push(&tq, &item)
	}
	ts := int64(1)
	expectedItem := Item{Value: "1", ExpirationTimestamp: &ts, index: -1}
	item := heap.Pop(&tq).(*Item)
	assert.Equal(t, *item, expectedItem)

	expectedItems := []string{"2", "3", "4", "5"}
	var gotItems []string
	for tq.Len() > 0 {
		item := heap.Pop(&tq).(*Item)
		gotItems = append(gotItems, item.Value)
	}
	assert.Equal(t, gotItems, expectedItems)
}

func TestDeleteFromQueue(t *testing.T) {
	items := map[string]int{"5": 5, "2": 2, "4": 4, "3": 3, "1": 1}
	tq := TTLQueue{}
	heap.Init(&tq)
	for key, val := range items {
		ts := int64(val)
		item := Item{Value: key, ExpirationTimestamp: &ts}
		heap.Push(&tq, &item)
	}

	err := tq.Delete("3")
	assert.Equal(t, err, nil)

	expectedItems := []string{"1", "2", "4", "5"}
	var gotItems []string
	for tq.Len() > 0 {
		item := heap.Pop(&tq).(*Item)
		gotItems = append(gotItems, item.Value)
	}
	assert.Equal(t, gotItems, expectedItems)
}

func TestUpdateQueue(t *testing.T) {
	items := map[string]int{"5": 5, "2": 2, "4": 4, "3": 3, "1": 1}
	tq := TTLQueue{}
	heap.Init(&tq)
	for key, val := range items {
		ts := int64(val)
		item := Item{Value: key, ExpirationTimestamp: &ts}
		heap.Push(&tq, &item)
	}

	newTTL := int64(124)
	err := tq.Update("3", &newTTL)
	assert.Equal(t, err, nil)

	expectedItems := []string{"1", "2", "4", "5", "3"}
	var gotItems []string
	var item *Item
	for tq.Len() > 0 {
		item = heap.Pop(&tq).(*Item)
		gotItems = append(gotItems, item.Value)
	}
	assert.Equal(t, gotItems, expectedItems)
	assert.Equal(t, item.ExpirationTimestamp, &newTTL)
}
