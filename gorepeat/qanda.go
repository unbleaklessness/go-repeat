package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func checkUnitUnixTime(unixTime int64) ierrori {
	if unixTime > now() {
		message := fmt.Sprintf("No Q&A available, next one will be at %s", fromUnix(unixTime))
		return ierror{m: message}
	}
	return nil
}

func openQOrA(isQuestion bool) ierrori {

	units, ie := findUnits()
	if ie != nil {
		return ie
	}

	unit, ok := unitWithLeastUnixTime(units)
	if !ok {
		return ierror{m: "Could not find unit with least time"}
	}

	ie = checkUnitUnixTime(unit.data.UnixTime)
	if ie != nil {
		return ie
	}

	if isQuestion {
		for _, questionFilePath := range unit.questionFilePaths {
			e := open(questionFilePath)
			if e != nil {
				return ierror{m: "Could not open question file", e: e}
			}
		}
	} else {
		for _, answerFilePath := range unit.answerFilePaths {
			e := open(answerFilePath)
			if e != nil {
				return ierror{m: "Could not open answer file", e: e}
			}
		}
	}

	return nil
}

func yesOrNo(isYes bool) ierrori {

	units, ie := findUnits()
	if ie != nil {
		return ie
	}

	unit, ok := unitWithLeastUnixTime(units)
	if !ok {
		return ierror{m: "Could not find unit with least time"}
	}

	ie = checkUnitUnixTime(unit.data.UnixTime)
	if ie != nil {
		return ie
	}

	newStage := 0

	if isYes {
		newStage = nextStage(unit.data.Stage)
	}

	unit.data.UnixTime = unixTimeForStage(newStage)
	unit.data.Stage = newStage

	unitDataBytes, e := json.Marshal(unit.data)
	if e != nil {
		return ierror{m: "Could not update unit data file", e: e}
	}

	e = ioutil.WriteFile(unit.dataFilePath, unitDataBytes, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not update unit data file", e: e}
	}

	return nil
}
