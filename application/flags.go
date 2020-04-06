package main

import "flag"

type flags struct {
	new      string
	move     string
	rename   string
	delete   string
	question bool
	rest     []string
}

func initializeFlags() flags {

	var (
		newFlag      *string
		moveFlag     *string
		renameFlag   *string
		deleteFlag   *string
		questionFlag *bool

		f flags
	)

	newFlag = flag.String("n", "", "New unit")
	moveFlag = flag.String("m", "", "Move unit")
	renameFlag = flag.String("r", "", "Rename unit")
	deleteFlag = flag.String("d", "", "Delete unit")
	questionFlag = flag.Bool("q", false, "Show a question")

	flag.Parse()

	f.new = *newFlag
	f.move = *moveFlag
	f.rename = *renameFlag
	f.delete = *deleteFlag
	f.question = *questionFlag
	f.rest = flag.Args()

	return f
}
