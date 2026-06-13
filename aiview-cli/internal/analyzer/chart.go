package analyzer

import (
	"fmt"
	"strings"
)

// RenderASCIIChart renders a trend result as an ASCII bar chart.
func RenderASCIIChart(result *TrendResult, width int) string {
	if len(result.Points) == 0 {
		return "No data available for the specified period."
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n📊 Trend: %s/%s (%d days)\n", result.Platform, result.Type, len(result.Points)))
	sb.WriteString(strings.Repeat("─", width+20) + "\n")

	// Find max value for scaling
	maxVal := result.Max
	if maxVal == 0 {
		maxVal = 1
	}

	barWidth := width - 12 // Reserve space for label and value

	for _, p := range result.Points {
		barLen := int((p.Value / maxVal) * float64(barWidth))
		if barLen < 0 {
			barLen = 0
		}

		label := p.Label[5:] // Show MM-DD only
		bar := strings.Repeat("█", barLen)
		if barLen == 0 && p.Value > 0 {
			bar = "▏"
		}

		sb.WriteString(fmt.Sprintf("  %s │%-*s %5.0f\n", label, barWidth, bar, p.Value))
	}

	sb.WriteString(strings.Repeat("─", width+20) + "\n")
	sb.WriteString(fmt.Sprintf("  Min: %.0f  Max: %.0f  Avg: %.1f  Change: %+.1f%%\n",
		result.Min, result.Max, result.Average, result.Change))

	return sb.String()
}

// RenderCompareTable renders cross-platform comparison as a table.
func RenderCompareTable(results []CompareResult, keyword string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n🔍 Cross-platform comparison: \"%s\"\n", keyword))
	sb.WriteString(strings.Repeat("─", 60) + "\n")
	sb.WriteString(fmt.Sprintf("  %-15s %-10s %-10s %s\n", "Platform", "Matches", "Score", "Top Items"))
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	for _, r := range results {
		topItems := ""
		if len(r.TopItems) > 0 {
			limit := 3
			if len(r.TopItems) < limit {
				limit = len(r.TopItems)
			}
			topItems = strings.Join(r.TopItems[:limit], ", ")
			if len(r.TopItems) > 3 {
				topItems += "..."
			}
		}
		sb.WriteString(fmt.Sprintf("  %-15s %-10d %-10.1f %s\n", r.Platform, r.Count, r.Score, topItems))
	}

	sb.WriteString(strings.Repeat("─", 60) + "\n")
	return sb.String()
}
