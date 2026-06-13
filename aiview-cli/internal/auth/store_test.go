package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := &Store{
		Dir:  dir,
		File: filepath.Join(dir, "test_credential.json"),
	}

	cred := &Credential{
		Sessdata:    "test_sessdata",
		BiliJct:     "test_bili_jct",
		AcTimeValue: "test_ac_time",
		Buvid3:      "test_buvid3",
		Buvid4:      "test_buvid4",
		DedeUserID:  "test_dedeuserid",
	}

	if err := store.Save(cred); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil credential")
	}
	if loaded.Sessdata != "test_sessdata" {
		t.Errorf("expected 'test_sessdata', got '%s'", loaded.Sessdata)
	}
	if loaded.BiliJct != "test_bili_jct" {
		t.Errorf("expected 'test_bili_jct', got '%s'", loaded.BiliJct)
	}
	if loaded.SavedAt == 0 {
		t.Error("expected non-zero SavedAt")
	}
}

func TestStore_LoadNonExistent(t *testing.T) {
	dir := t.TempDir()
	store := &Store{
		Dir:  dir,
		File: filepath.Join(dir, "nonexistent.json"),
	}

	cred, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cred != nil {
		t.Error("expected nil credential for non-existent file")
	}
}

func TestStore_Clear(t *testing.T) {
	dir := t.TempDir()
	store := &Store{
		Dir:  dir,
		File: filepath.Join(dir, "test_credential.json"),
	}

	store.Save(&Credential{Sessdata: "test"})

	if err := store.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// File should not exist after clear
	if _, err := os.Stat(store.File); !os.IsNotExist(err) {
		t.Error("credential file should not exist after clear")
	}
}

func TestStore_ClearNonExistent(t *testing.T) {
	dir := t.TempDir()
	store := &Store{
		Dir:  dir,
		File: filepath.Join(dir, "nonexistent.json"),
	}

	if err := store.Clear(); err != nil {
		t.Errorf("Clear on non-existent file should not error: %v", err)
	}
}

func TestCredential_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    *Credential
		want bool
	}{
		{"valid", &Credential{Sessdata: "abc"}, true},
		{"empty", &Credential{}, false},
		{"only_jct", &Credential{BiliJct: "abc"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredential_HasWriteCapability(t *testing.T) {
	tests := []struct {
		name string
		c    *Credential
		want bool
	}{
		{"has", &Credential{BiliJct: "abc"}, true},
		{"no", &Credential{Sessdata: "abc"}, false},
		{"empty", &Credential{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.HasWriteCapability(); got != tt.want {
				t.Errorf("HasWriteCapability() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredential_IsStale(t *testing.T) {
	// Current time should not be stale
	c := &Credential{SavedAt: 0}
	if !c.IsStale(7) {
		t.Error("SavedAt 0 should be stale")
	}
}