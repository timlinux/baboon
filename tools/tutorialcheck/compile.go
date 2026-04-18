// tools/tutorialcheck/compile.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CompileExample represents a standalone Go example to compile
type CompileExample struct {
	File      string
	StartLine int
	EndLine   int
	Code      string
}

// ExtractCompileExamples finds all ```go compile blocks in markdown
func ExtractCompileExamples(path string) ([]CompileExample, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var examples []CompileExample
	scanner := bufio.NewScanner(file)

	lineNum := 0
	inCompileBlock := false
	var currentLines []string
	var startLine int

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if compileBlockRe.MatchString(line) {
			inCompileBlock = true
			startLine = lineNum + 1
			currentLines = nil
			continue
		}

		if inCompileBlock && strings.HasPrefix(line, "```") {
			examples = append(examples, CompileExample{
				File:      path,
				StartLine: startLine,
				EndLine:   lineNum - 1,
				Code:      strings.Join(currentLines, "\n"),
			})
			inCompileBlock = false
			continue
		}

		if inCompileBlock {
			currentLines = append(currentLines, line)
		}
	}

	return examples, scanner.Err()
}

// CompileExamples compiles all examples and returns errors
func CompileExamples(examples []CompileExample, verbose bool) []ValidationError {
	var errors []ValidationError

	tmpDir, err := os.MkdirTemp("", "tutorialcheck-*")
	if err != nil {
		return []ValidationError{{Message: fmt.Sprintf("creating temp dir: %v", err)}}
	}
	defer os.RemoveAll(tmpDir)

	for i, ex := range examples {
		tmpFile := filepath.Join(tmpDir, fmt.Sprintf("example_%d.go", i))

		if err := os.WriteFile(tmpFile, []byte(ex.Code), 0644); err != nil {
			errors = append(errors, ValidationError{
				File:    ex.File,
				Line:    ex.StartLine,
				Message: fmt.Sprintf("writing temp file: %v", err),
			})
			continue
		}

		cmd := exec.Command("go", "build", "-o", "/dev/null", tmpFile)
		output, err := cmd.CombinedOutput()
		if err != nil {
			errors = append(errors, ValidationError{
				File:    ex.File,
				Line:    ex.StartLine,
				Message: fmt.Sprintf("compile error: %s", strings.TrimSpace(string(output))),
			})
		} else if verbose {
			fmt.Printf("  %s:%d: compile OK\n", ex.File, ex.StartLine)
		}
	}

	return errors
}

// ValidateCompileExamples extracts and compiles all examples in a directory
func ValidateCompileExamples(tutorialDir string, verbose bool) []ValidationError {
	var allErrors []ValidationError
	var allExamples []CompileExample

	err := filepath.Walk(tutorialDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		examples, err := ExtractCompileExamples(path)
		if err != nil {
			allErrors = append(allErrors, ValidationError{
				File:    path,
				Line:    0,
				Message: fmt.Sprintf("extracting examples: %v", err),
			})
			return nil
		}

		allExamples = append(allExamples, examples...)
		return nil
	})

	if err != nil {
		allErrors = append(allErrors, ValidationError{
			File:    tutorialDir,
			Line:    0,
			Message: err.Error(),
		})
	}

	if verbose && len(allExamples) > 0 {
		fmt.Printf("Compiling %d standalone examples...\n", len(allExamples))
	}

	compileErrors := CompileExamples(allExamples, verbose)
	allErrors = append(allErrors, compileErrors...)

	return allErrors
}
