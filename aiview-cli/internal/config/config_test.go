package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Platform != "bilibili" {
		t.Errorf("expected 'bilibili', got '%s'", cfg.Platform)
	}
	if cfg.CacheTTL != 300 {
		t.Errorf("expected 300, got %d", cfg.CacheTTL)
	}
	if cfg.Platforms.Bilibili.Timeout != 30 {
		t.Errorf("expected 30, got %d", cfg.Platforms.Bilibili.Timeout)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Switch to temp dir to avoid picking up real config
	origDir, _ := os.Getwd()
	dir := t.TempDir()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Platform != "bilibili" {
		t.Errorf("expected 'bilibili', got '%s'", cfg.Platform)
	}
}