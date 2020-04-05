package main

func stringsContains(xs []string, x string) bool {

	var (
		s string
	)

	for _, s = range xs {
		if s == x {
			return true
		}
	}

	return false
}

func stringsIndex(xs []string, x string) int {

	var (
		s string
		i int
	)

	for i, s = range xs {
		if s == x {
			return i
		}
	}

	return -1
}
