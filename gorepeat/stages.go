package main

var stages = []int64{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377}

func unixTimeForStage(stage int) int64 {
	return now() + stages[stage]*secondsInADay
}

func nextStage(stage int) int {
	stage++
	if stage >= len(stages) {
		stage--
	}
	return stage
}
