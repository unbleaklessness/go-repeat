package main

import (
	"database/sql"
)

func dispatch(db *sql.DB, f flags) ierrori {

	var ie ierrori

	if len(f.new) > 0 {
		ie = newUnit(db, f.new)
		if ie != nil {
			return ie
		}
	} else {
		// Question user.
	}

	return nil
}