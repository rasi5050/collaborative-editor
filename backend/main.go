package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
    err := initDocumentStorage()
    if err != nil {
        log.Fatalf("Failed to initialize document storage: %v", err)
    }

    http.HandleFunc("/ws", handleConnections) // Only WebSocket handler for the single file

    log.Println("Server started on :8082")
    log.Fatal(http.ListenAndServe(":8082", nil))
}

func initDocumentStorage() error {
    // Ensure the "documents" directory exists
    return os.MkdirAll("documents", 0755)
}
