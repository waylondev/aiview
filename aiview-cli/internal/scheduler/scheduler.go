// Package scheduler provides a simple cron-like task scheduler.
package scheduler

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"

	aiverr "github.com/jackwener/aiview/internal/errors"
)

// Job represents a scheduled task.
type Job struct {
	ID       string
	Interval time.Duration
	Command  string
	LastRun  time.Time
	NextRun  time.Time
	Running  bool
}

// Scheduler manages periodic task execution.
type Scheduler struct {
	mu      sync.Mutex
	jobs    map[string]*Job
	running bool
	stopCh  chan struct{}
}

// New creates a new Scheduler.
func New() *Scheduler {
	return &Scheduler{
		jobs:   make(map[string]*Job),
		stopCh: make(chan struct{}),
	}
}

// AddJob adds a new scheduled job.
func (s *Scheduler) AddJob(id string, interval time.Duration, command string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[id]; exists {
		return aiverr.InvalidInput("scheduler", fmt.Sprintf("job %q already exists", id))
	}

	now := time.Now()
	s.jobs[id] = &Job{
		ID:       id,
		Interval: interval,
		Command:  command,
		NextRun:  now.Add(interval),
	}

	return nil
}

// RemoveJob removes a scheduled job.
func (s *Scheduler) RemoveJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[id]; !exists {
		return aiverr.NotFound("scheduler", fmt.Sprintf("job %q not found", id))
	}

	delete(s.jobs, id)
	return nil
}

// ListJobs returns all scheduled jobs.
func (s *Scheduler) ListJobs() []*Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	var jobs []*Job
	for _, j := range s.jobs {
		jobs = append(jobs, j)
	}
	return jobs
}

// Start begins the scheduler loop.
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go s.loop()
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}
	s.running = false
	close(s.stopCh)
}

func (s *Scheduler) loop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case now := <-ticker.C:
			s.mu.Lock()
			for _, job := range s.jobs {
				if now.After(job.NextRun) && !job.Running {
					job.Running = true
					job.LastRun = now
					job.NextRun = now.Add(job.Interval)
					go s.executeJob(job)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Scheduler) executeJob(job *Job) {
	defer func() {
		s.mu.Lock()
		job.Running = false
		s.mu.Unlock()
	}()

	// Execute the command using the aiview binary itself
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", job.Command)
	} else {
		cmd = exec.Command("sh", "-c", job.Command)
	}
	if err := cmd.Run(); err != nil {
		log.Printf("Command execution failed: %v", err)
	}
}

// ParseInterval parses a human-readable interval string.
func ParseInterval(s string) (time.Duration, error) {
	// Support formats: "30s", "5m", "1h", "24h"
	return time.ParseDuration(s)
}
