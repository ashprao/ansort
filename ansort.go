// Package ansort provides natural sorting capabilities for alphanumeric strings.
// It implements intelligent sorting where numeric parts are sorted numerically
// rather than lexicographically, resulting in more intuitive ordering.
//
// The package supports two complementary use cases:
//
//  1. Direct natural sorting - Sort data in-memory using intelligent alphanumeric comparison
//  2. External system integration - Generate lexicographically sortable keys that maintain
//     natural order when sorted by external systems (databases, search engines)
//
// Example of direct natural sorting:
//
//	Standard sort: ["item1", "item10", "item2", "item20", "item3"]
//	Natural sort:  ["item1", "item2", "item3", "item10", "item20"]
//
// Basic usage for direct sorting:
//
//	data := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
//	ansort.SortStrings(data)
//	// Result: ["file1.txt", "file2.txt", "file10.txt", "file20.txt"]
//
// Example of external system integration:
//
//	sortKey := ansort.ToNaturalSortKey("file10.txt")
//	// Result: "file0000000010.txt" (lexicographically sortable by external systems)
package ansort

import (
	"errors"
	"fmt"
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

// ExternalSortKeyConfig holds configuration options for external sort key generation
type ExternalSortKeyConfig struct {
	// CaseSensitive determines whether alphabetic comparisons are case-sensitive
	// Default: true (case-sensitive)
	CaseSensitive bool
	// MaxNumericLength is the maximum length to pad numeric segments
	// Default: 10 (supports numbers up to 9,999,999,999)
	MaxNumericLength int
}

// DefaultExternalSortKeyConfig returns an ExternalSortKeyConfig with default settings
func DefaultExternalSortKeyConfig() ExternalSortKeyConfig {
	return ExternalSortKeyConfig{
		CaseSensitive:    true,
		MaxNumericLength: 10,
	}
}

// ExternalSortKeyOption is a functional option for configuring external sort key generation
type ExternalSortKeyOption func(*ExternalSortKeyConfig)

// WithMaxNumericLength sets the padding length for numeric segments in external sort keys.
// The default padding length is 10, which supports numbers up to 9,999,999,999.
//
// Example:
//
//	sortKey := ansort.ToNaturalSortKey("item5", ansort.WithMaxNumericLength(3))
//	// Result: "item005"
func WithMaxNumericLength(length int) ExternalSortKeyOption {
	return func(c *ExternalSortKeyConfig) {
		c.MaxNumericLength = length
	}
}

// WithExternalCaseSensitive sets the case sensitivity option for external sort key generation.
// Pass true for case-sensitive (default) or false for case-insensitive behavior.
// This is the external equivalent of WithCaseSensitive() for direct sorting.
//
// Example:
//
//	sortKey := ansort.ToNaturalSortKey("File10.TXT", ansort.WithExternalCaseSensitive(false))
//	// Result: "file0000000010.txt" (case-insensitive)
//
//	sortKey := ansort.ToNaturalSortKey("File10.TXT", ansort.WithExternalCaseSensitive(true))
//	// Result: "File0000000010.TXT" (case-sensitive, default)
func WithExternalCaseSensitive(caseSensitive bool) ExternalSortKeyOption {
	return func(c *ExternalSortKeyConfig) {
		c.CaseSensitive = caseSensitive
	}
}

// WithExternalCaseInsensitive is a convenience option for case-insensitive external sort key generation.
// This is equivalent to WithExternalCaseSensitive(false).
// This is the external equivalent of WithCaseInsensitive() for direct sorting.
//
// Example:
//
//	sortKey := ansort.ToNaturalSortKey("File10.TXT", ansort.WithExternalCaseInsensitive())
//	// Result: "file0000000010.txt"
func WithExternalCaseInsensitive() ExternalSortKeyOption {
	return WithExternalCaseSensitive(false)
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

// validateExternalSortKeyConfig validates the external sort key configuration options
// Returns an error if the configuration is invalid
func validateExternalSortKeyConfig(config ExternalSortKeyConfig) error {
	if config.MaxNumericLength <= 0 {
		return &ValidationError{
			Field:   "MaxNumericLength",
			Message: "must be greater than 0",
		}
	}
	if config.MaxNumericLength > 50 {
		return &ValidationError{
			Field:   "MaxNumericLength",
			Message: "must be 50 or less to prevent excessive memory usage",
		}
	}
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

// WithCaseSensitive sets the case sensitivity option for direct natural sorting.
// Pass true for case-sensitive (default) or false for case-insensitive behavior.
//
// Example:
//
//	ansort.SortStrings(data, ansort.WithCaseSensitive(false))
//	// Same as: ansort.SortStrings(data, ansort.WithCaseInsensitive())
func WithCaseSensitive(caseSensitive bool) Option {
	return func(c *Config) {
		c.CaseSensitive = caseSensitive
	}
}

// WithCaseInsensitive is a convenience option for case-insensitive direct natural sorting.
// This is equivalent to WithCaseSensitive(false).
//
// Example:
//
//	data := []string{"File2.txt", "file10.txt", "FILE1.txt"}
//	ansort.SortStrings(data, ansort.WithCaseInsensitive())
//	// Result: case-insensitive natural ordering
func WithCaseInsensitive() Option {
	return WithCaseSensitive(false)
}

// DefaultConfig returns a Config with default settings for direct natural sorting.
// The default configuration uses case-sensitive comparison.
//
// This function is primarily used internally but is exported for advanced use cases
// where you need to inspect or modify the default configuration.
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

// buildExternalSortKeyConfig creates an external sort key configuration from functional options
func buildExternalSortKeyConfig(options ...ExternalSortKeyOption) ExternalSortKeyConfig {
	config := DefaultExternalSortKeyConfig()
	for _, option := range options {
		option(&config)
	}
	return config
}

// NewSorter creates a new AlphanumericSorter that implements sort.Interface for natural sorting.
// The sorter can be used with Go's standard sort package functions like sort.Sort().
//
// Accepts nil slices for convenience (they will result in empty sorters).
// Use NewSorterValidated for strict validation with error reporting.
//
// Example:
//
//	data := []string{"file10.txt", "file2.txt", "file1.txt"}
//	sorter := ansort.NewSorter(data, ansort.WithCaseInsensitive())
//	sort.Sort(sorter)
//	// data is now sorted case-insensitively in natural order
func NewSorter(data []string, options ...Option) *AlphanumericSorter {
	config := buildConfig(options...)
	return &AlphanumericSorter{
		data:   data,
		config: config,
	}
}

// NewSorterValidated creates a new AlphanumericSorter with comprehensive validation.
// Unlike NewSorter, this function returns an error if validation fails, making it
// suitable for cases where you need strict input validation.
//
// Returns an error if the input data is nil or if configuration validation fails.
//
// Example:
//
//	data := []string{"file10.txt", "file2.txt", "file1.txt"}
//	sorter, err := ansort.NewSorterValidated(data, ansort.WithCaseInsensitive())
//	if err != nil {
//		log.Fatal(err)
//	}
//	sort.Sort(sorter)
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

// Len returns the number of elements in the collection.
// This method implements the sort.Interface.
func (s AlphanumericSorter) Len() int {
	return len(s.data)
}

// Swap swaps the elements with indexes i and j.
// This method implements the sort.Interface.
func (s AlphanumericSorter) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

// Less reports whether the element with index i should sort before
// the element with index j using natural alphanumeric comparison.
// This method implements the sort.Interface.
func (s AlphanumericSorter) Less(i, j int) bool {
	return Compare(s.data[i], s.data[j], WithCaseSensitive(s.config.CaseSensitive)) < 0
}

// Compare compares two strings using natural alphanumeric sorting rules.
// This is the core comparison function that handles intelligent numeric comparison.
//
// Accepts functional options to customize behavior such as case sensitivity.
// Returns:
//
//	-1 if a < b (a should come before b)
//	 0 if a == b (strings are equivalent)
//	+1 if a > b (a should come after b)
//
// Example:
//
//	result := ansort.Compare("file1.txt", "file10.txt")
//	// Returns: -1 (file1.txt comes before file10.txt)
//
//	result := ansort.Compare("File1.txt", "file1.txt", ansort.WithCaseInsensitive())
//	// Returns: 0 (equivalent when case-insensitive)
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
// with comprehensive validation. Unlike Compare, this function returns an error
// if configuration validation fails.
//
// Returns the comparison result and an error if validation fails.
// Returns:
//
//	-1 if a < b (a should come before b)
//	 0 if a == b (strings are equivalent)
//	+1 if a > b (a should come after b)
//
// Example:
//
//	result, err := ansort.CompareValidated("file1.txt", "file10.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// result: -1
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
// This is the main convenience function for direct natural sorting.
//
// The slice is modified in-place. Accepts functional options to customize behavior
// such as case sensitivity. For nil slices, the function returns early without error
// for convenience. Use SortStringsValidated for strict validation with error reporting.
//
// Example:
//
//	data := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
//	ansort.SortStrings(data)
//	// data is now: ["file1.txt", "file2.txt", "file10.txt", "file20.txt"]
//
//	data2 := []string{"File2.txt", "file10.txt", "FILE1.txt"}
//	ansort.SortStrings(data2, ansort.WithCaseInsensitive())
//	// data2 is now sorted case-insensitively
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
// with comprehensive validation and error reporting. Unlike SortStrings, this function
// returns an error if validation fails, making it suitable for cases where you need
// strict input validation.
//
// Returns an error if the input data is nil or if configuration validation fails.
//
// Example:
//
//	data := []string{"file10.txt", "file2.txt", "file1.txt"}
//	err := ansort.SortStringsValidated(data, ansort.WithCaseInsensitive())
//	if err != nil {
//		log.Fatal(err)
//	}
//	// data is now sorted case-insensitively
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

// ToNaturalSortKey generates a lexicographically sortable string from an alphanumeric input
// that maintains natural alphanumeric ordering when sorted by external systems.
// This is the main function for external system integration.
//
// The function parses the input string into numeric and non-numeric segments,
// pads numeric segments with leading zeros to ensure correct lexicographical sorting,
// and optionally normalizes case for consistent ordering.
//
// When external systems (databases, search engines like Elasticsearch) sort the generated
// keys lexicographically, the result maintains the same natural order that ansort's
// direct sorting would produce.
//
// Options:
//   - WithMaxNumericLength(int): Sets padding length (default: 10, supports up to 9,999,999,999)
//   - WithExternalCaseSensitive(bool): Explicitly sets case sensitivity (true = sensitive, false = insensitive)
//   - WithExternalCaseInsensitive(): Convenience option for case-insensitive behavior
//
// Example:
//
//	// Basic usage
//	sortKey := ansort.ToNaturalSortKey("file10.txt")
//	// Result: "file0000000010.txt"
//
//	// Custom padding and case-insensitive
//	sortKey := ansort.ToNaturalSortKey("Item5",
//		ansort.WithMaxNumericLength(3),
//		ansort.WithExternalCaseInsensitive())
//	// Result: "item005"
//
//	// External system integration workflow:
//	// 1. Generate keys: originalValue -> sortKey
//	// 2. Store both in external system
//	// 3. Query ordered by sortKey
//	// 4. Results maintain natural order
func ToNaturalSortKey(input string, options ...ExternalSortKeyOption) string {
	if input == "" {
		return ""
	}

	// Build configuration from options
	config := buildExternalSortKeyConfig(options...)

	// Note: For backward compatibility, this function doesn't validate configuration.
	// Use ToNaturalSortKeyValidated for strict validation with error reporting.

	// Tokenize the input string using existing parseString function
	tokens := parseString(input)

	// Build the sort key by processing each token
	var result strings.Builder

	for _, token := range tokens {
		if token.Type == NumericToken {
			// Pad numeric tokens with leading zeros
			paddedNumber := padNumericToken(token.Value, config.MaxNumericLength)
			result.WriteString(paddedNumber)
		} else {
			// Handle alphabetic tokens based on case sensitivity
			alphaValue := token.Value
			if !config.CaseSensitive {
				alphaValue = strings.ToLower(alphaValue)
			}
			result.WriteString(alphaValue)
		}
	}

	return result.String()
}

// ToNaturalSortKeyValidated generates a lexicographically sortable string with comprehensive validation.
// Unlike ToNaturalSortKey, this function returns an error if configuration validation fails,
// making it suitable for cases where you need strict input validation.
//
// Returns an error if configuration validation fails (e.g., invalid MaxNumericLength).
//
// Example:
//
//	sortKey, err := ansort.ToNaturalSortKeyValidated("file10.txt",
//		ansort.WithMaxNumericLength(100))
//	if err != nil {
//		log.Fatal(err) // Will fail due to excessive padding length
//	}
func ToNaturalSortKeyValidated(input string, options ...ExternalSortKeyOption) (string, error) {
	if input == "" {
		return "", nil
	}

	// Build and validate configuration from options
	config := buildExternalSortKeyConfig(options...)
	if err := validateExternalSortKeyConfig(config); err != nil {
		return "", err
	}

	// Tokenize the input string using existing parseString function
	tokens := parseString(input)

	// Build the sort key by processing each token
	var result strings.Builder

	for _, token := range tokens {
		if token.Type == NumericToken {
			// Pad numeric tokens with leading zeros
			paddedNumber := padNumericToken(token.Value, config.MaxNumericLength)
			result.WriteString(paddedNumber)
		} else {
			// Handle alphabetic tokens based on case sensitivity
			alphaValue := token.Value
			if !config.CaseSensitive {
				alphaValue = strings.ToLower(alphaValue)
			}
			result.WriteString(alphaValue)
		}
	}

	return result.String(), nil
}

// padNumericToken pads a numeric string with leading zeros to the specified length
func padNumericToken(numStr string, maxLength int) string {
	// If the number is already longer than maxLength, return as-is
	// This prevents truncation of valid numbers
	if len(numStr) >= maxLength {
		return numStr
	}

	// Pad with leading zeros
	return fmt.Sprintf("%0*s", maxLength, numStr)
}
