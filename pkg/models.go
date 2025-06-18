package pkg

import "time"

type JobPosting struct {
	ID          string    `bson:"_id,omitempty"`
	Title       string    `bson:"title"`
	Company     string    `bson:"company"`
	Location    string    `bson:"location"`
	Salary      string    `bson:"salary"`
	PostedOn    time.Time `bson:"postedOn"`
	Description string    `bson:"description"`
	URL         string    `bson:"url"`
	Source      string    `bson:"source"`
	ApplyURL    string    `bson:"applyUrl"`
	Skills      []string  `bson:"skills"`
	Experience  string    `bson:"experience"`

	Hash            string    `bson:"hash"`            // hash of Title+Company+Location+PostedOn
	DescriptionHash string    `bson:"descriptionHash"` // hash of Description + Salary
	LastUpdated     time.Time `bson:"lastUpdated"`
	CreatedAt       time.Time `bson:"createdAt"`
	ExpireAt        time.Time `bson:"expireAt"` // TTL field (PostedOn + 30 days)
}
