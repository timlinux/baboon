// tools/tutorialcheck/extract.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// CodeSnippet represents extracted code from a TUTORIAL marker
type CodeSnippet struct {
	ID        string
	FilePath  string
	StartLine int
	EndLine   int
	Code      string
}

var (
	markerStartRe = regexp.MustCompile(`//\s*TUTORIAL:([^:]+):start`)
	markerEndRe   = regexp.MustCompile(`//\s*TUTORIAL:([^:]+):end`)
)

// ExtractSnippets finds all TUTORIAL markers in Go files under the given root
func ExtractSnippets(root string, verbose bool) (map[string]CodeSnippet, error) {
	snippets := make(map[string]CodeSnippet)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and vendor/tools directories
		if info.IsDir() {
			if info.Name() == "vendor" || info.Name() == "tools" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fileSnippets, err := extractFromFile(path, verbose)
		if err != nil {
			return fmt.Errorf("extracting from %s: %w", path, err)
		}

		for id, snippet := range fileSnippets {
			if existing, ok := snippets[id]; ok {
				return fmt.Errorf("duplicate marker ID %q in %s and %s", id, existing.FilePath, path)
			}
			snippets[id] = snippet
		}

		return nil
	})

	return snippets, err
}

func extractFromFile(path string, verbose bool) (map[string]CodeSnippet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	snippets := make(map[string]CodeSnippet)
	scanner := bufio.NewScanner(file)

	var currentID string
	var currentLines []string
	var startLine int
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if matches := markerStartRe.FindStringSubmatch(line); matches != nil {
			if currentID != "" {
				return nil, fmt.Errorf("line %d: nested marker %q inside %q", lineNum, matches[1], currentID)
			}
			currentID = matches[1]
			startLine = lineNum + 1
			currentLines = nil
			if verbose {
				fmt.Printf("  Found marker start: %s at line %d\n", currentID, lineNum)
			}
			continue
		}

		if matches := markerEndRe.FindStringSubmatch(line); matches != nil {
			if currentID == "" {
				return nil, fmt.Errorf("line %d: end marker %q without start", lineNum, matches[1])
			}
			if matches[1] != currentID {
				return nil, fmt.Errorf("line %d: end marker %q doesn't match start %q", lineNum, matches[1], currentID)
			}

			snippets[currentID] = CodeSnippet{
				ID:        currentID,
				FilePath:  path,
				StartLine: startLine,
				EndLine:   lineNum - 1,
				Code:      strings.Join(currentLines, "\n"),
			}

			if verbose {
				fmt.Printf("  Found marker end: %s at line %d (%d lines)\n", currentID, lineNum, len(currentLines))
			}

			currentID = ""
			currentLines = nil
			continue
		}

		if currentID != "" {
			currentLines = append(currentLines, line)
		}
	}

	if currentID != "" {
		return nil, fmt.Errorf("unclosed marker %q starting at line %d", currentID, startLine-1)
	}

	return snippets, scanner.Err()
}
