package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/tlehman/git-llama/vdb"
)

const LLM_MODEL_NAME = "llama3.2"

const ERR_NOT_SINGLE_PROMPT = 1
const ERR_OLLAMA_API_FAIL = 2
const ERR_OLLAMA_NOT_RUNNING = 3
const ERR_VECTORDB_OPEN_FAIL = 4

func usage() {
	fmt.Printf("Usage:\n  git-llama [your prompt, delimited by quotes]\n")
	os.Exit(ERR_NOT_SINGLE_PROMPT)
}

// wrap the prompt with git-specific data for the LLM
func wrap(prompt string) string {
	return fmt.Sprintf("git command for %s just the command, no text", prompt)
}

func isOllamaRunning() (bool, error) {
	// Create a client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Use the default Ollama server address
	serverAddress := "http://localhost:11434"

	// Attempt a simple request to the root endpoint
	resp, err := client.Get(serverAddress)
	if err != nil {
		// If the connection fails (e.g., server not running), return false
		return false, nil
	}
	defer resp.Body.Close()

	// Check if the status code indicates the server is responding
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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

func main() {
	// open or create the vector database
	vectordb, err := vdb.Open(dbfilename(), LLM_MODEL_NAME)
	if err != nil {
		fmt.Printf("failed to open vector db: %s\n", err)
		os.Exit(ERR_VECTORDB_OPEN_FAIL)
	}
	defer vectordb.Close()
	_ = vectordb.Get("foo")

	// check if a single prompt is passed in
	if len(os.Args) != 2 {
		usage()
	}
	// check if ollama is running
	running, err := isOllamaRunning()
	if err != nil {
		fmt.Printf("failed to check if ollama is running: %s\n", err)
		os.Exit(ERR_OLLAMA_API_FAIL)
	}
	if !running {
		fmt.Printf("Ollama is not running! Please run `ollama serve` in another window\n")
		os.Exit(ERR_OLLAMA_NOT_RUNNING)
	}
	// the first argument is assumed to be the prompt input
	prompt := os.Args[1]

	// create the ollama client
	client, err := api.ClientFromEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed creating a client: %s", err)
		os.Exit(1)
	}
	req := &api.GenerateRequest{
		Model:  LLM_MODEL_NAME,
		Prompt: wrap(prompt),
	}
	// Context for the request
	ctx := context.Background()

	// Function to handle the response
	respond := func(resp api.GenerateResponse) error {
		fmt.Print(resp.Response) // Print the response as it streams
		return nil
	}

	// Send the prompt and get the response
	err = client.Generate(ctx, req, respond)
}
