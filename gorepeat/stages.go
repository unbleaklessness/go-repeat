package main

var stages = []int64{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765, 10946, 17711, 28657, 46368}

func unixTimeForStage(stage int) int64 {
	return now() + stages[stage]*secondsInDay
}

func inverseUnixTimeForStage(stage int) int64 {
	next := unixTimeForStage(stage)
	if stage > 0 {
		return next + secondsInHour
	}
	return next
}
