---
title: "Basic Types and Custom Types"
weight: 1
chapter: "02-settings"
concepts:
  - Basic types
  - Custom types
  - Constants with iota
  - Methods on types
---

# Basic Types and Custom Types

Go has a simple but powerful type system. In this page, you'll learn about Go's built-in types and how to create your own custom types with special behaviors.

## Go's Basic Types

Go provides several built-in types:

```go compile
package main

import "fmt"

func main() {
    // Numeric types
    var age int = 25
    var price float64 = 19.99

    // String type
    var name string = "Baboon"

    // Boolean type
    var isActive bool = true

    fmt.Printf("Age: %d, Price: %.2f, Name: %s, Active: %t\n",
        age, price, name, isActive)
}
```

{{% notice tip %}}
Go has multiple numeric types like `int8`, `int16`, `int32`, `int64`, `uint`, `float32`, etc. For most cases, `int` and `float64` are good defaults.
{{% /notice %}}

## Creating Custom Types

You can create your own types based on existing ones using the `type` keyword. This is useful when you want to add meaning to a simple type:

{{< code-ref id="settings/advancekey-type" >}}

Here, `AdvanceKey` is a custom type based on `int`. This gives us several advantages:

1. **Type Safety**: You can't accidentally pass any `int` where an `AdvanceKey` is expected
2. **Documentation**: The type name tells us what the value represents
3. **Methods**: We can attach methods to our custom type (more on this below)

## Constants with iota

The `iota` keyword is Go's enumerator. It starts at 0 and increments by 1 for each constant in a `const` block:

```go compile
package main

import "fmt"

type Status int

const (
    StatusPending  Status = iota // 0
    StatusActive                 // 1
    StatusComplete               // 2
)

func main() {
    fmt.Println("Pending:", StatusPending)    // 0
    fmt.Println("Active:", StatusActive)      // 1
    fmt.Println("Complete:", StatusComplete)  // 2
}
```

In Baboon's settings, we use this pattern to define the three advance key options:

- `AdvanceKeySpace` (0)
- `AdvanceKeyEnter` (1)
- `AdvanceKeyEither` (2)

{{% notice warning %}}
The actual numeric values don't matter for our code logic. What matters is that each constant has a unique value. Never rely on the specific numbers.
{{% /notice %}}

## Methods on Custom Types

One of the most powerful features of custom types is the ability to add methods to them. Here's how Baboon defines methods on `AdvanceKey`:

{{< code-ref id="settings/advancekey-methods" >}}

The syntax `func (a AdvanceKey) String() string` means:
- `(a AdvanceKey)` is the **receiver** - like `self` in Python or `this` in JavaScript
- `String()` is the method name
- `string` is the return type

Let's see these methods in action:

```go compile
package main

import "fmt"

type AdvanceKey int

const (
    AdvanceKeySpace  AdvanceKey = iota
    AdvanceKeyEnter
    AdvanceKeyEither
)

func (a AdvanceKey) String() string {
    switch a {
    case AdvanceKeySpace:
        return "Space"
    case AdvanceKeyEnter:
        return "Enter"
    case AdvanceKeyEither:
        return "Either"
    default:
        return "Space"
    }
}

func (a AdvanceKey) KeyHint() string {
    switch a {
    case AdvanceKeySpace:
        return "SPACE"
    case AdvanceKeyEnter:
        return "ENTER"
    case AdvanceKeyEither:
        return "SPACE or ENTER"
    default:
        return "SPACE"
    }
}

func main() {
    key := AdvanceKeyEnter
    fmt.Println("Display name:", key.String())
    fmt.Println("Hint:", key.KeyHint())
}
```

## Type Conversions

Sometimes you need to convert between types. Go requires explicit conversions:

```go compile
package main

import "fmt"

type AdvanceKey int

const (
    AdvanceKeySpace  AdvanceKey = iota
    AdvanceKeyEnter
)

func main() {
    // Convert int to AdvanceKey
    var num int = 1
    key := AdvanceKey(num)
    fmt.Println("Key:", key)  // 1

    // Convert AdvanceKey to int
    var k AdvanceKey = AdvanceKeyEnter
    value := int(k)
    fmt.Println("Value:", value)  // 1
}
```

{{% notice tip %}}
Go doesn't do implicit type conversions. This prevents bugs where you accidentally mix incompatible types.
{{% /notice %}}

## Why Use Custom Types?

You might wonder: why not just use `int` everywhere? Here's why custom types are better:

1. **Self-Documenting Code**: `func SetKey(key AdvanceKey)` is clearer than `func SetKey(key int)`
2. **Type Safety**: The compiler prevents you from passing the wrong kind of value
3. **Methods**: You can attach behavior specific to that type
4. **Future-Proofing**: If you need to change the underlying type later, you only change it in one place

{{< exercise >}}

Create a custom type called `Difficulty` with three levels: Easy, Medium, Hard. Add a `String()` method that returns the name, and a `MaxErrors()` method that returns how many errors are allowed (10 for Easy, 5 for Medium, 2 for Hard).

{{< solution >}}

```go
package main

import "fmt"

type Difficulty int

const (
    DifficultyEasy   Difficulty = iota
    DifficultyMedium
    DifficultyHard
)

func (d Difficulty) String() string {
    switch d {
    case DifficultyEasy:
        return "Easy"
    case DifficultyMedium:
        return "Medium"
    case DifficultyHard:
        return "Hard"
    default:
        return "Easy"
    }
}

func (d Difficulty) MaxErrors() int {
    switch d {
    case DifficultyEasy:
        return 10
    case DifficultyMedium:
        return 5
    case DifficultyHard:
        return 2
    default:
        return 10
    }
}

func main() {
    level := DifficultyMedium
    fmt.Printf("Playing on %s difficulty (max %d errors)\n",
        level.String(), level.MaxErrors())
}
```

{{< /solution >}}

{{< /exercise >}}

---

**Next:** [Structs](../structs/) - Learn how to group related data together
