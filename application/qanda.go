package main

import (
	"database/sql"
	"os/exec"
	"path/filepath"
)

func showQOrA(db *sql.DB, isQ bool) ierrori {

	var (
		e         error
		ie        ierrori
		rows      *sql.Rows
		thisError func(e error) ierrori
		unitPath  string
		qOrAPath  string
		files     []string
		file      string
		date      int64
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Could not ask a question", e: e}
	}

	rows, e = db.Query(`select path, min(date) from units`)
	if e != nil {
		return thisError(e)
	}
	defer rows.Close()

	if !rows.Next() {
		return thisError(nil)
	}

	e = rows.Scan(&unitPath, &date)
	if e != nil {
		return thisError(e)
	}

	if isQ {
		qOrAPath = filepath.Join(unitPath, questionsName)
	} else {
		qOrAPath = filepath.Join(unitPath, answersName)
	}

	files, ie = listFiles(qOrAPath)
	if ie != nil {
		return thisError(ie)
	}

	if len(files) < 1 {
		return thisError(nil)
	}

	for _, file = range files {
		_, e = exec.Command("cmd", "/c", "start", "", file).Output()
		if e != nil {
			return thisError(e)
		}
	}

	return nil
}
