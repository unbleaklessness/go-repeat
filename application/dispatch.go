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
		if f.questionText {
			ie = openDefault(db, f.new, true, textType)
			if ie != nil {
				return ie
			}
		}
		if f.questionImage {
			ie = openDefault(db, f.new, true, imageType)
			if ie != nil {
				return ie
			}
		}
		if f.answerText {
			ie = openDefault(db, f.new, false, textType)
			if ie != nil {
				return ie
			}
		}
		if f.answerImage {
			ie = openDefault(db, f.new, false, imageType)
			if ie != nil {
				return ie
			}
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
		ie = showQOrA(db, true)
		if ie != nil {
			return ie
		}
	} else if f.answer {
		ie = showQOrA(db, false)
		if ie != nil {
			return ie
		}
	} else if f.yes {
		ie = yesOrNo(db, true)
		if ie != nil {
			return ie
		}
	} else if f.no {
		ie = yesOrNo(db, false)
		if ie != nil {
			return ie
		}
	} else if len(f.defaultText) > 0 && len(f.rest) > 0 {
		ie = setDefault(db, textType, f.defaultText, f.rest[0])
		if ie != nil {
			return ie
		}
	} else if len(f.defaultImage) > 0 && len(f.rest) > 0 {
		ie = setDefault(db, imageType, f.defaultImage, f.rest[0])
		if ie != nil {
			return ie
		}
	} else {
		return ierror{m: "Unknown flag combination"}
	}

	return nil
}
