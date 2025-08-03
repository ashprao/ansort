# Alphanumeric Sorting Package

A Go package that provides natural sorting capabilities for alphanumeric strings, where numeric parts are sorted numerically rather than lexicographically. The package supports two complementary use cases:

1. **Direct natural sorting** - Sort data in-memory using intelligent alphanumeric comparison
2. **External system integration** - Generate lexicographically sortable keys that maintain natural order when sorted by external systems (databases, search engines)

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
    // Direct natural sorting (in-memory)
    data := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
    ansort.SortStrings(data)
    fmt.Println(data)
    // Output: [file1.txt file2.txt file10.txt file20.txt]
    
    // External system integration (generate sortable keys)
    sortKey := ansort.ToNaturalSortKey("file10.txt")
    fmt.Println(sortKey)
    // Output: file0000000010.txt (lexicographically sortable)
    
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

This package solves natural sorting challenges in two ways:

### 1. Direct Natural Sorting (In-Memory)
For applications that can control the sorting process directly, the package provides intelligent comparison functions that handle numeric parts correctly, including:
- **Multi-segment numbers**: `v1.2.10` < `v1.10.1` (semantic versioning)
- **Decimal numbers**: `file3.14.txt` < `file10.1.txt`
- **Mixed patterns**: IP addresses, version strings, file names with complex numbering

### 2. External System Integration (Pre-sorted Keys)
For systems that must rely on lexicographic sorting (databases, search engines like Elasticsearch), the package generates specially formatted keys that maintain natural order when sorted lexicographically:
- **Input**: `["item1", "item10", "item2"]`
- **Generated keys**: `["item0000000001", "item0000000010", "item0000000002"]`
- **External lexicographic sort result**: Natural order maintained!

## Project Structure

```
ansort/
├── README.md                   # Project documentation and usage guide
├── LICENSE                     # MIT License
├── go.mod                      # Go module definition
├── ansort.go                   # Main package implementation
├── ansort_test.go              # Core sorting functionality tests
├── config_test.go              # Configuration and functional options tests
├── external_sort_key_test.go   # External system integration tests
└── examples/                   # Usage examples and demos
    ├── basic/
    │   └── main.go             # Basic usage example
    └── external_sort_key/
        └── main.go             # External system integration example
```

### Key Artifacts

#### Core Files
- **`ansort.go`** - The main package implementation containing all sorting logic, types, and functions
- **`ansort_test.go`** - Core sorting functionality tests including natural sorting, tokenization, and basic API tests
- **`config_test.go`** - Configuration and functional options tests covering case sensitivity, leading zeros, and validation
- **`external_sort_key_test.go`** - External system integration tests for sort key generation and consistency verification
- **`go.mod`** - Go module file defining the module path and Go version requirements

#### Examples
- **`examples/`** - Directory containing practical usage examples
  - **`examples/basic/main.go`** - Demonstrates basic package usage and functionality
  - **`examples/external_sort_key/main.go`** - Demonstrates external system integration with sort key generation

This structure follows Go package conventions and provides clear separation of concerns between implementation, testing, documentation, and examples.

## Verification and Testing

Verify the package works correctly on your system:

```bash
# Run all tests and examples
make all            # Run both tests and examples (recommended)

# Individual verification
make test           # Run all tests with coverage
make examples       # Run all examples 

# Development workflow
make verify         # Comprehensive verification (fmt + vet + test + examples)
make help           # Show all available targets
```

Or run individual examples:
```bash
go run examples/basic/main.go
go run examples/external_sort_key/main.go
go run examples/batch_processing/main.go
go run examples/validation/main.go
```

## Usage Guide

The package provides functionality for two distinct use cases at opposite ends of a data pipeline:

### Use Case 1: Direct Natural Sorting
When you control the sorting process and can perform natural comparison directly on your data.

### Use Case 2: External System Integration
When you need to store data in external systems (databases, search engines) that only support lexicographic sorting, but you want natural ordering to be preserved.

---

## Direct Natural Sorting

Use these approaches when you can control the sorting process directly: **convenience functions** for simple use cases and **sorter objects** for integration with Go's standard library.

### When to Use Direct Natural Sorting

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

---

## External System Integration

Use these functions when you need to store naturally-sortable data in external systems that only support lexicographic sorting.

### When to Use External Sort Keys

**Use external sort keys** when:
- Storing data in databases with lexicographic sorting constraints
- Indexing data in search engines like Elasticsearch
- Working with systems that don't support custom sorting logic
- You need consistent natural ordering across different systems
- Pre-computing sort keys for performance optimization

**Use batch processing** (`ToNaturalSortKeys`) when:
- Processing multiple items with the same configuration (reduces overhead)
- Importing large datasets into external systems
- Performance-critical applications with many sort key generations
- Memory-efficient processing with pre-allocated result slices

### External Sort Key Generation

```go
import "github.com/ashprao/ansort"

// External sort key generation for databases/search engines
sortKey := ansort.ToNaturalSortKey("file10.txt")
// Result: "file0000000010.txt" (padded for lexicographic sorting)

// Custom padding length and case-insensitive keys
sortKey2 := ansort.ToNaturalSortKey("Item5", ansort.WithMaxNumericLength(5), ansort.WithExternalCaseInsensitive())
// Result: "item00005" (lowercase with 5-digit padding)

// Batch processing for multiple inputs (Performance optimized)
data := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
sortKeys := ansort.ToNaturalSortKeys(data, ansort.WithMaxNumericLength(5))
// Result: ["file00010.txt", "file00002.txt", "file00001.txt", "file00020.txt"]

// Batch processing with validation
sortKeys, err := ansort.ToNaturalSortKeysValidated(data, ansort.WithMaxNumericLength(5))
if err != nil {
    log.Fatal(err)
}
```

### Typical External Integration Workflow

```go
// Generate sort keys efficiently using batch processing
data := []string{"file10.txt", "file2.txt", "file1.txt"}
sortKeys := ansort.ToNaturalSortKeys(data) // Single config build, optimized processing

// Store both original and sort key in external system
for i, item := range data {
    // Insert into database/search engine with both values
    // External system sorts by sortKeys[i] lexicographically,
    // maintaining natural order of original item values
}

// Query results from external system ordered by sort key
// Returns original values in natural alphanumeric order
```

---

## Go Standard Library Integration

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

### Direct Natural Sorting Functions

- `SortStrings(data []string, options ...Option)` - Sorts a slice of strings in-place using natural ordering
- `Compare(a, b string, options ...Option) int` - Compares two strings using natural ordering rules

### External System Integration Functions

- `ToNaturalSortKey(input string, options ...ExternalSortKeyOption) string` - Generates lexicographically sortable keys for external systems (databases, Elasticsearch, etc.)
- `ToNaturalSortKeyValidated(input string, options ...ExternalSortKeyOption) (string, error)` - Generates keys with comprehensive validation and error reporting
- `ToNaturalSortKeys(inputs []string, options ...ExternalSortKeyOption) []string` - Batch processing for multiple inputs with performance optimization
- `ToNaturalSortKeysValidated(inputs []string, options ...ExternalSortKeyOption) ([]string, error)` - Batch processing with comprehensive validation

### Functional Options (Direct Sorting)

- `WithCaseInsensitive()` - Makes sorting/comparison case-insensitive
- `WithCaseSensitive(sensitive bool)` - Explicitly sets case sensitivity (true = sensitive, false = insensitive)

### External Sort Key Options

- `WithMaxNumericLength(int)` - Sets numeric padding length for external sort keys (default: 10)
- `WithExternalCaseSensitive(sensitive bool)` - Explicitly sets case sensitivity for external keys (true = sensitive, false = insensitive)
- `WithExternalCaseInsensitive()` - Convenience option for case-insensitive external sort key generation

### Standard Library Integration

- `NewSorter(data []string, options ...Option) *AlphanumericSorter` - Creates a sorter implementing `sort.Interface`

### Input Validation and Error Handling

The package provides "validated" variants of core functions that perform comprehensive input and configuration validation with detailed error reporting. Use these when you need strict validation and want to handle errors explicitly.

#### When to Use Validation Functions

**Use validation functions** (`*Validated` variants) when:
- Building production systems that require robust error handling
- Working with user input that might be invalid
- You need detailed error messages for debugging
- Integrating with systems that expect explicit error handling
- Writing libraries or frameworks where input validation is critical

**Use convenience functions** (non-validated variants) when:
- Building prototypes or internal tools with trusted input
- You prefer simpler APIs without error handling overhead
- Performance is critical and input is guaranteed to be valid
- You want graceful degradation instead of explicit errors

#### Validated Function Variants

- `SortStringsValidated(data []string, options ...Option) error` - Sorts with comprehensive validation and error reporting
- `NewSorterValidated(data []string, options ...Option) (*AlphanumericSorter, error)` - Creates sorter with validation
- `CompareValidated(a, b string, options ...Option) (int, error)` - Compares with validation
- `ToNaturalSortKeyValidated(input string, options ...ExternalSortKeyOption) (string, error)` - Generates external sort keys with validation
- `ToNaturalSortKeysValidated(inputs []string, options ...ExternalSortKeyOption) ([]string, error)` - Batch generates external sort keys with validation

#### Error Types

- `ValidationError` - Detailed validation errors with field-specific messages
- `ErrInvalidConfig` - Configuration validation failures
- `ErrNilInput` - Nil input provided where non-nil expected

#### Example: Production-Ready Error Handling

```go
import (
    "fmt"
    "errors"
    "github.com/ashprao/ansort"
)

func ProcessUserData(userInput []string) error {
    // Use validated function for robust error handling
    err := ansort.SortStringsValidated(userInput, ansort.WithCaseInsensitive())
    if err != nil {
        // Handle specific error types
        var validationErr *ansort.ValidationError
        if errors.As(err, &validationErr) {
            return fmt.Errorf("invalid input in field %s: %s", 
                validationErr.Field, validationErr.Message)
        }
        return fmt.Errorf("sorting failed: %w", err)
    }
    
    return nil
}

func GenerateExternalSortKeys(userInput []string, paddingLength int) ([]string, error) {
    // Use batch processing for efficiency and error handling
    return ansort.ToNaturalSortKeysValidated(userInput, 
        ansort.WithMaxNumericLength(paddingLength))
}

// For trusted internal data, use convenience functions  
func ProcessInternalData(data []string) {
    ansort.SortStrings(data) // No error handling needed
    keys := ansort.ToNaturalSortKeys(data) // Batch processing for efficiency
    // ... process keys
}
```

#### Validation vs Convenience Function Comparison

| Scenario | Recommended Function | Reason |
|----------|---------------------|---------|
| User input processing | `SortStringsValidated` | Need explicit error handling |
| API endpoint data | `NewSorterValidated` | Robust validation required |
| Internal data processing | `SortStrings` | Simpler API, trusted input |
| Library/framework building | `*Validated` variants | Proper error propagation |
| Performance-critical loops | Non-validated | Avoid validation overhead |
| Configuration validation | `*Validated` variants | Detailed error messages |

### Types

- `AlphanumericSorter` - Implements `sort.Interface` for integration with Go's sort package
- `Config` - Internal configuration structure (used by functional options)
- `Option` - Function type for configuring sorting behavior
- `ExternalSortKeyConfig` - Configuration for external sort key generation
- `ExternalSortKeyOption` - Function type for configuring external sort key behavior
- `ValidationError` - Detailed validation error with field-specific information

## Examples

See the `examples/` directory for comprehensive usage examples:
- `examples/basic/` - Basic natural sorting functionality
- `examples/external_sort_key/` - External system integration with sort key generation
- `examples/batch_processing/` - Batch processing optimization for multiple items
- `examples/validation/` - Validation functions with error handling

## Current Status

This package is currently at **v0.2.0** - a stable release with comprehensive natural sorting and batch processing features.

### ✅ Implemented Features (v0.2.0)
- **Core natural sorting**: Alphanumeric parsing and intelligent comparison
- **Case sensitivity options**: Case-sensitive/insensitive modes with functional options
- **Multi-segment number support**: Semantic versioning, IP addresses, complex patterns
- **Leading zero handling**: Proper numeric comparison with preserved formatting
- **Standard library integration**: Full `sort.Interface` implementation
- **External system integration**: Generate lexicographically sortable keys for external systems
- **Batch processing**: Optimized performance for processing multiple items efficiently
- **Comprehensive validation**: Input validation with detailed error reporting
- **High test coverage**: 95%+ coverage with comprehensive test suite

### 🔄 Future Versions
- **v0.3.0**: Unicode and special character support
- **v0.4.0**: Performance optimizations and API polish
- **v1.0.0**: Stable API guarantee after community feedback

### 📦 Version Compatibility
- **Pre-v1.0.0**: API may evolve based on feedback
- **Post-v1.0.0**: Semantic versioning with backward compatibility guarantees

### Key Discovery
The natural tokenization approach handles complex patterns automatically without special cases:
- Multi-segment numbers work through token-by-token comparison
- Decimal numbers sort correctly via natural numeric token handling
- Version strings, IP addresses, and mixed patterns all work seamlessly

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
