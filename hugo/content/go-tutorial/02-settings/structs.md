---
title: "Structs"
weight: 2
chapter: "02-settings"
concepts:
  - Struct definition
  - Struct fields
  - Exported vs unexported names
  - Constructor functions
  - Pointer receivers
---

# Structs

Structs are Go's way of grouping related data together. They're similar to classes in other languages, but simpler - they're just data containers.

## Defining a Struct

Here's how Baboon defines the `Settings` struct:

{{< code-ref id="settings/structs" >}}

A struct definition has:
- The `type` keyword
- The struct name (`Settings`)
- The `struct` keyword
- A list of fields in curly braces

Each field has:
- A name (`AdvanceKey`, `PerfectMode`, `PracticeMode`)
- A type (`AdvanceKey`, `bool`, `PracticeMode`)
- Optional struct tags (those backtick strings - we'll cover these in the next page)

## Creating Struct Instances

There are several ways to create a struct:

```go compile
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    // Method 1: Specify all fields by name
    p1 := Person{Name: "Alice", Age: 30}

    // Method 2: Use positional values (not recommended)
    p2 := Person{"Bob", 25}

    // Method 3: Create with zero values, then set fields
    var p3 Person
    p3.Name = "Charlie"
    p3.Age = 35

    fmt.Println(p1)
    fmt.Println(p2)
    fmt.Println(p3)
}
```

{{% notice warning %}}
Using positional values (Method 2) is fragile. If someone reorders the struct fields, your code breaks. Always use field names.
{{% /notice %}}

## Zero Values

If you don't initialize a field, it gets a "zero value":
- Numbers: `0`
- Booleans: `false`
- Strings: `""`
- Pointers: `nil`
- Structs: all fields are zero-valued

```go compile
package main

import "fmt"

type Config struct {
    Timeout int
    Debug   bool
    Name    string
}

func main() {
    var c Config
    fmt.Printf("Timeout: %d, Debug: %t, Name: '%s'\n",
        c.Timeout, c.Debug, c.Name)
}
```

## Exported vs Unexported Names

In Go, names that start with a capital letter are **exported** (public), while names starting with a lowercase letter are **unexported** (private to the package):

```go compile
package main

import "fmt"

type User struct {
    Name     string // Exported - visible outside the package
    Email    string // Exported
    password string // Unexported - only visible within the package
}

func main() {
    u := User{
        Name:     "Alice",
        Email:    "alice@example.com",
        password: "secret123",
    }
    fmt.Printf("User: %s (%s)\n", u.Name, u.Email)
    // Outside this package, you couldn't access u.password
}
```

All fields in Baboon's `Settings` struct are exported because we need to serialize them to JSON (more on that in the next page).

## Constructor Functions

Go doesn't have constructors like other languages, but by convention we create functions that return initialized structs. Here's how Baboon does it:

{{< code-ref id="settings/defaults" >}}

This function:
- Has a descriptive name: `DefaultSettings`
- Returns a pointer to `Settings` (`*Settings`)
- Initializes all fields with sensible defaults

Why return a pointer? Because in Go:
- Passing structs by value copies all the data
- Passing pointers is more efficient for larger structs
- Pointers let you modify the original struct

```go compile
package main

import "fmt"

type Counter struct {
    Value int
}

// Constructor returns a pointer
func NewCounter() *Counter {
    return &Counter{Value: 0}
}

// Method with pointer receiver can modify the struct
func (c *Counter) Increment() {
    c.Value++
}

func main() {
    counter := NewCounter()
    counter.Increment()
    counter.Increment()
    fmt.Println("Count:", counter.Value)
}
```

{{% notice tip %}}
The `&` operator takes the address of a value, creating a pointer. The `*` in the type means "pointer to".
{{% /notice %}}

## Accessing Struct Fields

Use the dot notation to access fields:

```go compile
package main

import "fmt"

type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    fmt.Println("Width:", rect.Width)
    fmt.Println("Height:", rect.Height)
    fmt.Println("Area:", rect.Area())
}
```

Go automatically dereferences pointers, so you can use the dot notation on pointers too:

```go compile
package main

import "fmt"

type Point struct {
    X, Y int
}

func main() {
    p := &Point{X: 10, Y: 20}
    // No need to write (*p).X - Go does it automatically
    fmt.Println("X:", p.X)
    fmt.Println("Y:", p.Y)
}
```

## Methods vs Functions

Methods are functions attached to a type. They have a receiver (the thing before the method name):

```go compile
package main

import "fmt"

type Temperature float64

// Method - has a receiver
func (t Temperature) ToFahrenheit() float64 {
    return float64(t)*9/5 + 32
}

// Regular function - no receiver
func CelsiusToFahrenheit(celsius float64) float64 {
    return celsius*9/5 + 32
}

func main() {
    temp := Temperature(25.0)

    // Method call
    fmt.Printf("%.1f°C = %.1f°F\n", temp, temp.ToFahrenheit())

    // Function call
    fmt.Printf("%.1f°C = %.1f°F\n", 25.0, CelsiusToFahrenheit(25.0))
}
```

## Value vs Pointer Receivers

Methods can have either value receivers or pointer receivers:

```go compile
package main

import "fmt"

type Counter struct {
    Count int
}

// Value receiver - gets a copy
func (c Counter) GetCount() int {
    return c.Count
}

// Pointer receiver - gets the original
func (c *Counter) Increment() {
    c.Count++
}

func main() {
    counter := Counter{Count: 0}

    counter.Increment()  // Modifies the original
    counter.Increment()

    fmt.Println("Count:", counter.GetCount())
}
```

**Rule of thumb:**
- Use pointer receivers when you need to modify the struct
- Use pointer receivers for large structs (to avoid copying)
- Use value receivers for small, immutable structs
- Be consistent: if some methods need pointer receivers, use them for all methods

{{< exercise >}}

Create a `Book` struct with fields for Title, Author, Pages, and Read (boolean). Create a constructor `NewBook` that takes title, author, and pages, and returns a pointer to a Book with Read defaulting to false. Add a method `MarkAsRead()` that sets Read to true.

{{< solution >}}

```go
package main

import "fmt"

type Book struct {
    Title  string
    Author string
    Pages  int
    Read   bool
}

func NewBook(title, author string, pages int) *Book {
    return &Book{
        Title:  title,
        Author: author,
        Pages:  pages,
        Read:   false,
    }
}

func (b *Book) MarkAsRead() {
    b.Read = true
}

func (b Book) String() string {
    status := "unread"
    if b.Read {
        status = "read"
    }
    return fmt.Sprintf("%s by %s (%d pages, %s)",
        b.Title, b.Author, b.Pages, status)
}

func main() {
    book := NewBook("The Go Programming Language", "Donovan & Kernighan", 380)
    fmt.Println(book)

    book.MarkAsRead()
    fmt.Println(book)
}
```

{{< /solution >}}

{{< /exercise >}}

---

**Next:** [JSON Tags](../json-tags/) - Learn how to serialize structs to JSON
