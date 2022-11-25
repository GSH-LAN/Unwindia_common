package mongodb

import "time"

type MongoBaseDocument struct {
	ID        string    `bson:"_id"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty"`
}
