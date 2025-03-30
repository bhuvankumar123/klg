package crud

import (
	"context"
	"strconv"
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
	timeRange := bson.M{}
	limit := int64(50) // default limit

	for k, v := range filter {
		if k == "level" {
			query["level"] = v
		} else if k == "message" {
			// Use regex for case-insensitive message matching
			query["message"] = bson.M{"$regex": v, "$options": "i"}
		} else if k == "starttime" {
			timeRange["$gte"] = v
		} else if k == "endtime" {
			timeRange["$lte"] = v
		} else if k == "recent" {
			if count, ok := v.(float64); ok {
				limit = int64(count)
			}
		} else {
			query["metadata."+k] = v
		}
	}

	// Add time range if either start or end time is present
	if len(timeRange) > 0 {
		query["timestamp"] = timeRange
	}

	// Set up options for sorting by timestamp in descending order and limit
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, query, opts)
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

func (s *mongoService) Delete(ctx context.Context, filter map[string]interface{}) error {
	collection := s.client.Database(s.database).Collection("logs")

	// If ID is present, delete specific document
	if id, ok := filter["id"].(string); ok && id != "" {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.Wrap(errBadRequest, "invalid log ID format")
		}
		_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			return errors.Wrap(err, "failed to delete log entry")
		}
		return nil
	}

	// If before timestamp is present, delete all documents before that time
	if beforeTime, ok := filter["before"].(string); ok && beforeTime != "" {
		timestamp, err := strconv.ParseInt(beforeTime, 10, 64)
		if err != nil {
			return errors.Wrap(errBadRequest, "invalid epoch timestamp format")
		}
		_, err = collection.DeleteMany(ctx, bson.M{"timestamp": bson.M{"$lt": timestamp}})
		if err != nil {
			return errors.Wrap(err, "failed to delete log entries")
		}
		return nil
	}

	return errors.Wrap(errBadRequest, "either id or before timestamp must be provided")
}
