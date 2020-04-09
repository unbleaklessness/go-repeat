package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func directoryExists(p string) bool {
	info, e := os.Stat(p)
	return e == nil && info.IsDir()
}

func fileExists(p string) bool {
	info, e := os.Stat(p)
	return e == nil && !info.IsDir()
}

func listDirectories(p string) ([]string, ierrori) {

	infos, e := ioutil.ReadDir(p)
	if e != nil {
		return nil, ierror{m: "Could not list directories", e: e}
	}

	directoryPaths := make([]string, 0)

	for _, info := range infos {
		if info.IsDir() {
			directoryPaths = append(directoryPaths, filepath.Join(p, info.Name()))
		}
	}

	return directoryPaths, nil
}

func listFiles(p string) ([]string, ierrori) {

	infos, e := ioutil.ReadDir(p)
	if e != nil {
		return nil, ierror{m: "Could not list files", e: e}
	}

	filePaths := make([]string, 0)

	for _, info := range infos {
		if !info.IsDir() {
			filePaths = append(filePaths, filepath.Join(p, info.Name()))
		}
	}

	return filePaths, nil
}
