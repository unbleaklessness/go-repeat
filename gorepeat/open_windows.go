// +build windows, !linux

package main

import "os/exec"

func open(p string) error {
	return exec.Command("cmd", "/c", "start", "", p).Start()
}
