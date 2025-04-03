package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tlehman/git-llama/ollm"
	"github.com/tlehman/git-llama/vdb"
)

const BUFSIZE = 1024

const ERR_NOT_SINGLE_PROMPT = 1
const ERR_OLLAMA_NOT_INSTALLED = 2
const ERR_VECTORDB_OPEN_FAIL = 3
const ERR_GIT_ERROR = 4

/*
  - start ollama
  - check git/vec delta
    -- update vec db if changed
  - semantic search
  - stop ollama
*/

func startOllama(stopChan chan bool) {
	ollamaPath := Which("ollama")
	cmd := exec.Command(ollamaPath, "serve")
	cmd.Start()
	stop := <-stopChan
	if stop {
		cmd.Process.Kill()
	}
}

func main() {
	stopChan := make(chan bool, 1)
	go startOllama(stopChan)
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

	stopOllama := true
	stopChan <- stopOllama
}
