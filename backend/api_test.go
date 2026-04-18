package backend

import (
	"testing"

	"github.com/timlinux/baboon/settings"
)

func TestDefaultConfig_PracticeModeIsWords(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.PracticeMode != settings.ModeWords {
		t.Errorf("DefaultConfig().PracticeMode = %v, want ModeWords", cfg.PracticeMode)
	}
}
