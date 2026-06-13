package scheduler

import (
	"testing"
	"time"
)

func TestScheduler_AddJob(t *testing.T) {
	s := New()

	err := s.AddJob("test1", time.Hour, "echo hello")
	if err != nil {
		t.Fatalf("AddJob failed: %v", err)
	}

	jobs := s.ListJobs()
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].ID != "test1" {
		t.Errorf("expected job ID 'test1', got '%s'", jobs[0].ID)
	}
}

func TestScheduler_AddDuplicateJob(t *testing.T) {
	s := New()

	s.AddJob("test1", time.Hour, "echo hello")
	err := s.AddJob("test1", time.Hour, "echo world")
	if err == nil {
		t.Error("expected error for duplicate job")
	}
}

func TestScheduler_RemoveJob(t *testing.T) {
	s := New()

	s.AddJob("test1", time.Hour, "echo hello")
	err := s.RemoveJob("test1")
	if err != nil {
		t.Fatalf("RemoveJob failed: %v", err)
	}

	jobs := s.ListJobs()
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestScheduler_RemoveNonExistent(t *testing.T) {
	s := New()

	err := s.RemoveJob("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent job")
	}
}

func TestParseInterval(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{"30s", 30 * time.Second, false},
		{"5m", 5 * time.Minute, false},
		{"1h", time.Hour, false},
		{"24h", 24 * time.Hour, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		d, err := ParseInterval(tt.input)
		if tt.wantErr && err == nil {
			t.Errorf("ParseInterval(%q) expected error", tt.input)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("ParseInterval(%q) unexpected error: %v", tt.input, err)
		}
		if d != tt.expected {
			t.Errorf("ParseInterval(%q) = %v, want %v", tt.input, d, tt.expected)
		}
	}
}

func TestScheduler_StartStop(t *testing.T) {
	s := New()
	s.Start()
	// Should not panic on double start
	s.Start()
	s.Stop()
}
