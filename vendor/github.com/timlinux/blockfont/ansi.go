// SPDX-FileCopyrightText: 2026 Tim Sutton / Kartoza
// SPDX-License-Identifier: MIT

package blockfont

import (
	"strings"
	"unicode/utf8"
)

// ANSI escape code constants for terminal styling
const (
	ANSIReset     = "\033[0m"
	ANSIRed       = "\033[1;31m"
	ANSIGreen     = "\033[1;32m"
	ANSIOrange    = "\033[1;33m"
	ANSIYellow    = "\033[1;33m" // Alias for orange
	ANSIBlue      = "\033[1;34m"
	ANSIMagenta   = "\033[1;35m"
	ANSICyan      = "\033[1;36m"
	ANSIWhite     = "\033[1;37m"
	ANSIInverse   = "\033[7m"
	ANSIDim       = "\033[2m"
	ANSIBold      = "\033[1m"
	ANSIItalic    = "\033[3m"
	ANSIUnderline = "\033[4m"
)

// RemoveANSI removes ANSI escape codes from a string.
// This is useful for calculating visible width of styled text.
func RemoveANSI(s string) string {
	var result strings.Builder
	inEscape := false

	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}

	return result.String()
}

// VisibleStringWidth calculates the visible width of a string, excluding ANSI escape codes.
// Returns the rune count of the visible text.
func VisibleStringWidth(s string) int {
	cleaned := RemoveANSI(s)
	return utf8.RuneCountInString(cleaned)
}

// InsertAt inserts overlay text at position x in the base string.
// Handles overlapping text composition while preserving ANSI codes.
func InsertAt(base, overlay string, x int) string {
	baseRunes := []rune(RemoveANSI(base))
	overlayRunes := []rune(overlay)

	// Extend base if needed
	if x > len(baseRunes) {
		padding := make([]rune, x-len(baseRunes))
		for i := range padding {
			padding[i] = ' '
		}
		baseRunes = append(baseRunes, padding...)
	}

	// Insert overlay
	result := make([]rune, 0, len(baseRunes)+len(overlayRunes))
	result = append(result, baseRunes[:x]...)
	result = append(result, overlayRunes...)
	if x+len(overlayRunes) < len(baseRunes) {
		result = append(result, baseRunes[x+len(overlayRunes):]...)
	}

	return string(result)
}

// WrapWithColor wraps text with ANSI color codes.
func WrapWithColor(text, color string) string {
	if color == "" {
		return text
	}
	return color + text + ANSIReset
}

// InvertLine inverts a single line using ANSI inverse video.
// This swaps foreground/background colors for the entire line.
func InvertLine(line string) string {
	return ANSIInverse + line + ANSIReset
}
