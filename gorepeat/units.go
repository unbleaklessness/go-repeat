package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

type unitData struct {
	UnixTime int64
	Stage    int
}

type unit struct {
	path                  string
	dataFilePath          string
	data                  unitData
	questionFilePaths     []string
	answerFilePaths       []string
	questionDirectoryPath string
	answerDirectoryPath   string
}

func newUnit(unitDirectoryPath string) ierrori {

	unitDirectoryPath = filepath.Clean(unitDirectoryPath)

	e := os.MkdirAll(unitDirectoryPath, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not create unit directory", e: e}
	}

	e = os.MkdirAll(filepath.Join(unitDirectoryPath, questionDirectoryName), os.ModePerm)
	if e != nil {
		return ierror{m: "Could not create unit question directory", e: e}
	}

	e = os.MkdirAll(filepath.Join(unitDirectoryPath, answerDirectoryName), os.ModePerm)
	if e != nil {
		return ierror{m: "Could not create unit answer directory", e: e}
	}

	stage := 0

	unitData := unitData{
		UnixTime: unixTimeForStage(stage),
		Stage:    stage,
	}

	unitDataBytes, e := json.Marshal(unitData)

	unitDataFilePath := filepath.Join(unitDirectoryPath, unitDataFileName)

	e = ioutil.WriteFile(unitDataFilePath, unitDataBytes, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not write unit data file", e: e}
	}

	return nil
}

func findUnits() ([]unit, ierrori) {

	currentDirectoryPath, e := os.Getwd()
	if e != nil {
		return nil, ierror{m: "Could not get current directory", e: e}
	}

	var inner func(p string) []unit

	inner = func(p string) []unit {

		units := make([]unit, 0)

		directoryPaths, ie := listDirectories(p)
		if ie != nil {
			return units
		}

		findSubUnits := func(p string) []unit {
			return append(units, inner(p)...)
		}

		for _, directoryPath := range directoryPaths {

			questionDirectoryPath := filepath.Join(directoryPath, questionDirectoryName)
			answerDirectoryPath := filepath.Join(directoryPath, answerDirectoryName)
			unitDataFilePath := filepath.Join(directoryPath, unitDataFileName)

			if !directoryExists(questionDirectoryPath) || !directoryExists(answerDirectoryPath) || !fileExists(unitDataFilePath) {
				return findSubUnits(directoryPath)
			}

			unitData := unitData{}

			unitDataBytes, e := ioutil.ReadFile(unitDataFilePath)
			if e != nil {
				return findSubUnits(directoryPath)
			}

			e = json.Unmarshal(unitDataBytes, &unitData)
			if e != nil {
				return findSubUnits(directoryPath)
			}

			if unitData.Stage < 0 || unitData.Stage >= len(stages) || unitData.UnixTime < 0 {
				return findSubUnits(directoryPath)
			}

			questionFilePaths, ie := listFiles(questionDirectoryPath)
			if ie != nil {
				return findSubUnits(directoryPath)
			}

			answerFilePaths, ie := listFiles(answerDirectoryPath)
			if ie != nil {
				return findSubUnits(directoryPath)
			}

			if len(questionFilePaths) < 1 || len(answerFilePaths) < 1 {
				return findSubUnits(directoryPath)
			}

			unit := unit{
				path:                  directoryPath,
				dataFilePath:          unitDataFilePath,
				data:                  unitData,
				questionFilePaths:     questionFilePaths,
				answerFilePaths:       answerFilePaths,
				questionDirectoryPath: questionDirectoryPath,
				answerDirectoryPath:   answerDirectoryPath,
			}
			units = append(units, unit)
		}

		return units
	}

	units := inner(currentDirectoryPath)

	if len(units) < 1 {
		return nil, ierror{m: "No units are found"}
	}

	return units, nil
}

func unitWithLeastUnixTime(units []unit) (unit, bool) {

	if len(units) < 1 {
		return unit{}, false
	}

	unitIndex := 0
	leastUnixTime := int64(math.MaxInt64)
	for i, unit := range units {
		unixTime := unit.data.UnixTime
		if unixTime < leastUnixTime {
			leastUnixTime = unixTime
			unitIndex = i
		}
	}

	return units[unitIndex], true
}
