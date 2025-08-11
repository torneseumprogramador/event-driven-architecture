package idempotency

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository implementação usando MongoDB
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository cria um novo repositório MongoDB
func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("processed_events"),
	}
}

// IsProcessed verifica se um evento já foi processado
func (r *MongoRepository) IsProcessed(ctx context.Context, eventID, serviceName string) (bool, error) {
	filter := bson.M{
		"event_id":     eventID,
		"service_name": serviceName,
	}
	
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// MarkAsProcessed marca um evento como processado
func (r *MongoRepository) MarkAsProcessed(ctx context.Context, eventID, serviceName string) error {
	processedEvent := bson.M{
		"event_id":     eventID,
		"service_name": serviceName,
		"processed_at": time.Now(),
	}
	
	_, err := r.collection.InsertOne(ctx, processedEvent)
	return err
}
