package mover

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestProcessMoveWithDeps_MkdirAllError(t *testing.T) {
	tests := []struct {
		name       string
		mkdirError error
		wantError  bool
		errorMsg   string
	}{
		{
			name:       "permission denied",
			mkdirError: os.ErrPermission,
			wantError:  true,
			errorMsg:   "failed to create directory",
		},
		{
			name:       "disk full",
			mkdirError: errors.New("no space left on device"),
			wantError:  true,
			errorMsg:   "failed to create directory",
		},
		{
			name:       "generic error",
			mkdirError: errors.New("unknown error"),
			wantError:  true,
			errorMsg:   "failed to create directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemOps{
				MkdirAllFunc: func(path string, perm os.FileMode) error {
					return tt.mkdirError
				},
			}
			writer := &bytes.Buffer{}

			err := ProcessMoveWithDeps("/source/file.txt", "documents", 0.8, mockFS, writer)

			if tt.wantError {
				if err == nil {
					t.Error("ProcessMoveWithDeps() expected error but got nil")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("ProcessMoveWithDeps() error = %v, want error containing %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ProcessMoveWithDeps() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestProcessMoveWithDeps_RenameError(t *testing.T) {
	tests := []struct {
		name        string
		renameError error
		expectLog   string
	}{
		{
			name:        "file not found",
			renameError: os.ErrNotExist,
			expectLog:   "Error moving file.txt:",
		},
		{
			name:        "permission denied",
			renameError: os.ErrPermission,
			expectLog:   "Error moving file.txt:",
		},
		{
			name:        "cross-device link",
			renameError: errors.New("invalid cross-device link"),
			expectLog:   "Error moving file.txt:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &MockFileSystemOps{
				RenameFunc: func(oldpath, newpath string) error {
					return tt.renameError
				},
			}
			writer := &bytes.Buffer{}

			err := ProcessMoveWithDeps("/source/file.txt", "documents", 0.8, mockFS, writer)

			// Should return the error
			if err == nil {
				t.Error("ProcessMoveWithDeps() expected error but got nil")
				return
			}

			// Should log the error
			output := writer.String()
			if !strings.Contains(output, tt.expectLog) {
				t.Errorf("ProcessMoveWithDeps() output = %q, want output containing %q", output, tt.expectLog)
			}
		})
	}
}