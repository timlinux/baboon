package words

import "strings"

// CommonWords contains the 1000 most common English words (British English)
// All words are lowercase only
var CommonWords = []string{
	"the", "be", "to", "of", "and", "a", "in", "that", "have", "it",
	"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
	"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
	"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
	"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
	"when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
	"people", "into", "year", "your", "good", "some", "could", "them", "see", "other",
	"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
	"back", "after", "use", "two", "how", "our", "work", "first", "well", "way",
	"even", "new", "want", "because", "any", "these", "give", "day", "most", "us",
	"is", "was", "are", "been", "has", "had", "were", "said", "each", "she",
	"may", "find", "long", "down", "did", "get", "made", "before", "might", "many",
	"write", "must", "water", "word", "such", "call", "side", "where", "help", "through",
	"much", "still", "same", "while", "great", "last", "small", "own", "found", "those",
	"never", "show", "under", "little", "every", "house", "world", "put", "old", "being",
	"once", "went", "next", "hand", "high", "off", "end", "live", "night", "school",
	"another", "away", "home", "something", "need", "study", "keep", "might", "point", "start",
	"head", "story", "city", "play", "spell", "young", "few", "enough", "always", "watch",
	"three", "letter", "until", "far", "children", "got", "walk", "example", "ease", "paper",
	"always", "music", "often", "run", "late", "hard", "set", "food", "both", "between",
	"name", "line", "right", "boy", "soon", "grow", "state", "left", "near", "kind",
	"together", "thought", "father", "mother", "white", "seem", "began", "country", "family", "fact",
	"earth", "part", "place", "life", "open", "read", "north", "south", "change", "question",
	"tell", "turn", "move", "sea", "light", "self", "face", "number", "group", "idea",
	"woman", "bring", "book", "problem", "better", "learn", "ask", "should", "britain", "plant",
	"above", "girl", "sometimes", "without", "stop", "four", "second", "later", "programme", "order",
	"man", "men", "woman", "women", "money", "eye", "thing", "room", "feel", "since",
	"area", "today", "already", "during", "minute", "morning", "within", "hear", "moment", "body",
	"close", "nothing", "certain", "along", "stand", "against", "government", "system", "car", "power",
	"company", "believe", "hold", "week", "person", "case", "service", "become", "possible", "mind",
	"however", "member", "pay", "law", "meet", "month", "true", "door", "early", "course",
	"real", "continue", "nice", "public", "include", "whether", "big", "half", "sure", "reason",
	"free", "black", "though", "level", "view", "important", "using", "centre", "best", "sense",
	"across", "business", "known", "working", "study", "health", "result", "development", "national", "social",
	"least", "whole", "support", "war", "control", "child", "given", "red", "present", "table",
	"market", "office", "action", "care", "perhaps", "perhaps", "window", "report", "decide", "building",
	"death", "yes", "behind", "reach", "local", "couple", "job", "position", "produce", "effect",
	"political", "interest", "five", "ago", "history", "land", "different", "economic", "international", "process",
	"voice", "art", "finally", "cost", "strong", "plan", "party", "available", "full", "class",
	"base", "common", "information", "research", "human", "love", "story", "court", "air", "record",
	"happen", "provide", "friend", "several", "field", "role", "appear", "age", "street", "likely",
	"police", "technology", "future", "inside", "paper", "clear", "special", "front", "rate", "value",
	"language", "among", "stay", "six", "period", "policy", "computer", "million", "rather", "create",
	"tax", "similar", "relationship", "personal", "single", "return", "experience", "form", "price", "british",
	"image", "accept", "doctor", "require", "type", "military", "agree", "actually", "theory", "step",
	"choose", "foreign", "involve", "hospital", "usually", "cause", "student", "private", "news", "game",
	"offer", "degree", "term", "material", "able", "low", "performance", "either", "cover", "behaviour",
	"energy", "building", "piece", "address", "describe", "player", "third", "current", "expect", "design",
	"training", "model", "team", "staff", "risk", "choice", "individual", "wrong", "data", "range",
	"simple", "activity", "film", "bit", "particular", "test", "size", "consider", "site", "pass",
	"star", "argue", "measure", "specific", "management", "direction", "indeed", "cut", "sign", "population",
	"director", "suddenly", "recently", "analysis", "project", "effort", "method", "physical", "discussion", "evidence",
	"major", "stock", "whatever", "opportunity", "scene", "environment", "practice", "financial", "season", "science",
	"music", "pressure", "natural", "treatment", "section", "college", "reduce", "response", "concern", "ten",
	"activity", "traditional", "budget", "culture", "amount", "seven", "growth", "figure", "factor", "success",
	"general", "responsibility", "subject", "role", "allow", "beyond", "deal", "rest", "community", "society",
	"recent", "article", "conference", "standard", "quickly", "interview", "character", "break", "resource", "approach",
	"eight", "increase", "strategy", "statement", "structure", "surface", "movement", "short", "add", "nature",
	"legal", "product", "medical", "situation", "whose", "event", "poor", "author", "reality", "anyone",
	"although", "expert", "ready", "final", "series", "production", "speech", "leader", "board", "drug",
	"unit", "election", "worker", "condition", "beautiful", "peace", "century", "trade", "attention", "nine",
	"education", "property", "detail", "goal", "attack", "operation", "skill", "picture", "station", "loss",
	"manager", "trouble", "brother", "summer", "campaign", "focus", "everything", "benefit", "everyone", "trial",
	"cell", "religious", "network", "apply", "audience", "knowledge", "memory", "administration", "bank", "purpose",
	"benefit", "protect", "coach", "forward", "security", "feeling", "capital", "everything", "challenge", "race",
	"fill", "share", "access", "artist", "debate", "quality", "chance", "feature", "agency", "difference",
	"central", "pull", "note", "establish", "explain", "identify", "organisation", "hair", "notice", "crime",
	"hot", "arm", "travel", "total", "average", "democratic", "space", "follow", "pattern", "professor",
	"weight", "marriage", "dead", "list", "senior", "affect", "professional", "hotel", "maintain", "miss",
	"kitchen", "blue", "wonder", "represent", "tree", "contain", "heart", "various", "ball", "green",
	"dark", "eat", "doctor", "colour", "stage", "shot", "trouble", "animal", "region", "assume",
	"develop", "popular", "collection", "fund", "drop", "perform", "hit", "bed", "remain", "original",
	"floor", "official", "deep", "foot", "carry", "pain", "behaviour", "determine", "fight", "success",
	"bottom", "basic", "civil", "cultural", "safe", "degree", "garden", "necessary", "claim", "complete",
	"account", "seek", "district", "generation", "threat", "exist", "employee", "exactly", "charge", "draw",
	"institution", "candidate", "former", "entire", "technology", "participate", "rule", "enter", "science", "disease",
	"tend", "effort", "date", "potential", "sport", "replace", "source", "edge", "card", "wide",
	"movement", "release", "fear", "prepare", "smile", "material", "oil", "modern", "cry", "rich",
	"dream", "former", "finger", "stuff", "imagine", "glass", "defence", "corner", "nearly", "impact",
	"wall", "video", "basic", "express", "conference", "independent", "appropriate", "western", "hang", "opinion",
	"promise", "save", "trouble", "mention", "drive", "wind", "tiny", "seat", "critical", "item",
	"decade", "pretty", "realise", "rise", "yard", "refer", "stone", "reflect", "limited", "huge",
	"connection", "count", "associate", "style", "extend", "absolutely", "damage", "variety", "spring", "complex",
	"truth", "collect", "commission", "executive", "recognise", "reform", "category", "significant", "citizen", "majority",
	"serve", "magazine", "demand", "prove", "content", "financial", "magazine", "dinner", "hotel", "lack",
	"spot", "reform", "track", "press", "ticket", "sleep", "train", "fly", "sell", "store",
	"influence", "teach", "mouth", "immediately", "afternoon", "extremely", "version", "conflict", "forget", "achieve",
	"temperature", "majority", "relatively", "exercise", "customer", "approve", "survey", "mass", "male", "directly",
	"female", "clearly", "primary", "remove", "credit", "income", "grant", "restaurant", "serve", "title",
	"chief", "drink", "overall", "positive", "cool", "sample", "bus", "speed", "native", "prevent",
	"battle", "address", "largely", "option", "element", "mental", "judge", "construction", "domestic", "handle",
	"acquire", "moral", "observer", "volume", "decade", "chief", "failure", "speech", "assessment", "balance",
	"waste", "difficulty", "tradition", "separate", "respect", "examine", "touch", "rely", "weapon", "initial",
	"reform", "trip", "generally", "slightly", "publish", "lift", "attempt", "shift", "profit", "reform",
	"discover", "circle", "struggle", "literature", "destroy", "newspaper", "proportion", "provision", "strike", "crowd",
	"essential", "sector", "length", "seriously", "hate", "spirit", "quiet", "equipment", "violence", "favour",
	"labour", "honour", "neighbour", "travelling", "catalogue", "theatre", "metre", "litre", "fibre", "analyse",
	"apologise", "capitalise", "categorise", "characterise", "civilise", "colonise", "criticise", "customise", "emphasise", "energise",
	"equalise", "finalise", "generalise", "globalise", "harmonise", "hospitalise", "hypothesise", "idealise", "immunise", "initialise",
	"legalise", "liberalise", "localise", "materialise", "maximise", "memorise", "minimise", "modernise", "naturalise", "neutralise",
}

// GetRandomWords returns n random words from the common words list
// All words are converted to lowercase and empty words are skipped
func GetRandomWords(n int, rng func(int) int) []string {
	result := make([]string, 0, n)
	for len(result) < n {
		word := strings.ToLower(CommonWords[rng(len(CommonWords))])
		// Skip empty words or words with only spaces
		if len(strings.TrimSpace(word)) > 0 {
			result = append(result, word)
		}
	}
	return result
}

// WordsPerRound is the fixed number of words per round
const WordsPerRound = 30

// TargetCharacters is the exact number of characters per round
const TargetCharacters = 150

// LetterStats represents frequency and accuracy data for a letter
type LetterStats struct {
	Presented int // Number of times presented
	Correct   int // Number of times typed correctly
}

// LetterData represents the combined frequency and accuracy data for letters
type LetterData map[string]LetterStats

// scoreWordByLetterData scores a word based on:
// 1. How many underrepresented letters it contains (frequency balancing)
// 2. How many letters the user has low accuracy on (practice weak letters)
// Higher score = more underrepresented + lower accuracy letters = should be preferred
func scoreWordByLetterData(word string, letterData LetterData) float64 {
	if len(letterData) == 0 {
		return 1.0 // No history, equal weight
	}

	// Find max frequency to normalize
	var maxFreq int
	for _, stats := range letterData {
		if stats.Presented > maxFreq {
			maxFreq = stats.Presented
		}
	}

	if maxFreq == 0 {
		return 1.0 // No data yet
	}

	// Calculate combined score for each letter
	var score float64
	for _, char := range word {
		letter := string(char)
		stats := letterData[letter]

		// Frequency score: low frequency = high score (range 0 to 1)
		normalizedFreq := float64(stats.Presented) / float64(maxFreq)
		freqScore := 1.0 - normalizedFreq

		// Accuracy score: low accuracy = high score (range 0 to 1)
		var accScore float64
		if stats.Presented > 0 {
			accuracy := float64(stats.Correct) / float64(stats.Presented)
			accScore = 1.0 - accuracy // Invert: low accuracy = high score
		} else {
			accScore = 0.5 // No data, neutral score
		}

		// Combined score: weight both factors equally
		// Letters that are both rare AND have low accuracy get highest scores
		letterScore := (freqScore + accScore) / 2.0
		score += letterScore
	}

	// Normalize by word length to avoid bias toward longer words
	if len(word) > 0 {
		score = score / float64(len(word))
	}

	// Ensure minimum score to avoid zero probability
	if score < 0.1 {
		score = 0.1
	}

	return score
}

// weightedRandomSelect selects a word from candidates weighted by letter data score
func weightedRandomSelect(candidates []string, letterData LetterData, rng func(int) int) string {
	if len(candidates) == 0 {
		return ""
	}
	if len(candidates) == 1 {
		return candidates[0]
	}

	// Calculate scores for all candidates
	scores := make([]float64, len(candidates))
	var totalScore float64
	for i, word := range candidates {
		scores[i] = scoreWordByLetterData(word, letterData)
		totalScore += scores[i]
	}

	// Weighted random selection
	target := float64(rng(1000000)) / 1000000.0 * totalScore
	var cumulative float64
	for i, score := range scores {
		cumulative += score
		if cumulative >= target {
			return candidates[i]
		}
	}

	// Fallback to last candidate
	return candidates[len(candidates)-1]
}

// GetRandomWordsFixedCount returns exactly numWords random words totalling exactly targetChars characters
// Uses a stratified selection approach to ensure both constraints are met
// Words containing underrepresented and low-accuracy letters are weighted higher
func GetRandomWordsFixedCount(numWords, targetChars int, rng func(int) int, letterData LetterData) []string {
	// Group words by length for efficient selection
	wordsByLength := make(map[int][]string)
	for _, word := range CommonWords {
		word = strings.ToLower(word)
		if len(strings.TrimSpace(word)) > 0 {
			length := len(word)
			wordsByLength[length] = append(wordsByLength[length], word)
		}
	}

	// Try multiple times to find a valid combination
	for attempt := 0; attempt < 100; attempt++ {
		result := make([]string, 0, numWords)
		currentChars := 0

		for i := 0; i < numWords; i++ {
			wordsRemaining := numWords - i
			charsRemaining := targetChars - currentChars

			if wordsRemaining == 0 {
				break
			}

			// Calculate ideal length for this word
			idealLength := float64(charsRemaining) / float64(wordsRemaining)

			// For the last word, we need exact match
			if wordsRemaining == 1 {
				if words, exists := wordsByLength[charsRemaining]; exists && len(words) > 0 {
					word := weightedRandomSelect(words, letterData, rng)
					result = append(result, word)
					currentChars += charsRemaining
				}
				break
			}

			// Find words close to the ideal length
			// Allow some variance to keep it interesting
			minLen := int(idealLength) - 2
			maxLen := int(idealLength) + 2
			if minLen < 1 {
				minLen = 1
			}

			// Ensure we can still reach target with remaining words
			// Min possible: remaining words * 1 char each (but we need words of at least length 1)
			// Max possible: remaining words * max word length
			minPossibleRemaining := wordsRemaining - 1 // Minimum 1 char per remaining word
			maxPossibleRemaining := (wordsRemaining - 1) * 15 // Assuming max word length ~15

			// Adjust bounds to ensure feasibility
			if charsRemaining-maxLen < minPossibleRemaining {
				minLen = charsRemaining - maxPossibleRemaining
				if minLen < 1 {
					minLen = 1
				}
			}
			if charsRemaining-minLen > maxPossibleRemaining {
				maxLen = charsRemaining - minPossibleRemaining
			}

			// Collect valid words within range
			var candidates []string
			for length := minLen; length <= maxLen; length++ {
				if words, exists := wordsByLength[length]; exists {
					candidates = append(candidates, words...)
				}
			}

			if len(candidates) == 0 {
				// Fallback: try any word that keeps us feasible
				for length, words := range wordsByLength {
					remaining := charsRemaining - length
					if remaining >= wordsRemaining-1 && remaining <= (wordsRemaining-1)*15 {
						candidates = append(candidates, words...)
					}
				}
			}

			if len(candidates) == 0 {
				// This attempt failed, try again
				break
			}

			// Use weighted selection based on letter data (frequency + accuracy)
			word := weightedRandomSelect(candidates, letterData, rng)
			result = append(result, word)
			currentChars += len(word)
		}

		// Check if we got exactly what we need
		if len(result) == numWords && currentChars == targetChars {
			return result
		}
	}

	// Fallback: just return random words (shouldn't happen with reasonable params)
	result := make([]string, 0, numWords)
	for len(result) < numWords {
		word := strings.ToLower(CommonWords[rng(len(CommonWords))])
		if len(strings.TrimSpace(word)) > 0 {
			result = append(result, word)
		}
	}
	return result
}
