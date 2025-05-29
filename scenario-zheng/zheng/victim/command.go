package main

import (
	"os/exec"
	"strings"
)

func execute(command string) string {
	cmd := exec.Command(command)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err.Error()
	} else {
		return out.String()
	}
}
