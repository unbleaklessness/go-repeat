package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	questionDirectoryName      = "Q"
	answerDirectoryName        = "A"
	unitDataFileName           = "D.json"
	logFileName                = "log.txt"
	templatesFileName          = "templates.json"
	configurationDirectoryName = ".go-repeat"
)

type template struct {
	Name     string
	FileName string
	Bytes    []byte
}

func main() {

	flags := parseFlags()

	configurationDirectoryPath, ie := getConfigurationDirectoryPath()
	if ie != nil {
		panic(ie.Message())
	}

	e := os.MkdirAll(configurationDirectoryPath, os.ModePerm)
	if e != nil {
		panic("Could not create project configuration directory")
	}

	templatesFilePath, ie := getTemplatesFilePath()
	if ie != nil {
		panic(ie.Message())
	}

	templateFile, e := os.OpenFile(templatesFilePath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if e != nil {
		panic("Could not create templates file")
	}
	templateFile.Close()

	logFileAPath := filepath.Join(configurationDirectoryPath, logFileName)

	logFile, e := os.OpenFile(logFileAPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)

	logger := log.New(logFile, "", log.Ldate|log.Ltime)

	ie = dispatch(flags)
	if ie != nil {
		fmt.Println(ie.Message())
		logger.Println(ie.Error())
	}
}
