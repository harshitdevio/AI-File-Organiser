package main

import (
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/gabriel-vasile/mimetype"
)

type FileInfo struct {
    Path     string
    MIMEType string
}

func ProcessDirectory(dirPath string) ([]FileInfo, error) {
    if valid, err := isValidDirectory(dirPath); !valid {
        return nil, err
    }

    var files []FileInfo

    err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        mtype, err := mimetype.DetectFile(path)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error detecting MIME for %s: %v\n", path, err)
            return nil
        }

        files = append(files, FileInfo{
            Path:     path,
            MIMEType: mtype.String(),
        })

        return nil
    })

    if err != nil {
        return nil, fmt.Errorf("error walking directory: %w", err)
    }

    return files, nil
}

func isValidDirectory(path string) (bool, error) {
    info, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return false, fmt.Errorf("directory does not exist: %s", path)
        }
        return false, fmt.Errorf("error accessing path: %w", err)
    }
    
    if !info.IsDir() {
        return false, fmt.Errorf("path is not a directory: %s", path)
    }
    
    return true, nil
}
