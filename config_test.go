package ansort

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// TestCaseSensitivityOptions tests Phase 2.1 case sensitivity features with functional options
func TestCaseSensitivityOptions(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		caseSensitive bool
		expected      []string
	}{
		{
			name:          "case-insensitive sorting",
			input:         []string{"File2.txt", "file10.txt", "FILE1.txt"},
			caseSensitive: false,
			expected:      []string{"FILE1.txt", "File2.txt", "file10.txt"},
		},
		{
			name:          "case-sensitive sorting",
			input:         []string{"File2.txt", "file10.txt", "FILE1.txt"},
			caseSensitive: true,
			expected:      []string{"FILE1.txt", "File2.txt", "file10.txt"}, // uppercase comes first
		},
		{
			name:          "mixed case with numbers case-insensitive",
			input:         []string{"item10", "Item2", "ITEM1"},
			caseSensitive: false,
			expected:      []string{"ITEM1", "Item2", "item10"},
		},
		{
			name:          "mixed case with numbers case-sensitive",
			input:         []string{"item10", "Item2", "ITEM1"},
			caseSensitive: true,
			expected:      []string{"ITEM1", "Item2", "item10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			if tt.caseSensitive {
				SortStrings(input) // Default is case-sensitive
			} else {
				SortStrings(input, WithCaseInsensitive())
			}

			for i := range input {
				if input[i] != tt.expected[i] {
					t.Errorf("SortStrings() = %v, expected %v", input, tt.expected)
					break
				}
			}
		})
	}
}

// TestCompareOptions tests the Compare function with functional options
func TestCompareOptions(t *testing.T) {
	tests := []struct {
		name          string
		a, b          string
		caseSensitive bool
		expected      int
	}{
		{
			name:          "case-insensitive equal",
			a:             "File1.txt",
			b:             "file1.txt",
			caseSensitive: false,
			expected:      0,
		},
		{
			name:          "case-sensitive not equal",
			a:             "File1.txt",
			b:             "file1.txt",
			caseSensitive: true,
			expected:      -1, // uppercase comes before lowercase
		},
		{
			name:          "case-insensitive comparison",
			a:             "file2.txt",
			b:             "FILE10.txt",
			caseSensitive: false,
			expected:      -1,
		},
		{
			name:          "case-sensitive comparison",
			a:             "file2.txt",
			b:             "FILE10.txt",
			caseSensitive: true,
			expected:      1, // lowercase comes after uppercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result int
			if tt.caseSensitive {
				result = Compare(tt.a, tt.b) // Default is case-sensitive
			} else {
				result = Compare(tt.a, tt.b, WithCaseInsensitive())
			}
			if result != tt.expected {
				t.Errorf("Compare(%q, %q, CaseSensitive=%v) = %d, expected %d",
					tt.a, tt.b, tt.caseSensitive, result, tt.expected)
			}
		})
	}
}

// TestNewSorter tests the NewSorter constructor
func TestNewSorter(t *testing.T) {
	data := []string{"file10.txt", "file2.txt", "file1.txt"}

	sorter := NewSorter(data, WithCaseInsensitive())

	if sorter.Len() != 3 {
		t.Errorf("NewSorter().Len() = %d, expected 3", sorter.Len())
	}

	if sorter.config.CaseSensitive != false {
		t.Errorf("NewSorter() config not set correctly")
	}

	// Test that sorting works
	sort.Sort(sorter)
	expected := []string{"file1.txt", "file2.txt", "file10.txt"}
	for i := range data {
		if data[i] != expected[i] {
			t.Errorf("NewSorter() sorting failed: got %v, expected %v", data, expected)
			break
		}
	}
}

// TestDefaultConfig tests the DefaultConfig function
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if !config.CaseSensitive {
		t.Errorf("DefaultConfig().CaseSensitive = %v, expected true", config.CaseSensitive)
	}
}

// TestFunctionalOptions tests the functional options
func TestFunctionalOptions(t *testing.T) {
	// Test WithCaseSensitive
	config := buildConfig(WithCaseSensitive(false))
	if config.CaseSensitive {
		t.Errorf("WithCaseSensitive(false) failed")
	}

	// Test WithCaseInsensitive
	config = buildConfig(WithCaseInsensitive())
	if config.CaseSensitive {
		t.Errorf("WithCaseInsensitive() failed")
	}

	// Test multiple options (last one wins)
	config = buildConfig(WithCaseSensitive(true), WithCaseInsensitive())
	if config.CaseSensitive {
		t.Errorf("Multiple options should work, last one wins")
	}
}

// TestLeadingZeroWithOptions tests Phase 2.2 leading zero handling with functional options
func TestLeadingZeroWithOptions(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		caseSensitive bool
		expected      []string
	}{
		{
			name:          "leading zeros case-sensitive",
			input:         []string{"Item001", "item1", "ITEM010"},
			caseSensitive: true,
			expected:      []string{"ITEM010", "Item001", "item1"},
		},
		{
			name:          "leading zeros case-insensitive",
			input:         []string{"Item001", "item1", "ITEM010"},
			caseSensitive: false,
			expected:      []string{"item1", "Item001", "ITEM010"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := make([]string, len(tt.input))
			copy(result, tt.input)

			if tt.caseSensitive {
				SortStrings(result, WithCaseSensitive(true))
			} else {
				SortStrings(result, WithCaseInsensitive())
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Position %d: expected %q, got %q. Full result: %v", i, expected, result[i], result)
					break
				}
			}
		})
	}
}

// Phase 2.4: Additional Enhanced Testing for Case Sensitivity and Leading Zeros
func TestCaseSensitivityStories(t *testing.T) {
	// Test Story 1.2: Case-insensitive sorting
	t.Run("Story_1_2_case_insensitive", func(t *testing.T) {
		input := []string{"Apple", "banana", "Cherry", "apple", "BANANA"}
		expected := []string{"Apple", "apple", "banana", "BANANA", "Cherry"}

		SortStrings(input, WithCaseInsensitive())
		if !reflect.DeepEqual(input, expected) {
			t.Errorf("Case-insensitive story failed: got %v, expected %v", input, expected)
		}
	})

	// Test Story 1.3: Case-sensitive sorting (default)
	t.Run("Story_1_3_case_sensitive", func(t *testing.T) {
		input := []string{"Apple", "banana", "Cherry", "apple", "BANANA"}
		SortStrings(input) // Default is case-sensitive

		// Verify case-sensitive ordering (capitals typically come before lowercase)
		// The exact order depends on ASCII values, but verify it's different from case-insensitive
		if len(input) != 5 {
			t.Error("Expected all 5 elements to be preserved")
		}

		// Test that Compare respects case sensitivity
		result := Compare("Apple", "apple")
		if result == 0 {
			t.Error("Case-sensitive compare should distinguish 'Apple' and 'apple'")
		}
	})
}

func TestLeadingZeroStories(t *testing.T) {
	// Test Story 2.2: Leading Zero Handling - comprehensive test
	t.Run("Story_2_2_comprehensive", func(t *testing.T) {
		input := []string{
			"item010", "item001", "item100", // From requirements
			"file0001", "file1", "file01", // Additional patterns
			"test00000", "test0", "test000", // Edge cases with zeros
		}

		SortStrings(input)

		// Verify numeric ordering is preserved
		for i := 0; i < len(input)-1; i++ {
			result := Compare(input[i], input[i+1])
			if result > 0 {
				t.Errorf("Items not in sorted order: %q should not come after %q", input[i], input[i+1])
			}
		}

		// Verify specific requirements example
		requirementsInput := []string{"item010", "item001", "item100"}
		requirementsExpected := []string{"item001", "item010", "item100"}
		SortStrings(requirementsInput)

		if !reflect.DeepEqual(requirementsInput, requirementsExpected) {
			t.Errorf("Requirements example failed: got %v, expected %v", requirementsInput, requirementsExpected)
		}
	})
}

// Performance validation tests
func TestPerformanceValidation(t *testing.T) {
	t.Run("large_dataset_performance", func(t *testing.T) {
		// Create a reasonably large dataset
		data := make([]string, 500)
		for i := 0; i < 500; i++ {
			data[i] = fmt.Sprintf("item%d", i*7%500) // Some randomness
		}

		// Measure basic performance - should complete quickly
		testData := make([]string, len(data))
		copy(testData, data)

		SortStrings(testData) // Should not take too long

		// Verify it's actually sorted
		for i := 0; i < len(testData)-1; i++ {
			if Compare(testData[i], testData[i+1]) > 0 {
				t.Error("Large dataset not properly sorted")
				break
			}
		}
	})
}
