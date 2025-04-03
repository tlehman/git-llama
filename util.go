package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

var foo io.Writer

const whichPath = "/usr/bin/which"

func Which(cmdPath string) string {
	var buf bytes.Buffer
	cmd := exec.Command("which", cmdPath)
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		fmt.Printf("error running which: %s\n", err)
		return ""
	}

	return strings.TrimSpace(buf.String())
}
