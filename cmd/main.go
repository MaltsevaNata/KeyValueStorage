package main

import (
	"KeyValueStorage/pkg/storage"
	"log"
	"time"
)

func main() {
	kvStorage := storage.New()

	// set value
	ttl := 2
	firstKey := "foo"
	firstVal := "bar"
	err := kvStorage.Set(firstKey, firstVal, &ttl)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Set {key: `%s`, value: `%s`} with TTL %d seconds\n", firstKey, firstVal, ttl)
	err = nil

	// get value
	val, err := kvStorage.Get(firstKey)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Got {key: `%s`, value: `%s`}", firstKey, val)
	err = nil

	// check deleting by TTL
	log.Println("Waiting for TTL to expire...")
	time.Sleep(3 * time.Second) // wait 1 sec more to let item be removed
	_, err = kvStorage.Get(firstKey)
	if err != nil {
		log.Printf("First value was deleted by TTL: %s \n", err)
	} else {
		log.Fatalf("Error: TTL didn't work")
	}
	err = nil

	// set new value
	secondKey := "abcde123"
	secondVal := "87443643209"
	err = kvStorage.Set(secondKey, secondVal, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Set {key: `%s`, value: `%s`} with no TTL", secondKey, secondVal)
	val, err = kvStorage.Get(secondKey)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Got {key: `%s`, value: `%s`}", secondKey, val)
	err = nil

	// delete value
	err = kvStorage.Delete(secondKey)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Deleted value by {key: `%s`}", secondKey)
	err = nil

	// check deleted
	_, err = kvStorage.Get(secondKey)
	if err != nil {
		log.Printf("Second value was deleted by user request: %s \n", err)
	} else {
		log.Fatalf("Error: Deleting didn't work")
	}
}
