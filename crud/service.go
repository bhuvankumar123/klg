package crud

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("log entry not found")
	ErrEmptyKey = errors.New("validation failed, required fields missing")
)

type (
	Service interface {
		Create(ctx context.Context, level string, message string, metadata map[string]interface{}) error
		Get(ctx context.Context, id string) (*LogEntry, error)
		List(ctx context.Context, filter map[string]interface{}) ([]LogEntry, error)
		Close(ctx context.Context) error
	}

	LogEntry struct {
		ID        string                 `bson:"_id,omitempty"`
		Timestamp time.Time              `bson:"timestamp"`
		Level     string                 `bson:"level"`
		Message   string                 `bson:"message"`
		Metadata  map[string]interface{} `bson:"metadata,omitempty"`
	}

	defaultService struct {
		store map[string]*LogEntry
	}
)

func NewLogEntry(level, message string, metadata map[string]interface{}) *LogEntry {
	return &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Metadata:  metadata,
	}
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

func NewService() (Service, error) {
	return &defaultService{
		store: make(map[string]*LogEntry),
	}, nil
}
