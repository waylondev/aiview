// Package analyzer provides data analysis capabilities for collected platform data.
package analyzer

import (
	"fmt"
	"sort"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
	"github.com/jackwener/aiview/internal/storage"
)

// TrendPoint represents a single data point in a trend.
type TrendPoint struct {
	Time  time.Time
	Value float64
	Label string
}

// TrendResult represents the result of trend analysis.
type TrendResult struct {
	Points   []TrendPoint
	Min      float64
	Max      float64
	Average  float64
	Change   float64 // percentage change from first to last
	Platform string
	Type     string
}

// CompareResult represents cross-platform comparison.
type CompareResult struct {
	Platform string
	Count    int
	Score    float64
	TopItems []string
}

// Analyzer performs data analysis on stored records.
type Analyzer struct {
	store storage.Storage
}

// New creates a new Analyzer.
func New(store storage.Storage) *Analyzer {
	return &Analyzer{store: store}
}

// AnalyzeTrend analyzes trend data for a platform and type over the given days.
func (a *Analyzer) AnalyzeTrend(platform, recordType string, days int) (*TrendResult, error) {
	records, err := a.store.Query(platform, recordType, 0) // 0 = no limit
	if err != nil {
		return nil, aiverr.Wrap(fmt.Errorf("query records: %w", err), aiverr.CodeAPIError, "analyzer")
	}

	// Group records by date
	cutoff := time.Now().AddDate(0, 0, -days)
	dateCounts := make(map[string]int)
	for _, r := range records {
		if r.CollectedAt.Before(cutoff) {
			continue
		}
		dateKey := r.CollectedAt.Format("2006-01-02")
		dateCounts[dateKey]++
	}

	// Build trend points
	var points []TrendPoint
	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		count := float64(dateCounts[date])
		points = append(points, TrendPoint{
			Time:  time.Now().AddDate(0, 0, -i),
			Value: count,
			Label: date,
		})
	}

	// Calculate statistics
	result := &TrendResult{
		Points:   points,
		Platform: platform,
		Type:     recordType,
	}

	if len(points) > 0 {
		var sum float64
		result.Min = points[0].Value
		result.Max = points[0].Value
		for _, p := range points {
			sum += p.Value
			if p.Value < result.Min {
				result.Min = p.Value
			}
			if p.Value > result.Max {
				result.Max = p.Value
			}
		}
		result.Average = sum / float64(len(points))

		if points[0].Value > 0 {
			result.Change = ((points[len(points)-1].Value - points[0].Value) / points[0].Value) * 100
		}
	}

	return result, nil
}

// ComparePlatforms compares data across multiple platforms for a keyword.
func (a *Analyzer) ComparePlatforms(keyword string, platforms []string) ([]CompareResult, error) {
	var results []CompareResult

	for _, p := range platforms {
		records, err := a.store.Query(p, "search", 100)
		if err != nil {
			continue
		}

		count := 0
		var topItems []string
		for _, r := range records {
			// Simple keyword matching in data
			dataStr := fmt.Sprintf("%v", r.Data)
			if containsKeyword(dataStr, keyword) {
				count++
				if title, ok := extractTitle(r.Data); ok {
					topItems = append(topItems, title)
				}
			}
		}

		score := float64(count)
		results = append(results, CompareResult{
			Platform: p,
			Count:    count,
			Score:    score,
			TopItems: topItems,
		})
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func containsKeyword(data, keyword string) bool {
	return len(data) > 0 && len(keyword) > 0 &&
		(len(data) < 10000 && containsIgnoreCase(data, keyword))
}

func containsIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	subLower := toLower(substr)
	return contains(sLower, subLower)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func extractTitle(data map[string]interface{}) (string, bool) {
	// Try common title fields
	for _, key := range []string{"title", "keyword", "word", "name", "caption"} {
		if v, ok := data[key]; ok {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
	}
	return "", false
}
