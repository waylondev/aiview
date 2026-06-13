package output

import (
	"testing"
)

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