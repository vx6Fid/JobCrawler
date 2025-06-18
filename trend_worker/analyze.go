package trend_worker

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type CountResult struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type TrendReport struct {
	TopSkills              []CountResult `json:"top_skills"`
	TopLocations           []CountResult `json:"top_locations"`
	TopCompanies           []CountResult `json:"top_companies"`
	ExperienceDistribution []CountResult `json:"experience_distribution"`
}

func AnalyzeTrendsByRole(role string) (*TrendReport, error) {
	coll, err := getMongoCollection()
	if err != nil {
		return nil, err
	}

	match := bson.D{}
	if role != "" {
		// Case-insensitive partial match on title
		match = bson.D{{Key: "title", Value: bson.D{{Key: "$regex", Value: role}, {Key: "$options", Value: "i"}}}}
	}

	skills, err := countAggregation(coll, "skills", match)
	if err != nil {
		return nil, err
	}

	locations, err := countAggregation(coll, "location", match)
	if err != nil {
		return nil, err
	}

	experience, err := countAggregation(coll, "experience", match)
	if err != nil {
		return nil, err
	}

	companies, err := countAggregation(coll, "company", match)
	if err != nil {
		return nil, err
	}

	report := &TrendReport{
		TopSkills:              skills,
		TopLocations:           locations,
		TopCompanies:           companies,
		ExperienceDistribution: experience,
	}

	return report, nil
}

func getMongoCollection() (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		log.Fatal("DATABASE_URL not set in .env file")
	}

	clientOptions := options.Client().ApplyURI(database_url)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client.Database("job_scraper").Collection("jobs"), nil
}

func countAggregation(coll *mongo.Collection, field string, match bson.D) ([]CountResult, error) {
	agg := []bson.M{}
	if len(match) > 0 {
		agg = append(agg, bson.M{"$match": match})
	}
	agg = append(agg,
		bson.M{"$unwind": "$" + field},
		bson.M{"$group": bson.M{"_id": "$" + field, "count": bson.M{"$sum": 1}}},
		bson.M{"$sort": bson.M{"count": -1}},
		bson.M{"$limit": 20},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cursor, err := coll.Aggregate(ctx, agg)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []CountResult
	for cursor.Next(ctx) {
		var r struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&r); err != nil {
			return nil, err
		}
		results = append(results, CountResult{Value: r.ID, Count: r.Count})
	}

	return results, nil
}
