package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	questionsDirectoryName = "questions"
	answersDirectoryName   = "answers"
	informationFileName    = "information.txt"
)

func usage() {
	fmt.Println("Please, provide a directory as the first argument!")
}

func main() {

	var (
		e         error
		arguments []string
		workPath  string
		units     []unit
	)

	arguments = os.Args

	if len(arguments) < 2 {
		usage()
		return
	}

	workPath = arguments[1]
	workPath, e = filepath.Abs(workPath)
	if e != nil {
		panic("Could not get current directory")
	}

	units = findUnits(workPath)

	var u unit
	for _, u = range units {
		fmt.Printf("%+v \n\n", u)
	}
}
