// SPDX-FileCopyrightText: 2026 Tim Sutton / Kartoza
// SPDX-License-Identifier: MIT

package blockfont

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WidgetOptions configures the Widget behavior
type WidgetOptions struct {
	// Width is the target width for rendering
	Width int
	// Height is the target height (optional)
	Height int
	// Alignment controls text alignment
	Alignment Alignment
	// VimMode enables vim-style editing
	VimMode bool
	// Animate enables animations
	Animate bool
	// WordWrap enables word wrapping
	WordWrap bool
	// CursorStyle controls cursor appearance
	CursorStyle CursorStyle
	// Theme provides colors and styles
	Theme Theme
}

// DefaultWidgetOptions returns sensible defaults
func DefaultWidgetOptions() WidgetOptions {
	return WidgetOptions{
		Width:       80,
		Height:      0,
		Alignment:   AlignLeft,
		VimMode:     false,
		Animate:     false,
		WordWrap:    false,
		CursorStyle: CursorBlock,
		Theme:       DefaultTheme,
	}
}

// Widget is a high-level component for rendering block text in bubbletea applications.
// It implements tea.Model for easy integration.
type Widget struct {
	text       string
	buffer     *Buffer
	animator   *Animator
	options    WidgetOptions
	highlights []CharHighlight
	colors     []lipgloss.Color
	focused    bool
}

// NewWidget creates a new Widget with the given options
func NewWidget(opts WidgetOptions) *Widget {
	w := &Widget{
		text:     "",
		buffer:   NewBuffer(""),
		animator: NewAnimator(),
		options:  opts,
		focused:  true,
	}
	return w
}

// SetText sets the widget text
func (w *Widget) SetText(text string) {
	w.text = text
	if w.options.VimMode {
		w.buffer.SetText(text)
	}
	w.highlights = nil
	w.colors = nil

	if w.options.Animate {
		w.animator.TriggerTransition(TransitionFadeIn)
	}
}

// Text returns the current text
func (w *Widget) Text() string {
	if w.options.VimMode {
		return w.buffer.Text()
	}
	return w.text
}

// SetHighlights sets character-level highlights
func (w *Widget) SetHighlights(highlights []CharHighlight) {
	w.highlights = highlights
}

// ColorCharacter sets a color for a specific character index
func (w *Widget) ColorCharacter(index int, color lipgloss.Color) {
	// Ensure colors slice is large enough
	if w.colors == nil || len(w.colors) <= index {
		newColors := make([]lipgloss.Color, index+1)
		if w.colors != nil {
			copy(newColors, w.colors)
		}
		w.colors = newColors
	}
	w.colors[index] = color
}

// ColorRange sets a color for a range of characters
func (w *Widget) ColorRange(start, end int, color lipgloss.Color) {
	for i := start; i < end; i++ {
		w.ColorCharacter(i, color)
	}
}

// ResetColors clears all character colors
func (w *Widget) ResetColors() {
	w.colors = nil
}

// Focus sets the widget focus state
func (w *Widget) Focus() {
	w.focused = true
}

// Blur removes focus from the widget
func (w *Widget) Blur() {
	w.focused = false
}

// IsFocused returns whether the widget is focused
func (w *Widget) IsFocused() bool {
	return w.focused
}

// Buffer returns the underlying buffer (for vim mode)
func (w *Widget) Buffer() *Buffer {
	return w.buffer
}

// Mode returns the current vim mode
func (w *Widget) Mode() Mode {
	return w.buffer.Mode()
}

// Init implements tea.Model
func (w *Widget) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (w *Widget) Update(msg tea.Msg) (*Widget, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !w.focused {
			return w, nil
		}

		if w.options.VimMode {
			cmd = w.handleVimKey(msg)
		}

	case AnimationTickMsg:
		if w.animator.Update() {
			cmd = w.tickAnimation()
		}
	}

	return w, cmd
}

// AnimationTickMsg is sent when animation should update
type AnimationTickMsg struct{}

// tickAnimation returns a command to continue animation
func (w *Widget) tickAnimation() tea.Cmd {
	return tea.Tick(AnimationInterval, func(t time.Time) tea.Msg {
		return AnimationTickMsg{}
	})
}

// handleVimKey handles keyboard input in vim mode
func (w *Widget) handleVimKey(msg tea.KeyMsg) tea.Cmd {
	switch w.buffer.Mode() {
	case ModeNormal:
		return w.handleNormalMode(msg)
	case ModeInsert:
		return w.handleInsertMode(msg)
	}
	return nil
}

// handleNormalMode handles keys in normal mode
func (w *Widget) handleNormalMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "h", "left":
		w.buffer.MoveLeft(1)
	case "l", "right":
		w.buffer.MoveRight(1)
	case "j", "down":
		w.buffer.MoveDown(1)
	case "k", "up":
		w.buffer.MoveUp(1)
	case "0", "home":
		w.buffer.MoveToLineStart()
	case "$", "end":
		w.buffer.MoveToLineEnd()
	case "g":
		w.buffer.MoveToFirstLine()
	case "G":
		w.buffer.MoveToLastLine()
	case "i":
		w.buffer.SetMode(ModeInsert)
	case "a":
		w.buffer.MoveRight(1)
		w.buffer.SetMode(ModeInsert)
	case "A":
		w.buffer.MoveToLineEnd()
		w.buffer.SetMode(ModeInsert)
	case "I":
		w.buffer.MoveToLineStart()
		w.buffer.SetMode(ModeInsert)
	case "o":
		w.buffer.MoveToLineEnd()
		w.buffer.Insert("\n")
		w.buffer.SetMode(ModeInsert)
	case "O":
		w.buffer.MoveToLineStart()
		w.buffer.Insert("\n")
		w.buffer.MoveUp(1)
		w.buffer.SetMode(ModeInsert)
	case "x":
		w.buffer.Delete(1)
	case "dd":
		w.buffer.SetRegister(w.buffer.DeleteLine())
	case "D":
		w.buffer.SetRegister(w.buffer.DeleteToEndOfLine())
	case "r":
		// Would need next key for replace
	case "u":
		// Undo not implemented
	}
	return nil
}

// handleInsertMode handles keys in insert mode
func (w *Widget) handleInsertMode(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		w.buffer.SetMode(ModeNormal)
		w.buffer.MoveLeft(1)
	case "backspace":
		if x, _ := w.buffer.CursorPosition(); x > 0 {
			w.buffer.MoveLeft(1)
			w.buffer.Delete(1)
		}
	case "delete":
		w.buffer.Delete(1)
	case "left":
		w.buffer.MoveLeft(1)
	case "right":
		w.buffer.MoveRight(1)
	case "up":
		w.buffer.MoveUp(1)
	case "down":
		w.buffer.MoveDown(1)
	case "home":
		w.buffer.MoveToLineStart()
	case "end":
		w.buffer.MoveToLineEnd()
	case "enter":
		w.buffer.Insert("\n")
	default:
		// Insert character
		if len(msg.Runes) > 0 {
			w.buffer.Insert(string(msg.Runes))
		}
	}
	return nil
}

// View implements tea.Model
func (w *Widget) View() string {
	text := w.Text()
	if text == "" {
		return ""
	}

	lines := w.Render()

	// Apply alignment
	if w.options.Width > 0 {
		lines = AlignLines(lines, w.options.Alignment, w.options.Width)
	}

	return strings.Join(lines, "\n")
}

// Render returns the rendered block text as lines
func (w *Widget) Render() []string {
	text := w.Text()
	if text == "" {
		return []string{}
	}

	// Use the full-featured RenderWithCursor for vim mode
	if w.options.VimMode && w.focused {
		cursorIdx := w.buffer.CursorIndex()
		isInsert := w.buffer.IsInsertMode()
		maxWidth := 0
		if w.options.WordWrap {
			maxWidth = w.options.Width
		}
		return RenderWithCursor(text, cursorIdx, w.highlights, isInsert, maxWidth, w.options.Theme)
	}

	// Simple rendering without cursor
	if w.highlights != nil {
		// Use RenderWithCursor but with cursor at -1 (off-screen)
		maxWidth := 0
		if w.options.WordWrap {
			maxWidth = w.options.Width
		}
		return RenderWithCursor(text, -1, w.highlights, false, maxWidth, w.options.Theme)
	}

	// Basic rendering without highlights or cursor
	lines := RenderPlainText(text, "")
	return lines
}

// RenderCentered renders the text centered within the given width
func (w *Widget) RenderCentered(width int) []string {
	lines := w.Render()
	return CenterLines(lines, width)
}

// SetWidth sets the widget width
func (w *Widget) SetWidth(width int) {
	w.options.Width = width
}

// SetHeight sets the widget height
func (w *Widget) SetHeight(height int) {
	w.options.Height = height
}

// SetAlignment sets the text alignment
func (w *Widget) SetAlignment(alignment Alignment) {
	w.options.Alignment = alignment
}

// EnableVimMode enables vim-style editing
func (w *Widget) EnableVimMode() {
	w.options.VimMode = true
	w.buffer.SetText(w.text)
}

// DisableVimMode disables vim-style editing
func (w *Widget) DisableVimMode() {
	w.options.VimMode = false
	w.text = w.buffer.Text()
}

// EnableAnimations enables animations
func (w *Widget) EnableAnimations() {
	w.options.Animate = true
}

// DisableAnimations disables animations
func (w *Widget) DisableAnimations() {
	w.options.Animate = false
}

// TriggerAnimation starts an animation
func (w *Widget) TriggerAnimation(t TransitionType) tea.Cmd {
	w.animator.TriggerTransition(t)
	return w.tickAnimation()
}

// IsAnimating returns whether an animation is in progress
func (w *Widget) IsAnimating() bool {
	return w.animator.IsAnimating
}

// SetTheme sets the color theme
func (w *Widget) SetTheme(theme Theme) {
	w.options.Theme = theme
}
