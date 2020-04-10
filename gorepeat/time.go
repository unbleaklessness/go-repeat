package main

import "time"

const (
	secondsInDay  = 86400
	secondsInHour = 3600
	timeLayout    = "2006-01-02 15:04:05"
)

func now() int64 {
	return time.Now().Unix()
}

func fromUnix(unixTime int64) string {
	return time.Unix(unixTime, 0).Format(timeLayout)
}
