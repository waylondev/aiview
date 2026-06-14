package pipeline

import (
	"errors"
	"testing"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/storage"
)

// mockCollector implements Collector for testing.
type mockCollector struct {
	data map[string]interface{}
	err  error
}

func (m *mockCollector) Collect(recordType string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

// mockStorage implements storage.Storage for testing.
type mockStorage struct {
	saved      []storage.Record
	saveErr    error
	queryData  []storage.Record
	queryErr   error
}

func (m *mockStorage) Save(record storage.Record) error {
	m.saved = append(m.saved, record)
	return m.saveErr
}

func (m *mockStorage) Query(platform, recordType string, limit int) ([]storage.Record, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return m.queryData, nil
}

func (m *mockStorage) QueryAll(limit int) ([]storage.Record, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return m.queryData, nil
}

func (m *mockStorage) Close() error {
	return nil
}

func TestNewPipeline(t *testing.T) {
	col := &mockCollector{data: map[string]interface{}{}}
	store := &mockStorage{}
	p := New("test-platform", col, store)

	if p == nil {
		t.Fatal("expected non-nil Pipeline")
	}
	if p.platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", p.platform)
	}
}

func TestAddStep(t *testing.T) {
	// Pipeline doesn't have AddStep method; the concept is tested via CollectAndStore
	// which effectively runs a sequence of collection steps.
	col := &mockCollector{data: map[string]interface{}{"key": "value"}}
	store := &mockStorage{}
	p := New("test-platform", col, store)

	err := p.CollectAndStore([]string{"hot", "trending"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.saved) != 2 {
		t.Errorf("expected 2 saved records, got %d", len(store.saved))
	}
}

func TestExecute(t *testing.T) {
	t.Run("successful collection and storage", func(t *testing.T) {
		col := &mockCollector{data: map[string]interface{}{"count": 42}}
		store := &mockStorage{}
		p := New("test-platform", col, store)

		err := p.CollectAndStore([]string{"hot"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(store.saved) != 1 {
			t.Fatalf("expected 1 saved record, got %d", len(store.saved))
		}
		rec := store.saved[0]
		if rec.Platform != "test-platform" {
			t.Errorf("expected platform %q, got %q", "test-platform", rec.Platform)
		}
		if rec.Type != "hot" {
			t.Errorf("expected type %q, got %q", "hot", rec.Type)
		}
		val, ok := rec.Data["count"]
		if !ok || val != 42 {
			t.Errorf("expected data count=42, got %v", rec.Data)
		}
	})

	t.Run("collector error wraps to PlatformError", func(t *testing.T) {
		col := &mockCollector{err: errors.New("collect failed")}
		store := &mockStorage{}
		p := New("test-platform", col, store)

		err := p.CollectAndStore([]string{"hot"})
		if err == nil {
			t.Fatal("expected error from collector")
		}
		pe, ok := err.(*aiverr.PlatformError)
		if !ok {
			t.Fatalf("expected *PlatformError, got %T", err)
		}
		if pe.Code != aiverr.CodeAPIError {
			t.Errorf("expected code %q, got %q", aiverr.CodeAPIError, pe.Code)
		}
	})

	t.Run("storage error wraps to PlatformError", func(t *testing.T) {
		col := &mockCollector{data: map[string]interface{}{"x": 1}}
		store := &mockStorage{saveErr: errors.New("save failed")}
		p := New("test-platform", col, store)

		err := p.CollectAndStore([]string{"hot"})
		if err == nil {
			t.Fatal("expected error from storage")
		}
		pe, ok := err.(*aiverr.PlatformError)
		if !ok {
			t.Fatalf("expected *PlatformError, got %T", err)
		}
		if pe.Code != aiverr.CodeAPIError {
			t.Errorf("expected code %q, got %q", aiverr.CodeAPIError, pe.Code)
		}
	})

	t.Run("empty types produces no error", func(t *testing.T) {
		col := &mockCollector{data: map[string]interface{}{}}
		store := &mockStorage{}
		p := New("test-platform", col, store)

		err := p.CollectAndStore([]string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(store.saved) != 0 {
			t.Errorf("expected 0 saved records, got %d", len(store.saved))
		}
	})
}