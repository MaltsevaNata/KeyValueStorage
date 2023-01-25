## Redis-like Key-value storage 
-  Supports methods set, get, delete.
-  Supports TTL

## Details
For effective implementation of deleting expired elements by TTL, priority queue is used based on https://pkg.go.dev/container/heap 

Expired items are deleted every `DeleteExpiredItemsPeriodMs` (500 ms by default)


## Local run

Example of storage usage is shown in ./cmd/app/main.py
```go run ./cmd/app```
