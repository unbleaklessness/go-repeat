package main

import "flag"

type flags struct {
	new           string
	move          string
	rename        string
	delete        string
	question      bool
	answer        bool
	yes           bool
	no            bool
	questionText  bool
	questionImage bool
	answerText    bool
	answerImage   bool
	defaultText   string
	defaultImage  string
	setRoot       string
	rest          []string
}

func initializeFlags() flags {

	var (
		newFlag           *string
		moveFlag          *string
		renameFlag        *string
		deleteFlag        *string
		questionFlag      *bool
		answerFlag        *bool
		yesFlag           *bool
		noFlag            *bool
		questionTextFlag  *bool
		questionImageFlag *bool
		answerTextFlag    *bool
		answerImageFlag   *bool
		defaultTextFlag   *string
		defaultImageFlag  *string
		setRootFlag       *string

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
	questionTextFlag = flag.Bool("q-text", false, "Open text editor for question")
	questionImageFlag = flag.Bool("q-image", false, "Open image editor for question")
	answerTextFlag = flag.Bool("a-text", false, "Open text editor for answer")
	answerImageFlag = flag.Bool("a-image", false, "Open image editor for answer")
	defaultTextFlag = flag.String("default-text", "", "Set default text file and command")
	defaultImageFlag = flag.String("default-image", "", "Set default image file and command")
	setRootFlag = flag.String("set-root", "", "Set root directory")

	flag.Parse()

	f.new = *newFlag
	f.move = *moveFlag
	f.rename = *renameFlag
	f.delete = *deleteFlag
	f.question = *questionFlag
	f.answer = *answerFlag
	f.yes = *yesFlag
	f.no = *noFlag
	f.questionText = *questionTextFlag
	f.questionImage = *questionImageFlag
	f.answerText = *answerTextFlag
	f.answerImage = *answerImageFlag
	f.defaultText = *defaultTextFlag
	f.defaultImage = *defaultImageFlag
	f.setRoot = *setRootFlag
	f.rest = flag.Args()

	return f
}
