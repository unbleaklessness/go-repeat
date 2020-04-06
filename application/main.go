package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	questionsName        = "Q"
	answersName          = "A"
	databaseName         = "data.sqlite"
	logFileName          = "log.txt"
	projectDirectoryName = ".go-repeat"
)

func main() {

	var (
		e            error
		ie           ierrori
		flags        flags
		db           *sql.DB
		logger       *log.Logger
		file         *os.File
		databasePath string
		projectPath  string
		logPath string
		home         string
	)

	flags = initializeFlags()

	home, e = os.UserHomeDir()
	if e != nil {
		panic("Could not get user's home directory")
	}

	projectPath = filepath.Join(home, projectDirectoryName)
	e = os.MkdirAll(projectPath, os.ModePerm)
	if e != nil {
		panic("Could not create project directory in user's home")
	}

	databasePath = filepath.Join(projectPath, databaseName)
	logPath = filepath.Join(projectPath, logFileName)

	db, e = sql.Open("sqlite3", databasePath)
	if e != nil {
		panic("Could not open database")
	}

	e = createTables(db)
	if e != nil {
		panic("Could not create database tables")
	}

	file, e = os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if e != nil {
		panic("Could not create a log file")
	}
	logger = log.New(file, "", log.Ldate|log.Ltime)

	ie = dispatch(db, flags)
	if ie != nil {
		fmt.Println(ie.Message())
		logger.Println(ie.Error())
	}
}
