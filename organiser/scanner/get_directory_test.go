package scanner

import (
	"testing"
)

func TestGetDirectory(t *testing.T) {
	tests := []struct {
		name      string
		pathFlag  string
		input     string
		want      string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "path flag provided",
			pathFlag:  "/home/user/documents",
			input:     "",
			want:      "/home/user/documents",
			wantError: false,
		},
		{
			name:      "path flag takes precedence over input",
			pathFlag:  "/home/user/flag-path",
			input:     "/home/user/input-path",
			want:      "/home/user/flag-path",
			wantError: false,
		},
		{
			name:      "input provided when flag is empty",
			pathFlag:  "",
			input:     "/home/user/documents",
			want:      "/home/user/documents",
			wantError: false,
		},
		{
			name:      "input with leading whitespace",
			pathFlag:  "",
			input:     "  /home/user/documents",
			want:      "/home/user/documents",
			wantError: false,
		},
		{
			name:      "input with trailing whitespace",
			pathFlag:  "",
			input:     "/home/user/documents  ",
			want:      "/home/user/documents",
			wantError: false,
		},
		{
			name:      "input with leading and trailing whitespace",
			pathFlag:  "",
			input:     "  /home/user/documents  \n",
			want:      "/home/user/documents",
			wantError: false,
		},
		{
			name:      "both empty returns error",
			pathFlag:  "",
			input:     "",
			want:      "",
			wantError: true,
			errorMsg:  "No directory provided",
		},
		{
			name:      "input is only whitespace",
			pathFlag:  "",
			input:     "   \n\t  ",
			want:      "",
			wantError: true,
			errorMsg:  "No directory provided",
		},
		{
			name:      "input is only newline",
			pathFlag:  "",
			input:     "\n",
			want:      "",
			wantError: true,
			errorMsg:  "No directory provided",
		},
	}
 
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDirectory(tt.pathFlag, tt.input)
 
			if tt.wantError {
				if err == nil {
					t.Errorf("GetDirectory() expected error but got nil")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("GetDirectory() error = %v, want %v", err.Error(), tt.errorMsg)
				}
				return
			}
 
			if err != nil {
				t.Errorf("GetDirectory() unexpected error: %v", err)
				return
			}
 
			if got != tt.want {
				t.Errorf("GetDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}
 