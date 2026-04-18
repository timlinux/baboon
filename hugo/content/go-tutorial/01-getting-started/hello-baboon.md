---
title: "Hello Baboon"
weight: 1
chapter: "01-getting-started"
concepts:
  - go-modules
  - project-structure
  - go-run
---

# Hello Baboon

Let's explore the Baboon codebase and understand how a Go project is organized.

## Installing Go

If you haven't already, install Go from [go.dev/dl](https://go.dev/dl/). Verify your installation:

```bash
go version
```

You should see something like `go version go1.21.0 linux/amd64`.

## Cloning Baboon

```bash
git clone https://github.com/timlinux/baboon.git
cd baboon
```

## Project Structure

Baboon follows standard Go project conventions:

```
baboon/
├── main.go           # Application entry point
├── go.mod            # Module definition
├── go.sum            # Dependency checksums
├── backend/          # Game engine and REST API
├── frontend/         # Terminal UI (Bubble Tea)
├── settings/         # User preferences
├── stats/            # Statistics tracking
├── words/            # Word dictionary
├── auth/             # Authentication
├── database/         # Database layer
└── web/              # React web frontend
```

## The go.mod File

Every Go project has a `go.mod` file that defines the module path and dependencies:

{{< file-ref path="go.mod" >}}

The module path `github.com/timlinux/baboon` is how other packages import this code.

## Running Baboon

Let's run the application:

```bash
go run .
```

This compiles and runs `main.go`. You should see the Baboon typing practice interface!

Press `Esc` to exit.

## Building Baboon

To create a standalone binary:

```bash
go build -o baboon .
./baboon
```

The `-o baboon` flag names the output binary.

## Understanding main.go

Every Go program starts execution in the `main` function of the `main` package. Let's look at the entry point:

{{< file-ref path="main.go" lines="1-20" >}}

{{% notice tip %}}
Go uses `package main` to indicate an executable program. Library packages use other names like `package settings`.
{{% /notice %}}

## Key Takeaways

1. **Modules** - Go uses `go.mod` to manage dependencies
2. **Packages** - Code is organized into packages (directories)
3. **main package** - Executables must have a `main` package with a `main()` function
4. **go run** - Compiles and runs in one step
5. **go build** - Creates a standalone binary

## Next Steps

In [Chapter 02: Settings](../02-settings/), we'll dive into the `settings/` package to learn about Go's type system, structs, and file I/O.

{{< exercise title="Explore the codebase" >}}
Use `go doc` to explore Baboon's packages:

```bash
go doc ./settings
go doc ./backend
```

What exported types and functions do you see?

{{< solution >}}
The `settings` package exports:
- `Settings` struct
- `AdvanceKey` type
- `Load()` and `Save()` functions
- Constants like `AdvanceKeySpace`

The `backend` package exports:
- `GameAPI` interface
- `Engine` struct
- Various result types
{{< /solution >}}
{{< /exercise >}}
