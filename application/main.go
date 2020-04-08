package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	questionDirectoryName      = "Q"
	answerDirectoryName        = "A"
	databaseFileName           = "Data.sqlite"
	logFileName                = "log.txt"
	configurationDirectoryName = ".go-repeat"
	rootConfigurationFileName  = "root.txt"
)

func createConfigurationDirectory() ierrori {

	var (
		configurationDirectoryPath string
		ie                         ierrori
		e                          error
	)

	configurationDirectoryPath, ie = getConfigurationDirectoryPath()
	if ie != nil {
		return ie
	}

	e = os.MkdirAll(configurationDirectoryPath, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not create configuration directory"}
	}

	return nil
}

func getDatabase() (*sql.DB, ierrori) {

	var (
		e                error
		ie               ierrori
		db               *sql.DB
		databaesFilePath string
	)

	databaesFilePath, ie = getDatabaseFilePath()
	if ie != nil {
		return db, ie
	}

	db, e = sql.Open("sqlite3", databaesFilePath)
	if e != nil {
		return db, ierror{m: "Could not open database", e: e}
	}

	return db, nil
}

func getLogger() (*log.Logger, ierrori) {

	var (
		logger      *log.Logger
		ie          ierrori
		logFilePath string
		logFile     *os.File
		e           error
	)

	logFilePath, ie = getLogFilePath()
	if ie != nil {
		return logger, ie
	}

	logFile, e = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if e != nil {
		return logger, ierror{m: "Could not create log file", e: e}
	}

	logger = log.New(logFile, "", log.Ldate|log.Ltime)

	return logger, nil
}

func setRoot(f flags) ierrori {

	var (
		rootConfigurationFilePath string
		ie                        ierrori
		e                         error
		rootConfigurationFile     *os.File
	)

	if len(f.setRoot) < 1 {
		return nil
	}

	rootConfigurationFilePath, ie = getRootConfigurationFilePath()
	if ie != nil {
		return ie
	}

	if fileExists(rootConfigurationFilePath) {
		e = os.Remove(rootConfigurationFilePath)
		if e != nil {
			return ierror{m: "Could not remove old root configuration file", e: e}
		}
	}

	rootConfigurationFile, e = os.Create(rootConfigurationFilePath)
	if e != nil {
		return ierror{m: "Could not create a new root configuration file", e: e}
	}

	_, e = rootConfigurationFile.WriteString(f.setRoot)
	if e != nil {
		return ierror{m: "Could not save a new root directory", e: e}
	}

	e = rootConfigurationFile.Close()
	if e != nil {
		return ierror{m: "Could not close root configuration file", e: e}
	}

	return nil
}

func main() {

	var (
		ie     ierrori
		flags  flags
		db     *sql.DB
		logger *log.Logger
	)

	flags = initializeFlags()

	ie = createConfigurationDirectory()
	if ie != nil {
		fmt.Println(ie.Message())
		return
	}

	ie = setRoot(flags)
	if ie != nil {
		fmt.Println(ie.Message())
		return
	}

	db, ie = getDatabase()
	if ie != nil {
		fmt.Println(ie.Message())
		return
	}

	logger, ie = getLogger()
	if ie != nil {
		fmt.Println(ie.Message())
		return
	}

	ie = createTables(db)
	if ie != nil {
		fmt.Println(ie.Message())
		return
	}

	ie = dispatch(db, flags)
	if ie != nil {
		fmt.Println(ie.Message())
		logger.Println(ie.Error())
	}
}
