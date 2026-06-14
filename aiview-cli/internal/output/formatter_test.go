package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

type testVideo struct {
	Title string
	Views int
	URL   string
}

func TestResolveFormat_JSON(t *testing.T) {
	f := ResolveFormat(true, false, false, false)
	if f != FormatJSON {
		t.Errorf("expected FormatJSON, got %d", f)
	}
}

func TestResolveFormat_YAML(t *testing.T) {
	f := ResolveFormat(false, true, false, false)
	if f != FormatYAML {
		t.Errorf("expected FormatYAML, got %d", f)
	}
}

func TestResolveFormat_Default(t *testing.T) {
	// When no format is specified and not a TTY, defaults to YAML
	f := ResolveFormat(false, false, false, false)
	if f != FormatYAML && f != FormatTable {
		t.Errorf("expected FormatYAML or FormatTable, got %d", f)
	}
}

func TestResolveFormat_Table(t *testing.T) {
	f := ResolveFormat(false, false, true, false)
	if f != FormatTable {
		t.Errorf("expected FormatTable, got %d", f)
	}
}

func TestResolveFormat_CSV(t *testing.T) {
	f := ResolveFormat(false, false, false, true)
	if f != FormatCSV {
		t.Errorf("expected FormatCSV, got %d", f)
	}
}

func TestFormatCount(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{100, "100"},
		{9999, "9999"},
		{10000, "1.0万"},
		{15000, "1.5万"},
		{100000, "10.0万"},
	}
	for _, tt := range tests {
		got := FormatCount(tt.input)
		if got != tt.expected {
			t.Errorf("FormatCount(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "00:00"},
		{59, "00:59"},
		{60, "01:00"},
		{3599, "59:59"},
		{3600, "1:00:00"},
		{3661, "1:01:01"},
	}
	for _, tt := range tests {
		got := FormatDuration(tt.seconds)
		if got != tt.expected {
			t.Errorf("FormatDuration(%d) = %q, want %q", tt.seconds, got, tt.expected)
		}
	}
}

func TestFormatDuration_Negative(t *testing.T) {
	got := FormatDuration(-10)
	if got != "00:00" {
		t.Errorf("expected '00:00' for negative, got %q", got)
	}
}

// captureStdout redirects os.Stdout and returns a function that restores it
// along with a channel that receives the captured output.
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestJSONOutput(t *testing.T) {
	data := []testVideo{
		{Title: "测试视频", Views: 1000, URL: "https://example.com/1"},
	}

	output := captureStdout(func() {
		EmitSuccess(data, FormatJSON)
	})

	if !strings.Contains(output, `"ok": true`) {
		t.Error("JSON output should contain ok: true")
	}
	if !strings.Contains(output, `"data"`) {
		t.Error("JSON output should contain data field")
	}
	if !strings.Contains(output, `"测试视频"`) {
		t.Error("JSON output should contain the video title")
	}
}

func TestTableOutput(t *testing.T) {
	data := []testVideo{
		{Title: "视频A", Views: 100, URL: "https://a.com"},
		{Title: "视频B", Views: 200, URL: "https://b.com"},
	}

	output := captureStdout(func() {
		RenderTable(data)
	})

	if !strings.Contains(output, "Title") {
		t.Error("table output should contain header 'Title'")
	}
	if !strings.Contains(output, "Views") {
		t.Error("table output should contain header 'Views'")
	}
	if !strings.Contains(output, "视频A") {
		t.Error("table output should contain '视频A'")
	}
	if !strings.Contains(output, "视频B") {
		t.Error("table output should contain '视频B'")
	}
}

func TestCSVOutput(t *testing.T) {
	data := []testVideo{
		{Title: "视频X", Views: 500, URL: "https://x.com"},
		{Title: "视频Y", Views: 800, URL: "https://y.com"},
	}

	output := captureStdout(func() {
		RenderCSV(data)
	})

	if !strings.Contains(output, "Title,Views,URL") {
		t.Error("CSV output should contain header row")
	}
	if !strings.Contains(output, "视频X") {
		t.Error("CSV output should contain '视频X'")
	}
	if !strings.Contains(output, "视频Y") {
		t.Error("CSV output should contain '视频Y'")
	}
	if !strings.Contains(output, "500") {
		t.Error("CSV output should contain view count")
	}
}

func TestPrettyOutput(t *testing.T) {
	data := []testVideo{
		{Title: "精选视频", Views: 9999, URL: "https://premium.com"},
	}

	output := captureStdout(func() {
		EmitSuccess(data, FormatYAML)
	})

	if !strings.Contains(output, "ok: true") {
		t.Error("YAML output should contain ok: true")
	}
	if !strings.Contains(output, "data:") {
		t.Error("YAML output should contain data field")
	}
	if !strings.Contains(output, "精选视频") {
		t.Error("YAML output should contain the video title")
	}
}

func TestEmptyOutput(t *testing.T) {
	// Test RenderTable with empty slice
	emptySlice := []testVideo{}
	output := captureStdout(func() {
		RenderTable(emptySlice)
	})
	if strings.TrimSpace(output) != "" {
		t.Error("empty slice should produce no output for table")
	}

	// Test RenderCSV with empty slice
	output = captureStdout(func() {
		RenderCSV(emptySlice)
	})
	if strings.TrimSpace(output) != "" {
		t.Error("empty slice should produce no output for CSV")
	}

	// Test EmitSuccess with nil data
	output = captureStdout(func() {
		EmitSuccess(nil, FormatJSON)
	})
	if !strings.Contains(output, `"ok": true`) {
		t.Error("EmitSuccess with nil should still emit ok: true")
	}

	// Test EmitError
	output = captureStdout(func() {
		EmitError("TEST_ERR", "test error message", FormatJSON)
	})
	if !strings.Contains(output, `"ok": false`) {
		t.Error("EmitError should emit ok: false")
	}
	if !strings.Contains(output, `TEST_ERR`) {
		t.Error("EmitError should contain the error code")
	}
}