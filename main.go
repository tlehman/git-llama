package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tlehman/git-llama/ollm"
	"github.com/tlehman/git-llama/vdb"
)

const ERR_NOT_SINGLE_PROMPT = 1
const ERR_OLLAMA_NOT_RUNNING = 2
const ERR_VECTORDB_OPEN_FAIL = 3

func usage() {
	fmt.Printf("Usage:\n  git-llama [your prompt, delimited by quotes]\n")
	os.Exit(ERR_NOT_SINGLE_PROMPT)
}

func dbfilename() string {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return ""
	}

	// Construct the full path by joining home directory with .git-llama.db
	dbPath := filepath.Join(homeDir, ".git-llama.db")
	return dbPath
}

func exitIfOllamaIsNotRunning() {
	if !ollm.IsOllamaRunning() {
		fmt.Printf("Ollama is not running! Please run `ollama serve` in another window\n")
		os.Exit(ERR_OLLAMA_NOT_RUNNING)
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
}
