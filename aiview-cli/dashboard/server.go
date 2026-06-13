// Package dashboard provides a web-based monitoring dashboard for aiview.
package dashboard

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

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
func NewServer(port int, store storage.Storage, sched *scheduler.Scheduler) *Server {
	s := &Server{
		port:      port,
		storage:   store,
		analyzer:  analyzer.New(store),
		scheduler: sched,
		mux:       http.NewServeMux(),
	}
	s.setupRoutes()
	return s
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) setupRoutes() {
	// API endpoints
	s.mux.HandleFunc("/api/trend", s.handleTrend)
	s.mux.HandleFunc("/api/platforms", s.handlePlatforms)
	s.mux.HandleFunc("/api/schedule", s.handleSchedule)
	s.mux.HandleFunc("/api/history", s.handleHistory)

	// Static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(fmt.Sprintf("failed to get static files: %v", err))
	}
	fileServer := http.FileServer(http.FS(staticFS))
	s.mux.Handle("/", fileServer)
}
