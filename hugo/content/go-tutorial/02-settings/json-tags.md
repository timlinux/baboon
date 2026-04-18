---
title: "JSON Tags and Serialization"
weight: 3
chapter: "02-settings"
concepts:
  - Struct tags
  - JSON marshaling
  - JSON unmarshaling
  - encoding/json package
---

# JSON Tags and Serialization

Go makes it easy to convert structs to and from JSON. This is how Baboon saves settings to disk and loads them back.

## What Are Struct Tags?

Look at the `Settings` struct again:

{{< code-ref id="settings/structs" >}}

Those strings in backticks (`` `json:"advance_key"` ``) are **struct tags**. They're metadata attached to fields that other code can read.

The `encoding/json` package reads these tags to know:
- What JSON field name to use
- Whether to skip empty fields
- Other serialization options

## JSON Tag Syntax

The basic syntax is `` `json:"field_name"` ``:

```go compile
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Age       int    `json:"age"`
}

func main() {
    user := User{
        FirstName: "Jane",
        LastName:  "Doe",
        Age:       28,
    }

    data, _ := json.Marshal(user)
    fmt.Println(string(data))
}
```

Without the tags, the JSON would use the exact field names: `{"FirstName":"Jane","LastName":"Doe","Age":28}`. With tags, we get snake_case: `{"first_name":"Jane","last_name":"Doe","age":28}`.

{{% notice tip %}}
Use snake_case in JSON tags to follow common JSON conventions. Go uses PascalCase for exported names, but JSON typically uses snake_case.
{{% /notice %}}

## Tag Options

You can add options after a comma:

```go compile
package main

import (
    "encoding/json"
    "fmt"
)

type Config struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Password string `json:"password,omitempty"`  // Omit if empty
    Debug    bool   `json:"debug,omitempty"`     // Omit if false
    Internal string `json:"-"`                   // Never include
}

func main() {
    config := Config{
        Host:     "localhost",
        Port:     8080,
        Password: "",      // Will be omitted
        Debug:    false,   // Will be omitted
        Internal: "secret",
    }

    data, _ := json.MarshalIndent(config, "", "  ")
    fmt.Println(string(data))
}
```

Common options:
- `omitempty` - Omit field if it's empty (zero value)
- `-` - Never include this field in JSON

## Marshaling: Go to JSON

"Marshaling" means converting a Go value to JSON:

```go compile
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    p := Person{Name: "Alice", Age: 30}

    // Marshal to compact JSON
    data, err := json.Marshal(p)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Compact:", string(data))

    // Marshal to pretty-printed JSON
    data, err = json.MarshalIndent(p, "", "  ")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Pretty:\n" + string(data))
}
```

Here's how Baboon marshals settings:

{{< code-ref id="settings/save" >}}

The `MarshalIndent` function takes:
- The value to marshal
- A prefix string (usually empty)
- An indent string (two spaces for readability)

## Unmarshaling: JSON to Go

"Unmarshaling" means converting JSON to a Go value:

```go compile
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    jsonStr := `{"name":"Bob","age":25}`

    var p Person
    err := json.Unmarshal([]byte(jsonStr), &p)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
}
```

Key points:
- Pass a byte slice (`[]byte`) to `Unmarshal`
- Pass a pointer to the struct (`&p`) so it can be modified
- The JSON field names must match the struct tags

Here's how Baboon unmarshals settings:

{{< code-ref id="settings/load" >}}

## Handling Missing Fields

If a JSON field is missing, the struct field gets its zero value:

```go compile
package main

import (
    "encoding/json"
    "fmt"
)

type Config struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

func main() {
    // JSON is missing the "port" field
    jsonStr := `{"host":"localhost"}`

    var config Config
    json.Unmarshal([]byte(jsonStr), &config)

    fmt.Printf("Host: %s, Port: %d\n", config.Host, config.Port)
    // Output: Host: localhost, Port: 0
}
```

This is why Baboon's `Load` function returns `DefaultSettings()` if the file doesn't exist - it ensures all fields have sensible values.

## Handling Extra Fields

If JSON has extra fields that aren't in the struct, they're ignored:

```go compile
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    Name string `json:"name"`
}

func main() {
    // JSON has extra fields "age" and "email"
    jsonStr := `{"name":"Alice","age":30,"email":"alice@example.com"}`

    var user User
    json.Unmarshal([]byte(jsonStr), &user)

    fmt.Printf("Name: %s\n", user.Name)
    // Extra fields are silently ignored
}
```

This makes your code forward-compatible: if you add new fields later, old code won't break.

## Custom Types and JSON

When you have custom types like `AdvanceKey`, the JSON package marshals them as their underlying type:

```go compile
package main

import (
    "encoding/json"
    "fmt"
)

type AdvanceKey int

const (
    AdvanceKeySpace  AdvanceKey = iota
    AdvanceKeyEnter
    AdvanceKeyEither
)

type Settings struct {
    AdvanceKey AdvanceKey `json:"advance_key"`
}

func main() {
    settings := Settings{AdvanceKey: AdvanceKeyEnter}

    data, _ := json.MarshalIndent(settings, "", "  ")
    fmt.Println(string(data))
    // Output: {"advance_key": 1}
}
```

The `AdvanceKey` (which has value 1) is marshaled as the number `1`. When unmarshaling, the number is converted back to an `AdvanceKey`.

{{% notice warning %}}
If you change the order of constants, old JSON files will map to different values. This is why you should never rely on specific iota values for persistence.
{{% /notice %}}

## Error Handling

Always check errors when marshaling/unmarshaling:

```go compile
package main

import (
    "encoding/json"
    "fmt"
)

type Config struct {
    Timeout int `json:"timeout"`
}

func main() {
    // Invalid JSON (missing closing brace)
    invalidJSON := `{"timeout": 30`

    var config Config
    err := json.Unmarshal([]byte(invalidJSON), &config)
    if err != nil {
        fmt.Println("Failed to parse JSON:", err)
        return
    }

    fmt.Println("Config:", config)
}
```

{{< exercise >}}

Create a `GameStats` struct with fields for PlayerName, Score, Level, and HighScore. Add appropriate JSON tags. Write a function that marshals a GameStats to pretty JSON, then unmarshals it back and prints the values. Make HighScore optional (omit if zero).

{{< solution >}}

```go
package main

import (
    "encoding/json"
    "fmt"
)

type GameStats struct {
    PlayerName string `json:"player_name"`
    Score      int    `json:"score"`
    Level      int    `json:"level"`
    HighScore  int    `json:"high_score,omitempty"`
}

func main() {
    // Create some stats
    stats := GameStats{
        PlayerName: "Alice",
        Score:      1500,
        Level:      5,
        HighScore:  0, // Will be omitted
    }

    // Marshal to JSON
    data, err := json.MarshalIndent(stats, "", "  ")
    if err != nil {
        fmt.Println("Marshal error:", err)
        return
    }
    fmt.Println("JSON:")
    fmt.Println(string(data))

    // Unmarshal back
    var stats2 GameStats
    err = json.Unmarshal(data, &stats2)
    if err != nil {
        fmt.Println("Unmarshal error:", err)
        return
    }

    fmt.Printf("\nUnmarshaled: %s, Score=%d, Level=%d, HighScore=%d\n",
        stats2.PlayerName, stats2.Score, stats2.Level, stats2.HighScore)
}
```

{{< /solution >}}

{{< /exercise >}}

---

**Next:** [File I/O](../file-io/) - Learn how to read and write files
