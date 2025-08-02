// Package ansort provides natural sorting capabilities for alphanumeric strings.
// It implements intelligent sorting where numeric parts are sorted numerically
// rather than lexicographically, resulting in more intuitive ordering.
//
// Example:
//
//	Standard sort: ["item1", "item10", "item2", "item20", "item3"]
//	Natural sort:  ["item1", "item2", "item3", "item10", "item20"]
//
// Basic usage:
//
//	data := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
//	ansort.SortStrings(data)
//	// Result: ["file1.txt", "file2.txt", "file10.txt", "file20.txt"]
package ansort

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// TokenType represents the type of token in a parsed string
type TokenType int

const (
	// AlphaToken represents alphabetic characters
	AlphaToken TokenType = iota
	// NumericToken represents numeric characters
	NumericToken
)

// Token represents a segment of a string (either alphabetic or numeric)
type Token struct {
	Type  TokenType
	Value string
}

// Config holds configuration options for alphanumeric sorting
type Config struct {
	// CaseSensitive determines whether alphabetic comparisons are case-sensitive
	// Default: true (case-sensitive)
	CaseSensitive bool
}

// ErrInvalidConfig is returned when configuration validation fails
var ErrInvalidConfig = errors.New("invalid configuration")

// ErrNilInput is returned when a nil input is provided where non-nil is expected
var ErrNilInput = errors.New("nil input provided")

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return "validation error in field '" + e.Field + "': " + e.Message
}

// validateConfig validates the configuration options
// Returns an error if the configuration is invalid
func validateConfig(config Config) error {
	// Currently all Config fields are valid by design (bool can't be invalid)
	// This function provides a framework for future validation needs
	return nil
}

// validateSlice validates that a slice is not nil for operations that require non-nil input
// Returns an error if validation fails
func validateSlice(data []string, operationName string) error {
	if data == nil {
		return &ValidationError{
			Field:   "data",
			Message: "slice cannot be nil for " + operationName,
		}
	}
	return nil
}

// Option is a functional option for configuring sorting behavior
type Option func(*Config)

// WithCaseSensitive sets the case sensitivity option
func WithCaseSensitive(caseSensitive bool) Option {
	return func(c *Config) {
		c.CaseSensitive = caseSensitive
	}
}

// WithCaseInsensitive is a convenience option for case-insensitive sorting
func WithCaseInsensitive() Option {
	return WithCaseSensitive(false)
}

// DefaultConfig returns a Config with default settings
func DefaultConfig() Config {
	return Config{
		CaseSensitive: true,
	}
}

// buildConfig creates a configuration from functional options
func buildConfig(options ...Option) Config {
	config := DefaultConfig()
	for _, option := range options {
		option(&config)
	}
	return config
}

// NewSorter creates a new AlphanumericSorter with the specified options.
// Accepts nil slices for convenience (they will result in empty sorters).
// Use NewSorterValidated for strict validation with error reporting.
func NewSorter(data []string, options ...Option) *AlphanumericSorter {
	config := buildConfig(options...)
	return &AlphanumericSorter{
		data:   data,
		config: config,
	}
}

// NewSorterValidated creates a new AlphanumericSorter with comprehensive validation.
// Returns an error if validation fails.
func NewSorterValidated(data []string, options ...Option) (*AlphanumericSorter, error) {
	if err := validateSlice(data, "NewSorterValidated"); err != nil {
		return nil, err
	}

	config := buildConfig(options...)
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return &AlphanumericSorter{
		data:   data,
		config: config,
	}, nil
}

// AlphanumericSorter implements sort.Interface for natural alphanumeric sorting
type AlphanumericSorter struct {
	data   []string
	config Config
}

// Len returns the number of elements in the collection
func (s AlphanumericSorter) Len() int {
	return len(s.data)
}

// Swap swaps the elements with indexes i and j
func (s AlphanumericSorter) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

// Less reports whether the element with index i should sort before
// the element with index j using natural alphanumeric comparison
func (s AlphanumericSorter) Less(i, j int) bool {
	return Compare(s.data[i], s.data[j], WithCaseSensitive(s.config.CaseSensitive)) < 0
}

// Compare compares two strings using natural alphanumeric sorting rules.
// Accepts functional options to customize behavior.
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	+1 if a > b
func Compare(a, b string, options ...Option) int {
	// Handle identical strings quickly
	if a == b {
		return 0
	}

	// Build configuration from options
	config := buildConfig(options...)

	// Tokenize both strings
	tokensA := parseString(a)
	tokensB := parseString(b)

	// Compare token by token
	minLen := len(tokensA)
	if len(tokensB) < minLen {
		minLen = len(tokensB)
	}

	for i := 0; i < minLen; i++ {
		result := compareTokensWithConfig(tokensA[i], tokensB[i], config)
		if result != 0 {
			return result
		}
	}

	// If all compared tokens are equal, the shorter string comes first
	if len(tokensA) < len(tokensB) {
		return -1
	} else if len(tokensA) > len(tokensB) {
		return 1
	}

	return 0
}

// CompareValidated compares two strings using natural alphanumeric sorting rules
// with comprehensive validation.
// Returns the comparison result and an error if validation fails.
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	+1 if a > b
func CompareValidated(a, b string, options ...Option) (int, error) {
	// Handle identical strings quickly
	if a == b {
		return 0, nil
	}

	// Build and validate configuration from options
	config := buildConfig(options...)
	if err := validateConfig(config); err != nil {
		return 0, err
	}

	// Tokenize both strings
	tokensA := parseString(a)
	tokensB := parseString(b)

	// Compare token by token
	minLen := len(tokensA)
	if len(tokensB) < minLen {
		minLen = len(tokensB)
	}

	for i := 0; i < minLen; i++ {
		result := compareTokensWithConfig(tokensA[i], tokensB[i], config)
		if result != 0 {
			return result, nil
		}
	}

	// If all compared tokens are equal, the shorter string comes first
	if len(tokensA) < len(tokensB) {
		return -1, nil
	} else if len(tokensA) > len(tokensB) {
		return 1, nil
	}

	return 0, nil
}

// SortStrings sorts a slice of strings using natural alphanumeric ordering.
// Accepts functional options to customize behavior.
// The slice is modified in-place.
// For nil slices, the function returns early without error for convenience.
// Use SortStringsValidated for strict validation with error reporting.
func SortStrings(data []string, options ...Option) {
	if data == nil {
		return
	}
	config := buildConfig(options...)
	if err := validateConfig(config); err != nil {
		// Log or handle config error gracefully - for backward compatibility,
		// we don't return errors from this function
		return
	}
	sorter := AlphanumericSorter{data: data, config: config}
	sort.Sort(sorter)
}

// SortStringsValidated sorts a slice of strings using natural alphanumeric ordering
// with comprehensive validation and error reporting.
// Returns an error if validation fails.
func SortStringsValidated(data []string, options ...Option) error {
	if err := validateSlice(data, "SortStringsValidated"); err != nil {
		return err
	}

	config := buildConfig(options...)
	if err := validateConfig(config); err != nil {
		return err
	}

	sorter := AlphanumericSorter{data: data, config: config}
	sort.Sort(sorter)
	return nil
}

// parseString tokenizes a string into alternating alphabetic and numeric segments.
// It separates the input string into tokens where each token is either purely
// alphabetic or purely numeric characters.
func parseString(s string) []Token {
	if len(s) == 0 {
		return []Token{}
	}

	var tokens []Token
	runes := []rune(s)

	for i := 0; i < len(runes); {
		start := i

		// Check if current character is a digit
		if unicode.IsDigit(runes[i]) {
			// Collect all consecutive digits
			for i < len(runes) && unicode.IsDigit(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{
				Type:  NumericToken,
				Value: string(runes[start:i]),
			})
		} else {
			// Collect all consecutive non-digits
			for i < len(runes) && !unicode.IsDigit(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{
				Type:  AlphaToken,
				Value: string(runes[start:i]),
			})
		}
	}

	return tokens
}

// compareTokens compares two tokens according to their types and values.
// Uses default configuration (case-sensitive).
// Numeric tokens are compared numerically, alphabetic tokens lexicographically.
// Numeric tokens always sort before alphabetic tokens when types differ.
func compareTokens(a, b Token, options ...Option) int {
	config := buildConfig(options...)
	return compareTokensWithConfig(a, b, config)
}

// compareTokensWithConfig compares two tokens according to their types and values
// with the specified configuration.
// Numeric tokens are compared numerically, alphabetic tokens lexicographically.
// Numeric tokens always sort before alphabetic tokens when types differ.
func compareTokensWithConfig(a, b Token, config Config) int {
	// If types are different, numeric comes before alphabetic
	if a.Type != b.Type {
		if a.Type == NumericToken {
			return -1
		}
		return 1
	}

	// Both tokens are the same type
	if a.Type == NumericToken {
		// Compare numerically - convert to integers
		aNum, aErr := strconv.Atoi(a.Value)
		bNum, bErr := strconv.Atoi(b.Value)

		// If either conversion fails, fall back to string comparison
		if aErr != nil || bErr != nil {
			if a.Value < b.Value {
				return -1
			} else if a.Value > b.Value {
				return 1
			}
			return 0
		}

		// Compare as integers
		if aNum < bNum {
			return -1
		} else if aNum > bNum {
			return 1
		}

		// Phase 2.2: Leading Zero Handling
		// When numeric values are equal, use string comparison as tie-breaker
		// This means "1" comes before "001" (shorter first)
		if a.Value != b.Value {
			if len(a.Value) < len(b.Value) {
				return -1
			} else if len(a.Value) > len(b.Value) {
				return 1
			}
			// If same length but different strings, use lexicographic order
			if a.Value < b.Value {
				return -1
			} else if a.Value > b.Value {
				return 1
			}
		}

		return 0
	} else {
		// Both are alphabetic - compare based on case sensitivity setting
		aValue := a.Value
		bValue := b.Value

		if !config.CaseSensitive {
			aValue = strings.ToLower(aValue)
			bValue = strings.ToLower(bValue)
		}

		if aValue < bValue {
			return -1
		} else if aValue > bValue {
			return 1
		}
		return 0
	}
}
