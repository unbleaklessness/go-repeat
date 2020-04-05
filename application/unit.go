package main

import (
	"database/sql"
	"os"
	"path"
)

func newUnit(db *sql.DB, name string) ierrori {

	var (
		e                error
		unitPath         string
		currentDirectory string
		thisError        func(e error) ierrori
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Cound not create a new unit", e: e}
	}

	currentDirectory, e = os.Getwd()
	if e != nil {
		return thisError(e)
	}
	unitPath = path.Join(currentDirectory, name)

	e = os.MkdirAll(unitPath, os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	e = os.MkdirAll(path.Join(unitPath, questionsName), os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	e = os.MkdirAll(path.Join(unitPath, answersName), os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`insert into units
		(path, date, score)
		values
		($1, $2, $3)
	`, unitPath, now(), 0)
	if e != nil {
		return thisError(e)
	}

	return nil
}
