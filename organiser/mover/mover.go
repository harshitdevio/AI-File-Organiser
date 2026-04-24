package mover

import (
    "fmt"
    "os"
    "path/filepath"
)

func MoveFiles(results []any) {
}

func ProcessMove(path string, topic string, confidence float64) {
    targetDir := "misc"
    if confidence >= 0.60 {
        targetDir = topic
    }

    os.MkdirAll(targetDir, os.ModePerm)

    newName := filepath.Base(path)
    err := os.Rename(path, filepath.Join(targetDir, newName))
    
    if err != nil {
        fmt.Printf("Error moving %s: %v\n", newName, err)
    } else {
        fmt.Printf("Organized: %s -> %s (%.2f)\n", newName, targetDir, confidence)
    }
}