package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/vx6fid/job-crawler/internal/crawler"
)

func CrawlHandler(w http.ResponseWriter, r *http.Request) {
	roles := r.URL.Query()["role"] // allows multiple ?role=dev&role=ml

	var validRoles []string
	for _, role := range roles {
		trimmed := strings.TrimSpace(role)
		if crawler.IsRoleAllowed(trimmed) {
			validRoles = append(validRoles, trimmed)
		}
	}

	if len(validRoles) == 0 {
		http.Error(w, "No valid roles provided", http.StatusBadRequest)
		return
	}

	// run crawl in background
	go func() {
		_ = crawler.StartCrawling(validRoles, 50, 60*time.Second)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Crawling started",
		"roles":   validRoles,
	})
}
