package main

import (
	"log"
	"minivault/interfaces"
	"net/http"
)

func main() {
	// Set up the HTTP handler for the /generate endpoint.
	// When a POST request is made to /generate, the GenerateHandler function from the interfaces package is invoked.
	http.HandleFunc("/generate", interfaces.GenerateHandler)
	// Register the streaming endpoint for SSE responses
	http.HandleFunc("/generate/streaming", interfaces.StreamingGenerateHandler)

	// Log a message to indicate that the MiniVault API server is running and listening on port 8080.
	log.Println("MiniVault API running on :8080")

	// Start the HTTP server on port 8080. This function blocks and will run indefinitely,
	// serving incoming requests to the /generate endpoint.
	http.ListenAndServe(":8080", nil)
}
