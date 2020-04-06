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

func moveUnit(db *sql.DB, name string, newPath string) ierrori {

	var (
		e           error
		current     string
		unitPath    string
		newUnitPath string
		thisError   func(e error) ierrori
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Could not move unit", e: e}
	}

	unitPath = filepath.Join(current, name)

	if !filepath.IsAbs(newPath) {
		current, e = os.Getwd()
		if e != nil {
			return thisError(e)
		}
		newUnitPath = filepath.Join(current, newPath, name)
	} else {
		newUnitPath = filepath.Join(newPath, name)
	}

	if !directoryExists(newPath) {
		e = os.MkdirAll(newPath, os.ModePerm)
		if e != nil {
			return thisError(e)
		}
	}

	e = os.Rename(unitPath, newUnitPath)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`update units set path = $1 where path = $2`, newUnitPath, unitPath)
	if e != nil {
		return thisError(e)
	}

	return nil
}

func renameUnit(db *sql.DB, oldName string, newName string) ierrori {

	var (
		e           error
		current     string
		unitPath    string
		newUnitPath string
		thisError   func(e error) ierrori
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Could not rename unit", e: e}
	}

	current, e = os.Getwd()
	if e != nil {
		return thisError(e)
	}

	unitPath = filepath.Join(current, oldName)
	newUnitPath = filepath.Join(current, newName)

	e = os.Rename(unitPath, newUnitPath)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`update units set path = $1 where path = $2`, newUnitPath, unitPath)
	if e != nil {
		return thisError(e)
	}

	return nil
}
