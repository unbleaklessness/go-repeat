package main

import (
	"bufio"
	"database/sql"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type fileType = int

const (
	imageType fileType = iota + 1
	textType
)

func setDefault(db *sql.DB, t fileType, fileName string, command string) ierrori {

	var (
		currentDirectory string
		filePath         string
		data             []byte
		thisError        func(e error) ierrori
		e                error
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Could not set default file and command", e: e}
	}

	currentDirectory, e = os.Getwd()
	if e != nil {
		return thisError(e)
	}
	filePath = filepath.Join(currentDirectory, fileName)

	data, e = ioutil.ReadFile(filePath)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`delete from defaults where type = $1`, t)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`insert into defaults
		(type, command, data, name)
		values
		($1, $2, $3, $4)`, t, command, data, fileName)
	if e != nil {
		return thisError(e)
	}

	return nil
}

func openDefault(db *sql.DB, unitName string, isQ bool, t fileType) ierrori {

	var (
		e                error
		thisError        func(e error) ierrori
		rows             *sql.Rows
		data             []byte
		command          string
		currentDirectory string
		filePath         string
		fileName         string
		file             *os.File
		writer           *bufio.Writer
		commandSplit     []string
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Could not open default file", e: e}
	}

	rows, e = db.Query(`select command, data, name from defaults where type = $1`, t)
	if e != nil {
		return thisError(e)
	}
	defer rows.Close()

	if !rows.Next() {
		return thisError(nil)
	}

	e = rows.Scan(&command, &data, &fileName)
	if e != nil {
		return thisError(e)
	}

	if len(command) < 1 {
		return thisError(nil)
	}

	currentDirectory, e = os.Getwd()
	if e != nil {
		return thisError(e)
	}

	if isQ {
		filePath = filepath.Join(currentDirectory, unitName, questionsName, fileName)
	} else {
		filePath = filepath.Join(currentDirectory, unitName, answersName, fileName)
	}

	file, e = os.Create(filePath)
	if e != nil {
		return thisError(e)
	}

	writer = bufio.NewWriter(file)

	_, e = writer.Write(data)
	if e != nil {
		return thisError(e)
	}
	e = writer.Flush()
	if e != nil {
		return thisError(e)
	}
	file.Close()

	commandSplit = strings.Split(command, " ")
	commandSplit = append(commandSplit, filePath)

	e = exec.Command(commandSplit[0], commandSplit[1:]...).Start()
	if e != nil {
		return thisError(e)
	}

	return nil
}
