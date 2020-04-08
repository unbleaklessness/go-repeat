package main

import "database/sql"

func createTables(db *sql.DB) ierrori {

	var e error

	_, e = db.Exec(`create table if not exists units (
		id integer primary key autoincrement,
		path text unique not null,
		date integer not null,
		stage integer not null
	)`)
	if e != nil {
		return ierror{m: "Could not create `units` table", e: e}
	}

	_, e = db.Exec(`create table if not exists defaults (
		type integer unique not null,
		command text not null,
		data blob not null,
		name text not null
	)`)
	if e != nil {
		return ierror{m: "Could not create `defaults` table", e: e}
	}

	return nil
}
