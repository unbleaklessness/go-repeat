package main

import (
	"os"
	"path/filepath"
)

func getConfigurationDirectoryPath() (string, ierrori) {
	homeDirectoryPath, e := os.UserHomeDir()
	if e != nil {
		return homeDirectoryPath, ierror{m: "Could not get user's home directory", e: e}
	}
	return filepath.Join(homeDirectoryPath, configurationDirectoryName), nil
}

func getTemplatesFilePath() (string, ierrori) {
	configurationDirectoryPath, ie := getConfigurationDirectoryPath()
	if ie != nil {
		return configurationDirectoryPath, ie
	}
	return filepath.Join(configurationDirectoryPath, templatesFileName), nil
}

func getLogFilePath() (string, ierrori) {
	configurationDirectoryPath, ie := getConfigurationDirectoryPath()
	if ie != nil {
		return configurationDirectoryPath, ie
	}
	return filepath.Join(configurationDirectoryPath, logFileName), nil
}
