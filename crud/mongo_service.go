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

	// Build query
	query := bson.M{}

	// Add level filter if present
	if level, ok := filter["level"].(string); ok && level != "" {
		query["level"] = level
	}

	// Add message filter if present
	if message, ok := filter["message"].(string); ok && message != "" {
		query["message"] = bson.M{"$regex": message, "$options": "i"}
	}

	// Add time range filters if present
	if startTime, ok := filter["starttime"].(string); ok && startTime != "" {
		startTimestamp, err := strconv.ParseInt(startTime, 10, 64)
		if err != nil {
			return nil, errors.Wrap(errBadRequest, "invalid start time format")
		}
		query["timestamp"] = bson.M{"$gte": startTimestamp}
	}

	if endTime, ok := filter["endtime"].(string); ok && endTime != "" {
		endTimestamp, err := strconv.ParseInt(endTime, 10, 64)
		if err != nil {
			return nil, errors.Wrap(errBadRequest, "invalid end time format")
		}
		if _, exists := query["timestamp"]; exists {
			query["timestamp"].(bson.M)["$lte"] = endTimestamp
		} else {
			query["timestamp"] = bson.M{"$lte": endTimestamp}
		}
	}

	// Set up options for sorting and limiting
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Sort by timestamp in descending order

	// Handle recent parameter
	if recent, ok := filter["recent"].(string); ok && recent != "" {
		limit, err := strconv.ParseInt(recent, 10, 64)
		if err != nil {
			return nil, errors.Wrap(errBadRequest, "invalid recent value")
		}
		opts.SetLimit(limit)
	}

	// Execute query
	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query logs")
	}
	defer cursor.Close(ctx)

	var logs []LogEntry
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, errors.Wrap(err, "failed to decode logs")
	}

	return logs, nil
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

		// Delete all documents with timestamp less than the specified time
		result, err := collection.DeleteMany(ctx, bson.M{"timestamp": bson.M{"$lt": timestamp}})
		if err != nil {
			return errors.Wrap(err, "failed to delete log entries")
		}

		if result.DeletedCount == 0 {
			return errors.Wrap(errBadRequest, "no logs found before the specified timestamp")
		}

		return nil
	}

	return errors.Wrap(errBadRequest, "either id or before timestamp must be provided")
}
