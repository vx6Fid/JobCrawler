package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vx6fid/job-crawler/internal/crawler"
	"github.com/vx6fid/job-crawler/trend_worker"
)

func TrendReportHandler(w http.ResponseWriter, r *http.Request) {
	role := strings.TrimSpace(r.URL.Query().Get("role"))

	if !crawler.IsRoleAllowed(role) {
		http.Error(w, "Invalid or unsupported role", http.StatusBadRequest)
		return
	}

	report, err := trend_worker.AnalyzeTrendsByRole(role)
	if err != nil {
		http.Error(w, "Failed to generate trend report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
