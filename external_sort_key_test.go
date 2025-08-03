package ansort

import (
	"errors"
	"fmt"
	"sort"
	"testing"
)

// TestToNaturalSortKey tests the basic external sort key generation functionality
func TestToNaturalSortKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  []ExternalSortKeyOption
		expected string
	}{
		{
			name:     "basic numeric padding",
			input:    "file10.txt",
			options:  nil, // default options
			expected: "file0000000010.txt",
		},
		{
			name:     "empty string",
			input:    "",
			options:  nil,
			expected: "",
		},
		{
			name:     "custom padding length",
			input:    "item5",
			options:  []ExternalSortKeyOption{WithMaxNumericLength(3)},
			expected: "item005",
		},
		{
			name:     "case insensitive",
			input:    "File10.TXT",
			options:  []ExternalSortKeyOption{WithExternalCaseInsensitive()},
			expected: "file0000000010.txt",
		},
		{
			name:     "multiple numeric segments",
			input:    "v1.2.10",
			options:  []ExternalSortKeyOption{WithMaxNumericLength(3)},
			expected: "v001.002.010",
		},
		{
			name:     "number longer than padding",
			input:    "file12345678901",
			options:  []ExternalSortKeyOption{WithMaxNumericLength(5)},
			expected: "file12345678901", // no truncation, preserves original
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToNaturalSortKey(tt.input, tt.options...)
			if result != tt.expected {
				t.Errorf("ToNaturalSortKey(%q, %v) = %q; want %q",
					tt.input, tt.options, result, tt.expected)
			}
		})
	}
}

// TestExternalSortKeyConsistency tests that external sort keys maintain the same
// ordering as natural sort when lexicographically sorted
func TestExternalSortKeyConsistency(t *testing.T) {
	testData := []string{
		"file1.txt",
		"file2.txt",
		"file10.txt",
		"file20.txt",
		"item1",
		"item10",
		"item2",
		"v1.2.3",
		"v1.10.2",
		"v1.2.10",
	}

	// Create a copy for natural sorting
	naturalSorted := make([]string, len(testData))
	copy(naturalSorted, testData)
	SortStrings(naturalSorted)

	// Generate sort keys and create pairs
	type keyValuePair struct {
		key   string
		value string
	}

	var pairs []keyValuePair
	for _, item := range testData {
		key := ToNaturalSortKey(item)
		pairs = append(pairs, keyValuePair{key: key, value: item})
	}

	// Sort by external sort keys (lexicographically)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].key < pairs[j].key
	})

	// Extract the externally sorted values
	externallySorted := make([]string, len(pairs))
	for i, pair := range pairs {
		externallySorted[i] = pair.value
	}

	// Compare the results
	if len(naturalSorted) != len(externallySorted) {
		t.Fatalf("Length mismatch: natural=%d, external=%d",
			len(naturalSorted), len(externallySorted))
	}

	for i := 0; i < len(naturalSorted); i++ {
		if naturalSorted[i] != externallySorted[i] {
			t.Errorf("Order mismatch at position %d: natural=%q, external=%q",
				i, naturalSorted[i], externallySorted[i])
		}
	}
}

// TestExternalSortKeyConfigValidation tests the validation of external sort key configuration
func TestExternalSortKeyConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      ExternalSortKeyConfig
		expectError bool
		errorField  string
	}{
		{
			name: "valid default config",
			config: ExternalSortKeyConfig{
				CaseSensitive:    true,
				MaxNumericLength: 10,
			},
			expectError: false,
		},
		{
			name: "zero padding length",
			config: ExternalSortKeyConfig{
				CaseSensitive:    true,
				MaxNumericLength: 0,
			},
			expectError: true,
			errorField:  "MaxNumericLength",
		},
		{
			name: "negative padding length",
			config: ExternalSortKeyConfig{
				CaseSensitive:    true,
				MaxNumericLength: -1,
			},
			expectError: true,
			errorField:  "MaxNumericLength",
		},
		{
			name: "excessive padding length",
			config: ExternalSortKeyConfig{
				CaseSensitive:    true,
				MaxNumericLength: 100,
			},
			expectError: true,
			errorField:  "MaxNumericLength",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExternalSortKeyConfig(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for config %+v, but got none", tt.config)
					return
				}

				if valErr, ok := err.(*ValidationError); ok {
					if valErr.Field != tt.errorField {
						t.Errorf("Expected error field %q, got %q", tt.errorField, valErr.Field)
					}
				} else {
					t.Errorf("Expected ValidationError, got %T: %v", err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for config %+v, but got: %v", tt.config, err)
				}
			}
		})
	}
}

// TestExternalSortKeyOptions tests the functional options for external sort key configuration
func TestExternalSortKeyOptions(t *testing.T) {
	// Test default configuration
	defaultConfig := buildExternalSortKeyConfig()
	if !defaultConfig.CaseSensitive {
		t.Errorf("Default case sensitivity should be true, got false")
	}
	if defaultConfig.MaxNumericLength != 10 {
		t.Errorf("Default MaxNumericLength should be 10, got %d", defaultConfig.MaxNumericLength)
	}

	// Test WithMaxNumericLength option
	config := buildExternalSortKeyConfig(WithMaxNumericLength(5))
	if config.MaxNumericLength != 5 {
		t.Errorf("Expected MaxNumericLength=5, got %d", config.MaxNumericLength)
	}

	// Test WithExternalCaseSensitive option
	config = buildExternalSortKeyConfig(WithExternalCaseSensitive(false))
	if config.CaseSensitive {
		t.Errorf("WithExternalCaseSensitive(false) failed: expected CaseSensitive=false, got true")
	}

	config = buildExternalSortKeyConfig(WithExternalCaseSensitive(true))
	if !config.CaseSensitive {
		t.Errorf("WithExternalCaseSensitive(true) failed: expected CaseSensitive=true, got false")
	}

	// Test WithExternalCaseInsensitive option (should be equivalent to WithExternalCaseSensitive(false))
	config = buildExternalSortKeyConfig(WithExternalCaseInsensitive())
	if config.CaseSensitive {
		t.Errorf("Expected CaseSensitive=false, got true")
	}

	// Test that WithExternalCaseInsensitive() and WithExternalCaseSensitive(false) are equivalent
	config1 := buildExternalSortKeyConfig(WithExternalCaseInsensitive())
	config2 := buildExternalSortKeyConfig(WithExternalCaseSensitive(false))
	if config1.CaseSensitive != config2.CaseSensitive {
		t.Errorf("WithExternalCaseInsensitive() and WithExternalCaseSensitive(false) should be equivalent")
	}

	// Test multiple options
	config = buildExternalSortKeyConfig(
		WithMaxNumericLength(8),
		WithExternalCaseInsensitive(),
	)
	if config.MaxNumericLength != 8 {
		t.Errorf("Expected MaxNumericLength=8, got %d", config.MaxNumericLength)
	}
	if config.CaseSensitive {
		t.Errorf("Expected CaseSensitive=false, got true")
	}
}

// TestToNaturalSortKeyValidated tests the validated version of ToNaturalSortKey
func TestToNaturalSortKeyValidated(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		result, err := ToNaturalSortKeyValidated("file10.txt", WithMaxNumericLength(5))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		expected := "file00010.txt"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("invalid configuration - zero padding", func(t *testing.T) {
		_, err := ToNaturalSortKeyValidated("file10.txt", WithMaxNumericLength(0))
		if err == nil {
			t.Error("Expected error for zero padding length, got nil")
		}

		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		} else if validationErr.Field != "MaxNumericLength" {
			t.Errorf("Expected error field 'MaxNumericLength', got %q", validationErr.Field)
		}
	})

	t.Run("invalid configuration - excessive padding", func(t *testing.T) {
		_, err := ToNaturalSortKeyValidated("file10.txt", WithMaxNumericLength(100))
		if err == nil {
			t.Error("Expected error for excessive padding length, got nil")
		}

		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		} else if validationErr.Field != "MaxNumericLength" {
			t.Errorf("Expected error field 'MaxNumericLength', got %q", validationErr.Field)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := ToNaturalSortKeyValidated("")
		if err != nil {
			t.Errorf("Expected no error for empty input, got %v", err)
		}
		if result != "" {
			t.Errorf("Expected empty result for empty input, got %q", result)
		}
	})

	t.Run("comparison with non-validated function", func(t *testing.T) {
		input := "file10.txt"
		options := []ExternalSortKeyOption{WithMaxNumericLength(5), WithExternalCaseInsensitive()}

		validated, err := ToNaturalSortKeyValidated(input, options...)
		if err != nil {
			t.Errorf("Validated function failed: %v", err)
		}

		nonValidated := ToNaturalSortKey(input, options...)

		if validated != nonValidated {
			t.Errorf("Results should be identical: validated=%q, non-validated=%q", validated, nonValidated)
		}
	})
}

// TestToNaturalSortKeys tests the batch processing functionality
func TestToNaturalSortKeys(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result := ToNaturalSortKeys(nil)
		if result != nil {
			t.Errorf("Expected nil result for nil input, got %v", result)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		result := ToNaturalSortKeys([]string{})
		if result == nil {
			t.Error("Expected empty slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got length %d", len(result))
		}
	})

	t.Run("single item", func(t *testing.T) {
		input := []string{"file10.txt"}
		result := ToNaturalSortKeys(input)
		expected := []string{"file0000000010.txt"}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		if result[0] != expected[0] {
			t.Errorf("Expected %q, got %q", expected[0], result[0])
		}
	})

	t.Run("multiple items", func(t *testing.T) {
		input := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
		result := ToNaturalSortKeys(input)
		expected := []string{"file0000000010.txt", "file0000000002.txt", "file0000000001.txt", "file0000000020.txt"}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
			return
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("Index %d: expected %q, got %q", i, exp, result[i])
			}
		}
	})

	t.Run("with options", func(t *testing.T) {
		input := []string{"File10.TXT", "file2.txt", "FILE1.txt"}
		result := ToNaturalSortKeys(input, WithMaxNumericLength(3), WithExternalCaseInsensitive())
		expected := []string{"file010.txt", "file002.txt", "file001.txt"}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
			return
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("Index %d: expected %q, got %q", i, exp, result[i])
			}
		}
	})

	t.Run("empty strings in input", func(t *testing.T) {
		input := []string{"file1.txt", "", "file2.txt"}
		result := ToNaturalSortKeys(input)
		expected := []string{"file0000000001.txt", "", "file0000000002.txt"}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
			return
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("Index %d: expected %q, got %q", i, exp, result[i])
			}
		}
	})

	t.Run("consistency with single function", func(t *testing.T) {
		input := []string{"file10.txt", "file2.txt", "file1.txt"}
		options := []ExternalSortKeyOption{WithMaxNumericLength(5), WithExternalCaseInsensitive()}

		// Batch processing
		batchResult := ToNaturalSortKeys(input, options...)

		// Individual processing
		var individualResult []string
		for _, item := range input {
			individualResult = append(individualResult, ToNaturalSortKey(item, options...))
		}

		if len(batchResult) != len(individualResult) {
			t.Errorf("Length mismatch: batch=%d, individual=%d", len(batchResult), len(individualResult))
			return
		}

		for i, batchItem := range batchResult {
			if batchItem != individualResult[i] {
				t.Errorf("Index %d: batch=%q, individual=%q", i, batchItem, individualResult[i])
			}
		}
	})
}

// TestToNaturalSortKeysValidated tests the validated batch processing functionality
func TestToNaturalSortKeysValidated(t *testing.T) {
	t.Run("nil input validation", func(t *testing.T) {
		_, err := ToNaturalSortKeysValidated(nil)
		if err == nil {
			t.Error("Expected error for nil input, got nil")
		}

		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		} else if validationErr.Field != "data" {
			t.Errorf("Expected error field 'data', got %q", validationErr.Field)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := ToNaturalSortKeysValidated([]string{})
		if err != nil {
			t.Errorf("Expected no error for empty input, got %v", err)
		}
		if result == nil {
			t.Error("Expected empty slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got length %d", len(result))
		}
	})

	t.Run("valid configuration", func(t *testing.T) {
		input := []string{"file10.txt", "file2.txt", "file1.txt"}
		result, err := ToNaturalSortKeysValidated(input, WithMaxNumericLength(5))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		expected := []string{"file00010.txt", "file00002.txt", "file00001.txt"}

		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
			return
		}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("Index %d: expected %q, got %q", i, exp, result[i])
			}
		}
	})

	t.Run("invalid configuration - zero padding", func(t *testing.T) {
		input := []string{"file10.txt", "file2.txt"}
		_, err := ToNaturalSortKeysValidated(input, WithMaxNumericLength(0))
		if err == nil {
			t.Error("Expected error for zero padding length, got nil")
		}

		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		} else if validationErr.Field != "MaxNumericLength" {
			t.Errorf("Expected error field 'MaxNumericLength', got %q", validationErr.Field)
		}
	})

	t.Run("invalid configuration - excessive padding", func(t *testing.T) {
		input := []string{"file10.txt", "file2.txt"}
		_, err := ToNaturalSortKeysValidated(input, WithMaxNumericLength(100))
		if err == nil {
			t.Error("Expected error for excessive padding length, got nil")
		}

		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		} else if validationErr.Field != "MaxNumericLength" {
			t.Errorf("Expected error field 'MaxNumericLength', got %q", validationErr.Field)
		}
	})

	t.Run("consistency with single validated function", func(t *testing.T) {
		input := []string{"file10.txt", "file2.txt", "file1.txt"}
		options := []ExternalSortKeyOption{WithMaxNumericLength(5), WithExternalCaseInsensitive()}

		// Batch processing
		batchResult, err := ToNaturalSortKeysValidated(input, options...)
		if err != nil {
			t.Errorf("Batch processing failed: %v", err)
		}

		// Individual processing
		var individualResult []string
		for _, item := range input {
			individual, err := ToNaturalSortKeyValidated(item, options...)
			if err != nil {
				t.Errorf("Individual processing failed for %q: %v", item, err)
			}
			individualResult = append(individualResult, individual)
		}

		if len(batchResult) != len(individualResult) {
			t.Errorf("Length mismatch: batch=%d, individual=%d", len(batchResult), len(individualResult))
			return
		}

		for i, batchItem := range batchResult {
			if batchItem != individualResult[i] {
				t.Errorf("Index %d: batch=%q, individual=%q", i, batchItem, individualResult[i])
			}
		}
	})

	t.Run("consistency with non-validated batch function", func(t *testing.T) {
		input := []string{"file10.txt", "file2.txt", "file1.txt"}
		options := []ExternalSortKeyOption{WithMaxNumericLength(5), WithExternalCaseInsensitive()}

		// Validated batch processing
		validatedResult, err := ToNaturalSortKeysValidated(input, options...)
		if err != nil {
			t.Errorf("Validated batch processing failed: %v", err)
		}

		// Non-validated batch processing
		nonValidatedResult := ToNaturalSortKeys(input, options...)

		if len(validatedResult) != len(nonValidatedResult) {
			t.Errorf("Length mismatch: validated=%d, non-validated=%d", len(validatedResult), len(nonValidatedResult))
			return
		}

		for i, validatedItem := range validatedResult {
			if validatedItem != nonValidatedResult[i] {
				t.Errorf("Index %d: validated=%q, non-validated=%q", i, validatedItem, nonValidatedResult[i])
			}
		}
	})
}

// TestBatchProcessingPerformance tests the performance characteristics of batch processing
func TestBatchProcessingPerformance(t *testing.T) {
	// Generate a reasonable-sized test dataset
	input := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		input[i] = fmt.Sprintf("file%d.txt", i)
	}

	t.Run("batch processing correctness", func(t *testing.T) {
		result := ToNaturalSortKeys(input)
		if len(result) != len(input) {
			t.Errorf("Expected length %d, got %d", len(input), len(result))
		}

		// Verify that the results are consistent
		for i, original := range input {
			expected := ToNaturalSortKey(original)
			if result[i] != expected {
				t.Errorf("Index %d: expected %q, got %q", i, expected, result[i])
				break // Don't spam the output
			}
		}
	})

	t.Run("order preservation", func(t *testing.T) {
		// Test that batch processing preserves input order
		testInput := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"}
		result := ToNaturalSortKeys(testInput)

		// The output should be in the same order as input (not sorted)
		expected := []string{"file0000000010.txt", "file0000000002.txt", "file0000000001.txt", "file0000000020.txt"}

		for i, exp := range expected {
			if result[i] != exp {
				t.Errorf("Order not preserved. Index %d: expected %q, got %q", i, exp, result[i])
			}
		}
	})
}
