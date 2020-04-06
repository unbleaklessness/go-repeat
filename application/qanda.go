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

func yesOrNo(db *sql.DB, isYes bool) ierrori {

	var (
		e         error
		rows      *sql.Rows
		date      int64
		newDate   int64
		stage     int
		newStage  int
		stages    []int64
		id        int
		thisError func(e error) ierrori
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Could not answer a question", e: e}
	}

	rows, e = db.Query(`select id, min(date), stage from units`)
	if e != nil {
		return thisError(e)
	}

	if !rows.Next() {
		return thisError(nil)
	}

	e = rows.Scan(&id, &date, &stage)
	if e != nil {
		return thisError(e)
	}

	stages = getStages()

	if stage < 0 || stage >= len(stages) {
		return thisError(nil)
	}

	if isYes {
		if stage >= len(stages)-1 {
			newStage = stage
		} else {
			newStage = stage + 1
		}
	} else {
		newStage = 0
	}

	newDate = now() + stages[newStage]*secondsInDay

	rows.Close()
	_, e = db.Exec(`update units set date = $1, stage = $2 where id = $3`, newDate, newStage, id)
	if e != nil {
		return thisError(e)
	}

	return nil
}
