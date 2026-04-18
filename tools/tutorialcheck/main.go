// tools/tutorialcheck/main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	flagAll     = flag.Bool("all", false, "Validate all tutorial content")
	flagVerbose = flag.Bool("verbose", false, "Verbose output")
	flagExtract = flag.Bool("extract-only", false, "Only extract code snippets, don't validate")
)

func main() {
	flag.Parse()

	if *flagAll {
		fmt.Println("Validating all tutorial content...")
		if err := validateAll(*flagVerbose); err != nil {
			fmt.Fprintf(os.Stderr, "Validation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All validations passed!")
		return
	}

	if *flagExtract {
		fmt.Println("Extracting code snippets...")
		if err := extractAll(*flagVerbose); err != nil {
			fmt.Fprintf(os.Stderr, "Extraction failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Validate specific files passed as arguments
	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Usage: tutorialcheck [--all] [--verbose] [--extract-only] [files...]")
		os.Exit(0)
	}

	if err := validateFiles(files, *flagVerbose); err != nil {
		fmt.Fprintf(os.Stderr, "Validation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Validation passed!")
}

func validateAll(verbose bool) error {
	// Extract all snippets from source
	snippets, err := ExtractSnippets(".", verbose)
	if err != nil {
		return fmt.Errorf("extracting snippets: %w", err)
	}

	if verbose {
		fmt.Printf("Found %d code snippets\n", len(snippets))
	}

	// Validate markdown references
	tutorialDir := "hugo/content/go-tutorial"
	if _, err := os.Stat(tutorialDir); os.IsNotExist(err) {
		fmt.Println("Tutorial directory not found, skipping markdown validation")
		return nil
	}

	errors := ValidateAllMarkdown(tutorialDir, snippets, verbose)
	if len(errors) > 0 {
		fmt.Printf("\nFound %d validation errors:\n", len(errors))
		for _, e := range errors {
			fmt.Printf("  %s\n", e.Error())
		}
		return fmt.Errorf("%d validation errors", len(errors))
	}

	return nil
}

func extractAll(verbose bool) error {
	snippets, err := ExtractSnippets(".", verbose)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d code snippets:\n", len(snippets))
	for id, snippet := range snippets {
		fmt.Printf("  %s: %s:%d-%d (%d lines)\n",
			id, snippet.FilePath, snippet.StartLine, snippet.EndLine,
			snippet.EndLine-snippet.StartLine+1)
	}

	return nil
}

func validateFiles(files []string, verbose bool) error {
	snippets, err := ExtractSnippets(".", verbose)
	if err != nil {
		return fmt.Errorf("extracting snippets: %w", err)
	}

	var allErrors []ValidationError
	for _, file := range files {
		if strings.HasSuffix(file, ".md") {
			errors := ValidateMarkdown(file, snippets, verbose)
			allErrors = append(allErrors, errors...)
		}
	}

	if len(allErrors) > 0 {
		for _, e := range allErrors {
			fmt.Printf("  %s\n", e.Error())
		}
		return fmt.Errorf("%d validation errors", len(allErrors))
	}

	return nil
}
