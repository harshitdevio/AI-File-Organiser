package scanner

import (
	"testing"
	"errors"
	"os"
)

func Test_GetDirectory(t *testing.T) {
    tempPath := t.TempDir()

    _, err := os.Stat(tempPath)
    
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            t.Errorf("Expected path %s to exist, but it didn't", tempPath)
        } else {
            t.Errorf("Unexpected error checking path: %v", err)
        }
    }
}