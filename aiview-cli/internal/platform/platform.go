// Package platform defines the core Platform interface and registry for multi-platform CLI applications.
package platform

import (
	"github.com/jackwener/aiview/internal/config"
	"github.com/spf13/cobra"
)

// Client is the interface that all platform API clients must implement.
type Client interface {
	// PlatformName returns the platform name.
	PlatformName() string
}

// Platform defines the interface for a platform implementation.
type Platform interface {
	// Name returns the platform identifier (e.g., "bilibili", "twitter").
	Name() string
	// Commands returns all cobra commands for this platform.
	Commands() []*cobra.Command
	// NewClient creates a new API client for this platform.
	NewClient(cfg *config.Config) (Client, error)
}