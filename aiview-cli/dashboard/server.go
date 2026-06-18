// Package dashboard provides a web-based monitoring dashboard for aiview.
package dashboard

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackwener/aiview/internal/analyzer"
	"github.com/jackwener/aiview/internal/scheduler"
	"github.com/jackwener/aiview/internal/storage"
)

//go:embed static/*
var staticFiles embed.FS

// Server represents the dashboard HTTP server.
type Server struct {
	port      int
	storage   storage.Storage
	analyzer  *analyzer.Analyzer
	scheduler *scheduler.Scheduler
	mux       *http.ServeMux
}

// NewServer creates a new dashboard server.
func NewServer(port int, store storage.Storage, sched *scheduler.Scheduler) (*Server, error) {
	s := &Server{
		port:      port,
		storage:   store,
		analyzer:  analyzer.New(store),
		scheduler: sched,
		mux:       http.NewServeMux(),
	}
	if err := s.setupRoutes(); err != nil {
		return nil, err
	}
	return s, nil
}

// Start starts the HTTP server with graceful shutdown.
func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.mux,
	}

	// Channel to listen for errors from the server
	errCh := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-quit:
		// Create a deadline for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	}
}

func (s *Server) setupRoutes() error {
	// API endpoints
	s.mux.HandleFunc("/api/trend", s.handleTrend)
	s.mux.HandleFunc("/api/platforms", s.handlePlatforms)
	s.mux.HandleFunc("/api/schedule", s.handleSchedule)
	s.mux.HandleFunc("/api/history", s.handleHistory)

	// Static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to get static files: %w", err)
	}
	fileServer := http.FileServer(http.FS(staticFS))
	s.mux.Handle("/", fileServer)
	return nil
}
