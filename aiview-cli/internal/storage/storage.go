// Package storage provides data persistence backends for collected data.
package storage

import "time"

// Record represents a collected data record.
type Record struct {
	Platform    string                 `json:"platform"`
	Type        string                 `json:"type"`
	Data        map[string]interface{} `json:"data"`
	CollectedAt time.Time              `json:"collected_at"`
}

// Storage defines the interface for data persistence.
type Storage interface {
	Save(record Record) error
	Query(platform, recordType string, limit int) ([]Record, error)
	Close() error
}
