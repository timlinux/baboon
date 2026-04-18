package backend

import (
	"os"
	"testing"

	"github.com/timlinux/baboon/settings"
)

func TestDefaultConfig_PracticeModeIsWords(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.PracticeMode != settings.ModeWords {
		t.Errorf("DefaultConfig().PracticeMode = %v, want ModeWords", cfg.PracticeMode)
	}
}

func TestEngine_NgramMode_StartsRound(t *testing.T) {
	// Set HOME to temp directory for sandboxed test environment
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	cfg := DefaultConfig()
	cfg.PracticeMode = settings.ModeNgrams

	engine, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	state := engine.GetGameState()

	// Should have n-grams (2-3 letter sequences)
	if len(state.Words) == 0 {
		t.Error("Expected non-empty word list")
	}

	for _, ngram := range state.Words {
		if len(ngram) < 2 || len(ngram) > 3 {
			t.Errorf("Expected n-gram (2-3 chars), got %q", ngram)
		}
	}
}

func TestEngine_WordMode_StartsRound(t *testing.T) {
	// Set HOME to temp directory for sandboxed test environment
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	cfg := DefaultConfig()
	cfg.PracticeMode = settings.ModeWords

	engine, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	state := engine.GetGameState()

	// Should have words (30 by default)
	if len(state.Words) != 30 {
		t.Errorf("Expected 30 words, got %d", len(state.Words))
	}
}
