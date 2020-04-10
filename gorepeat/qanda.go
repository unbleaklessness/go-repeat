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

func openQuestionOrAnswer(isQuestion bool) ierrori {

	units, ie := findUnits()
	if ie != nil {
		return ie
	}

	unit, ok := unitWithLeastUnixTime(units)
	if !ok {
		return ierror{m: "Could not find unit with appropriate time"}
	}

	ie = unitUnixTimeIsInFuture(unit)
	if ie != nil {
		return ie
	}

	aOrB := isQuestion

	if unit.data.getInverse() {
		aOrB = !aOrB
	}

	if aOrB {
		for _, aFilePath := range unit.aFilePaths {
			e := open(aFilePath)
			if e != nil {
				return ierror{m: "Could not open A association file", e: e}
			}
		}
	} else {
		for _, bFilePath := range unit.bFilePaths {
			e := open(bFilePath)
			if e != nil {
				return ierror{m: "Could not open B association file", e: e}
			}
		}
	}

	return nil
}

func answer(correct bool) ierrori {

	units, ie := findUnits()
	if ie != nil {
		return ie
	}

	unit, ok := unitWithLeastUnixTime(units)
	if !ok {
		return ierror{m: "Could not find unit with appropriate time"}
	}

	ie = unitUnixTimeIsInFuture(unit)
	if ie != nil {
		return ie
	}

	if correct {
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
