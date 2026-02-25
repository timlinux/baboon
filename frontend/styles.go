// Package frontend provides the terminal user interface for the typing practice application.
// This package handles all rendering, user input, and visual presentation.
// It communicates with the backend exclusively through the GameAPI interface.
package frontend

import "github.com/charmbracelet/lipgloss"

// Colour constants for consistent styling
const (
	ColourCorrect   = "10"  // Bright green
	ColourIncorrect = "9"   // Bright red
	ColourUntyped   = "8"   // Gray
	ColourTitle     = "14"  // Cyan
	ColourLabel     = "7"   // Light gray
	ColourValue     = "15"  // White
	ColourSession   = "6"   // Cyan
	ColourNewBest   = "226" // Yellow
	ColourHelp      = "8"   // Gray
	ColourEmptyBar  = "236" // Dark gray
	ColourPrevNext  = "240" // Dim gray
)

// Gradient colours from red (slow) through yellow to green (fast)
var GradientColours = []string{
	"196", "202", "208", "214", "220", "226",
	"190", "154", "118", "82", "46", "47",
}

// Styles holds all the lipgloss styles used in the application
type Styles struct {
	// Typing screen styles
	Correct   lipgloss.Style
	Incorrect lipgloss.Style
	Untyped   lipgloss.Style
	Progress  lipgloss.Style
	PrevNext  lipgloss.Style
	Help      lipgloss.Style

	// Results screen styles
	Title        lipgloss.Style
	NewBest      lipgloss.Style
	SessionLabel lipgloss.Style
	SessionValue lipgloss.Style
	Label        lipgloss.Style
	Value        lipgloss.Style
	LetterLabel  lipgloss.Style
	LetterHeader lipgloss.Style
	ErrorStyle   lipgloss.Style
	CountStyle   lipgloss.Style
}

// NewStyles creates a new Styles instance with all styles initialised
func NewStyles() Styles {
	return Styles{
		// Typing screen
		Correct:   lipgloss.NewStyle().Foreground(lipgloss.Color(ColourCorrect)),
		Incorrect: lipgloss.NewStyle().Foreground(lipgloss.Color(ColourIncorrect)),
		Untyped:   lipgloss.NewStyle().Foreground(lipgloss.Color(ColourUntyped)),
		Progress:  lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true),
		PrevNext:  lipgloss.NewStyle().Foreground(lipgloss.Color(ColourPrevNext)),
		Help:      lipgloss.NewStyle().Foreground(lipgloss.Color(ColourHelp)),

		// Results screen
		Title:        lipgloss.NewStyle().Foreground(lipgloss.Color(ColourTitle)).Bold(true),
		NewBest:      lipgloss.NewStyle().Foreground(lipgloss.Color(ColourNewBest)).Bold(true),
		SessionLabel: lipgloss.NewStyle().Foreground(lipgloss.Color(ColourSession)).Width(18).Align(lipgloss.Right),
		SessionValue: lipgloss.NewStyle().Foreground(lipgloss.Color(ColourValue)).Width(8).Align(lipgloss.Right),
		Label:        lipgloss.NewStyle().Foreground(lipgloss.Color(ColourLabel)).Width(18).Align(lipgloss.Right),
		Value:        lipgloss.NewStyle().Foreground(lipgloss.Color(ColourValue)).Width(8).Align(lipgloss.Right),
		LetterLabel:  lipgloss.NewStyle().Foreground(lipgloss.Color(ColourLabel)).Width(18).Align(lipgloss.Right),
		LetterHeader: lipgloss.NewStyle().Foreground(lipgloss.Color(ColourValue)).Bold(true),
		ErrorStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		CountStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color(ColourHelp)),
	}
}

// GetAccuracyColour returns a colour on red-yellow-green gradient based on accuracy (0-100)
// This uses absolute thresholds - for relative comparisons, use GetRelativeColour
func GetAccuracyColour(accuracy float64) string {
	if accuracy >= 95 {
		return "46" // Bright green
	} else if accuracy >= 90 {
		return "82"
	} else if accuracy >= 85 {
		return "118"
	} else if accuracy >= 80 {
		return "154"
	} else if accuracy >= 75 {
		return "190"
	} else if accuracy >= 70 {
		return "226" // Yellow
	} else if accuracy >= 65 {
		return "220"
	} else if accuracy >= 60 {
		return "214"
	} else if accuracy >= 50 {
		return "208"
	} else if accuracy >= 40 {
		return "202"
	}
	return "196" // Red
}

// GetRelativeColour returns a colour based on where a value falls within a min-max range.
// This ensures visual differentiation even when all values are clustered.
// A value at min gets red, at max gets green, with gradient in between.
func GetRelativeColour(value, minVal, maxVal float64) string {
	if maxVal <= minVal {
		return "46" // All same - show green
	}

	// Normalize to 0-1 range
	normalized := (value - minVal) / (maxVal - minVal)

	// Clamp to 0-1
	if normalized < 0 {
		normalized = 0
	}
	if normalized > 1 {
		normalized = 1
	}

	// Map to color gradient (index 0-11 in GradientColours)
	// Index 0 = red (worst), Index 11 = green (best)
	colourIdx := int(normalized * 11)
	if colourIdx > 11 {
		colourIdx = 11
	}

	return GradientColours[colourIdx]
}

// GetFrequencyColour returns a colour based on frequency (0-1)
// Higher frequency = more red (problematic letters appear more), lower = green
func GetFrequencyColour(frequency float64) string {
	// Invert: high frequency should be red (more practice needed), low should be green
	return GetAccuracyColour((1.0 - frequency) * 100)
}

// GetSeekTimeColour returns a colour based on seek time (inverted: fast=green, slow=red)
func GetSeekTimeColour(seekTimeMs, maxSeekTimeMs float64) string {
	if maxSeekTimeMs == 0 {
		return "46" // Green if no data
	}
	normalized := seekTimeMs / maxSeekTimeMs
	accuracy := (1.0 - normalized) * 100
	return GetAccuracyColour(accuracy)
}
