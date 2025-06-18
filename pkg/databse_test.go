package pkg

import (
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestUpsertJob(t *testing.T) {
	/* Debugging: Print current working directory
	dir, _ := os.Getwd()
	log.Println("Current working directory:", dir)
	*/

	err := godotenv.Load(filepath.Join("..", ".env"))
	if err != nil {
		log.Println("Could not load .env file for test:", err)
	}

	err = ConnectMongo()
	if err != nil {
		t.Fatalf("MongoDB connection failed: %v", err)
	}

	job := JobPosting{
		Title:       "DevOps Engineer",
		Company:     "Chipcolate",
		Location:    "Remote",
		Salary:      "$75,000 - $99,999 USD",
		PostedOn:    time.Now().AddDate(0, 0, -1),
		Description: "Work with Docker, Kubernetes, Terraform...",
		URL:         "https://weworkremotely.com/jobs/devops-engineer",
		Source:      "WeWorkRemotely",
		ApplyURL:    "https://apply-link.com",
		Skills:      []string{"Docker", "Kubernetes", "Terraform"},
		Experience:  "3+ years",
	}

	err = UpsertJob(job)
	if err != nil {
		t.Fatalf("Failed to upsert job: %v", err)
	}
}
