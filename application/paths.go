package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	errorNoRootConfigurationFile = iota + 64824075
)

func getConfigurationDirectoryPath() (string, ierrori) {

	var (
		e                          error
		homeDirectoryPath          string
		configurationDirectoryPath string
	)

	homeDirectoryPath, e = os.UserHomeDir()
	if e != nil {
		return configurationDirectoryPath, ierror{m: "Could not user's home directory", e: e}
	}

	configurationDirectoryPath = filepath.Join(homeDirectoryPath, configurationDirectoryName)

	return configurationDirectoryPath, nil
}

func getRootConfigurationFilePath() (string, ierrori) {

	var (
		ie                         ierrori
		configurationDirectoryPath string
		rootConfigurationFilePath  string
	)

	configurationDirectoryPath, ie = getConfigurationDirectoryPath()
	if ie != nil {
		return rootConfigurationFilePath, ie
	}

	rootConfigurationFilePath = filepath.Join(configurationDirectoryPath, rootConfigurationFileName)

	return rootConfigurationFilePath, nil
}

func getRootDirectoryPath() (string, ierrori) {

	var (
		e                         error
		ie                        ierrori
		rootConfigurationFilePath string
		rootFile                  *os.File
		root                      string
		rootBytes                 []byte
	)

	rootConfigurationFilePath, ie = getRootConfigurationFilePath()
	if ie != nil {
		return root, ie
	}

	rootFile, e = os.Open(rootConfigurationFilePath)
	if e != nil {
		return root, ierror{
			m: "Root directory is not set up. Please, set root directory with `-set-root` flag",
			e: e,
			c: errorNoRootConfigurationFile,
		}
	}

	rootBytes, e = ioutil.ReadAll(rootFile)
	if e != nil {
		return root, ierror{m: "Could not read root configuration file", e: e}
	}

	root = string(rootBytes)

	return root, nil
}

func getDatabaseFilePath() (string, ierrori) {

	var (
		ie                ierrori
		rootDirectoryPath string
		databaesFilePath  string
	)

	rootDirectoryPath, ie = getRootDirectoryPath()
	if ie != nil {
		return databaesFilePath, ie
	}

	databaesFilePath = filepath.Join(rootDirectoryPath, databaseFileName)

	return databaesFilePath, nil
}

func getLogFilePath() (string, ierrori) {

	var (
		logFilePath                string
		configurationDirectoryPath string
		ie                         ierrori
	)

	configurationDirectoryPath, ie = getConfigurationDirectoryPath()
	if ie != nil {
		return logFilePath, ie
	}

	logFilePath = filepath.Join(configurationDirectoryPath, logFileName)

	return logFilePath, nil
}
