package utils

import "time"

func GetTimestamp() int64 {
	now := time.Now()
	return now.UnixMilli()
}
