package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SendToService hides the HTTP implementation from the rest of the app
func SendToService(targetURL string, data interface{}) error {
	client := &http.Client{Timeout: 30 * time.Second}

	jsonData, err := json.Marshal(data)
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