package main

import (
	"database/sql"
	"os"
	"path"
)

func newUnit(db *sql.DB, name string) ierrori {

	var (
		e         error
		p         string
		w         string
		thisError func(e error) ierrori
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Cound not create a new unit", e: e}
	}

	w, e = os.Getwd()
	if e != nil {
		return thisError(e)
	}
	p = path.Join(w, name)

	e = os.MkdirAll(p, os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	e = os.MkdirAll(path.Join(p, questionsName), os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	e = os.MkdirAll(path.Join(p, answersName), os.ModePerm)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`insert into units
		(path, date, score)
		values
		($1, $2, $3)
	`, p, now(), 0)
	if e != nil {
		return thisError(e)
	}

	return nil
}
