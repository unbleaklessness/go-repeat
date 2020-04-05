package main

import (
	"path"
)

type unit struct {
	directory       string
	informationFile string
	qAndAs          map[string]string
}

func getQAndAs(questionsDirectory string, answersDirectory string) map[string]string {

	var (
		ie ierrori

		qAndAs map[string]string

		questions []string
		answers   []string

		question string

		questionDirectory string
		answerDirectory   string

		questionFiles []string
		answerFiles   []string

		answerIndex int
	)

	qAndAs = make(map[string]string)

	questions, ie = listDirectoryNames(questionsDirectory)
	if ie != nil {
		return qAndAs
	}

	answers, ie = listDirectoryNames(answersDirectory)
	if ie != nil {
		return qAndAs
	}

	for _, question = range questions {
		answerIndex = stringsIndex(answers, question)
		if answerIndex != -1 {
			questionDirectory = path.Join(questionsDirectory, question)
			questionFiles, ie = listFiles(questionDirectory)
			if ie != nil {
				continue
			}
			if len(questionFiles) < 1 {
				continue
			}

			answerDirectory = path.Join(answersDirectory, answers[answerIndex])
			answerFiles, ie = listFiles(answerDirectory)
			if ie != nil {
				continue
			}
			if len(answerFiles) < 1 {
				continue
			}

			qAndAs[questionDirectory] = answerDirectory
		}
	}

	return qAndAs
}

func toUnit(p string) (unit, bool) {

	var (
		u unit

		questionsDirectory string
		answersDirectory   string
		informationFile    string

		qAndAs map[string]string
	)

	questionsDirectory = path.Join(p, questionsDirectoryName)
	answersDirectory = path.Join(p, answersDirectoryName)
	informationFile = path.Join(p, informationFileName)

	if !fileExists(informationFile) {
		return u, false
	}

	if !directoryExists(questionsDirectory) {
		return u, false
	}

	if !directoryExists(answersDirectory) {
		return u, false
	}

	qAndAs = getQAndAs(questionsDirectory, answersDirectory)
	if len(qAndAs) < 1 {
		return u, false
	}

	u.qAndAs = qAndAs
	u.directory = p
	u.informationFile = informationFile

	return u, true
}

func findUnits(p string) []unit {

	var (
		units       []unit
		u           unit
		directories []string
		directory   string
		ie          ierrori
		ok          bool
	)

	directories, ie = listDirectories(p)
	if ie != nil {
		return units
	}

	for _, directory = range directories {
		if u, ok = toUnit(directory); ok {
			units = append(units, u)
		} else {
			units = append(units, findUnits(directory)...)
		}
	}

	return units
}
