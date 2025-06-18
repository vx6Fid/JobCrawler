package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/vx6fid/job-crawler/api_server/routes"
	"github.com/vx6fid/job-crawler/pkg"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	if err := pkg.ConnectMongo(); err != nil {
		log.Fatal("[mongo] MongoDB connection failed:", err)
	}

	routes.RegisterRoutes()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("api_server/static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api_server/templates/index.html")
	})
	log.Println("[api] API server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
