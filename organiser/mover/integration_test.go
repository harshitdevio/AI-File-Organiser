package mover

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessMove_RealFilesystem(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "mover-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create source file
	sourceFile := filepath.Join(tempDir, "test-file.txt")
	err = os.WriteFile(sourceFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create target directory inside temp dir
	targetDir := filepath.Join(tempDir, "documents")

	// Use real filesystem operations
	realFS := &RealFileSystemOps{}
	writer := &bytes.Buffer{}

	// Override the target directory to be inside temp dir
	// We need to test the logic, so we'll use ProcessMoveWithDeps
	err = ProcessMoveWithDeps(sourceFile, targetDir, 0.85, realFS, writer)

	if err != nil {
		t.Errorf("ProcessMoveWithDeps() unexpected error: %v", err)
		return
	}

	// Verify source file no longer exists
	if _, err := os.Stat(sourceFile); !os.IsNotExist(err) {
		t.Error("Source file still exists after move")
	}

	// Verify target file exists
	targetFile := filepath.Join(targetDir, "test-file.txt")
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		t.Error("Target file does not exist after move")
	}

	// Verify content
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Errorf("Failed to read target file: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("Target file content = %q, want %q", string(content), "test content")
	}

	// Verify output message
	output := writer.String()
	if !strings.Contains(output, "Organized:") {
		t.Errorf("Output missing 'Organized:' message: %q", output)
	}
}

func TestProcessMove_RealFilesystem_LowConfidence(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mover-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create source file
	sourceFile := filepath.Join(tempDir, "unknown.dat")
	err = os.WriteFile(sourceFile, []byte("unknown content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	realFS := &RealFileSystemOps{}
	writer := &bytes.Buffer{}

	// Save current directory and change to temp dir
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldWd)

	// Low confidence should move to "misc"
	err = ProcessMoveWithDeps(sourceFile, "documents", 0.40, realFS, writer)

	if err != nil {
		t.Errorf("ProcessMoveWithDeps() unexpected error: %v", err)
		return
	}

	// Verify file moved to misc directory
	miscFile := filepath.Join(tempDir, "misc", "unknown.dat")
	if _, err := os.Stat(miscFile); os.IsNotExist(err) {
		t.Error("File was not moved to misc directory")
	}

	// Verify output indicates misc directory
	output := writer.String()
	if !strings.Contains(output, "-> misc") {
		t.Errorf("Output does not indicate misc directory: %q", output)
	}
}