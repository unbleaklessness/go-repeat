package main

import "flag"

type flags struct {
	new    string
	rename string
	delete string
	rest   []string
}

func initializeFlags() flags {

	var (
		newFlag    *string
		renameFlag *string
		deleteFlag *string

		f flags
	)

	newFlag = flag.String("n", "", "New unit")
	renameFlag = flag.String("r", "", "Rename unit")
	deleteFlag = flag.String("d", "", "Delete unit")

	flag.Parse()

	f.new = *newFlag
	f.rename = *renameFlag
	f.delete = *deleteFlag
	f.rest = flag.Args()

	return f
}
