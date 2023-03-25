package main

import (
	"log"
	"net/http"
	"os"
)

// The main function handles ports assignment, sets up handler endpoints and starts the HTTP-server
func main() {

	// Handle port assignment (either based on environment variable, or local override)
	port := os.Getenv("PORT")
	defaultPort := 8080
	if port == "" {
		log.Println("PORT has not been set. Default: 8080", defaultPort)
	}

	// Set up handler endpoints

	// Starting HTTP-server
	log.Println("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
