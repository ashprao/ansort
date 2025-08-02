# Alphanumeric Sorting Package

A Go package that provides natural sorting capabilities for alphanumeric strings, where numeric parts are sorted numerically rather than lexicographically.

## Installation

```bash
go get github.com/ashprao/ansort
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/ashprao/ansort"
)

func main() {
    // Basic natural sorting
    data := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
    ansort.SortStrings(data)
    fmt.Println(data)
    // Output: [file1.txt file2.txt file10.txt file20.txt]
    
    // Works with complex patterns too
    versions := []string{"v1.10.1", "v1.2.10", "v1.2.2"}
    ansort.SortStrings(versions)
    fmt.Println(versions)
    // Output: [v1.2.2 v1.2.10 v1.10.1]
}
```

## Problem Solved

Standard string sorting treats numbers as individual characters, leading to unintuitive results:

- **Standard sort**: `["item1", "item10", "item2", "item20", "item3"]`
- **Natural sort**: `["item1", "item2", "item3", "item10", "item20"]`

This package implements natural (human-friendly) sorting that handles numeric parts intelligently, including:
- **Multi-segment numbers**: `v1.2.10` < `v1.10.1` (semantic versioning)
- **Decimal numbers**: `file3.14.txt` < `file10.1.txt`
- **Mixed patterns**: IP addresses, version strings, file names with complex numbering

## Project Structure

```
ansort/
├── README.md                   # Project documentation and usage guide
├── LICENSE                     # MIT License
├── go.mod                      # Go module definition
├── ansort.go                   # Main package implementation
├── ansort_test.go              # Unit tests for the package
└── examples/                   # Usage examples and demos
    └── basic/
        └── main.go             # Basic usage example
```

### Key Artifacts

#### Core Files
- **`ansort.go`** - The main package implementation containing all sorting logic, types, and functions
- **`ansort_test.go`** - Comprehensive test suite ensuring code quality and correctness
- **`go.mod`** - Go module file defining the module path and Go version requirements

#### Examples
- **`examples/`** - Directory containing practical usage examples
  - **`examples/basic/main.go`** - Demonstrates basic package usage and functionality

This structure follows Go package conventions and provides clear separation of concerns between implementation, testing, documentation, and examples.

## Usage Guide

The package provides two main approaches for sorting: **convenience functions** for simple use cases and **sorter objects** for integration with Go's standard library.

### When to Use Which Approach

**Use convenience functions** (`SortStrings`, `Compare`) when:
- You want to sort data directly with minimal setup
- You're doing one-time sorting operations
- You prefer a simple, clean API for basic use cases

**Use sorter objects** (`NewSorter` + `sort.Sort`) when:
- Integrating with existing code that expects `sort.Interface`
- You need to use Go's specialized sorting functions (`sort.Stable`, `sort.IsSorted`)
- Building frameworks or libraries that work with different sorting implementations
- Performance-critical code where you want to reuse the sorter object
- You need more control over the sorting process

### Convenience Functions (Recommended for Most Cases)

```go
import "github.com/ashprao/ansort"

// Simple natural sorting
data := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
ansort.SortStrings(data)
// data is now: ["file1.txt", "file2.txt", "file10.txt", "file20.txt"]

// Case-insensitive sorting
data2 := []string{"File2.txt", "file10.txt", "FILE1.txt"}
ansort.SortStrings(data2, ansort.WithCaseInsensitive())
// Result: ["FILE1.txt", "File2.txt", "file10.txt"]

// String comparison
result := ansort.Compare("file1.txt", "file10.txt")
// result is -1 (file1.txt comes before file10.txt)

// Multi-segment numbers (semantic versioning)
versions := []string{"v1.10.1", "v1.2.10", "v1.2.2"}
ansort.SortStrings(versions)
// Result: ["v1.2.2", "v1.2.10", "v1.10.1"]

// Decimal numbers
prices := []string{"price$19.99", "price$5.50", "price$100.00"}
ansort.SortStrings(prices)
// Result: ["price$5.50", "price$19.99", "price$100.00"]
```

### Integration with Go's sort Package

Use this approach when integrating with existing code that expects `sort.Interface`:

```go
import (
    "sort"
    "github.com/ashprao/ansort"
)

data := []string{"item10", "item2", "item1"}

// Create a sorter and use with sort.Sort
sorter := ansort.NewSorter(data, ansort.WithCaseInsensitive())
sort.Sort(sorter)
// data is now sorted case-insensitively
```

#### Common Integration Scenarios

**1. Working with Existing Infrastructure**
```go
// Existing function that works with any sort.Interface
func TimedSort(data sort.Interface) time.Duration {
    start := time.Now()
    sort.Sort(data)
    return time.Since(start)
}

// Your code can plug into existing infrastructure
files := []string{"doc10.txt", "doc2.txt", "doc1.txt"}
sorter := ansort.NewSorter(files)
duration := TimedSort(sorter) // Works seamlessly
```

**2. Using Different Sort Algorithms**
```go
data := []string{"file10.txt", "file2.txt", "file1.txt"}
sorter := ansort.NewSorter(data)

// Use stable sort instead of regular sort
sort.Stable(sorter)

// Check if data is already sorted
if sort.IsSorted(sorter) {
    fmt.Println("Data is already naturally sorted!")
}
```

**3. Performance-Critical Reusable Sorting**
```go
// Reuse sorter for multiple sort operations
data := []string{"file10.txt", "file2.txt", "file1.txt"}
sorter := ansort.NewSorter(data, ansort.WithCaseInsensitive())

for i := 0; i < 1000; i++ {
    modifyData(data) // Data gets modified between sorts
    sort.Sort(sorter) // Reuse the same sorter object
}
```

> **💡 Tip**: Run `go run examples/basic/main.go` to see natural sorting in action with detailed output examples.

## API Reference

### Convenience Functions

- `SortStrings(data []string, options ...Option)` - Sorts a slice of strings in-place using natural ordering
- `Compare(a, b string, options ...Option) int` - Compares two strings using natural ordering rules

### Functional Options

- `WithCaseInsensitive()` - Makes sorting/comparison case-insensitive
- `WithCaseSensitive(sensitive bool)` - Explicitly sets case sensitivity (true = sensitive, false = insensitive)

### Standard Library Integration

- `NewSorter(data []string, options ...Option) *AlphanumericSorter` - Creates a sorter implementing `sort.Interface`

### Types

- `AlphanumericSorter` - Implements `sort.Interface` for integration with Go's sort package
- `Config` - Internal configuration structure (used by functional options)
- `Option` - Function type for configuring sorting behavior

## Current Status

This package is currently at **v0.1.0** - a stable MVP with comprehensive natural sorting features.

### ✅ Implemented Features (v0.1.0)
- **Core natural sorting**: Basic alphanumeric parsing and comparison
- **Case sensitivity options**: Case-sensitive/insensitive modes with functional options
- **Multi-segment number support**: Semantic versioning (`v1.2.10`), IP addresses, complex patterns
- **Decimal number support**: Real decimals (`3.14`), prices (`$19.99`), mixed patterns
- **Leading zero handling**: Proper numeric comparison with preserved formatting
- **Standard library integration**: Full `sort.Interface` implementation
- **Comprehensive error handling**: Input validation and graceful error handling
- **High test coverage**: 90%+ coverage with comprehensive test suite

### 🔄 Future Versions
- **v0.2.0**: Unicode and special character support (Phase 3.3)
- **v0.3.0**: Performance optimizations and API polish (Phase 4)
- **v1.0.0**: Stable API guarantee after community feedback

### 📦 Version Compatibility
- **Pre-v1.0.0**: API may evolve based on feedback
- **Post-v1.0.0**: Semantic versioning with backward compatibility guarantees

### Key Discovery
The natural tokenization approach handles complex patterns automatically without special cases:
- Multi-segment numbers work through token-by-token comparison
- Decimal numbers sort correctly via natural numeric token handling
- Version strings, IP addresses, and mixed patterns all work seamlessly

## Contributing

This project follows a phased development approach. Please see `implementation_plan.md` for detailed development phases and contribution guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Examples

See the `examples/` directory for comprehensive usage examples:
- `examples/basic/` - Basic natural sorting functionality including multi-segment and decimal number examples

Run the example:
```bash
go run examples/basic/main.go
```
