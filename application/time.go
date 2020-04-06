package main

import "time"

const (
	secondsInDay = 86400
)

func now() int64 {
	return time.Now().Unix()
}
