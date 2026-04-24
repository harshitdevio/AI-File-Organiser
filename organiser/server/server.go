package server

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "organiser/mover" 
)

type ClassificationResult struct {
    Filepath   string             `json:"filepath"`
    TopTopic   string             `json:"top_topic"`
    Confidence float64            `json:"confidence"`
    AllScores  map[string]float64 `json:"all_scores"`
    Error      *string            `json:"error"`
}

func IngestHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading body", http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()

    var results []any
    if err := json.Unmarshal(body, &results); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Passing the results to the Mover module
    mover.MoveFiles(results)

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"status":"processed","count":%d}`, len(results))
}