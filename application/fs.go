package main

import (
	"io/ioutil"
	"os"
	"path"
)

func directoryExists(p string) bool {

	var (
		e    error
		info os.FileInfo
	)

	info, e = os.Stat(p)

	return e == nil && info.IsDir()
}

func fileExists(p string) bool {

	var (
		e    error
		info os.FileInfo
	)

	info, e = os.Stat(p)

	return e == nil && !info.IsDir()
}

func listDirectories(p string) ([]string, ierrori) {

	var (
		infos       []os.FileInfo
		e           error
		directories []string
	)

	infos, e = ioutil.ReadDir(p)
	if e != nil {
		return directories, ierror{m: "Could not list directories", e: e}
	}

	for _, info := range infos {
		if info.IsDir() {
			directories = append(directories, path.Join(p, info.Name()))
		}
	}

	return directories, nil
}

func listDirectoryNames(p string) ([]string, ierrori) {

	var (
		infos       []os.FileInfo
		e           error
		directories []string
	)

	infos, e = ioutil.ReadDir(p)
	if e != nil {
		return directories, ierror{m: "Could not list directories", e: e}
	}

	for _, info := range infos {
		if info.IsDir() {
			directories = append(directories, info.Name())
		}
	}

	return directories, nil
}

func listFiles(p string) ([]string, ierrori) {

	var (
		infos []os.FileInfo
		e     error
		files []string
	)

	infos, e = ioutil.ReadDir(p)
	if e != nil {
		return files, ierror{m: "Could not list files", e: e}
	}

	for _, info := range infos {
		if !info.IsDir() {
			files = append(files, path.Join(p, info.Name()))
		}
	}

	return files, nil
}

func listFileNames(p string) ([]string, ierrori) {

	var (
		infos []os.FileInfo
		e     error
		files []string
	)

	infos, e = ioutil.ReadDir(p)
	if e != nil {
		return files, ierror{m: "Could not list files", e: e}
	}

	for _, info := range infos {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
	}

	return files, nil
}
