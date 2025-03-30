package crud

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoService struct {
	client   *mongo.Client
	database string
}

func NewMongoService(uri, database string) (Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to MongoDB")
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(err, "failed to ping MongoDB")
	}

	return &mongoService{
		client:   client,
		database: database,
	}, nil
}

func (s *mongoService) Create(ctx context.Context, level string, message string, metadata map[string]interface{}) error {
	collection := s.client.Database(s.database).Collection("logs")

	entry := NewLogEntry(level, message, metadata)

	_, err := collection.InsertOne(ctx, entry)
	if err != nil {
		return errors.Wrap(err, "failed to insert log entry")
	}

	return nil
}

func (s *mongoService) Get(ctx context.Context, id string) (*LogEntry, error) {
	collection := s.client.Database(s.database).Collection("logs")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(errBadRequest, "invalid log ID format")
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&entry)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get log entry")
	}

	return &entry, nil
}

func (s *mongoService) List(ctx context.Context, filter map[string]interface{}) ([]LogEntry, error) {
	collection := s.client.Database(s.database).Collection("logs")

	// Convert filter to MongoDB query
	query := bson.M{}
	for k, v := range filter {
		if k == "level" {
			query["level"] = v
		} else if k == "message" {
			query["message"] = v
		} else {
			query["metadata."+k] = v
		}
	}

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find log entries")
	}
	defer cursor.Close(ctx)

	var entries []LogEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, errors.Wrap(err, "failed to decode log entries")
	}

	return entries, nil
}

func (s *mongoService) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
