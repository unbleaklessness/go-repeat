package main

import "flag"

type flags struct {
	new            string
	question       bool
	answer         bool
	yes            bool
	no             bool
	addTemplate    bool
	deleteTemplate string
	questionAs     string
	answerAs       string
	questionInline string
	answerInline   string
	listTemplates  bool
	renameTemplate bool

	rest []string
}

func parseFlags() flags {

	newFlag := flag.String("n", "", "New unit")
	questionFlag := flag.Bool("q", false, "Show a question")
	answerFlag := flag.Bool("a", false, "Show an answer")
	yesFlag := flag.Bool("yes", false, "You answered correclty")
	noFlag := flag.Bool("no", false, "You answered incorreclty")
	addTemplateFlag := flag.Bool("add-template", false, "Add a template file")
	deleteTemplateFlag := flag.String("delete-template", "", "Delete template file")
	questionAsFlag := flag.String("q-as", "", "Create template file for question and open it")
	answerAsFlag := flag.String("a-as", "", "Create template file for answer and open it")
	questionInlineFlag := flag.String("q-is", "", "Inline question to template file")
	answerInlineFlag := flag.String("a-is", "", "Inline answer to template file")
	listTemplatesFlag := flag.Bool("list-templates", false, "List template files")
	renameTemplateFlag := flag.Bool("rename-template", false, "Rename a template")

	flag.Parse()

	f := flags{}

	f.new = *newFlag
	f.question = *questionFlag
	f.answer = *answerFlag
	f.yes = *yesFlag
	f.no = *noFlag
	f.addTemplate = *addTemplateFlag
	f.deleteTemplate = *deleteTemplateFlag
	f.questionAs = *questionAsFlag
	f.answerAs = *answerAsFlag
	f.questionInline = *questionInlineFlag
	f.answerInline = *answerInlineFlag
	f.listTemplates = *listTemplatesFlag
	f.renameTemplate = *renameTemplateFlag
	f.rest = flag.Args()

	return f
}
