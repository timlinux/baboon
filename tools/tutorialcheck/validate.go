// tools/tutorialcheck/validate.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents a single validation failure
type ValidationError struct {
	File    string
	Line    int
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
}

var (
	// Matches {{< code-ref id="settings/structs" >}} or {{<code-ref id="settings/structs">}}
	codeRefRe = regexp.MustCompile(`\{\{<\s*code-ref\s+id="([^"]+)"`)

	// Matches {{< file-ref path="settings/settings.go" lines="47-55" >}}
	fileRefRe = regexp.MustCompile(`\{\{<\s*file-ref\s+path="([^"]+)"(?:\s+lines="(\d+)(?:-(\d+))?")?`)

	// Matches ```go compile
	compileBlockRe = regexp.MustCompile("^```go\\s+compile\\s*$")
)

// ValidateMarkdown checks all references in a markdown file
func ValidateMarkdown(path string, snippets map[string]CodeSnippet, verbose bool) []ValidationError {
	var errors []ValidationError

	file, err := os.Open(path)
	if err != nil {
		return []ValidationError{{File: path, Line: 0, Message: err.Error()}}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check code-ref shortcodes
		if matches := codeRefRe.FindStringSubmatch(line); matches != nil {
			id := matches[1]
			if _, ok := snippets[id]; !ok {
				errors = append(errors, ValidationError{
					File:    path,
					Line:    lineNum,
					Message: fmt.Sprintf("code-ref id %q not found in source files", id),
				})
			} else if verbose {
				fmt.Printf("  %s:%d: code-ref %q OK\n", path, lineNum, id)
			}
		}

		// Check file-ref shortcodes
		if matches := fileRefRe.FindStringSubmatch(line); matches != nil {
			filePath := matches[1]
			if err := validateFileRef(filePath, matches[2], matches[3]); err != nil {
				errors = append(errors, ValidationError{
					File:    path,
					Line:    lineNum,
					Message: err.Error(),
				})
			} else if verbose {
				fmt.Printf("  %s:%d: file-ref %q OK\n", path, lineNum, filePath)
			}
		}
	}

	return errors
}

func validateFileRef(path, startLine, endLine string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file-ref path %q: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("file-ref path %q is a directory", path)
	}

	if startLine == "" {
		return nil
	}

	// Count lines in file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	start, _ := strconv.Atoi(startLine)
	if start > lineCount {
		return fmt.Errorf("file-ref line %d exceeds file length %d", start, lineCount)
	}

	if endLine != "" {
		end, _ := strconv.Atoi(endLine)
		if end > lineCount {
			return fmt.Errorf("file-ref line %d exceeds file length %d", end, lineCount)
		}
	}

	return nil
}

// ValidateAllMarkdown validates all markdown files in the tutorial directory
func ValidateAllMarkdown(tutorialDir string, snippets map[string]CodeSnippet, verbose bool) []ValidationError {
	var allErrors []ValidationError

	err := filepath.Walk(tutorialDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		if verbose {
			fmt.Printf("Validating %s\n", path)
		}

		errors := ValidateMarkdown(path, snippets, verbose)
		allErrors = append(allErrors, errors...)

		return nil
	})

	if err != nil {
		allErrors = append(allErrors, ValidationError{
			File:    tutorialDir,
			Line:    0,
			Message: err.Error(),
		})
	}

	return allErrors
}
