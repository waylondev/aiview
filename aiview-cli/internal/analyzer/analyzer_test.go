package analyzer

import (
	"testing"
	"time"

	"github.com/jackwener/aiview/internal/storage"
)

// mockStorage implements storage.Storage for testing.
type mockStorage struct {
	records []storage.Record
}

func (m *mockStorage) Save(record storage.Record) error {
	m.records = append(m.records, record)
	return nil
}

func (m *mockStorage) Query(platform, recordType string, limit int) ([]storage.Record, error) {
	var result []storage.Record
	for _, r := range m.records {
		if r.Platform == platform && r.Type == recordType {
			result = append(result, r)
		}
	}
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (m *mockStorage) QueryAll(limit int) ([]storage.Record, error) {
	var result []storage.Record
	for _, r := range m.records {
		result = append(result, r)
	}
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}

func (m *mockStorage) Close() error { return nil }

func TestAnalyzeTrend_Empty(t *testing.T) {
	store := &mockStorage{}
	a := New(store)

	result, err := a.AnalyzeTrend("bilibili", "hot", 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Points) != 7 {
		t.Errorf("expected 7 points, got %d", len(result.Points))
	}
}

func TestAnalyzeTrend_WithData(t *testing.T) {
	store := &mockStorage{}
	// Add records for today and yesterday
	store.Save(storage.Record{
		Platform:    "bilibili",
		Type:        "hot",
		Data:        map[string]interface{}{"items": []interface{}{}},
		CollectedAt: time.Now(),
	})
	store.Save(storage.Record{
		Platform:    "bilibili",
		Type:        "hot",
		Data:        map[string]interface{}{"items": []interface{}{}},
		CollectedAt: time.Now().AddDate(0, 0, -1),
	})

	a := New(store)
	result, err := a.AnalyzeTrend("bilibili", "hot", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Average < 0 {
		t.Errorf("expected non-negative average, got %f", result.Average)
	}
}

func TestRenderASCIIChart_Empty(t *testing.T) {
	result := &TrendResult{Platform: "test", Type: "hot"}
	chart := RenderASCIIChart(result, 60)
	if chart == "" {
		t.Error("expected non-empty chart string")
	}
}

func TestRenderCompareTable(t *testing.T) {
	results := []CompareResult{
		{Platform: "bilibili", Count: 10, Score: 10, TopItems: []string{"item1"}},
		{Platform: "douyin", Count: 5, Score: 5, TopItems: []string{"item2"}},
	}
	table := RenderCompareTable(results, "test")
	if table == "" {
		t.Error("expected non-empty table string")
	}
}
