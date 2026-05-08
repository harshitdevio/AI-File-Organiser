package mover 
import (
	"bytes"
	"path/filepath"
	"os"
	"testing"
)
func TestProcessMoveWithDeps_ConfidenceThreshold(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		topic          string
		confidence     float64
		expectedTarget string
	}{
		{
			name:           "confidence below threshold uses misc",
			path:           "/source/file.txt",
			topic:          "documents",
			confidence:     0.59,
			expectedTarget: "misc",
		},
		{
			name:           "confidence at threshold uses topic",
			path:           "/source/file.txt",
			topic:          "documents",
			confidence:     0.60,
			expectedTarget: "documents",
		},
		{
			name:           "confidence above threshold uses topic",
			path:           "/source/file.txt",
			topic:          "images",
			confidence:     0.95,
			expectedTarget: "images",
		},
		{
			name:           "zero confidence uses misc",
			path:           "/source/file.txt",
			topic:          "videos",
			confidence:     0.0,
			expectedTarget: "misc",
		},
		{
			name:           "perfect confidence uses topic",
			path:           "/source/file.txt",
			topic:          "music",
			confidence:     1.0,
			expectedTarget: "music",
		},
		{
			name:           "just below threshold uses misc",
			path:           "/source/important.doc",
			topic:          "work",
			confidence:     0.5999,
			expectedTarget: "misc",
		},
		{
			name:           "just above threshold uses topic",
			path:           "/source/important.doc",
			topic:          "work",
			confidence:     0.6001,
			expectedTarget: "work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemOps{}
			writer := &bytes.Buffer{}

			ProcessMoveWithDeps(tt.path, tt.topic, tt.confidence, mockFS, writer)

			// Verify correct directory was created
			if len(mockFS.MkdirAllCalls) != 1 {
				t.Fatalf("Expected 1 MkdirAll call, got %d", len(mockFS.MkdirAllCalls))
			}

			if mockFS.MkdirAllCalls[0].Path != tt.expectedTarget {
				t.Errorf("MkdirAll called with %s, want %s", mockFS.MkdirAllCalls[0].Path, tt.expectedTarget)
			}

			// Verify file was moved to correct directory
			if len(mockFS.RenameCalls) != 1 {
				t.Fatalf("Expected 1 Rename call, got %d", len(mockFS.RenameCalls))
			}

			expectedNewPath := filepath.Join(tt.expectedTarget, filepath.Base(tt.path))
			if mockFS.RenameCalls[0].NewPath != expectedNewPath {
				t.Errorf("Rename called with newpath %s, want %s", mockFS.RenameCalls[0].NewPath, expectedNewPath)
			}

			if mockFS.RenameCalls[0].OldPath != tt.path {
				t.Errorf("Rename called with oldpath %s, want %s", mockFS.RenameCalls[0].OldPath, tt.path)
			}
		})
	}
}

// TestProcessMoveWithDeps_FilePathHandling tests various file path scenarios
func TestProcessMoveWithDeps_FilePathHandling(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		topic            string
		confidence       float64
		expectedBasename string
		expectedOldPath  string
	}{
		{
			name:             "simple filename",
			path:             "document.pdf",
			topic:            "docs",
			confidence:       0.8,
			expectedBasename: "document.pdf",
			expectedOldPath:  "document.pdf",
		},
		{
			name:             "absolute path",
			path:             "/home/user/downloads/photo.jpg",
			topic:            "images",
			confidence:       0.9,
			expectedBasename: "photo.jpg",
			expectedOldPath:  "/home/user/downloads/photo.jpg",
		},
		{
			name:             "relative path",
			path:             "folder/subfolder/video.mp4",
			topic:            "videos",
			confidence:       0.75,
			expectedBasename: "video.mp4",
			expectedOldPath:  "folder/subfolder/video.mp4",
		},
		{
			name:             "filename with spaces",
			path:             "/tmp/my document.txt",
			topic:            "text",
			confidence:       0.65,
			expectedBasename: "my document.txt",
			expectedOldPath:  "/tmp/my document.txt",
		},
		{
			name:             "filename with special characters",
			path:             "/data/file-name_v2.0.txt",
			topic:            "data",
			confidence:       0.7,
			expectedBasename: "file-name_v2.0.txt",
			expectedOldPath:  "/data/file-name_v2.0.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemOps{}
			writer := &bytes.Buffer{}

			ProcessMoveWithDeps(tt.path, tt.topic, tt.confidence, mockFS, writer)

			// Verify basename extraction
			if len(mockFS.RenameCalls) != 1 {
				t.Fatalf("Expected 1 Rename call, got %d", len(mockFS.RenameCalls))
			}

			actualBasename := filepath.Base(mockFS.RenameCalls[0].NewPath)
			if actualBasename != tt.expectedBasename {
				t.Errorf("Basename = %s, want %s", actualBasename, tt.expectedBasename)
			}

			// Verify original path is preserved
			if mockFS.RenameCalls[0].OldPath != tt.expectedOldPath {
				t.Errorf("OldPath = %s, want %s", mockFS.RenameCalls[0].OldPath, tt.expectedOldPath)
			}
		})
	}
}

func TestProcessMoveWithDeps_DirectoryPermissions(t *testing.T) {
	mockFS := &MockFileSystemOps{}
	writer := &bytes.Buffer{}

	ProcessMoveWithDeps("/source/file.txt", "documents", 0.8, mockFS, writer)

	if len(mockFS.MkdirAllCalls) != 1 {
		t.Fatalf("Expected 1 MkdirAll call, got %d", len(mockFS.MkdirAllCalls))
	}

	if mockFS.MkdirAllCalls[0].Perm != os.ModePerm {
		t.Errorf("MkdirAll called with perm %v, want %v", mockFS.MkdirAllCalls[0].Perm, os.ModePerm)
	}
}

func TestProcessMoveWithDeps_SuccessOutput(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		topic          string
		confidence     float64
		expectedOutput string
	}{
		{
			name:           "high confidence to topic",
			path:           "/source/report.pdf",
			topic:          "documents",
			confidence:     0.95,
			expectedOutput: "Organized: report.pdf -> documents (0.95)\n",
		},
		{
			name:           "low confidence to misc",
			path:           "/source/unknown.dat",
			topic:          "data",
			confidence:     0.45,
			expectedOutput: "Organized: unknown.dat -> misc (0.45)\n",
		},
		{
			name:           "exact threshold",
			path:           "/source/image.png",
			topic:          "images",
			confidence:     0.60,
			expectedOutput: "Organized: image.png -> images (0.60)\n",
		},
		{
			name:           "confidence with many decimals",
			path:           "/source/file.txt",
			topic:          "text",
			confidence:     0.876543,
			expectedOutput: "Organized: file.txt -> text (0.88)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemOps{}
			writer := &bytes.Buffer{}

			err := ProcessMoveWithDeps(tt.path, tt.topic, tt.confidence, mockFS, writer)

			if err != nil {
				t.Errorf("ProcessMoveWithDeps() unexpected error: %v", err)
				return
			}

			output := writer.String()
			if output != tt.expectedOutput {
				t.Errorf("ProcessMoveWithDeps() output = %q, want %q", output, tt.expectedOutput)
			}
		})
	}
}