package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BatchRequest wraps file data in the structure expected by Python API
type BatchRequest struct {
	Files interface{} `json:"files"`
}

// SendToService hides the HTTP implementation from the rest of the app
func SendToService(targetURL string, data interface{}) error {
	client := &http.Client{Timeout: 30 * time.Second}

	// Wrap data in {files: [...]} structure
	batchReq := BatchRequest{Files: data}

	jsonData, err := json.Marshal(batchReq)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	resp, err := client.Post(targetURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned error: %d", resp.StatusCode)
	}

	return nil
}