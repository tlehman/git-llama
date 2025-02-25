package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ollama/ollama/api"
)

func usage() {
	fmt.Printf("Usage:\n  git-llama [your prompt, delimited by quotes]\n")
	os.Exit(1)
}

// wrap the prompt with git-specific data for the LLM
func wrap(prompt string) string {
	return fmt.Sprintf("git command for %s just the command, no text", prompt)
}

func main() {
	// check if a single prompt is passed in
	if len(os.Args) != 2 {
		usage()
	}
	// the first argument is assumed to be the prompt input
	prompt := os.Args[1]

	// check if ollama is running
	client, err := api.ClientFromEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
		os.Exit(1)
	}
	req := &api.GenerateRequest{
		Model:  "llama3.2",
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
