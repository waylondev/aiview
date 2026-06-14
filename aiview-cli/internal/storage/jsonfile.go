package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
)

// JSONFileStorage stores records as JSON files.
type JSONFileStorage struct {
	dir string
	mu  sync.Mutex
}

func NewJSONFileStorage(dir string) (*JSONFileStorage, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}
	return &JSONFileStorage{dir: dir}, nil
}

func (s *JSONFileStorage) Save(record Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filename := filepath.Join(s.dir, fmt.Sprintf("%s_%s_%s.json",
		record.Platform, record.Type,
		time.Now().Format("20060102_150405")))

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return aiverr.Wrap(fmt.Errorf("marshal record: %w", err), aiverr.CodeParseError, "storage")
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *JSONFileStorage) Query(platform, recordType string, limit int) ([]Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, aiverr.Wrap(fmt.Errorf("read dir: %w", err), aiverr.CodeAPIError, "storage")
	}

	var records []Record
	prefix := platform + "_" + recordType + "_"

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) < len(prefix) || name[:len(prefix)] != prefix {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.dir, name))
		if err != nil {
			continue
		}

		var record Record
		if err := json.Unmarshal(data, &record); err != nil {
			continue
		}
		records = append(records, record)

		if limit > 0 && len(records) >= limit {
			break
		}
	}
	return records, nil
}

func (s *JSONFileStorage) QueryAll(limit int) ([]Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, aiverr.Wrap(fmt.Errorf("read dir: %w", err), aiverr.CodeAPIError, "storage")
	}

	var records []Record
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.dir, entry.Name()))
		if err != nil {
			continue
		}
		var record Record
		if err := json.Unmarshal(data, &record); err != nil {
			continue
		}
		records = append(records, record)
		if limit > 0 && len(records) >= limit {
			break
		}
	}
	return records, nil
}

func (s *JSONFileStorage) Close() error {
	return nil
}
