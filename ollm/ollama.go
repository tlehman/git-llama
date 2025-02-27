// Package ollm wraps the ollama api
package ollm

import (
	"fmt"
	"net/http"
	"time"
)

// Use the default Ollama server address
const serverAddress = "http://localhost:11434"

func IsOllamaRunning() (bool, error) {
	// Create a client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

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
