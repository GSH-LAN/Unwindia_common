package workitemLock

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

const (
	lockCollectionPrefix = "lock"
)

// MongoWorkItemLock is a workitem lock implementation based on MongoDB
type MongoWorkItemLock struct {
	lockName     string
	lockId       string
	dbCollection *mongo.Collection
}

func NewMongoWorkItemLock(lockName, lockId string, db *mongo.Database) *MongoWorkItemLock {
	lock := &MongoWorkItemLock{
		lockName:     lockName,
		lockId:       lockId,
		dbCollection: db.Collection(fmt.Sprintf("%s_%s", lockCollectionPrefix, lockName)),
	}

	go lock.StartHousekeeping()

	return lock
}

func (w *MongoWorkItemLock) housekeeping() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := w.dbCollection.DeleteMany(ctx, bson.M{"expiresAt": bson.M{"$lt": time.Now()}})
	if err != nil {
		log.Error().Err(err).Msg("Error while housekeeping")
	}
}

func (w *MongoWorkItemLock) Lock(ctx context.Context, workitemID string, ttl *time.Duration) error {

	expiresAfter := defaultTtl
	if ttl != nil {
		expiresAfter = *ttl
	}

	expiresAt := time.Now().Add(expiresAfter)

	doc := &WorkItemLockEntry{
		ID:        workitemID,
		LockedBy:  w.lockId,
		CreatedAt: time.Now(),
		ExpiresAt: &expiresAt,
	}

	// TTL index
	index := mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "createdAt", Value: bsonx.Int32(1)}},
		Options: options.Index().SetExpireAfterSeconds(int32(expiresAfter.Seconds())),
	}

	_, err := w.dbCollection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	_, err = w.dbCollection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

func (w *MongoWorkItemLock) Unlock(ctx context.Context, workitemID string) error {
	deleteResult, err := w.dbCollection.DeleteOne(ctx, bson.M{"_id": workitemID})
	if err != nil {
		return err
	}
	log.Trace().Str("workitemID", workitemID).Int64("deletedCount", deleteResult.DeletedCount).Msg("Removed lock for workitem")
	return nil
}

func (w *MongoWorkItemLock) StartHousekeeping() {
	ticker := time.NewTicker(housekeepingInterval)
	for {
		select {
		case <-ticker.C:
			w.housekeeping()
		}
	}
}
