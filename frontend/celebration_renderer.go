// Package frontend provides the terminal user interface for the typing practice application.
package frontend

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/timlinux/blockfont"
)

// RenderCelebrationScreen renders the celebration animation
func (r *Renderer) RenderCelebrationScreen(state *CelebrationState) string {
	if state.Phase == PhaseMessage {
		return r.renderCelebrationMessage(state)
	}
	return r.renderFireworks(state)
}

// renderFireworks renders the fireworks particle effect phase
func (r *Renderer) renderFireworks(state *CelebrationState) string {
	// Create a screen buffer
	buffer := make([][]screenCell, r.height)
	for y := range r.height {
		buffer[y] = make([]screenCell, r.width)
		for x := range r.width {
			buffer[y][x] = screenCell{char: ' ', color: 0}
		}
	}

	// Render text grid (with destroyed cells as spaces)
	if state.TextGrid != nil {
		for y, row := range state.TextGrid {
			if y >= r.height {
				break
			}
			for x, cell := range row {
				if x >= r.width {
					break
				}
				if !cell.Destroyed && cell.Char != ' ' && cell.Char != 0 {
					buffer[y][x] = screenCell{
						char:  cell.Char,
						color: cell.Color,
					}
				}
			}
		}
	}

	// Render particles on top
	for _, p := range state.Particles {
		x := int(p.X)
		y := int(p.Y)

		if x < 0 || x >= r.width || y < 0 || y >= r.height {
			continue
		}

		// Fade color based on brightness
		color := p.Color
		if p.Brightness < 0.5 {
			// Fade to darker version
			color = fadeColor(color, p.Brightness*2)
		}

		buffer[y][x] = screenCell{
			char:  p.Char,
			color: color,
		}
	}

	// Convert buffer to string
	return r.bufferToString(buffer, state.Width, state.Height)
}

// screenCell represents a single cell in the screen buffer
type screenCell struct {
	char  rune
	color int // 256-color code, 0 for default
}

// fadeColor fades a color toward black based on factor (0-1)
func fadeColor(color int, factor float64) int {
	if factor <= 0 {
		return 232 // Almost black in 256-color palette
	}
	if factor >= 1 {
		return color
	}

	// For grayscale range (232-255)
	if color >= 232 && color <= 255 {
		newGray := 232 + int(float64(color-232)*factor)
		newGray = max(newGray, 232)
		return newGray
	}

	// For the main 6x6x6 color cube (16-231)
	// Just shift toward the lower intensity version
	if color >= 16 && color <= 231 {
		color -= 16
		r := color / 36
		g := (color % 36) / 6
		b := color % 6

		r = int(float64(r) * factor)
		g = int(float64(g) * factor)
		b = int(float64(b) * factor)

		return 16 + r*36 + g*6 + b
	}

	// For basic 16 colors, just return a gray
	return 240
}

// bufferToString converts a screen buffer to a string
func (r *Renderer) bufferToString(buffer [][]screenCell, width, height int) string {
	var sb strings.Builder

	for y := range height {
		if y >= len(buffer) {
			sb.WriteString("\n")
			continue
		}

		row := buffer[y]
		var lineBuilder strings.Builder

		// Track runs of same color for efficiency
		currentColor := -1
		var runBuilder strings.Builder

		for x := range width {
			var cell screenCell
			if x < len(row) {
				cell = row[x]
			} else {
				cell = screenCell{char: ' ', color: 0}
			}

			char := cell.char
			if char == 0 {
				char = ' '
			}

			if cell.color != currentColor {
				// Flush previous run
				if runBuilder.Len() > 0 {
					if currentColor > 0 {
						style := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", currentColor)))
						lineBuilder.WriteString(style.Render(runBuilder.String()))
					} else {
						lineBuilder.WriteString(runBuilder.String())
					}
					runBuilder.Reset()
				}
				currentColor = cell.color
			}

			runBuilder.WriteRune(char)
		}

		// Flush final run
		if runBuilder.Len() > 0 {
			if currentColor > 0 {
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", currentColor)))
				lineBuilder.WriteString(style.Render(runBuilder.String()))
			} else {
				lineBuilder.WriteString(runBuilder.String())
			}
		}

		sb.WriteString(lineBuilder.String())
		if y < height-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// renderCelebrationMessage renders the "PERSONAL BEST!" block font message
func (r *Renderer) renderCelebrationMessage(state *CelebrationState) string {
	// Build the message based on what types of bests were achieved
	var message string
	var valueStr string

	if state.BestTypes&BestWPM != 0 {
		message = "PERSONAL BEST"
		valueStr = fmt.Sprintf("%.0f WPM", state.WPMValue)
	} else if state.BestTypes&BestAccuracy != 0 {
		message = "PERFECT ACCURACY"
		valueStr = fmt.Sprintf("%.1f%%", state.AccuracyValue)
	} else if state.BestTypes&BestTime != 0 {
		message = "FASTEST TIME"
		valueStr = fmt.Sprintf("%.1fs", state.TimeValue)
	} else {
		message = "NEW RECORD"
		valueStr = fmt.Sprintf("%.0f WPM", state.WPMValue)
	}

	// Render block font message
	messageLines := blockfont.RenderWord(message)

	// Style with gold color
	goldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)

	// Build centered content
	var content []string

	// Add some vertical padding
	for range 3 {
		content = append(content, "")
	}

	// Add block font message (each line)
	for lineIdx := range blockfont.LetterHeight {
		var lineBuilder strings.Builder
		for charIdx, letterLine := range messageLines[lineIdx] {
			lineBuilder.WriteString(letterLine)
			if charIdx < len(messageLines[lineIdx])-1 {
				lineBuilder.WriteString(" ")
			}
		}
		content = append(content, goldStyle.Render(lineBuilder.String()))
	}

	content = append(content, "")
	content = append(content, "")

	// Add value below
	content = append(content, valueStyle.Render(valueStr))

	content = append(content, "")
	content = append(content, "")

	// Add sparkle decoration
	sparkleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	sparkles := "✦  ·  ✧  ·  ✦  ·  ✧  ·  ✦"
	content = append(content, sparkleStyle.Render(sparkles))

	content = append(content, "")

	// Add "Press any key to continue" hint
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	content = append(content, hintStyle.Render("Press any key to continue"))

	// Center everything
	mainContent := lipgloss.JoinVertical(lipgloss.Center, content...)

	// Fixed header at top
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("14")).
		Bold(true)
	header := lipgloss.PlaceHorizontal(r.width, lipgloss.Center,
		headerStyle.Render("🐒 BABOON - Typing Practice"))

	// Calculate heights
	headerHeight := 1
	contentHeight := strings.Count(mainContent, "\n") + 1
	availableHeight := r.height - headerHeight - 2

	// Calculate top padding to center main content
	topPadding := max((availableHeight-contentHeight)/2, 0)

	// Build full screen layout
	var fullContent strings.Builder

	fullContent.WriteString(header)
	fullContent.WriteString("\n")

	for range topPadding {
		fullContent.WriteString("\n")
	}

	centeredMain := lipgloss.PlaceHorizontal(r.width, lipgloss.Center, mainContent)
	fullContent.WriteString(centeredMain)

	// Fill remaining space
	currentHeight := headerHeight + 1 + topPadding + contentHeight
	for i := currentHeight; i < r.height; i++ {
		fullContent.WriteString("\n")
	}

	return fullContent.String()
}
