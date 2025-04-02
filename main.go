package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tlehman/git-llama/ollm"
	"github.com/tlehman/git-llama/vdb"
)

const BUFSIZE = 1024

const ERR_NOT_SINGLE_PROMPT = 1
const ERR_OLLAMA_NOT_RUNNING = 2
const ERR_VECTORDB_OPEN_FAIL = 3
const ERR_GIT_ERROR = 4

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

func exitIfOllamaIsNotRunning() {
	if !ollm.IsOllamaRunning() {
		fmt.Printf("Ollama is not running! Please run `ollama serve` in another window\n")
		os.Exit(ERR_OLLAMA_NOT_RUNNING)
	}
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

func main() {
	exitIfOllamaIsNotRunning()
	// open or create the vector database
	vectordb, err := vdb.Open(dbfilename(), ollm.LLM_MODEL_NAME)
	if err != nil {
		fmt.Printf("failed to open vector db: %s\n", err)
		os.Exit(ERR_VECTORDB_OPEN_FAIL)
	}
	// check if .git-llama.db is in .git/info/exclude
	ensureDbIsGitExcluded()
	defer vectordb.Close()
	// Get the dimension of the model
	dim := ollm.ModelDimension(ollm.LLM_MODEL_NAME)
	// Create the table for the model
	vectordb.CreateTableIdempotent(dim)

	// check if a single prompt is passed in
	if len(os.Args) != 2 {
		usage()
	}

	// the first argument is assumed to be the prompt input
	prompt := os.Args[1]
	response := ollm.Generate(prompt)
	embedding := ollm.Embed(prompt)
	err = vectordb.Insert(response, embedding)
	if err != nil {
		fmt.Printf("failed inserting embedding vector: %s\n", err)
	}
	fmt.Println(response)
}
