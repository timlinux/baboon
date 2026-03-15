// SPDX-FileCopyrightText: 2026 Tim Sutton / Kartoza
// SPDX-License-Identifier: MIT

package blockfont

import (
	"strings"
	"unicode/utf8"
)

// Mode represents vim editing modes
type Mode int

const (
	// ModeNormal is the default mode for navigation and commands
	ModeNormal Mode = iota
	// ModeInsert is for typing text
	ModeInsert
	// ModeVisual is for selecting text
	ModeVisual
	// ModeVisualLine is for selecting whole lines
	ModeVisualLine
	// ModeVisualBlock is for selecting rectangular blocks
	ModeVisualBlock
	// ModeCommand is for entering ex commands
	ModeCommand
)

// String returns the mode name for display
func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeVisual:
		return "VISUAL"
	case ModeVisualLine:
		return "V-LINE"
	case ModeVisualBlock:
		return "V-BLOCK"
	case ModeCommand:
		return "COMMAND"
	default:
		return "UNKNOWN"
	}
}

// Buffer represents a text buffer with cursor position for vim-style editing
type Buffer struct {
	lines    []string
	cursorX  int // Column (character index in current line)
	cursorY  int // Row (line index)
	mode     Mode
	register string // Yank register
}

// NewBuffer creates a new buffer with the given text
func NewBuffer(text string) *Buffer {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	return &Buffer{
		lines:   lines,
		cursorX: 0,
		cursorY: 0,
		mode:    ModeNormal,
	}
}

// Text returns the full buffer content as a string
func (b *Buffer) Text() string {
	return strings.Join(b.lines, "\n")
}

// SetText replaces the buffer content
func (b *Buffer) SetText(text string) {
	b.lines = strings.Split(text, "\n")
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}
	b.clampCursor()
}

// Lines returns all lines
func (b *Buffer) Lines() []string {
	return b.lines
}

// CurrentLine returns the current line
func (b *Buffer) CurrentLine() string {
	if b.cursorY >= 0 && b.cursorY < len(b.lines) {
		return b.lines[b.cursorY]
	}
	return ""
}

// CursorPosition returns the cursor position (x=column, y=line)
func (b *Buffer) CursorPosition() (x, y int) {
	return b.cursorX, b.cursorY
}

// CursorIndex returns the absolute character index in the buffer
func (b *Buffer) CursorIndex() int {
	index := 0
	for i := 0; i < b.cursorY && i < len(b.lines); i++ {
		index += utf8.RuneCountInString(b.lines[i]) + 1 // +1 for newline
	}
	index += b.cursorX
	return index
}

// SetCursorPosition sets the cursor position
func (b *Buffer) SetCursorPosition(x, y int) {
	b.cursorX = x
	b.cursorY = y
	b.clampCursor()
}

// SetCursorIndex sets cursor by absolute character index
func (b *Buffer) SetCursorIndex(index int) {
	pos := 0
	for y, line := range b.lines {
		lineLen := utf8.RuneCountInString(line)
		if pos+lineLen >= index {
			b.cursorY = y
			b.cursorX = index - pos
			b.clampCursor()
			return
		}
		pos += lineLen + 1 // +1 for newline
	}
	// Past end of buffer
	if len(b.lines) > 0 {
		b.cursorY = len(b.lines) - 1
		b.cursorX = utf8.RuneCountInString(b.lines[b.cursorY])
	}
	b.clampCursor()
}

// Mode returns the current mode
func (b *Buffer) Mode() Mode {
	return b.mode
}

// SetMode sets the current mode
func (b *Buffer) SetMode(mode Mode) {
	b.mode = mode
	b.clampCursor()
}

// IsInsertMode returns true if in insert mode
func (b *Buffer) IsInsertMode() bool {
	return b.mode == ModeInsert
}

// clampCursor ensures cursor is within valid bounds
func (b *Buffer) clampCursor() {
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}

	if b.cursorY < 0 {
		b.cursorY = 0
	}
	if b.cursorY >= len(b.lines) {
		b.cursorY = len(b.lines) - 1
	}

	lineLen := utf8.RuneCountInString(b.lines[b.cursorY])
	if b.cursorX < 0 {
		b.cursorX = 0
	}

	// In normal mode, cursor can't be past last character
	// In insert mode, cursor can be at end of line (after last char)
	if b.mode == ModeNormal {
		if lineLen > 0 && b.cursorX >= lineLen {
			b.cursorX = lineLen - 1
		}
	} else {
		if b.cursorX > lineLen {
			b.cursorX = lineLen
		}
	}

	if b.cursorX < 0 {
		b.cursorX = 0
	}
}

// CharAt returns the character at the given position
func (b *Buffer) CharAt(x, y int) rune {
	if y < 0 || y >= len(b.lines) {
		return 0
	}
	line := b.lines[y]
	runes := []rune(line)
	if x < 0 || x >= len(runes) {
		return 0
	}
	return runes[x]
}

// CharUnderCursor returns the character under the cursor
func (b *Buffer) CharUnderCursor() rune {
	return b.CharAt(b.cursorX, b.cursorY)
}

// Insert inserts text at the cursor position
func (b *Buffer) Insert(text string) {
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}

	// Ensure cursorY is valid
	if b.cursorY >= len(b.lines) {
		b.cursorY = len(b.lines) - 1
	}
	if b.cursorY < 0 {
		b.cursorY = 0
	}

	line := b.lines[b.cursorY]
	runes := []rune(line)

	// Ensure cursorX is valid
	if b.cursorX > len(runes) {
		b.cursorX = len(runes)
	}
	if b.cursorX < 0 {
		b.cursorX = 0
	}

	// Split text by newlines
	parts := strings.Split(text, "\n")

	if len(parts) == 1 {
		// No newlines, simple insert
		newRunes := make([]rune, 0, len(runes)+len(text))
		newRunes = append(newRunes, runes[:b.cursorX]...)
		newRunes = append(newRunes, []rune(text)...)
		newRunes = append(newRunes, runes[b.cursorX:]...)
		b.lines[b.cursorY] = string(newRunes)
		b.cursorX += utf8.RuneCountInString(text)
	} else {
		// Multi-line insert
		before := string(runes[:b.cursorX])
		after := string(runes[b.cursorX:])

		// First line: before + first part
		b.lines[b.cursorY] = before + parts[0]

		// Middle lines
		newLines := make([]string, 0, len(b.lines)+len(parts)-1)
		newLines = append(newLines, b.lines[:b.cursorY+1]...)
		for i := 1; i < len(parts)-1; i++ {
			newLines = append(newLines, parts[i])
		}
		// Last line: last part + after
		newLines = append(newLines, parts[len(parts)-1]+after)
		newLines = append(newLines, b.lines[b.cursorY+1:]...)

		b.lines = newLines
		b.cursorY += len(parts) - 1
		b.cursorX = utf8.RuneCountInString(parts[len(parts)-1])
	}
}

// Delete deletes n characters starting at cursor
func (b *Buffer) Delete(n int) string {
	if len(b.lines) == 0 || n <= 0 {
		return ""
	}

	var deleted strings.Builder
	remaining := n

	for remaining > 0 && b.cursorY < len(b.lines) {
		line := b.lines[b.cursorY]
		runes := []rune(line)

		if b.cursorX >= len(runes) {
			// At end of line, delete newline (join with next line)
			if b.cursorY < len(b.lines)-1 {
				deleted.WriteRune('\n')
				b.lines[b.cursorY] = line + b.lines[b.cursorY+1]
				b.lines = append(b.lines[:b.cursorY+1], b.lines[b.cursorY+2:]...)
				remaining--
			} else {
				break
			}
		} else {
			// Delete characters from current line
			charsToDelete := remaining
			if b.cursorX+charsToDelete > len(runes) {
				charsToDelete = len(runes) - b.cursorX
			}

			deleted.WriteString(string(runes[b.cursorX : b.cursorX+charsToDelete]))
			newRunes := append(runes[:b.cursorX], runes[b.cursorX+charsToDelete:]...)
			b.lines[b.cursorY] = string(newRunes)
			remaining -= charsToDelete
		}
	}

	b.clampCursor()
	return deleted.String()
}

// DeleteLine deletes the current line
func (b *Buffer) DeleteLine() string {
	if len(b.lines) == 0 {
		return ""
	}

	deleted := b.lines[b.cursorY]

	if len(b.lines) == 1 {
		b.lines[0] = ""
	} else {
		b.lines = append(b.lines[:b.cursorY], b.lines[b.cursorY+1:]...)
	}

	b.clampCursor()
	return deleted
}

// DeleteToEndOfLine deletes from cursor to end of line
func (b *Buffer) DeleteToEndOfLine() string {
	if len(b.lines) == 0 {
		return ""
	}

	line := b.lines[b.cursorY]
	runes := []rune(line)

	if b.cursorX >= len(runes) {
		return ""
	}

	deleted := string(runes[b.cursorX:])
	b.lines[b.cursorY] = string(runes[:b.cursorX])
	b.clampCursor()
	return deleted
}

// ReplaceChar replaces the character under the cursor
func (b *Buffer) ReplaceChar(r rune) {
	if len(b.lines) == 0 {
		return
	}

	line := b.lines[b.cursorY]
	runes := []rune(line)

	if b.cursorX < len(runes) {
		runes[b.cursorX] = r
		b.lines[b.cursorY] = string(runes)
	}
}

// GetRegister returns the yank register content
func (b *Buffer) GetRegister() string {
	return b.register
}

// SetRegister sets the yank register content
func (b *Buffer) SetRegister(text string) {
	b.register = text
}

// Clone creates a copy of the buffer
func (b *Buffer) Clone() *Buffer {
	linesCopy := make([]string, len(b.lines))
	copy(linesCopy, b.lines)
	return &Buffer{
		lines:    linesCopy,
		cursorX:  b.cursorX,
		cursorY:  b.cursorY,
		mode:     b.mode,
		register: b.register,
	}
}

// MoveLeft moves cursor left by n characters
func (b *Buffer) MoveLeft(n int) {
	for i := 0; i < n; i++ {
		if b.cursorX > 0 {
			b.cursorX--
		} else if b.cursorY > 0 {
			b.cursorY--
			b.cursorX = utf8.RuneCountInString(b.lines[b.cursorY])
		}
	}
	b.clampCursor()
}

// MoveRight moves cursor right by n characters
func (b *Buffer) MoveRight(n int) {
	for i := 0; i < n; i++ {
		lineLen := utf8.RuneCountInString(b.lines[b.cursorY])
		if b.cursorX < lineLen {
			b.cursorX++
		} else if b.cursorY < len(b.lines)-1 {
			b.cursorY++
			b.cursorX = 0
		}
	}
	b.clampCursor()
}

// MoveUp moves cursor up by n lines
func (b *Buffer) MoveUp(n int) {
	b.cursorY -= n
	b.clampCursor()
}

// MoveDown moves cursor down by n lines
func (b *Buffer) MoveDown(n int) {
	b.cursorY += n
	b.clampCursor()
}

// MoveToLineStart moves cursor to start of line
func (b *Buffer) MoveToLineStart() {
	b.cursorX = 0
}

// MoveToLineEnd moves cursor to end of line
func (b *Buffer) MoveToLineEnd() {
	if b.cursorY < len(b.lines) {
		b.cursorX = utf8.RuneCountInString(b.lines[b.cursorY])
		b.clampCursor()
	}
}

// MoveToFirstLine moves cursor to first line
func (b *Buffer) MoveToFirstLine() {
	b.cursorY = 0
	b.clampCursor()
}

// MoveToLastLine moves cursor to last line
func (b *Buffer) MoveToLastLine() {
	b.cursorY = len(b.lines) - 1
	b.clampCursor()
}
