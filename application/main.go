package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	questionsName    = "Q"
	answersName      = "A"
	databaseFilePath = "data.sqlite"
	logFilePath      = "log.txt"
)

func main() {

	var (
		e      error
		ie     ierrori
		flags  flags
		db     *sql.DB
		logger *log.Logger
		file   *os.File
	)

	flags = initializeFlags()

	db, e = sql.Open("sqlite3", databaseFilePath)
	if e != nil {
		panic("Error opening database")
	}

	e = createTables(db)
	if e != nil {
		panic("Could not create database tables")
	}

	file, e = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
