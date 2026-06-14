// Package pipeline provides batch data collection and storage pipeline.
package pipeline

import (
	"fmt"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/storage"
)

// Collector defines the interface for data collection.
type Collector interface {
	Collect(recordType string) (map[string]interface{}, error)
}

// Pipeline orchestrates data collection and storage.
type Pipeline struct {
	collector Collector
	storage   storage.Storage
	platform  string
}

func New(platform string, collector Collector, store storage.Storage) *Pipeline {
	return &Pipeline{
		collector: collector,
		storage:   store,
		platform:  platform,
	}
}

// CollectAndStore collects data and stores it.
func (p *Pipeline) CollectAndStore(types []string) error {
	for _, t := range types {
		data, err := p.collector.Collect(t)
		if err != nil {
			return aiverr.Wrap(err, aiverr.CodeAPIError, p.platform)
		}

		record := storage.Record{
			Platform:    p.platform,
			Type:        t,
			Data:        data,
			CollectedAt: time.Now(),
		}

		if err := p.storage.Save(record); err != nil {
			return aiverr.Wrap(fmt.Errorf("save %s: %w", t, err), aiverr.CodeAPIError, "pipeline")
		}

		fmt.Printf("Collected and stored: %s/%s\n", p.platform, t)
	}
	return nil
}
