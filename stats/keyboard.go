package stats

// KeyboardMapping provides touch typing key assignments
// Standard QWERTY layout with proper finger assignments

// Finger assignments (0-9):
// Left hand:  0=pinky, 1=ring, 2=middle, 3=index
// Right hand: 6=index, 7=middle, 8=ring, 9=pinky

var KeyToFinger = map[rune]int{
	// Left pinky
	'q': 0, 'a': 0, 'z': 0,
	// Left ring
	'w': 1, 's': 1, 'x': 1,
	// Left middle
	'e': 2, 'd': 2, 'c': 2,
	// Left index (both regular and stretch positions)
	'r': 3, 'f': 3, 'v': 3, 't': 3, 'g': 3, 'b': 3,
	// Right index (both regular and stretch positions)
	'y': 6, 'h': 6, 'n': 6, 'u': 6, 'j': 6, 'm': 6,
	// Right middle
	'i': 7, 'k': 7,
	// Right ring
	'o': 8, 'l': 8,
	// Right pinky
	'p': 9,
}

// Hand assignments (0 = left, 1 = right)
var KeyToHand = map[rune]int{
	'q': 0, 'w': 0, 'e': 0, 'r': 0, 't': 0,
	'a': 0, 's': 0, 'd': 0, 'f': 0, 'g': 0,
	'z': 0, 'x': 0, 'c': 0, 'v': 0, 'b': 0,
	'y': 1, 'u': 1, 'i': 1, 'o': 1, 'p': 1,
	'h': 1, 'j': 1, 'k': 1, 'l': 1,
	'n': 1, 'm': 1,
}

// Row assignments (0 = top row, 1 = home row, 2 = bottom row)
var KeyToRow = map[rune]int{
	'q': 0, 'w': 0, 'e': 0, 'r': 0, 't': 0, 'y': 0, 'u': 0, 'i': 0, 'o': 0, 'p': 0,
	'a': 1, 's': 1, 'd': 1, 'f': 1, 'g': 1, 'h': 1, 'j': 1, 'k': 1, 'l': 1,
	'z': 2, 'x': 2, 'c': 2, 'v': 2, 'b': 2, 'n': 2, 'm': 2,
}

// Finger names for display
var FingerNames = map[int]string{
	0: "L Pinky",
	1: "L Ring",
	2: "L Middle",
	3: "L Index",
	6: "R Index",
	7: "R Middle",
	8: "R Ring",
	9: "R Pinky",
}

// HandNames for display
var HandNames = map[int]string{
	0: "Left",
	1: "Right",
}

// RowNames for display
var RowNames = map[int]string{
	0: "Top",
	1: "Home",
	2: "Bottom",
}

// IsSameFingerBigram returns true if two characters use the same finger
func IsSameFingerBigram(a, b rune) bool {
	fingerA, okA := KeyToFinger[a]
	fingerB, okB := KeyToFinger[b]
	if !okA || !okB {
		return false
	}
	return fingerA == fingerB
}

// IsSameHand returns true if two characters use the same hand
func IsSameHand(a, b rune) bool {
	handA, okA := KeyToHand[a]
	handB, okB := KeyToHand[b]
	if !okA || !okB {
		return false
	}
	return handA == handB
}

// GetFinger returns the finger assignment for a character (-1 if unknown)
func GetFinger(c rune) int {
	if finger, ok := KeyToFinger[c]; ok {
		return finger
	}
	return -1
}

// GetHand returns the hand assignment for a character (-1 if unknown)
func GetHand(c rune) int {
	if hand, ok := KeyToHand[c]; ok {
		return hand
	}
	return -1
}

// GetRow returns the row assignment for a character (-1 if unknown)
func GetRow(c rune) int {
	if row, ok := KeyToRow[c]; ok {
		return row
	}
	return -1
}
