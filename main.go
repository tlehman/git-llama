package main

import (
	"fmt"
	"os"

	"github.com/ollama/ollama/api"
)

func usage() {
	fmt.Printf("Usage:\n  git-llama [your prompt, delimited by quotes]\n")
	os.Exit(1)
}

func main() {
	// check if a single prompt is passed in
	if len(os.Args) != 2 {
		usage()
	}

	// check if ollama is running
	_, err := api.ClientFromEnvironment()
	if err != nil {
		fmt.Sprintf("err: %s", err)
		os.Exit(1)
	}

	// the first argument is assumed to be the prompt input
}
