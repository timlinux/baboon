// SPDX-FileCopyrightText: 2026 Tim Sutton / Kartoza
// SPDX-License-Identifier: MIT

package blockfont

import (
	"strings"
)

// insertCursor is a narrow vertical line cursor for insert mode
var insertCursor = []string{
	"|",
	"|",
	"|",
	"|",
	"|",
	"|",
}

// RenderWithCursor renders text with cursor and character-level highlighting support.
// Returns multiple lines (6 per row of characters + 1 underline row) that can be displayed.
// This is the full-featured rendering function used by typing/editing applications.
func RenderWithCursor(text string, cursorIdx int, highlights []CharHighlight, isInsertMode bool, maxWidth int, theme Theme) []string {
	runes := []rune(text)

	// Render each character in large font with appropriate styling
	// We need to build 6 lines (LetterHeight) for each row of text + 1 underline row
	letterLines := make([][]string, LetterHeight)
	underlineRow := make([]string, 0, len(runes)*2+1)
	for i := range letterLines {
		letterLines[i] = make([]string, 0, len(runes)*2+1)
	}

	for i, r := range runes {
		// Insert mode cursor: add narrow | line BEFORE the character at cursor position
		if isInsertMode && i == cursorIdx {
			for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
				letterLines[lineIdx] = append(letterLines[lineIdx], ANSIWhite+insertCursor[lineIdx]+ANSIReset)
			}
			underlineRow = append(underlineRow, " ")
		}

		// Get the block letter representation
		letter := BlockLetters[r]
		if letter == nil {
			letter = BlockLetters[' ']
		}

		// Determine styling for this character
		isCursorPos := i == cursorIdx && !isInsertMode

		// Calculate letter width for underline
		letterWidth := GetLetterWidth(r)

		// Add each line of the letter with styling
		for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
			var letterLine string
			if lineIdx < len(letter) {
				letterLine = letter[lineIdx]
			} else {
				letterLine = ""
			}

			if isCursorPos {
				// Normal mode cursor: highlight with cyan color
				letterLines[lineIdx] = append(letterLines[lineIdx], ANSICyan+letterLine+ANSIReset)
			} else if highlights != nil && i < len(highlights) {
				// Apply highlight coloring
				color := getHighlightColor(highlights[i], theme)
				if color != "" {
					letterLines[lineIdx] = append(letterLines[lineIdx], color+letterLine+ANSIReset)
				} else {
					letterLines[lineIdx] = append(letterLines[lineIdx], letterLine)
				}
			} else {
				letterLines[lineIdx] = append(letterLines[lineIdx], letterLine)
			}
		}

		// Add underline row: underlines under cursor, spaces elsewhere
		if isCursorPos {
			underlineRow = append(underlineRow, ANSICyan+strings.Repeat("▔", letterWidth)+ANSIReset)
		} else {
			underlineRow = append(underlineRow, strings.Repeat(" ", letterWidth))
		}
	}

	// Handle cursor at end of text
	if cursorIdx >= len(runes) {
		if isInsertMode {
			// Insert cursor at end: just add the | line
			for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
				letterLines[lineIdx] = append(letterLines[lineIdx], ANSIWhite+insertCursor[lineIdx]+ANSIReset)
			}
			underlineRow = append(underlineRow, " ")
		} else {
			// Normal mode cursor at end: show a block cursor placeholder with underline
			cursorBlock := BlockLetters[' ']
			if cursorBlock == nil {
				cursorBlock = []string{"   ", "   ", "   ", "   ", "   ", "   "}
			}
			cursorWidth := 3
			for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
				var letterLine string
				if lineIdx < len(cursorBlock) {
					letterLine = cursorBlock[lineIdx]
				} else {
					letterLine = "   "
				}
				letterLines[lineIdx] = append(letterLines[lineIdx], ANSICyan+letterLine+ANSIReset)
			}
			underlineRow = append(underlineRow, ANSICyan+strings.Repeat("▔", cursorWidth)+ANSIReset)
		}
	}

	// If text fits on one row or no wrapping needed, join all letters with spacing
	totalWidth := CalculateTotalWidth(runes, cursorIdx, isInsertMode)
	if maxWidth <= 0 || totalWidth <= maxWidth {
		result := make([]string, LetterHeight+1) // +1 for underline row
		spacing := strings.Repeat(" ", LetterSpacing)
		for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
			result[lineIdx] = strings.Join(letterLines[lineIdx], spacing)
		}
		// Add underline row
		result[LetterHeight] = strings.Join(underlineRow, spacing)
		return result
	}

	// Text needs to wrap on word boundaries - left justified
	return wrapWithCursor(text, cursorIdx, highlights, isInsertMode, maxWidth, theme)
}

// CalculateTotalWidth calculates the total display width including cursor
func CalculateTotalWidth(runes []rune, cursorIdx int, isInsertMode bool) int {
	total := 0
	for i, r := range runes {
		if isInsertMode && i == cursorIdx {
			total += 1 + LetterSpacing // Insert cursor width
		}
		total += GetLetterWidth(r)
		if i < len(runes)-1 {
			total += LetterSpacing
		}
	}
	if cursorIdx >= len(runes) {
		if isInsertMode {
			total += 1 + LetterSpacing // Insert cursor at end
		} else {
			total += GetLetterWidth(' ') + LetterSpacing // Block cursor
		}
	}
	return total
}

// wrapWithCursor wraps text at word boundaries (spaces) and returns left-justified output
// with proper cursor and highlighting tracking across lines
func wrapWithCursor(text string, cursorIdx int, highlights []CharHighlight, isInsertMode bool, maxWidth int, theme Theme) []string {
	// Split text into words
	words := strings.Split(text, " ")
	if len(words) == 0 {
		return []string{}
	}

	var result []string
	var currentLineWords []string
	currentWidth := 0
	charOffset := 0      // Track character position for cursor
	lineStartOffset := 0 // Track where current line starts in original text

	for _, word := range words {
		wordWidth := 0
		for _, r := range word {
			wordWidth += GetLetterWidth(r)
		}
		// Add spacing between letters within word
		if len([]rune(word)) > 1 {
			wordWidth += (len([]rune(word)) - 1) * LetterSpacing
		}

		// Check if we need insert cursor within or before this word
		wordStart := charOffset
		wordEnd := charOffset + len([]rune(word))
		if isInsertMode && cursorIdx >= wordStart && cursorIdx <= wordEnd {
			wordWidth += 1 + LetterSpacing // Account for insert cursor
		}

		// Add space width if not first word on line
		spaceWidth := 0
		if len(currentLineWords) > 0 {
			spaceWidth = GetLetterWidth(' ') + LetterSpacing
		}

		// Check if word fits on current line
		if currentWidth > 0 && currentWidth+spaceWidth+wordWidth > maxWidth {
			// Render current line and start new one
			lineText := strings.Join(currentLineWords, " ")
			lineResult := renderWrappedLine(lineText, lineStartOffset, cursorIdx, highlights, isInsertMode, theme)
			result = append(result, lineResult...)

			currentLineWords = []string{word}
			currentWidth = wordWidth
			lineStartOffset = charOffset // New line starts at current word position
		} else {
			currentLineWords = append(currentLineWords, word)
			currentWidth += spaceWidth + wordWidth
		}

		charOffset = wordEnd + 1 // +1 for the space after word
	}

	// Render remaining words
	if len(currentLineWords) > 0 {
		lineText := strings.Join(currentLineWords, " ")
		lineResult := renderWrappedLine(lineText, lineStartOffset, cursorIdx, highlights, isInsertMode, theme)
		result = append(result, lineResult...)
	}

	return result
}

// renderWrappedLine renders a single line of words with proper cursor and highlighting
func renderWrappedLine(lineText string, startOffset int, cursorIdx int, highlights []CharHighlight, isInsertMode bool, theme Theme) []string {
	runes := []rune(lineText)

	// Ensure startOffset is never negative
	if startOffset < 0 {
		startOffset = 0
	}

	// Get highlights for this portion of text
	var lineHighlights []CharHighlight
	if highlights != nil && startOffset < len(highlights) {
		endIdx := startOffset + len(runes)
		if endIdx > len(highlights) {
			endIdx = len(highlights)
		}
		lineHighlights = highlights[startOffset:endIdx]
	}

	// Build letter lines for this row + underline row
	letterLines := make([][]string, LetterHeight)
	underlineRow := make([]string, 0, len(runes)*2+1)
	for i := range letterLines {
		letterLines[i] = make([]string, 0, len(runes)*2+1)
	}

	for i, r := range runes {
		globalIdx := startOffset + i

		// Insert mode cursor
		if isInsertMode && globalIdx == cursorIdx {
			for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
				letterLines[lineIdx] = append(letterLines[lineIdx], ANSIWhite+insertCursor[lineIdx]+ANSIReset)
			}
			underlineRow = append(underlineRow, " ")
		}

		letter := BlockLetters[r]
		if letter == nil {
			letter = BlockLetters[' ']
		}

		isCursorPos := globalIdx == cursorIdx && !isInsertMode
		letterWidth := GetLetterWidth(r)

		for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
			var letterLine string
			if lineIdx < len(letter) {
				letterLine = letter[lineIdx]
			} else {
				letterLine = ""
			}

			if isCursorPos {
				// Normal mode cursor: highlight with cyan color
				letterLines[lineIdx] = append(letterLines[lineIdx], ANSICyan+letterLine+ANSIReset)
			} else if lineHighlights != nil && i < len(lineHighlights) {
				color := getHighlightColor(lineHighlights[i], theme)
				if color != "" {
					letterLines[lineIdx] = append(letterLines[lineIdx], color+letterLine+ANSIReset)
				} else {
					letterLines[lineIdx] = append(letterLines[lineIdx], letterLine)
				}
			} else {
				letterLines[lineIdx] = append(letterLines[lineIdx], letterLine)
			}
		}

		// Add underline row: underlines under cursor, spaces elsewhere
		if isCursorPos {
			underlineRow = append(underlineRow, ANSICyan+strings.Repeat("▔", letterWidth)+ANSIReset)
		} else {
			underlineRow = append(underlineRow, strings.Repeat(" ", letterWidth))
		}
	}

	// Handle cursor at end of this line (only if cursor is at end of full text)
	endGlobalIdx := startOffset + len(runes)
	if cursorIdx == endGlobalIdx {
		if isInsertMode {
			for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
				letterLines[lineIdx] = append(letterLines[lineIdx], ANSIWhite+insertCursor[lineIdx]+ANSIReset)
			}
			underlineRow = append(underlineRow, " ")
		} else {
			cursorBlock := BlockLetters[' ']
			if cursorBlock == nil {
				cursorBlock = []string{"   ", "   ", "   ", "   ", "   ", "   "}
			}
			cursorWidth := 3
			for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
				var letterLine string
				if lineIdx < len(cursorBlock) {
					letterLine = cursorBlock[lineIdx]
				} else {
					letterLine = "   "
				}
				letterLines[lineIdx] = append(letterLines[lineIdx], ANSICyan+letterLine+ANSIReset)
			}
			underlineRow = append(underlineRow, ANSICyan+strings.Repeat("▔", cursorWidth)+ANSIReset)
		}
	}

	// Join letters with spacing - LEFT JUSTIFIED (no centering) + underline row
	result := make([]string, LetterHeight+1)
	spacing := strings.Repeat(" ", LetterSpacing)
	for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
		result[lineIdx] = strings.Join(letterLines[lineIdx], spacing)
	}
	result[LetterHeight] = strings.Join(underlineRow, spacing)

	return result
}

// RenderPlainText renders plain text in large font without any highlighting.
// Used for reference/target text display.
func RenderPlainText(text string, color string) []string {
	runes := []rune(text)
	letterLines := make([][]string, LetterHeight)
	for i := range letterLines {
		letterLines[i] = make([]string, 0, len(runes))
	}

	for _, r := range runes {
		letter := BlockLetters[r]
		if letter == nil {
			letter = BlockLetters[' ']
		}

		for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
			var letterLine string
			if lineIdx < len(letter) {
				letterLine = letter[lineIdx]
			} else {
				letterLine = ""
			}

			if color != "" {
				letterLines[lineIdx] = append(letterLines[lineIdx], color+letterLine+ANSIReset)
			} else {
				letterLines[lineIdx] = append(letterLines[lineIdx], letterLine)
			}
		}
	}

	result := make([]string, LetterHeight)
	for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
		spacing := strings.Repeat(" ", LetterSpacing)
		result[lineIdx] = strings.Join(letterLines[lineIdx], spacing)
	}

	return result
}

// getHighlightColor returns the ANSI color code for a highlight type
func getHighlightColor(highlight CharHighlight, theme Theme) string {
	switch highlight {
	case HighlightCorrect:
		return ANSIGreen
	case HighlightIncorrect:
		return ANSIRed
	case HighlightDelete:
		return ANSIRed
	case HighlightChange:
		return ANSIOrange
	case HighlightTarget:
		return ANSIGreen
	case HighlightPending:
		return ANSIDim
	case HighlightCursor:
		return ANSICyan
	default:
		return ""
	}
}

// GetDisplayWidth returns the display width of text when rendered in large font
func GetDisplayWidth(text string) int {
	return GetTotalWidth(text)
}
