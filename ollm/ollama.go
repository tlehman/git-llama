// Package ollm wraps the ollama api
package ollm

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/tlehman/git-llama/vdb"
)

// Use the default Ollama server address
var serverAddress = os.Getenv("OLLAMA_HOST")

const serverAddressDefault = "http://localhost:11434"
const LLM_MODEL_NAME = "llama3.2"

func IsOllamaRunning() bool {
	if serverAddress == "" {
		serverAddress = serverAddressDefault
	}
	// Create a client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Attempt a simple request to the root endpoint
	resp, err := client.Get(serverAddress)
	if err != nil {
		// If the connection fails (e.g., server not running), return false
		return false
	}
	defer resp.Body.Close()

	// Check if the status code indicates the server is responding
	if resp.StatusCode == http.StatusOK {
		return true
	}

	return false
}

// wrap the prompt with git-specific data for the LLM
func wrap(prompt string) string {
	return fmt.Sprintf("git command for %s just the command, no text", prompt)
}

func Generate(prompt string) string {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed creating a client: %s", err)
		os.Exit(1)
	}

	req := &api.GenerateRequest{
		Model:  LLM_MODEL_NAME,
		Prompt: wrap(prompt),
		Stream: new(bool),
	}
	// Context for the request
	ctx := context.Background()

	// Create channel and wait group to get all the response back
	var responseChan chan string = make(chan string, 1)

	// Function to handle the response
	respond := func(resp api.GenerateResponse) error {
		responseChan <- resp.Response
		return nil
	}

	// Send the prompt and get the response
	err = client.Generate(ctx, req, respond)

	response := <-responseChan

	// Since the LLMs denote code with backticks `, we want to strip those:
	return strings.ReplaceAll(response, "`", "")
}

func Embed(prompt string) *vdb.Vector {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed creating a client: %s", err)
		os.Exit(1)
	}

	req := &api.EmbedRequest{
		Model: LLM_MODEL_NAME,
		Input: prompt,
	}

	// Context for the Request
	ctx := context.Background()

	response, err := client.Embed(ctx, req)
	if err != nil {
		fmt.Printf("failed calling ollama embed API: %s\n", err)
		return nil
	}

	return &vdb.Vector{Values: response.Embeddings[0]}
}

func ModelDimension(modelname string) int {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed creating a client: %s", err)
		os.Exit(1)
	}
	ctx := context.Background()
	req := &api.ShowRequest{
		Model: modelname,
	}
	response, err := client.Show(ctx, req)
	if err != nil {
		fmt.Printf("failed calling /api/show: %s\n", err)
		return -1
	}
	dim := (response.ModelInfo["llama.embedding_length"]).(float64)
	return int(dim)
}
