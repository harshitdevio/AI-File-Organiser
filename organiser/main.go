package main

import (
	"fmt"
	"log"
	"organiser/classifier"
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
        fmt.Printf("File: %s | Type: %s\n", file.Path, file.MIMEType)
    }

}