package main

import "flag"

type flags struct {
	new    string
	move   string
	delete string
	rest   []string
}

func initializeFlags() flags {

	var (
		newFlag    *string
		moveFlag   *string
		deleteFlag *string

		f flags
	)

	newFlag = flag.String("n", "", "New unit")
	moveFlag = flag.String("m", "", "Move unit")
	deleteFlag = flag.String("d", "", "Delete unit")

	flag.Parse()

	f.new = *newFlag
	f.move = *moveFlag
	f.delete = *deleteFlag
	f.rest = flag.Args()

	return f
}
