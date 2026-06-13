package bilibili

import "testing"

func TestExtractBVID_ValidURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"BV1xx411c7m9", "BV1xx411c7m9"},
		{"https://www.bilibili.com/video/BV1xx411c7m9", "BV1xx411c7m9"},
		{"https://b23.tv/BV1xx411c7m9", "BV1xx411c7m9"},
	}
	for _, tt := range tests {
		result, err := ExtractBVID(tt.input)
		if err != nil {
			t.Errorf("ExtractBVID(%q) unexpected error: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("ExtractBVID(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractBVID_Invalid(t *testing.T) {
	_, err := ExtractBVID("invalid_input")
	if err == nil {
		t.Error("ExtractBVID('invalid_input') expected error, got nil")
	}
}