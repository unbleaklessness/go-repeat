package main

import "flag"

type flags struct {
	new      string
	move     string
	rename   string
	delete   string
	question bool
	answer   bool
	yes      bool
	no       bool
	rest     []string
}

func initializeFlags() flags {

	var (
		newFlag      *string
		moveFlag     *string
		renameFlag   *string
		deleteFlag   *string
		questionFlag *bool
		answerFlag   *bool
		yesFlag      *bool
		noFlag       *bool

		f flags
	)

	newFlag = flag.String("n", "", "New unit")
	moveFlag = flag.String("m", "", "Move unit")
	renameFlag = flag.String("r", "", "Rename unit")
	deleteFlag = flag.String("d", "", "Delete unit")
	questionFlag = flag.Bool("q", false, "Show a question")
	answerFlag = flag.Bool("a", false, "Show an answer")
	yesFlag = flag.Bool("yes", false, "You answered correclty")
	noFlag = flag.Bool("no", false, "You answered incorreclty")

	flag.Parse()

	f.new = *newFlag
	f.move = *moveFlag
	f.rename = *renameFlag
	f.delete = *deleteFlag
	f.question = *questionFlag
	f.answer = *answerFlag
	f.yes = *yesFlag
	f.no = *noFlag
	f.rest = flag.Args()

	return f
}
