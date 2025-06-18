package cmd

import (
	"log"

	"github.com/vx6fid/job-crawler/trend_worker"
)

func main() {
	reports, err := trend_worker.AnalyzeTrendsByRole("")
	if err != nil {
		log.Fatalf("Trend analysis failed: %v", err)
	}
	log.Printf("Trend analysis completed successfully: %+v", reports)
}
