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
	} else if len(f.delete) > 0 {
		ie = deleteUnit(db, f.delete)
		if ie != nil {
			return ie
		}
	} else if len(f.move) > 0 && len(f.rest) > 0 {
		ie = moveUnit(db, f.move, f.rest[0])
		if ie != nil {
			return ie
		}
	} else if len(f.rename) > 0 && len(f.rest) > 0 {
		ie = renameUnit(db, f.rename, f.rest[0])
		if ie != nil {
			return ie
		}
	} else if f.question {
		ie = question(db)
		if ie != nil {
			return ie
		}
	} else {
		return ierror{m: "Unknown flag combination"}
	}

	return nil
}
