---
title: "File I/O"
weight: 4
chapter: "02-settings"
concepts:
  - os package
  - File reading
  - File writing
  - Error handling
  - defer keyword
  - Path manipulation
---

# File I/O

Reading and writing files is essential for any application that needs to persist data. In this page, you'll learn how Baboon saves and loads settings files.

## The os Package

Go's `os` package provides functions for interacting with the operating system, including file operations:

```go compile
package main

import (
    "fmt"
    "os"
)

func main() {
    // Get the current working directory
    dir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Working directory:", dir)

    // Get the user's home directory
    home, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Home directory:", home)
}
```

## Building File Paths

The `path/filepath` package provides cross-platform path manipulation:

```go compile
package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    // Join path elements with the OS-specific separator
    path := filepath.Join("config", "app", "settings.json")
    fmt.Println("Path:", path)
    // Unix: config/app/settings.json
    // Windows: config\app\settings.json

    // Extract directory and filename
    dir := filepath.Dir(path)
    base := filepath.Base(path)
    fmt.Println("Directory:", dir)
    fmt.Println("Filename:", base)
}
```

Here's how Baboon builds the settings path:

{{< code-ref id="settings/path" >}}

This function:
1. Gets the user's home directory
2. Joins it with `.config/baboon` to follow XDG conventions
3. Creates the directory if it doesn't exist (with `os.MkdirAll`)
4. Returns the full path to `settings.json`

{{% notice tip %}}
Always use `filepath.Join` instead of concatenating strings with `/`. It handles platform differences automatically.
{{% /notice %}}

## Creating Directories

`os.MkdirAll` creates a directory and all parent directories if they don't exist:

```go compile
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // Create nested directories
    path := filepath.Join("tmp", "config", "app")
    err := os.MkdirAll(path, 0755)
    if err != nil {
        fmt.Println("Error creating directories:", err)
        return
    }
    fmt.Println("Created directories:", path)

    // Clean up
    os.RemoveAll("tmp")
}
```

The second parameter (`0755`) is the permission mode:
- `7` (owner): read + write + execute
- `5` (group): read + execute
- `5` (others): read + execute

## Reading Files

Go provides several ways to read files. The simplest is `os.ReadFile`:

```go compile
package main

import (
    "fmt"
    "os"
)

func main() {
    // Create a test file
    os.WriteFile("test.txt", []byte("Hello, World!"), 0644)

    // Read the entire file into memory
    data, err := os.ReadFile("test.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }

    fmt.Println("Content:", string(data))

    // Clean up
    os.Remove("test.txt")
}
```

Here's how Baboon loads settings:

{{< code-ref id="settings/load" >}}

This function:
1. Gets the settings file path
2. Reads the entire file with `os.ReadFile`
3. Checks if the error is "file not found" with `os.IsNotExist`
4. Returns default settings if the file doesn't exist
5. Unmarshals the JSON data into a `Settings` struct

{{% notice warning %}}
`os.ReadFile` loads the entire file into memory. For large files (megabytes+), use `os.Open` and read in chunks instead.
{{% /notice %}}

## Checking Error Types

Not all errors are created equal. Sometimes you want to handle specific errors differently:

```go compile
package main

import (
    "fmt"
    "os"
)

func main() {
    data, err := os.ReadFile("nonexistent.txt")
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("File doesn't exist - that's okay!")
        } else {
            fmt.Println("Real error:", err)
        }
        return
    }
    fmt.Println("Data:", string(data))
}
```

Common error checks:
- `os.IsNotExist(err)` - File doesn't exist
- `os.IsPermission(err)` - Permission denied
- `os.IsTimeout(err)` - Operation timed out

## Writing Files

`os.WriteFile` writes data to a file, creating it if needed:

```go compile
package main

import (
    "fmt"
    "os"
)

func main() {
    data := []byte("Hello, File!")

    // Write to file (create or overwrite)
    err := os.WriteFile("output.txt", data, 0644)
    if err != nil {
        fmt.Println("Error writing file:", err)
        return
    }

    fmt.Println("File written successfully")

    // Read it back
    readData, _ := os.ReadFile("output.txt")
    fmt.Println("Content:", string(readData))

    // Clean up
    os.Remove("output.txt")
}
```

The permission `0644` means:
- `6` (owner): read + write
- `4` (group): read
- `4` (others): read

Here's how Baboon saves settings:

{{< code-ref id="settings/save" >}}

## File Operations with Open

For more control, use `os.Open`, `os.Create`, or `os.OpenFile`:

```go compile
package main

import (
    "fmt"
    "os"
)

func main() {
    // Create a file
    file, err := os.Create("example.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()

    // Write to it
    _, err = file.WriteString("Line 1\n")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    file.WriteString("Line 2\n")

    // Sync to disk
    file.Sync()

    fmt.Println("File created successfully")

    // Clean up
    os.Remove("example.txt")
}
```

## The defer Keyword

`defer` schedules a function call to run when the surrounding function returns:

```go compile
package main

import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Create("defer-example.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close() // Will be called when main() returns

    file.WriteString("Hello!")

    fmt.Println("File will be closed automatically")
    // file.Close() happens here
    os.Remove("defer-example.txt")
}
```

`defer` is especially useful for cleanup:
- Closing files
- Releasing locks
- Closing network connections

Multiple defers execute in LIFO (Last In, First Out) order:

```go compile
package main

import "fmt"

func main() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    fmt.Println("Start")
    // Output:
    // Start
    // 3
    // 2
    // 1
}
```

## Atomic Writes

Writing to a file can fail midway, leaving a corrupted file. For critical data, use atomic writes:

```go compile
package main

import (
    "fmt"
    "os"
)

func atomicWrite(filename string, data []byte) error {
    // Write to a temporary file
    tmpFile := filename + ".tmp"
    err := os.WriteFile(tmpFile, data, 0644)
    if err != nil {
        return err
    }

    // Atomically rename (replacing the original)
    return os.Rename(tmpFile, filename)
}

func main() {
    data := []byte("Important data")
    err := atomicWrite("important.txt", data)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Data written atomically")

    // Clean up
    os.Remove("important.txt")
}
```

{{% notice tip %}}
Baboon doesn't use atomic writes for settings because the data isn't critical. But for databases or important user data, always use atomic writes.
{{% /notice %}}

## Checking if Files Exist

There's no direct "file exists" function, but you can check by trying to stat it:

```go compile
package main

import (
    "fmt"
    "os"
)

func fileExists(filename string) bool {
    _, err := os.Stat(filename)
    return err == nil
}

func main() {
    // Create a test file
    os.WriteFile("test.txt", []byte("test"), 0644)

    fmt.Println("test.txt exists:", fileExists("test.txt"))
    fmt.Println("missing.txt exists:", fileExists("missing.txt"))

    // Clean up
    os.Remove("test.txt")
}
```

## Putting It All Together

Let's write a complete example that mimics Baboon's settings behavior:

```go compile
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

type Config struct {
    Theme string `json:"theme"`
    Port  int    `json:"port"`
}

func DefaultConfig() *Config {
    return &Config{
        Theme: "dark",
        Port:  8080,
    }
}

func getConfigPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    configDir := filepath.Join(homeDir, ".myapp")
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return "", err
    }
    return filepath.Join(configDir, "config.json"), nil
}

func Load() (*Config, error) {
    path, err := getConfigPath()
    if err != nil {
        return DefaultConfig(), err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return DefaultConfig(), nil
        }
        return DefaultConfig(), err
    }

    var c Config
    if err := json.Unmarshal(data, &c); err != nil {
        return DefaultConfig(), err
    }

    return &c, nil
}

func (c *Config) Save() error {
    path, err := getConfigPath()
    if err != nil {
        return err
    }

    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0644)
}

func main() {
    // Load config (will use defaults first time)
    config, err := Load()
    if err != nil {
        fmt.Println("Error loading:", err)
        return
    }
    fmt.Printf("Loaded: Theme=%s, Port=%d\n", config.Theme, config.Port)

    // Modify and save
    config.Theme = "light"
    config.Port = 9000
    if err := config.Save(); err != nil {
        fmt.Println("Error saving:", err)
        return
    }
    fmt.Println("Saved successfully")

    // Load again to verify
    config2, _ := Load()
    fmt.Printf("Reloaded: Theme=%s, Port=%d\n", config2.Theme, config2.Port)

    // Clean up
    path, _ := getConfigPath()
    os.Remove(path)
    os.Remove(filepath.Dir(path))
}
```

{{< exercise >}}

Create a simple note-taking app that saves notes to `~/.notes/notes.txt`. Implement:
- `SaveNote(note string)` - Appends a note to the file
- `LoadNotes()` - Returns all notes as a slice of strings (one per line)
- Handle the case where the file doesn't exist yet

{{< solution >}}

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func getNotesPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    notesDir := filepath.Join(homeDir, ".notes")
    if err := os.MkdirAll(notesDir, 0755); err != nil {
        return "", err
    }
    return filepath.Join(notesDir, "notes.txt"), nil
}

func SaveNote(note string) error {
    path, err := getNotesPath()
    if err != nil {
        return err
    }

    // Open file for appending (create if doesn't exist)
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(note + "\n")
    return err
}

func LoadNotes() ([]string, error) {
    path, err := getNotesPath()
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return []string{}, nil // Empty list if file doesn't exist
        }
        return nil, err
    }

    // Split by newlines and filter empty lines
    lines := strings.Split(string(data), "\n")
    notes := []string{}
    for _, line := range lines {
        if line != "" {
            notes = append(notes, line)
        }
    }

    return notes, nil
}

func main() {
    // Save some notes
    SaveNote("Buy milk")
    SaveNote("Call dentist")
    SaveNote("Finish Go tutorial")

    // Load and display
    notes, err := LoadNotes()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Your notes:")
    for i, note := range notes {
        fmt.Printf("%d. %s\n", i+1, note)
    }

    // Clean up
    path, _ := getNotesPath()
    os.Remove(path)
    os.Remove(filepath.Dir(path))
}
```

{{< /solution >}}

{{< /exercise >}}

---

**Congratulations!** You've completed Chapter 02 and learned the fundamentals of Go programming through Baboon's settings package. You now understand:

- Basic and custom types
- Structs and methods
- JSON serialization with struct tags
- File I/O and error handling

**Next:** In the next chapter, we'll explore the word list package and learn about slices, maps, and interfaces.
