package main

import "database/sql"

func createTables(db *sql.DB) ierrori {

	var (
		e         error
		thisError func(e error) ierrori
	)

	thisError = func(e error) ierrori {
		return ierror{m: "Could not create tables", e: e}
	}

	_, e = db.Exec(`create table if not exists units (
		id integer primary key autoincrement,
		path text unique not null,
		date integer not null,
		stage integer not null
	)`)
	if e != nil {
		return thisError(e)
	}

	_, e = db.Exec(`create table if not exists defaults (
		type integer unique not null,
		command text not null,
		data blob not null,
		name text not null
	)`)
	if e != nil {
		return thisError(e)
	}

	return nil
}
