package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

type unitData struct {
	UnixTime        int64
	Stage           int
	Inverse         bool
	InverseUnixTime int64
	InverseStage    int
}

type unit struct {
	path           string
	dataFilePath   string
	data           unitData
	aFilePaths     []string
	bFilePaths     []string
	aDirectoryPath string
	bDirectoryPath string
}

func (u *unitData) getInverse() bool {
	return u.InverseUnixTime < u.UnixTime && u.Inverse
}

func (u *unitData) getStage() int {
	if u.getInverse() {
		return u.InverseStage
	}
	return u.Stage
}

func (u *unitData) getUnixTime() int64 {
	if u.getInverse() {
		return u.InverseUnixTime
	}
	return u.UnixTime
}

func (u *unitData) setInverse(inverse bool) {
	u.Inverse = inverse
}

func (u *unitData) setStage(stage int) {
	if u.getInverse() {
		u.InverseStage = stage
		return
	}
	u.Stage = stage
}

func (u *unitData) setUnixTime(unixTime int64) {
	if u.getInverse() {
		u.InverseUnixTime = unixTime
		return
	}
	u.UnixTime = unixTime
}

func (u *unitData) nextStage() {
	if u.getInverse() {
		u.InverseStage++
		if u.InverseStage >= len(stages) {
			u.InverseStage--
		}
		u.InverseUnixTime = inverseUnixTimeForStage(u.InverseStage)
		return
	}
	u.Stage++
	if u.Stage >= len(stages) {
		u.Stage--
	}
	u.UnixTime = unixTimeForStage(u.Stage)
}

func (u *unitData) previousStage() {
	if u.getInverse() {
		u.InverseStage = 0
		u.InverseUnixTime = inverseUnixTimeForStage(u.InverseStage)
		return
	}
	u.Stage = 0
	u.UnixTime = unixTimeForStage(u.Stage)
}

func newUnit(unitDirectoryPath string, isInverse bool) ierrori {

	unitDirectoryPath = filepath.Clean(unitDirectoryPath)

	e := os.MkdirAll(unitDirectoryPath, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not create unit directory", e: e}
	}

	e = os.MkdirAll(filepath.Join(unitDirectoryPath, aDirectoryName), os.ModePerm)
	if e != nil {
		return ierror{m: "Could not create unit A association directory", e: e}
	}

	e = os.MkdirAll(filepath.Join(unitDirectoryPath, bDirectoryName), os.ModePerm)
	if e != nil {
		return ierror{m: "Could not create unit B association directory", e: e}
	}

	stage := 0

	unitData := unitData{
		UnixTime:        unixTimeForStage(stage),
		Stage:           stage,
		Inverse:         isInverse,
		InverseUnixTime: inverseUnixTimeForStage(stage),
		InverseStage:    stage,
	}

	unitDataBytes, e := json.Marshal(unitData)

	unitDataFilePath := filepath.Join(unitDirectoryPath, unitDataFileName)

	e = ioutil.WriteFile(unitDataFilePath, unitDataBytes, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not write unit data file", e: e}
	}

	return nil
}

func getUnit(unitPath string) (unit, ierrori) {

	aDirectoryPath := filepath.Join(unitPath, aDirectoryName)
	bDirectoryPath := filepath.Join(unitPath, bDirectoryName)
	unitDataFilePath := filepath.Join(unitPath, unitDataFileName)

	if !directoryExists(aDirectoryPath) {
		return unit{}, ierror{m: "Could not get unit, A association directory does not exists"}
	}

	if !directoryExists(bDirectoryPath) {
		return unit{}, ierror{m: "Could not get unit, B association directory does not exists"}
	}

	if !fileExists(unitDataFilePath) {
		return unit{}, ierror{m: "Could not get unit, data file does not exists"}
	}

	aFilePaths, ie := listFiles(aDirectoryPath)
	if ie != nil {
		return unit{}, ierror{m: "Could not list files in unit's A association directory", e: ie}
	}

	bFilePaths, ie := listFiles(bDirectoryPath)
	if ie != nil {
		return unit{}, ierror{m: "Could not list files in unit's B association directory", e: ie}
	}

	if len(aFilePaths) < 1 {
		return unit{}, ierror{m: "No files in unit's A association directory"}
	}

	if len(bFilePaths) < 1 {
		return unit{}, ierror{m: "No files in unit's B association directory"}
	}

	unitData := unitData{}

	unitDataBytes, e := ioutil.ReadFile(unitDataFilePath)
	if e != nil {
		return unit{}, ierror{m: "Could not read unit's data file", e: e}
	}

	e = json.Unmarshal(unitDataBytes, &unitData)
	if e != nil {
		return unit{}, ierror{m: "Could not read unit's data file", e: e}
	}

	if unitData.Stage < 0 || unitData.Stage >= len(stages) || unitData.UnixTime < 0 ||
		unitData.InverseStage < 0 || unitData.InverseStage >= len(stages) || unitData.InverseUnixTime < 0 {
		return unit{}, ierror{m: "Unit's data file is invalid"}
	}

	unit := unit{
		path:           unitPath,
		dataFilePath:   unitDataFilePath,
		data:           unitData,
		aFilePaths:     aFilePaths,
		bFilePaths:     bFilePaths,
		aDirectoryPath: aDirectoryPath,
		bDirectoryPath: bDirectoryPath,
	}

	return unit, nil
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

			unit, ie := getUnit(directoryPath)
			if ie != nil {
				return findSubUnits(directoryPath)
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
		unixTime := unit.data.getUnixTime()
		if unixTime < leastUnixTime {
			leastUnixTime = unixTime
			unitIndex = i
		}
	}

	return units[unitIndex], true
}

func toggleInverse(unitPath string) ierrori {

	unit, ie := getUnit(unitPath)
	if ie != nil {
		return ie
	}

	unit.data.Inverse = !unit.data.Inverse

	unitDataBytes, e := json.Marshal(unit.data)
	if e != nil {
		return ierror{m: "Could not update unit's data", e: e}
	}

	e = ioutil.WriteFile(unit.dataFilePath, unitDataBytes, os.ModePerm)
	if e != nil {
		return ierror{m: "Could not update unit's data", e: e}
	}

	return nil
}
