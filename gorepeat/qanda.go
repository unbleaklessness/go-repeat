package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func unitUnixTimeIsInFuture(unit unit) ierrori {
	unixTime := unit.data.getUnixTime()
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

	ie = unitUnixTimeIsInFuture(unit)
	if ie != nil {
		return ie
	}

	if unit.data.getInverse() {
		isQuestion = !isQuestion
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

	ie = unitUnixTimeIsInFuture(unit)
	if ie != nil {
		return ie
	}

	if isYes {
		unit.data.nextStage()
	} else {
		unit.data.previousStage()
	}

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
