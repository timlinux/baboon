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

func TestEngine_TracksTrigramSeekTime(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	cfg := DefaultConfig()
	engine, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("NewEngine() error = %v", err)
	}

	// Get the current word and type its first 3 letters
	state := engine.GetGameState()
	currentWord := state.CurrentWord

	// Skip if word is too short
	if len(currentWord) < 3 {
		t.Skip("First word too short for trigram test")
	}

	// Type first 3 letters of the current word with timing
	char1 := string(currentWord[0])
	char2 := string(currentWord[1])
	char3 := string(currentWord[2])
	expectedTrigram := char1 + char2 + char3

	engine.ProcessKeystrokeWithTiming(char1, 0)   // First letter, no seek time
	engine.ProcessKeystrokeWithTiming(char2, 100) // Second letter
	engine.ProcessKeystrokeWithTiming(char3, 100) // Third letter - should record trigram

	session := engine.GetSessionStats()

	// Check trigram was recorded
	if _, exists := session.TrigramSeekTime[expectedTrigram]; !exists {
		t.Errorf("Expected trigram '%s' to be tracked, got: %v", expectedTrigram, session.TrigramSeekTime)
	}
}
