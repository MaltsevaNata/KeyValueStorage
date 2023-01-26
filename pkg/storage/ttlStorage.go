package storage

import (
	"KeyValueStorage/internal/config"
	"KeyValueStorage/internal/ttlQueue"
	"KeyValueStorage/internal/utils"
	"container/heap"
	"errors"
	"sync"
	"time"
)

type TTLStorage struct {
	storage  map[string]string
	ttlQueue ttlQueue.TTLQueue
	lock     sync.RWMutex
}

func New() (storage *TTLStorage) {
	tq := ttlQueue.TTLQueue{}
	heap.Init(&tq)
	storage = &TTLStorage{storage: make(map[string]string), ttlQueue: tq}
	go func() {
		for range time.Tick(time.Millisecond * config.DeleteExpiredItemsPeriodMs) {
			storage.removeExpired()
		}
	}()
	return
}

func (storage *TTLStorage) removeExpired() {
	tsNow := utils.GetTimestamp()
	var oldestItem *ttlQueue.Item
	storage.lock.Lock()
	for {
		oldestItem = storage.ttlQueue.Peek()
		if oldestItem != nil && oldestItem.ExpirationTimestamp != nil && tsNow-*oldestItem.ExpirationTimestamp > 0 {
			delete(storage.storage, oldestItem.Value)
			heap.Pop(&storage.ttlQueue)
		} else {
			break // no expired items remained
		}
	}
	storage.lock.Unlock()
}

func (storage *TTLStorage) Set(key, val string, ttlSec *int) (err error) {
	storage.lock.Lock()
	var expirationPtrMs *int64
	if ttlSec != nil {
		nowTsMs := utils.GetTimestamp()
		expirationMs := nowTsMs + int64(*ttlSec*1000)
		expirationPtrMs = &expirationMs
	}
	_, ok := storage.storage[key]
	if !ok {
		queueItem := ttlQueue.Item{Value: key, ExpirationTimestamp: expirationPtrMs}
		storage.storage[key] = val
		heap.Push(&storage.ttlQueue, &queueItem)
	} else {
		err = storage.ttlQueue.Update(key, expirationPtrMs)
	}
	storage.lock.Unlock()
	return err
}

func (storage *TTLStorage) Get(key string) (val string, err error) {
	storage.lock.Lock()
	var ok bool
	if val, ok = storage.storage[key]; !ok {
		err = errors.New("key not found")
	}
	storage.lock.Unlock()
	return val, err
}

func (storage *TTLStorage) Delete(key string) (err error) {
	storage.lock.Lock()
	if _, ok := storage.storage[key]; !ok {
		err = errors.New("key not found")
	} else {
		delete(storage.storage, key)
		err = storage.ttlQueue.Delete(key)
	}
	storage.lock.Unlock()
	return err
}

func (storage *TTLStorage) GetAllItems() map[string]string {
	return storage.storage
}
