package main

import (
	"fmt"
	"log"
	"os"
)

const (
	aDirectoryName             = "A"
	bDirectoryName             = "B"
	unitDataFileName           = "D.json"
	logFileName                = "log.txt"
	templatesFileName          = "templates.json"
	defaultInlineFileName      = "inline.txt"
	configurationDirectoryName = ".go-repeat"
)

func createLogger() (*log.Logger, ierrori) {

	logFilePath, ie := getLogFilePath()
	if ie != nil {
		return nil, ie
	}

	logFile, e := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if e != nil {
		return nil, ierror{m: "Could not create log file", e: e}
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime)

	return logger, nil
}

func createConfigurationDirectory() ierrori {

	configurationDirectoryPath, ie := getConfigurationDirectoryPath()
	if ie != nil {
		return ie
	}

	e := os.MkdirAll(configurationDirectoryPath, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not create project configuration directory", e: e}
	}

	return nil
}

func main() {

	flags := parseFlags()

	ie := createConfigurationDirectory()
	if ie != nil {
		panic(ie.Message())
	}

	ie = createTemplates()
	if ie != nil {
		panic(ie.Message())
	}

	logger, ie := createLogger()
	if ie != nil {
		panic(ie.Message())
	}

	ie = dispatch(flags)
	if ie != nil {
		fmt.Println(ie.Message())
		logger.Println(ie.Error())
	}
}
