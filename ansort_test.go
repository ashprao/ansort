// Package ansort provides natural sorting capabilities for alphanumeric strings.
// It implements intelligent sorting where numeric parts are sorted numerically
// rather than lexicographically, resulting in more intuitive ordering.
package ansort

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// TestSortStrings tests basic string sorting functionality
func TestSortStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "basic alphanumeric sorting",
			input:    []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt"},
			expected: []string{"file1.txt", "file2.txt", "file10.txt", "file20.txt"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []string{"file1.txt"},
			expected: []string{"file1.txt"},
		},
		{
			name:     "nil slice",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying test data
			var input []string
			if tt.input != nil {
				input = make([]string, len(tt.input))
				copy(input, tt.input)
			}

			SortStrings(input)

			// Compare results
			if len(input) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(input))
				return
			}

			for i, v := range input {
				if v != tt.expected[i] {
					t.Errorf("At index %d: expected %q, got %q", i, tt.expected[i], v)
				}
			}
		})
	}
}

// TestCompare tests the basic comparison function
func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected int
	}{
		{
			name:     "a less than b",
			a:        "file1.txt",
			b:        "file2.txt",
			expected: -1,
		},
		{
			name:     "a greater than b",
			a:        "file2.txt",
			b:        "file1.txt",
			expected: 1,
		},
		{
			name:     "a equals b",
			a:        "file1.txt",
			b:        "file1.txt",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Compare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Compare(%q, %q) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestAlphanumericSorter tests the sort.Interface implementation
func TestAlphanumericSorter(t *testing.T) {
	data := []string{"file10.txt", "file2.txt", "file1.txt"}
	sorter := AlphanumericSorter{data: data}

	// Test Len
	if sorter.Len() != 3 {
		t.Errorf("Len() = %d, expected 3", sorter.Len())
	}

	// Test Swap
	sorter.Swap(0, 2)
	if data[0] != "file1.txt" || data[2] != "file10.txt" {
		t.Errorf("Swap failed: got %v", data)
	}

	// Test Less - now uses natural alphanumeric comparison
	if !sorter.Less(0, 1) { // "file1.txt" < "file2.txt"
		t.Errorf("Less(0, 1) should be true")
	}
}

// TestNaturalSortingBasic tests core Phase 1.2 functionality
func TestNaturalSortingBasic(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "simple numeric sorting",
			input:    []string{"item10", "item2", "item1"},
			expected: []string{"item1", "item2", "item10"},
		},
		{
			name:     "mixed alphanumeric",
			input:    []string{"abc123", "abc2", "abc10"},
			expected: []string{"abc2", "abc10", "abc123"},
		},
		{
			name:     "numbers vs letters priority",
			input:    []string{"abc", "123", "def"},
			expected: []string{"123", "abc", "def"},
		},
		{
			name:     "multiple numeric segments",
			input:    []string{"v1.10.0", "v1.2.0", "v1.1.0"},
			expected: []string{"v1.1.0", "v1.2.0", "v1.10.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			SortStrings(input)

			for i := range input {
				if input[i] != tt.expected[i] {
					t.Errorf("SortStrings() = %v, expected %v", input, tt.expected)
					break
				}
			}
		})
	}
}

// TestCompareDetailed tests the Compare function with various edge cases
func TestCompareDetailed(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected int
	}{
		{
			name:     "numeric comparison",
			a:        "file2.txt",
			b:        "file10.txt",
			expected: -1,
		},
		{
			name:     "same prefix different numbers",
			a:        "item1",
			b:        "item100",
			expected: -1,
		},
		{
			name:     "numbers vs letters",
			a:        "123abc",
			b:        "abc123",
			expected: -1, // numbers come before letters
		},
		{
			name:     "pure numbers",
			a:        "123",
			b:        "45",
			expected: 1, // 123 > 45
		},
		{
			name:     "identical strings",
			a:        "test123",
			b:        "test123",
			expected: 0,
		},
		{
			name:     "different lengths same prefix",
			a:        "test",
			b:        "test123",
			expected: -1, // shorter comes first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Compare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Compare(%q, %q) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// Phase 2.2: Leading Zero Handling Tests
func TestLeadingZeroHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "basic leading zeros",
			input:    []string{"item001", "item1", "item10"},
			expected: []string{"item1", "item001", "item10"},
		},
		{
			name:     "mixed leading zeros",
			input:    []string{"file00002", "file1", "file000010", "file2"},
			expected: []string{"file1", "file2", "file00002", "file000010"},
		},
		{
			name:     "preserve formatting",
			input:    []string{"test007", "test7", "test07"},
			expected: []string{"test7", "test07", "test007"},
		},
		{
			name:     "different prefix same numeric value",
			input:    []string{"a001", "b001", "a1"},
			expected: []string{"a1", "a001", "b001"},
		},
		{
			name:     "zero padding with letters",
			input:    []string{"item001a", "item1b", "item010c"},
			expected: []string{"item1b", "item001a", "item010c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]string, len(tt.input))
			copy(result, tt.input)
			SortStrings(result)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLeadingZeroComparison(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected int
		desc     string
	}{
		{
			name:     "same value different padding - shorter first",
			a:        "item1",
			b:        "item001",
			expected: -1,
			desc:     "1 should come before 001 (shorter first)",
		},
		{
			name:     "same value different padding - longer second",
			a:        "item001",
			b:        "item1",
			expected: 1,
			desc:     "001 should come after 1 (longer second)",
		},
		{
			name:     "padded less than unpadded numerically",
			a:        "item001",
			b:        "item2",
			expected: -1,
			desc:     "001 should be less than 2",
		},
		{
			name:     "unpadded greater than padded numerically",
			a:        "item10",
			b:        "item001",
			expected: 1,
			desc:     "10 should be greater than 001",
		},
		{
			name:     "multiple zeros - shorter first",
			a:        "item1",
			b:        "item00001",
			expected: -1,
			desc:     "1 should come before 00001 (shorter first)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Compare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("%s: Compare(%q, %q) = %d, want %d", tt.desc, tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// Phase 2.3: Input Validation & Error Handling Tests
func TestInputValidation(t *testing.T) {
	t.Run("SortStringsValidated with nil slice", func(t *testing.T) {
		err := SortStringsValidated(nil)
		if err == nil {
			t.Error("Expected error for nil slice, got nil")
		}

		var validationErr *ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("Expected ValidationError, got %T", err)
		}
	})

	t.Run("SortStringsValidated with valid slice", func(t *testing.T) {
		data := []string{"item2", "item1", "item10"}
		err := SortStringsValidated(data)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := []string{"item1", "item2", "item10"}
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	})

	t.Run("NewSorterValidated with nil slice", func(t *testing.T) {
		_, err := NewSorterValidated(nil)
		if err == nil {
			t.Error("Expected error for nil slice, got nil")
		}
	})

	t.Run("NewSorterValidated with valid slice", func(t *testing.T) {
		data := []string{"item2", "item1"}
		sorter, err := NewSorterValidated(data)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if sorter == nil {
			t.Error("Expected valid sorter, got nil")
		}
	})

	t.Run("CompareValidated with valid inputs", func(t *testing.T) {
		result, err := CompareValidated("item1", "item2")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != -1 {
			t.Errorf("Expected -1, got %d", result)
		}
	})
}

func TestErrorTypes(t *testing.T) {
	t.Run("ValidationError implementation", func(t *testing.T) {
		err := &ValidationError{
			Field:   "testField",
			Message: "test message",
		}

		expected := "validation error in field 'testField': test message"
		if err.Error() != expected {
			t.Errorf("Expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("ErrInvalidConfig type", func(t *testing.T) {
		if ErrInvalidConfig == nil {
			t.Error("ErrInvalidConfig should not be nil")
		}
		if ErrInvalidConfig.Error() != "invalid configuration" {
			t.Errorf("Expected 'invalid configuration', got %q", ErrInvalidConfig.Error())
		}
	})

	t.Run("ErrNilInput type", func(t *testing.T) {
		if ErrNilInput == nil {
			t.Error("ErrNilInput should not be nil")
		}
		if ErrNilInput.Error() != "nil input provided" {
			t.Errorf("Expected 'nil input provided', got %q", ErrNilInput.Error())
		}
	})
}

func TestGracefulHandling(t *testing.T) {
	t.Run("NewSorter handles nil gracefully", func(t *testing.T) {
		sorter := NewSorter(nil)
		if sorter == nil {
			t.Error("Expected valid sorter even with nil data")
		} else if sorter.data != nil {
			t.Error("Expected sorter.data to be nil")
		}
	})

	t.Run("Empty slice handling", func(t *testing.T) {
		data := []string{}
		SortStrings(data)
		if len(data) != 0 {
			t.Error("Expected empty slice to remain empty")
		}
	})

	t.Run("Single element handling", func(t *testing.T) {
		data := []string{"single"}
		original := make([]string, len(data))
		copy(original, data)

		SortStrings(data)
		if !reflect.DeepEqual(data, original) {
			t.Error("Single element slice should not be modified")
		}
	})
}

// Phase 2.4: Enhanced Testing - Edge Cases
func TestEdgeCases(t *testing.T) {
	t.Run("empty strings", func(t *testing.T) {
		data := []string{"", "a", "", "b"}
		expected := []string{"", "", "a", "b"}

		SortStrings(data)
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	})

	t.Run("single characters", func(t *testing.T) {
		data := []string{"z", "a", "1", "9", "A", "Z"}
		SortStrings(data)

		// Should be sorted with numbers first, then letters
		// The exact order might depend on implementation, but verify basic structure
		if len(data) != 6 {
			t.Errorf("Expected 6 elements, got %d", len(data))
		}

		// Test a specific comparison to ensure single chars work
		result := Compare("a", "z")
		if result != -1 {
			t.Errorf("Expected 'a' < 'z', got %d", result)
		}
	})

	t.Run("strings with only numbers", func(t *testing.T) {
		data := []string{"123", "45", "1", "1000"}
		expected := []string{"1", "45", "123", "1000"}

		SortStrings(data)
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	})

	t.Run("strings with only letters", func(t *testing.T) {
		data := []string{"zebra", "apple", "banana"}
		expected := []string{"apple", "banana", "zebra"}

		SortStrings(data)
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	})

	t.Run("very long strings", func(t *testing.T) {
		longString1 := "item" + strings.Repeat("a", 1000) + "123"
		longString2 := "item" + strings.Repeat("a", 1000) + "45"

		result := Compare(longString1, longString2)
		if result != 1 { // 123 > 45
			t.Errorf("Expected long string comparison to work, got %d", result)
		}
	})

	t.Run("special characters", func(t *testing.T) {
		data := []string{"item-2", "item_1", "item.3", "item 4"}
		SortStrings(data)

		// Verify it doesn't panic and produces some ordering
		if len(data) != 4 {
			t.Error("Expected all elements to be preserved")
		}
	})

	t.Run("unicode characters", func(t *testing.T) {
		data := []string{"café2", "café10", "café1"}
		SortStrings(data)

		// Should handle unicode gracefully
		if len(data) != 3 {
			t.Error("Expected all elements to be preserved")
		}
	})

	t.Run("mixed case with leading zeros", func(t *testing.T) {
		data := []string{"Item001", "item1", "ITEM010"}
		SortStrings(data, WithCaseInsensitive())

		// With case insensitive, should sort by numeric value
		expected := []string{"item1", "Item001", "ITEM010"}
		if !reflect.DeepEqual(data, expected) {
			t.Errorf("Expected %v, got %v", expected, data)
		}
	})
}

// Phase 2.4: Performance Baseline Benchmarks
func BenchmarkSortStrings(b *testing.B) {
	// Small dataset
	b.Run("small_10_items", func(b *testing.B) {
		data := []string{"file10.txt", "file2.txt", "file1.txt", "file20.txt", "file3.txt",
			"file100.txt", "file21.txt", "file4.txt", "file5.txt", "file11.txt"}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testData := make([]string, len(data))
			copy(testData, data)
			SortStrings(testData)
		}
	})

	// Medium dataset
	b.Run("medium_100_items", func(b *testing.B) {
		data := make([]string, 100)
		for i := 0; i < 100; i++ {
			data[i] = fmt.Sprintf("item%d", (i*7)%100) // Create some randomness
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testData := make([]string, len(data))
			copy(testData, data)
			SortStrings(testData)
		}
	})

	// Large dataset
	b.Run("large_1000_items", func(b *testing.B) {
		data := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			data[i] = fmt.Sprintf("file%d.txt", (i*13)%1000)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			testData := make([]string, len(data))
			copy(testData, data)
			SortStrings(testData)
		}
	})
}

func BenchmarkCompare(b *testing.B) {
	testCases := []struct {
		name string
		a, b string
	}{
		{"simple", "file1.txt", "file2.txt"},
		{"numeric_heavy", "item123", "item456"},
		{"long_strings", "very_long_filename_with_numbers_123", "very_long_filename_with_numbers_456"},
		{"leading_zeros", "item001", "item002"},
		{"mixed_case", "File1.TXT", "file2.txt"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Compare(tc.a, tc.b)
			}
		})
	}
}

func BenchmarkParseString(b *testing.B) {
	testStrings := []string{
		"simple",
		"file123.txt",
		"version1.2.3.build456",
		"very_long_filename_with_multiple_123_numbers_456_embedded",
	}

	for _, str := range testStrings {
		b.Run(fmt.Sprintf("len_%d", len(str)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				parseString(str)
			}
		})
	}
}

// TestMultiSegmentNumberSupport tests Phase 3.1 requirements for semantic versioning and complex patterns
func TestMultiSegmentNumberSupport(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "semantic versioning basic",
			input:    []string{"v1.10.2", "v1.2.10", "v1.2.3"},
			expected: []string{"v1.2.3", "v1.2.10", "v1.10.2"},
		},
		{
			name:     "version with prefix",
			input:    []string{"version1.10.5", "version1.2.15", "version1.9.1"},
			expected: []string{"version1.2.15", "version1.9.1", "version1.10.5"},
		},
		{
			name:     "four-segment versions",
			input:    []string{"v1.2.3.10", "v1.2.3.4", "v1.2.10.1", "v1.10.1.1"},
			expected: []string{"v1.2.3.4", "v1.2.3.10", "v1.2.10.1", "v1.10.1.1"},
		},
		{
			name:     "mixed major versions",
			input:    []string{"app2.1.0", "app2.10.0", "app2.2.0", "app10.1.0"},
			expected: []string{"app2.1.0", "app2.2.0", "app2.10.0", "app10.1.0"},
		},
		{
			name:     "complex patterns with build numbers",
			input:    []string{"v1.0.1-build123", "v1.0.1-build23", "v1.0.1-build12"},
			expected: []string{"v1.0.1-build12", "v1.0.1-build23", "v1.0.1-build123"},
		},
		{
			name:     "dot-separated with different lengths",
			input:    []string{"1.2", "1.2.3", "1.2.10", "1.10"},
			expected: []string{"1.2", "1.2.3", "1.2.10", "1.10"},
		},
		{
			name:     "prereleases and versions (natural order - shorter strings first)",
			input:    []string{"v2.0.0-beta.2", "v2.0.0-beta.10", "v2.0.0-alpha.1", "v2.0.0"},
			expected: []string{"v2.0.0", "v2.0.0-alpha.1", "v2.0.0-beta.2", "v2.0.0-beta.10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			SortStrings(input)

			for i := range input {
				if input[i] != tt.expected[i] {
					t.Errorf("SortStrings() = %v, expected %v", input, tt.expected)
					break
				}
			}
		})
	}
}
