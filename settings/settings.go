package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AdvanceKey defines the key(s) used to advance to the next word
type AdvanceKey int

const (
	AdvanceKeySpace AdvanceKey = iota // Default: space only
	AdvanceKeyEnter                   // Enter only
	AdvanceKeyEither                  // Either space or enter
)

// String returns the display name for the advance key setting
func (a AdvanceKey) String() string {
	switch a {
	case AdvanceKeySpace:
		return "Space"
	case AdvanceKeyEnter:
		return "Enter"
	case AdvanceKeyEither:
		return "Either"
	default:
		return "Space"
	}
}

// KeyHint returns the help text for which key to press
func (a AdvanceKey) KeyHint() string {
	switch a {
	case AdvanceKeySpace:
		return "SPACE"
	case AdvanceKeyEnter:
		return "ENTER"
	case AdvanceKeyEither:
		return "SPACE or ENTER"
	default:
		return "SPACE"
	}
}

// Settings holds user preferences
type Settings struct {
	AdvanceKey AdvanceKey `json:"advance_key"`
}

// DefaultSettings returns the default settings
func DefaultSettings() *Settings {
	return &Settings{
		AdvanceKey: AdvanceKeySpace,
	}
}

// GetSettingsPath returns the path to the settings file
func GetSettingsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	settingsDir := filepath.Join(homeDir, ".config", "baboon")
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(settingsDir, "settings.json"), nil
}

// Load loads settings from disk, returning defaults if not found
func Load() (*Settings, error) {
	path, err := GetSettingsPath()
	if err != nil {
		return DefaultSettings(), err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultSettings(), nil
		}
		return DefaultSettings(), err
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return DefaultSettings(), err
	}

	return &s, nil
}

// Save saves settings to disk
func (s *Settings) Save() error {
	path, err := GetSettingsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
