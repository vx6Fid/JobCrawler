package pkg

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var jobCollection *mongo.Collection

func EnsureTTLIndex(collection *mongo.Collection) error {
	index := mongo.IndexModel{
		Keys: bson.M{"CreatedAt": 1},
		Options: options.Index().
			SetExpireAfterSeconds(30 * 24 * 60 * 60). // 30 days
			SetName("CreatedAt_TTL"),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), index)
	if err != nil {
		log.Printf("[mongo] Failed to create TTL index: %v", err)
	} else {
		log.Println("[mongo] TTL index on 'CreatedAt' ensured.")
	}

	return err
}

func ConnectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		log.Fatal("DATABASE_URL not set in .env file")
	}

	clientOptions := options.Client().ApplyURI(database_url)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	jobCollection = client.Database("job_scraper").Collection("jobs")
	_ = EnsureTTLIndex(jobCollection) // Ensure TTL index on CreatedAt

	log.Println("[mongo] Connected to MongoDB")
	return nil
}

func UpsertJob(job JobPosting) error {
	log.Printf("[mongo] Upserting job: %s at %s", job.Title, job.Company)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	job.Hash = GenerateHash(job.Title, job.Company, job.Location, job.PostedOn.Format("2006-01-02"))
	job.DescriptionHash = GenerateHash(job.Description)

	filter := bson.M{"hash": job.Hash}

	var existing JobPosting
	err := jobCollection.FindOne(ctx, filter).Decode(&existing)

	job.LastUpdated = time.Now()

	// ToDo: this is a temporary fix. We should use the actual expiration date from the job posting.
	job.ExpireAt = job.PostedOn.AddDate(0, 0, 30)

	if err == nil {
		// Update only if description changed
		if existing.DescriptionHash != job.DescriptionHash {
			job.CreatedAt = existing.CreatedAt
			_, err = jobCollection.ReplaceOne(ctx, filter, job)
			return err
		}

		_, err = jobCollection.UpdateOne(ctx, filter, bson.M{
			"$set": bson.M{"lastUpdated": time.Now()},
		})
		log.Printf("[mongo] Job already exists, updated lastUpdated for: %s at %s", job.Title, job.Company)
		return err
	}

	job.CreatedAt = time.Now()
	_, err = jobCollection.InsertOne(ctx, job)
	log.Printf("[mongo] Job inserted: %s at %s", job.Title, job.Company)
	return err
}
