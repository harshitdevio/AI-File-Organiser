package main 

import (
	"fmt"
	"log"
	"net/http"
	"organiser/classifier"
	"organiser/network"
	"organiser/scanner"
	"organiser/server"
)

func main() {

	go func() {
		http.HandleFunc("/api/ingest", server.IngestHandler)
		fmt.Println("[Server] Listening for results on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Receiver server failed: %v", err)
		}
	}()

	dir, err := scanner.Scan()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	files, err := classifier.ProcessDirectory(dir)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}

	fmt.Printf("Found %d files:\n", len(files))
	for _, file := range files {
		fmt.Printf("File: %s | Type: %s\n", file.FilePath, file.MIMEType)
	}

	apiURL := "http://127.0.0.1:8000/detect-topics-batch"
	fmt.Printf("Sending %d files to Python LLM...\n", len(files))

	if err := network.SendToService(apiURL, files); err != nil {
		log.Fatalf("Dispatch failed: %v", err)
	}

	fmt.Println("Transfer complete. Waiting for Python to finish and send results back...")

	select {} 
}