package main

import "database/sql"

func createTables(db *sql.DB) ierrori {

	var e error

	_, e = db.Exec(`create table if not exists units (
		id integer primary key autoincrement,
		path text unique not null,
		date integer not null,
		score real not null,
		stage integer not null
	)`)
	if e != nil {
		return ierror{m: "Could not create tables", e: e}
	}

	return nil
}
