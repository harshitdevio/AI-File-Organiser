package main

import (
	"fmt"
	"log"
	"organiser/classifier"
	"organiser/network"
	"organiser/scanner"
)

func main() {
	dir, err := scanner.Scan()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	files, err := classifier.ProcessDirectory(dir)
	if err != nil {
		log.Fatalf("Something went wring: %v", err)
	}
	fmt.Printf("Found %d files:\n", len(files))
	for _, file := range files {
		fmt.Printf("File: %s | Type: %s\n", file.FilePath, file.MIMEType)
	}
	apiURL := "http://127.0.0.1:8000/detect-topics-batch"
	fmt.Printf("Sending %d files to LLM...\n", len(files))

	if err := network.SendToService(apiURL, files); err != nil {
		log.Fatalf("Dispatch failed: %v", err)
	}

	fmt.Println("Transfer complete.")
}
