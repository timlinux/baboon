// SPDX-FileCopyrightText: 2026 Tim Sutton / Kartoza
// SPDX-License-Identifier: MIT

package blockfont

import (
	"github.com/charmbracelet/lipgloss"
)

// CharHighlight represents different highlight states for characters
type CharHighlight int

const (
	// HighlightNone indicates no highlight
	HighlightNone CharHighlight = iota
	// HighlightCorrect indicates correctly typed/matched
	HighlightCorrect
	// HighlightIncorrect indicates incorrectly typed/mismatched
	HighlightIncorrect
	// HighlightCursor indicates cursor position
	HighlightCursor
	// HighlightDelete indicates text to be deleted
	HighlightDelete
	// HighlightChange indicates text to be changed
	HighlightChange
	// HighlightTarget indicates target/goal text
	HighlightTarget
	// HighlightPending indicates untyped/pending text
	HighlightPending
)

// CursorStyle represents how the cursor should be displayed
type CursorStyle int

const (
	// CursorBlock uses a block cursor (inverts character)
	CursorBlock CursorStyle = iota
	// CursorLine uses a vertical line cursor
	CursorLine
	// CursorUnderline uses an underline cursor
	CursorUnderline
)

// Theme holds colors and styles for block font rendering
type Theme struct {
	// Primary color for main text
	Primary lipgloss.Color
	// Secondary color for dimmed/inactive text
	Secondary lipgloss.Color
	// Accent color for highlights
	Accent lipgloss.Color
	// Correct color for correctly matched text
	Correct lipgloss.Color
	// Incorrect color for incorrectly matched text
	Incorrect lipgloss.Color
	// Cursor color for cursor highlights
	Cursor lipgloss.Color
	// Delete color for text to be deleted
	Delete lipgloss.Color
	// Change color for text to be changed
	Change lipgloss.Color
	// Target color for target/goal text
	Target lipgloss.Color
	// Pending color for untyped text
	Pending lipgloss.Color
}

// DefaultTheme provides a default color theme
var DefaultTheme = Theme{
	Primary:   lipgloss.Color("#FFFFFF"),
	Secondary: lipgloss.Color("#888888"),
	Accent:    lipgloss.Color("#FFB347"),
	Correct:   lipgloss.Color("#32CD32"),
	Incorrect: lipgloss.Color("#FF6B6B"),
	Cursor:    lipgloss.Color("#FFFFFF"),
	Delete:    lipgloss.Color("#FF4444"),
	Change:    lipgloss.Color("#FF8C00"),
	Target:    lipgloss.Color("#00FF00"),
	Pending:   lipgloss.Color("#666666"),
}

// KartozaTheme provides a Kartoza-branded color theme
var KartozaTheme = Theme{
	Primary:   lipgloss.Color("#FFFFFF"),
	Secondary: lipgloss.Color("#888888"),
	Accent:    lipgloss.Color("#FF6B35"),
	Correct:   lipgloss.Color("#32CD32"),
	Incorrect: lipgloss.Color("#FF4444"),
	Cursor:    lipgloss.Color("#FFB347"),
	Delete:    lipgloss.Color("#B22222"),
	Change:    lipgloss.Color("#FF6B35"),
	Target:    lipgloss.Color("#00FF00"),
	Pending:   lipgloss.Color("#555555"),
}

// GradientColors provides a gradient from red (slow) to green (fast)
var GradientColors = []string{
	"#8B0000", "#B22222", "#CD5C5C", "#F08080", // Reds
	"#FF6B35", "#FF8C00", "#FFA500", "#FFB347", // Oranges
	"#FFD700", "#ADFF2F", "#32CD32", "#00FF00", // Yellows to greens
}

// StyleConfig holds styling configuration for the widget
type StyleConfig struct {
	Theme       Theme
	CursorStyle CursorStyle
}

// DefaultStyleConfig returns the default style configuration
func DefaultStyleConfig() StyleConfig {
	return StyleConfig{
		Theme:       DefaultTheme,
		CursorStyle: CursorBlock,
	}
}

// NewStyle creates a lipgloss style for the given highlight
func (t Theme) NewStyle(highlight CharHighlight) lipgloss.Style {
	style := lipgloss.NewStyle()

	switch highlight {
	case HighlightCorrect:
		style = style.Foreground(t.Correct)
	case HighlightIncorrect:
		style = style.Foreground(t.Incorrect)
	case HighlightCursor:
		style = style.Reverse(true)
	case HighlightDelete:
		style = style.Foreground(t.Delete).Bold(true)
	case HighlightChange:
		style = style.Foreground(t.Change).Bold(true)
	case HighlightTarget:
		style = style.Foreground(t.Target)
	case HighlightPending:
		style = style.Foreground(t.Pending)
	default:
		style = style.Foreground(t.Primary)
	}

	return style
}

// GetGradientColor returns a color from the gradient based on a value (0.0 to 1.0)
func GetGradientColor(value float64, colors []string) string {
	if len(colors) == 0 {
		return ""
	}
	if value <= 0 {
		return colors[0]
	}
	if value >= 1 {
		return colors[len(colors)-1]
	}

	index := int(value * float64(len(colors)-1))
	if index >= len(colors) {
		index = len(colors) - 1
	}
	return colors[index]
}

// GetWPMColor returns a color based on words per minute value
func GetWPMColor(wpm int) string {
	// Map WPM to gradient (0 = slow/red, 11 = fast/green)
	// Typical range: 100-600 WPM
	index := (wpm - 100) / 50
	if index < 0 {
		index = 0
	}
	if index >= len(GradientColors) {
		index = len(GradientColors) - 1
	}
	return GradientColors[index]
}

// GetProgressColor returns a color based on progress value (0.0 to 1.0)
func GetProgressColor(progress float64) string {
	return GetGradientColor(progress, GradientColors)
}
