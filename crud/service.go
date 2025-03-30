package crud

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrNotFound     = errors.New("log entry not found")
	ErrEmptyKey     = errors.New("Bad Request, required fields missing")
	ErrInvalidLevel = errors.New("invalid log level")
)

// Valid log levels
var ValidLogLevels = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
	"fatal": true,
}

// Service interface defines the contract for log operations
type Service interface {
	Create(ctx context.Context, level string, message string, metadata map[string]interface{}) error
	Get(ctx context.Context, id string) (*LogEntry, error)
	List(ctx context.Context, filter map[string]interface{}) ([]LogEntry, error)
	Delete(ctx context.Context, filter map[string]interface{}) error
	Close(ctx context.Context) error
}

// LogEntry represents a log entry in the system
type LogEntry struct {
	ID        string                 `json:"id" bson:"_id,omitempty"`
	Timestamp int64                  `json:"timestamp" bson:"timestamp"`
	Level     string                 `json:"level" bson:"level"`
	Message   string                 `json:"message" bson:"message"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// NewLogEntry creates a new log entry with the current timestamp
func NewLogEntry(level string, message string, metadata map[string]interface{}) *LogEntry {
	return &LogEntry{
		Timestamp: time.Now().Unix(),
		Level:     level,
		Message:   message,
		Metadata:  metadata,
	}
}

// ValidateLogLevel checks if the given level is valid
func ValidateLogLevel(level string) error {
	if !ValidLogLevels[strings.ToLower(level)] {
		return errors.Wrapf(ErrInvalidLevel, "valid levels are: %v", getValidLevels())
	}
	return nil
}

// getValidLevels returns a slice of valid log levels
func getValidLevels() []string {
	levels := make([]string, 0, len(ValidLogLevels))
	for level := range ValidLogLevels {
		levels = append(levels, level)
	}
	return levels
}

// defaultService implements the Service interface using in-memory storage
type defaultService struct {
	store map[string]*LogEntry
}

func (s *defaultService) Create(
	ctx context.Context, level string, message string, metadata map[string]interface{},
) error {
	if level == "" || message == "" {
		return ErrEmptyKey
	}

	entry := NewLogEntry(level, message, metadata)
	s.store[entry.ID] = entry
	return nil
}

func (s *defaultService) Get(ctx context.Context, id string) (*LogEntry, error) {
	if entry, ok := s.store[id]; ok {
		return entry, nil
	}
	return nil, ErrNotFound
}

func (s *defaultService) List(ctx context.Context, filter map[string]interface{}) ([]LogEntry, error) {
	entries := make([]LogEntry, 0)
	for _, entry := range s.store {
		match := true
		for k, v := range filter {
			if entry.Metadata[k] != v {
				match = false
				break
			}
		}
		if match {
			entries = append(entries, *entry)
		}
	}
	return entries, nil
}

func (s *defaultService) Close(ctx context.Context) error {
	s.store = make(map[string]*LogEntry)
	return nil
}

func (s *defaultService) Delete(ctx context.Context, filter map[string]interface{}) error {
	if id, ok := filter["id"].(string); ok && id != "" {
		delete(s.store, id)
	}
	return nil
}

func NewService() (Service, error) {
	return &defaultService{
		store: make(map[string]*LogEntry),
	}, nil
}
