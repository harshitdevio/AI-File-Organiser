package mover

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileSystemOps interface {
	MkdirAll(path string, perm os.FileMode) error
	Rename(oldpath, newpath string) error
}

type RealFileSystemOps struct{}

func (rfs *RealFileSystemOps) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (rfs *RealFileSystemOps) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func ProcessMoveWithDeps(path string, topic string, confidence float64, fsOps FileSystemOps, writer io.Writer) error {
	targetDir := "misc"
	if confidence >= 0.60 {
		targetDir = topic
	}

	err := fsOps.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	newName := filepath.Base(path)
	newPath := filepath.Join(targetDir, newName)

	err = fsOps.Rename(path, newPath)

	if err != nil {
		fmt.Fprintf(writer, "Error moving %s: %v\n", newName, err)
		return err
	}

	fmt.Fprintf(writer, "Organized: %s -> %s (%.2f)\n", newName, targetDir, confidence)
	return nil
}

func ProcessMove(path string, topic string, confidence float64) {
	ProcessMoveWithDeps(path, topic, confidence, &RealFileSystemOps{}, os.Stdout)
}