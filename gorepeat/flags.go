package main

import "flag"

type flags struct {
	new              string
	question         bool
	answer           bool
	yes              bool
	no               bool
	addTemplate      bool
	deleteTemplate   string
	AAs              string
	BAs              string
	AInline          string
	BInline          string
	listTemplates    bool
	renameTemplate   bool
	inverse          bool
	setDefaultInline string

	rest []string
}

func parseFlags() flags {

	newFlag := flag.String("n", "", "New unit")
	questionFlag := flag.Bool("q", false, "Show a question")
	answerFlag := flag.Bool("a", false, "Show an answer")
	yesFlag := flag.Bool("yes", false, "Your answer is correct")
	noFlag := flag.Bool("no", false, "Your answer is incorrect")
	addTemplateFlag := flag.Bool("add-template", false, "Add a template")
	deleteTemplateFlag := flag.String("delete-template", "", "Delete a template")
	AAsFlag := flag.String("a-as", "", "Create template file for A association and open it")
	BAsFlag := flag.String("b-as", "", "Create template file for B association and open it")
	AInlineFlag := flag.String("a-is", "", "Inline A association")
	BInlineFlag := flag.String("b-is", "", "Inline B association")
	listTemplatesFlag := flag.Bool("list-templates", false, "List templates")
	renameTemplateFlag := flag.Bool("rename-template", false, "Rename a template")
	inverseFlag := flag.Bool("i", false, "Include inverse Q&A")
	setDefaultInlineFlag := flag.String("set-default-inline", "", "Set default template for inlining")

	flag.Parse()

	f := flags{}

	f.new = *newFlag
	f.question = *questionFlag
	f.answer = *answerFlag
	f.yes = *yesFlag
	f.no = *noFlag
	f.addTemplate = *addTemplateFlag
	f.deleteTemplate = *deleteTemplateFlag
	f.AAs = *AAsFlag
	f.BAs = *BAsFlag
	f.AInline = *AInlineFlag
	f.BInline = *BInlineFlag
	f.listTemplates = *listTemplatesFlag
	f.renameTemplate = *renameTemplateFlag
	f.inverse = *inverseFlag
	f.setDefaultInline = *setDefaultInlineFlag
	f.rest = flag.Args()

	return f
}
