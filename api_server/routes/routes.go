package routes

import (
	"net/http"

	"github.com/vx6fid/job-crawler/api_server/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/api/trends", handlers.TrendReportHandler)
	http.HandleFunc("/api/crawl", handlers.CrawlHandler)
}
