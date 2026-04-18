package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPracticeMode_String(t *testing.T) {
	tests := []struct {
		mode PracticeMode
		want string
	}{
		{ModeWords, "Words"},
		{ModeNgrams, "N-grams"},
	}
	for _, tt := range tests {
		if got := tt.mode.String(); got != tt.want {
			t.Errorf("PracticeMode.String() = %v, want %v", got, tt.want)
		}
	}
}

func TestSettings_PracticeModeDefaultsToWords(t *testing.T) {
	s := DefaultSettings()
	if s.PracticeMode != ModeWords {
		t.Errorf("DefaultSettings().PracticeMode = %v, want ModeWords", s.PracticeMode)
	}
}

func TestSettings_PracticeModePersists(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create config dir
	configDir := filepath.Join(tmpDir, ".config", "baboon")
	os.MkdirAll(configDir, 0755)

	// Save settings with N-gram mode
	s := DefaultSettings()
	s.PracticeMode = ModeNgrams
	if err := s.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load and verify
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.PracticeMode != ModeNgrams {
		t.Errorf("Loaded PracticeMode = %v, want ModeNgrams", loaded.PracticeMode)
	}
}
