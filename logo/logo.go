package logo

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"strings"
)

//go:embed baboon-logo.png
var logoPNG []byte

// Protocol represents the terminal graphics protocol to use
type Protocol int

const (
	ProtocolNone  Protocol = iota
	ProtocolKitty          // Kitty graphics protocol
	ProtocolSixel          // Sixel graphics
)

// DetectProtocol returns the best available terminal graphics protocol.
func DetectProtocol() Protocol {
	// Check for Kitty
	termProgram := os.Getenv("TERM_PROGRAM")
	term := os.Getenv("TERM")
	if termProgram == "kitty" || strings.Contains(term, "kitty") {
		return ProtocolKitty
	}

	// Check for terminals known to support Kitty graphics protocol
	if termProgram == "WezTerm" || termProgram == "wezterm" {
		return ProtocolKitty
	}

	// Ghostty supports Kitty graphics protocol
	if termProgram == "ghostty" {
		return ProtocolKitty
	}

	// Check for Sixel support via common terminals
	// foot, mlterm, xterm (with -ti vt340), contour, etc.
	if termProgram == "foot" || termProgram == "mlterm" || termProgram == "contour" {
		return ProtocolSixel
	}

	// TERM_FEATURES or VTE-based terminals with Sixel
	if os.Getenv("SIXEL_SUPPORT") == "1" {
		return ProtocolSixel
	}

	return ProtocolNone
}

// Render returns the logo as terminal escape sequences for the detected protocol.
// The targetWidth is in terminal cell columns. Returns empty string if no protocol supported.
func Render(targetCols int) string {
	protocol := DetectProtocol()
	switch protocol {
	case ProtocolKitty:
		return renderKitty(targetCols)
	case ProtocolSixel:
		return renderSixel(targetCols)
	default:
		return ""
	}
}

// renderKitty renders the logo using the Kitty graphics protocol.
// Sends the PNG directly (f=100) with chunked transmission.
func renderKitty(targetCols int) string {
	encoded := base64.StdEncoding.EncodeToString(logoPNG)

	var sb strings.Builder
	const chunkSize = 4096

	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunk := encoded[i:end]

		isFirst := i == 0
		isLast := end >= len(encoded)

		if isFirst && isLast {
			// Single chunk
			sb.WriteString(fmt.Sprintf("\x1b_Ga=T,f=100,c=%d,q=2;%s\x1b\\", targetCols, chunk))
		} else if isFirst {
			// First chunk of multi-chunk
			sb.WriteString(fmt.Sprintf("\x1b_Ga=T,f=100,c=%d,m=1,q=2;%s\x1b\\", targetCols, chunk))
		} else if isLast {
			// Last chunk
			sb.WriteString(fmt.Sprintf("\x1b_Gm=0;%s\x1b\\", chunk))
		} else {
			// Middle chunk
			sb.WriteString(fmt.Sprintf("\x1b_Gm=1;%s\x1b\\", chunk))
		}
	}

	// Calculate approximate height in rows (assuming ~2:1 aspect ratio for cells)
	// The image is 400x224, so at targetCols width, height is roughly targetCols*224/400/2 rows
	rows := (targetCols * 224) / (400 * 2)
	if rows < 1 {
		rows = 1
	}

	// Add empty lines to reserve space for the image
	for i := 0; i < rows; i++ {
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderSixel renders the logo using the Sixel graphics protocol.
func renderSixel(targetCols int) string {
	img, _, err := image.Decode(bytes.NewReader(logoPNG))
	if err != nil {
		return ""
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Scale to target width in pixels (assume ~8px per cell column)
	targetPx := targetCols * 8
	if targetPx > w {
		targetPx = w
	}
	scale := float64(targetPx) / float64(w)
	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)

	// Simple nearest-neighbor resize
	resized := image.NewRGBA(image.Rect(0, 0, newW, newH))
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			srcX := int(float64(x) / scale)
			srcY := int(float64(y) / scale)
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	return encodeSixel(resized)
}

// encodeSixel encodes an image as a Sixel escape sequence.
// Uses a simple 256-color quantization.
func encodeSixel(img *image.RGBA) string {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w == 0 || h == 0 {
		return ""
	}

	var sb strings.Builder

	// Build a simple color palette (6x6x6 color cube = 216 colors)
	type palEntry struct {
		r, g, b uint8
	}
	palette := make([]palEntry, 216)
	for i := 0; i < 216; i++ {
		r := (i / 36) * 51
		g := ((i / 6) % 6) * 51
		b := (i % 6) * 51
		palette[i] = palEntry{uint8(r), uint8(g), uint8(b)}
	}

	// Nearest palette color
	nearest := func(c color.Color) int {
		r, g, b, a := c.RGBA()
		if a < 128<<8 {
			return -1 // transparent
		}
		ri := int((r>>8)+25) / 51
		gi := int((g>>8)+25) / 51
		bi := int((b>>8)+25) / 51
		if ri > 5 {
			ri = 5
		}
		if gi > 5 {
			gi = 5
		}
		if bi > 5 {
			bi = 5
		}
		return ri*36 + gi*6 + bi
	}

	// Sixel header
	sb.WriteString("\x1bPq\n")

	// Define palette
	for i, p := range palette {
		// Sixel uses percentage 0-100 for RGB
		rp := int(p.r) * 100 / 255
		gp := int(p.g) * 100 / 255
		bp := int(p.b) * 100 / 255
		sb.WriteString(fmt.Sprintf("#%d;2;%d;%d;%d", i, rp, gp, bp))
	}
	sb.WriteString("\n")

	// Encode pixels in bands of 6 rows
	for band := 0; band*6 < h; band++ {
		startY := band * 6

		// For each color used in this band, emit a row
		usedColors := make(map[int]bool)
		for y := startY; y < startY+6 && y < h; y++ {
			for x := 0; x < w; x++ {
				ci := nearest(img.At(x, y))
				if ci >= 0 {
					usedColors[ci] = true
				}
			}
		}

		first := true
		for ci := range usedColors {
			if !first {
				sb.WriteByte('$') // carriage return (go back to start of band)
			}
			first = false

			sb.WriteString(fmt.Sprintf("#%d", ci))

			for x := 0; x < w; x++ {
				sixelByte := byte(0)
				for bit := 0; bit < 6; bit++ {
					y := startY + bit
					if y < h {
						pi := nearest(img.At(x, y))
						if pi == ci {
							sixelByte |= 1 << uint(bit)
						}
					}
				}
				sb.WriteByte(sixelByte + 63) // Sixel data byte = value + 63
			}
		}

		sb.WriteByte('-') // next band (newline)
	}

	// Sixel terminator
	sb.WriteString("\x1b\\")

	return sb.String()
}
