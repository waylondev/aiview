package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackwener/aiview/internal/platform"
)

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// handleTrend handles GET /api/trend requests.
func (s *Server) handleTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	platformName := r.URL.Query().Get("platform")
	recordType := r.URL.Query().Get("type")
	daysStr := r.URL.Query().Get("days")

	if platformName == "" {
		platformName = "bilibili"
	}
	if recordType == "" {
		recordType = "hot"
	}
	days := 7
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 && d <= 30 {
			days = d
		}
	}

	result, err := s.analyzer.AnalyzeTrend(platformName, recordType, days)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("analyze trend: %v", err))
		return
	}

	// Convert to JSON-friendly format
	type pointJSON struct {
		Date  string  `json:"date"`
		Value float64 `json:"value"`
	}
	points := make([]pointJSON, len(result.Points))
	for i, p := range result.Points {
		points[i] = pointJSON{Date: p.Label, Value: p.Value}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"platform": result.Platform,
		"type":     result.Type,
		"points":   points,
		"min":      result.Min,
		"max":      result.Max,
		"average":  result.Average,
		"change":   result.Change,
	})
}

// handlePlatforms handles GET /api/platforms requests.
func (s *Server) handlePlatforms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	platforms := platform.List()
	type platformStatus struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	statuses := make([]platformStatus, len(platforms))
	for i, name := range platforms {
		statuses[i] = platformStatus{
			Name:   name,
			Status: "active",
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"platforms": statuses,
		"count":     len(platforms),
	})
}

// handleSchedule handles GET /api/schedule requests.
func (s *Server) handleSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var jobs []map[string]interface{}
	if s.scheduler != nil {
		for _, job := range s.scheduler.ListJobs() {
			jobs = append(jobs, map[string]interface{}{
				"id":       job.ID,
				"interval": job.Interval.String(),
				"command":  job.Command,
				"last_run": job.LastRun.Format(time.RFC3339),
				"next_run": job.NextRun.Format(time.RFC3339),
				"running":  job.Running,
			})
		}
	}

	if jobs == nil {
		jobs = []map[string]interface{}{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"jobs":  jobs,
		"count": len(jobs),
	})
}

// handleHistory handles GET /api/history requests.
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	records, err := s.storage.QueryAll(limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("query history: %v", err))
		return
	}

	type historyItem struct {
		Platform    string `json:"platform"`
		Type        string `json:"type"`
		CollectedAt string `json:"collected_at"`
	}

	items := make([]historyItem, len(records))
	for i, r := range records {
		items[i] = historyItem{
			Platform:    r.Platform,
			Type:        r.Type,
			CollectedAt: r.CollectedAt.Format(time.RFC3339),
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": items,
		"count": len(items),
	})
}
