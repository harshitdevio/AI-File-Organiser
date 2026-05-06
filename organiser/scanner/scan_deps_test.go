package scanner 

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestScanWithDeps(t *testing.T) {
	tests := []struct {
		name         string
		pathFlag     string
		input        string
		want         string
		wantError    bool
		errorMsg     string
		wantOutput   string
	}{
		{
			name:       "path flag provided, no input needed",
			pathFlag:   "/home/user/documents",
			input:      "",
			want:       "/home/user/documents",
			wantError:  false,
			wantOutput: "Working on directory: /home/user/documents\n",
		},
		{
			name:       "no path flag, read from input",
			pathFlag:   "",
			input:      "/home/user/downloads\n",
			want:       "/home/user/downloads",
			wantError:  false,
			wantOutput: "Enter directory path: Working on directory: /home/user/downloads\n",
		},
		{
			name:       "input with extra whitespace",
			pathFlag:   "",
			input:      "  /home/user/music  \n",
			want:       "/home/user/music",
			wantError:  false,
			wantOutput: "Enter directory path: Working on directory: /home/user/music\n",
		},
		{
			name:       "empty input returns error",
			pathFlag:   "",
			input:      "\n",
			want:       "",
			wantError:  true,
			errorMsg:   "No directory provided",
			wantOutput: "Enter directory path: ",
		},
		{
			name:       "whitespace only input returns error",
			pathFlag:   "",
			input:      "   \n",
			want:       "",
			wantError:  true,
			errorMsg:   "No directory provided",
			wantOutput: "Enter directory path: ",
		},
		{
			name:       "input without newline",
			pathFlag:   "",
			input:      "/home/user/videos",
			want:       "/home/user/videos",
			wantError:  false,
			wantOutput: "Enter directory path: Working on directory: /home/user/videos\n",
		},
	}
 
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			writer := &bytes.Buffer{}
 
			got, err := ScanWithDeps(tt.pathFlag, reader, writer)
 
			if tt.wantError {
				if err == nil {
					t.Errorf("ScanWithDeps() expected error but got nil")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("ScanWithDeps() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ScanWithDeps() unexpected error: %v", err)
					return
				}
			}
 
			if got != tt.want {
				t.Errorf("ScanWithDeps() = %v, want %v", got, tt.want)
			}
 
			gotOutput := writer.String()
			if gotOutput != tt.wantOutput {
				t.Errorf("ScanWithDeps() output = %q, want %q", gotOutput, tt.wantOutput)
			}
		})
	}
}
 
// simulating a reader that returns an error via mockReader
type mockReader struct {
	err error
}
 
func (m *mockReader) Read(p []byte) (n int, err error) {
	return 0, m.err
}

func TestScanWithDeps_ReadError(t *testing.T) {
	mockErr := errors.New("mock read error")
	reader := &mockReader{err: mockErr}
	writer := &bytes.Buffer{}
 
	_, err := ScanWithDeps("", reader, writer)
 
	if err == nil {
		t.Fatal("ScanWithDeps() expected error but got nil")
	}
 
	if !strings.Contains(err.Error(), "failed to read input") {
		t.Errorf("ScanWithDeps() error = %v, want error containing 'failed to read input'", err)
	}
}
 
func TestScanWithDeps_PathFlagSkipsInput(t *testing.T) {
	reader := &mockReader{err: errors.New("should not read")}
	writer := &bytes.Buffer{}
 
	got, err := ScanWithDeps("/test/path", reader, writer)
 
	if err != nil {
		t.Errorf("ScanWithDeps() unexpected error: %v", err)
	}
 
	if got != "/test/path" {
		t.Errorf("ScanWithDeps() = %v, want %v", got, "/test/path")
	}
 
	output := writer.String()
	if strings.Contains(output, "Enter directory path") {
		t.Error("ScanWithDeps() should not prompt when path flag is provided")
	}
}