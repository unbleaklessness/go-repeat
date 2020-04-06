package main

import (
	"database/sql"
	"os"
	"path/filepath"
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
	unitPath = filepath.Join(currentDirectory, name)

	e = os.MkdirAll(unitPath, os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	e = os.MkdirAll(filepath.Join(unitPath, questionsName), os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	e = os.MkdirAll(filepath.Join(unitPath, answersName), os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`insert into units
		(path, date, score, stage)
		values
		($1, $2, $3, $4)
	`, unitPath, now(), 0, 1)
	if e != nil {
		return thisError(e)
	}

	return nil
}

func deleteUnit(db *sql.DB, name string) ierrori {

	var (
		e         error
		current   string
		unitPath  string
		thisError func(e error) ierrori
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Could not delete unit", e: e}
	}

	current, e = os.Getwd()
	if e != nil {
		return thisError(e)
	}

	unitPath = filepath.Join(current, name)

	e = os.RemoveAll(unitPath)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`delete from units where path = $1`, unitPath)
	if e != nil {
		return thisError(e)
	}

	return nil
}
