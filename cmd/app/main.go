package main

import (
	"KeyValueStorage/pkg/storage"
	"log"
	"time"
)

func main() {
	kvStorage := storage.New()
	items := map[string]string{"foo": "111", "bar": "222", "fdrkj": "333", "rtrn": "444", "ufecs": "555"}
	ttl := 2
	itemTTL := map[string]*int{"foo": &ttl, "bar": nil, "fdrkj": &ttl, "rtrn": nil, "ufecs": nil}

	// set value
	for key, val := range items {
		err := kvStorage.Set(key, val, itemTTL[key])
		if err != nil {
			log.Fatalf(err.Error())
		}
		if itemTTL[key] != nil {
			log.Printf("Set {key: `%s`, value: `%s`} with TTL %d seconds\n", key, val, *itemTTL[key])
		} else {
			log.Printf("Set {key: `%s`, value: `%s`} with no TTL\n", key, val)
		}
	}

	// get value
	val, err := kvStorage.Get("foo")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Got {key: `%s`, value: `%s`}", "foo", val)
	err = nil

	// check deleting by TTL
	log.Println("Waiting for TTL to expire...")
	time.Sleep(3 * time.Second) // wait 1 sec more to let items with TTL be removed

	_, err1 := kvStorage.Get("foo")
	_, err2 := kvStorage.Get("fdrkj")
	if err1 != nil && err2 != nil {
		log.Println("Values were deleted by TTL")
	} else {
		log.Fatalf("Error: TTL didn't work")
	}

	// get value without ttl
	val, err = kvStorage.Get("bar")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Got {key: `%s`, value: `%s`}", "bar", val)
	err = nil

	// delete value
	err = kvStorage.Delete("bar")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Deleted value by {key: `%s`}", "bar")
	err = nil

	// check deleted
	_, err = kvStorage.Get("bar")
	if err != nil {
		log.Printf("Value was deleted by user request: %s \n", err)
	} else {
		log.Fatalf("Error: Deleting didn't work")
	}
}
