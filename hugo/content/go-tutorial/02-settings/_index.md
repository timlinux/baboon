---
title: "Chapter 02: Settings Package"
weight: 2
chapter: "02-settings"
sources:
  - settings/settings.go
---

# The Settings Package

The `settings/` package is one of the simplest in Baboon, making it perfect for learning Go fundamentals. It handles loading and saving user preferences like which key advances to the next word.

## What You'll Learn

- Go's basic types (`int`, `string`, `bool`)
- Defining custom types with `type`
- Creating structs with fields
- JSON struct tags for serialization
- File I/O with `os` package
- Error handling patterns

## Package Overview

{{< file-ref path="settings/settings.go" >}}

This single file contains everything:
- A custom `AdvanceKey` type with constants
- Methods on custom types
- A `Settings` struct
- Functions for loading/saving to disk

## Pages

1. [Basic Types](./basic-types/) - Go's type system and custom types
2. [Structs](./structs/) - Defining and using structs
3. [JSON Tags](./json-tags/) - Serialization with struct tags
4. [File I/O](./file-io/) - Reading and writing files
