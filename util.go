package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

func usage() {
	fmt.Printf("Usage:\n  git-llama [your prompt, delimited by quotes]\n")
	os.Exit(ERR_NOT_SINGLE_PROMPT)
}

func dbfilename() string {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return ""
	}

	// Construct the full path by joining working directory with .git-llama.db
	dbPath := filepath.Join(wd, ".git-llama.db")
	return dbPath
}

func ensureDbIsGitExcluded() {
	excludeFilePath := ".git/info/exclude"
	dbfn := ".git-llama.db"
	data, err := os.ReadFile(excludeFilePath)
	if err != nil {
		fmt.Printf("git repo is invalid: %s\n", err)
		os.Exit(ERR_GIT_ERROR)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		// if the line matches dbfn, then the db filename is excluded from git, and we can return
		if line == dbfn {
			return
		}
	}
	if err := scanner.Err(); err != nil {
		// Handle any scanner errors here
		fmt.Printf("error reading file: %s\n", err)
	}
	// if you make it through the loop, then the dbfilename is NOT in the .git/info/exclude file
	excludeFile, err := os.OpenFile(excludeFilePath, os.O_RDWR, os.ModeAppend)
	_, err = excludeFile.WriteString(dbfn + "\n")
	if err != nil {
		fmt.Printf("error appending dbfilename to .git/info/exclude\n")
	}
}
