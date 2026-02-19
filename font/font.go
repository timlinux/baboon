package font

// BlockLetter represents a single letter as block characters
// Each letter is 6 lines tall and variable width
// Uses smooth block elements: █ (full), ▀ (upper), ▄ (lower), ▌ (left), ▐ (right)
var BlockLetters = map[rune][]string{
	'a': {
		" ▄██▄ ",
		"██  ██",
		"██████",
		"██  ██",
		"██  ██",
		"      ",
	},
	'b': {
		"████▄ ",
		"██  ██",
		"████▀ ",
		"██  ██",
		"████▀ ",
		"      ",
	},
	'c': {
		" ▄███▄",
		"██    ",
		"██    ",
		"██    ",
		" ▀███▀",
		"      ",
	},
	'd': {
		"███▄  ",
		"██ ▀██",
		"██  ██",
		"██ ▄██",
		"███▀  ",
		"      ",
	},
	'e': {
		"██████",
		"██    ",
		"████  ",
		"██    ",
		"██████",
		"      ",
	},
	'f': {
		"██████",
		"██    ",
		"████  ",
		"██    ",
		"██    ",
		"      ",
	},
	'g': {
		" ▄███▄",
		"██    ",
		"██ ▀██",
		"██  ██",
		" ▀███▀",
		"      ",
	},
	'h': {
		"██  ██",
		"██  ██",
		"██████",
		"██  ██",
		"██  ██",
		"      ",
	},
	'i': {
		"██████",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"██████",
		"      ",
	},
	'j': {
		"██████",
		"    ██",
		"    ██",
		"██  ██",
		" ▀██▀ ",
		"      ",
	},
	'k': {
		"██  ██",
		"██ ██ ",
		"████  ",
		"██ ██ ",
		"██  ██",
		"      ",
	},
	'l': {
		"██    ",
		"██    ",
		"██    ",
		"██    ",
		"██████",
		"      ",
	},
	'm': {
		"██▄ ▄██",
		"███▀███",
		"██ ▀ ██",
		"██   ██",
		"██   ██",
		"       ",
	},
	'n': {
		"██▄  ██",
		"███▄ ██",
		"██ ████",
		"██  ███",
		"██   ██",
		"       ",
	},
	'o': {
		" ▄██▄ ",
		"██  ██",
		"██  ██",
		"██  ██",
		" ▀██▀ ",
		"      ",
	},
	'p': {
		"████▄ ",
		"██  ██",
		"████▀ ",
		"██    ",
		"██    ",
		"      ",
	},
	'q': {
		" ▄██▄ ",
		"██  ██",
		"██  ██",
		"██ ▄██",
		" ▀████",
		"      ",
	},
	'r': {
		"████▄ ",
		"██  ██",
		"████▀ ",
		"██ ▀█ ",
		"██  ██",
		"      ",
	},
	's': {
		" ▄████",
		"██    ",
		" ▀██▄ ",
		"    ██",
		"████▀ ",
		"      ",
	},
	't': {
		"██████",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"      ",
	},
	'u': {
		"██  ██",
		"██  ██",
		"██  ██",
		"██  ██",
		" ▀██▀ ",
		"      ",
	},
	'v': {
		"██  ██",
		"██  ██",
		"██  ██",
		" ▀██▀ ",
		"  ██  ",
		"      ",
	},
	'w': {
		"██   ██",
		"██   ██",
		"██ ▄ ██",
		"███▀███",
		"██▀ ▀██",
		"       ",
	},
	'x': {
		"██  ██",
		" ▀██▀ ",
		"  ██  ",
		" ▄██▄ ",
		"██  ██",
		"      ",
	},
	'y': {
		"██  ██",
		" ▀██▀ ",
		"  ██  ",
		"  ██  ",
		"  ██  ",
		"      ",
	},
	'z': {
		"██████",
		"   ▄█▀",
		"  ██  ",
		" ▄█▀  ",
		"██████",
		"      ",
	},
	' ': {
		"   ",
		"   ",
		"   ",
		"   ",
		"   ",
		"   ",
	},
	',': {
		"    ",
		"    ",
		"    ",
		" ██ ",
		" ▄█ ",
		"    ",
	},
	'.': {
		"    ",
		"    ",
		"    ",
		"    ",
		" ██ ",
		"    ",
	},
	';': {
		"    ",
		" ██ ",
		"    ",
		" ██ ",
		" ▄█ ",
		"    ",
	},
	':': {
		"    ",
		" ██ ",
		"    ",
		" ██ ",
		"    ",
		"    ",
	},
	'!': {
		" ██ ",
		" ██ ",
		" ██ ",
		"    ",
		" ██ ",
		"    ",
	},
	'?': {
		" ▄██▄ ",
		"█▀  ██",
		"   ██ ",
		"      ",
		"  ██  ",
		"      ",
	},
}

const LetterHeight = 6
const LetterSpacing = 1

// RenderWord renders a word as block letters, returning each line separately
// Each rune in the returned slices corresponds to one character of the original word
func RenderWord(word string) [][]string {
	lines := make([][]string, LetterHeight)
	for i := range lines {
		lines[i] = make([]string, 0, len(word))
	}

	for _, char := range word {
		letter, ok := BlockLetters[char]
		if !ok {
			// Use space for unknown characters
			letter = BlockLetters[' ']
		}

		for lineIdx := 0; lineIdx < LetterHeight; lineIdx++ {
			if lineIdx < len(letter) {
				lines[lineIdx] = append(lines[lineIdx], letter[lineIdx])
			} else {
				lines[lineIdx] = append(lines[lineIdx], "")
			}
		}
	}

	return lines
}

// GetLetterWidth returns the width of a letter
func GetLetterWidth(char rune) int {
	letter, ok := BlockLetters[char]
	if !ok || len(letter) == 0 {
		return 3
	}
	maxWidth := 0
	for _, line := range letter {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	return maxWidth
}
