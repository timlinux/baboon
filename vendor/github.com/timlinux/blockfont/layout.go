// SPDX-FileCopyrightText: 2026 Tim Sutton / Kartoza
// SPDX-License-Identifier: MIT

package blockfont

import (
	"strings"
)

// Alignment represents text alignment options
type Alignment int

const (
	// AlignLeft aligns text to the left
	AlignLeft Alignment = iota
	// AlignCenter centers text
	AlignCenter
	// AlignRight aligns text to the right
	AlignRight
)

// CenterLines centers each line within the given width.
// ANSI escape codes are excluded from width calculations.
func CenterLines(lines []string, width int) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		visibleWidth := VisibleStringWidth(line)
		if visibleWidth < width {
			padding := (width - visibleWidth) / 2
			result[i] = strings.Repeat(" ", padding) + line
		} else {
			result[i] = line
		}
	}
	return result
}

// LeftJustify left-justifies lines with optional left margin.
func LeftJustify(lines []string, margin int) []string {
	if margin <= 0 {
		return lines
	}
	result := make([]string, len(lines))
	padding := strings.Repeat(" ", margin)
	for i, line := range lines {
		result[i] = padding + line
	}
	return result
}

// RightJustify right-justifies lines within the given width.
// ANSI escape codes are excluded from width calculations.
func RightJustify(lines []string, width int) []string {
	result := make([]string, len(lines))
	for i, line := range lines {
		visibleWidth := VisibleStringWidth(line)
		if visibleWidth < width {
			padding := width - visibleWidth
			result[i] = strings.Repeat(" ", padding) + line
		} else {
			result[i] = line
		}
	}
	return result
}

// AlignLines aligns lines according to the specified alignment and width.
func AlignLines(lines []string, alignment Alignment, width int) []string {
	switch alignment {
	case AlignCenter:
		return CenterLines(lines, width)
	case AlignRight:
		return RightJustify(lines, width)
	default:
		return lines
	}
}

// WrapOnWordBoundaries wraps text at word boundaries (spaces) to fit within maxWidth.
// Returns the text split into multiple lines.
func WrapOnWordBoundaries(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	var result []string
	var currentLine strings.Builder
	currentWidth := 0

	for i, word := range words {
		wordWidth := GetTotalWidth(word)

		// Check if word fits on current line
		spaceWidth := 0
		if currentWidth > 0 {
			spaceWidth = GetLetterWidth(' ') + LetterSpacing
		}

		if currentWidth > 0 && currentWidth+spaceWidth+wordWidth > maxWidth {
			// Word doesn't fit, start new line
			result = append(result, currentLine.String())
			currentLine.Reset()
			currentWidth = 0
		}

		// Add space if not first word on line
		if currentWidth > 0 {
			currentLine.WriteRune(' ')
			currentWidth += spaceWidth
		}

		currentLine.WriteString(word)
		currentWidth += wordWidth

		// Last word
		if i == len(words)-1 {
			result = append(result, currentLine.String())
		}
	}

	return result
}

// PadToWidth pads or truncates a string to the exact visible width.
// ANSI escape codes are preserved.
func PadToWidth(line string, width int) string {
	visibleWidth := VisibleStringWidth(line)
	if visibleWidth == width {
		return line
	}
	if visibleWidth < width {
		return line + strings.Repeat(" ", width-visibleWidth)
	}
	// Truncate - need to be careful with ANSI codes
	return truncateToWidth(line, width)
}

// truncateToWidth truncates a string to the given visible width.
// Preserves ANSI escape codes.
func truncateToWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}

	var result strings.Builder
	visibleCount := 0
	inEscape := false

	for _, r := range s {
		if r == '\033' {
			inEscape = true
			result.WriteRune(r)
			continue
		}
		if inEscape {
			result.WriteRune(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}

		if visibleCount >= width {
			break
		}
		result.WriteRune(r)
		visibleCount++
	}

	return result.String()
}

// MaxLineWidth returns the maximum visible width across all lines.
func MaxLineWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		w := VisibleStringWidth(line)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

// JoinBlockLines joins multiple lines from RenderWord with the specified spacing.
func JoinBlockLines(lines [][]string, spacing int) []string {
	result := make([]string, LetterHeight)
	spacer := strings.Repeat(" ", spacing)

	for lineIdx := range LetterHeight {
		result[lineIdx] = strings.Join(lines[lineIdx], spacer)
	}

	return result
}
