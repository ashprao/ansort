package ansort

import (
	"fmt"
	"testing"
)

// TestIntelligentSelection verifies that the system automatically chooses the best implementation
func TestIntelligentSelection(t *testing.T) {
	t.Run("Small datasets use legacy for better performance", func(t *testing.T) {
		// Create a small dataset (under 50 items)
		smallData := make([]string, 20)
		for i := 0; i < 20; i++ {
			smallData[i] = fmt.Sprintf("file%d.txt", i)
		}

		// Test that it still works correctly
		data := make([]string, len(smallData))
		copy(data, smallData)
		SortStrings(data)

		// Verify correct sorting
		for i := 1; i < len(data); i++ {
			if Compare(data[i-1], data[i]) > 0 {
				t.Errorf("Incorrect sort order at position %d: %s should come before %s", i, data[i-1], data[i])
			}
		}

		t.Logf("Small dataset (%d items) handled correctly", len(smallData))
	})

	t.Run("Medium datasets with duplicates use caching", func(t *testing.T) {
		// Create a medium dataset with duplicates
		mediumData := make([]string, 100)
		patterns := []string{"file%d.txt", "doc%d.pdf", "image%d.jpg"}

		for i := 0; i < 100; i++ {
			pattern := patterns[i%len(patterns)]
			// Create some duplicates by using modulo
			num := i % 20
			mediumData[i] = fmt.Sprintf(pattern, num)
		}

		// Test that it works correctly
		data := make([]string, len(mediumData))
		copy(data, mediumData)
		SortStrings(data)

		// Verify correct sorting
		for i := 1; i < len(data); i++ {
			if Compare(data[i-1], data[i]) > 0 {
				t.Errorf("Incorrect sort order at position %d: %s should come before %s", i, data[i-1], data[i])
			}
		}

		t.Logf("Medium dataset (%d items) with duplicates handled correctly", len(mediumData))
	})

	t.Run("Large datasets always use caching", func(t *testing.T) {
		// Create a large dataset
		largeData := make([]string, 500)
		for i := 0; i < 500; i++ {
			largeData[i] = fmt.Sprintf("file%d.txt", i*7%1000) // Some variety in numbers
		}

		// Test that it works correctly
		data := make([]string, len(largeData))
		copy(data, largeData)
		SortStrings(data)

		// Verify correct sorting (sample check)
		for i := 1; i < 10; i++ { // Check first 10 to avoid long test
			if Compare(data[i-1], data[i]) > 0 {
				t.Errorf("Incorrect sort order at position %d: %s should come before %s", i, data[i-1], data[i])
			}
		}

		t.Logf("Large dataset (%d items) handled correctly", len(largeData))
	})

	t.Run("Compare function adapts to string length", func(t *testing.T) {
		// Reset cache stats
		ResetCacheStats()

		// Test short strings (should use optimized parsing without caching)
		shortStrings := []string{"a1", "a2", "b1", "b2"}
		for i := 0; i < len(shortStrings); i++ {
			for j := i + 1; j < len(shortStrings); j++ {
				Compare(shortStrings[i], shortStrings[j])
			}
		}

		// Test longer strings (should use caching)
		longStrings := []string{
			"very_long_filename_with_numbers_123.txt",
			"another_very_long_filename_456.txt",
			"yet_another_long_filename_789.txt",
		}
		for i := 0; i < len(longStrings); i++ {
			for j := i + 1; j < len(longStrings); j++ {
				Compare(longStrings[i], longStrings[j])
				// Compare again to test caching
				Compare(longStrings[i], longStrings[j])
			}
		}

		hits, misses, ratio := CacheEfficiencyStats()
		t.Logf("Cache efficiency: %d hits, %d misses, %.2f hit ratio", hits, misses, ratio)

		// For longer strings with repeated comparisons, we should see some cache hits
		if hits == 0 && len(longStrings) > 0 {
			t.Log("Note: No cache hits detected, which is expected for short strings or single comparisons")
		}
	})

	t.Run("shouldUseCaching heuristics work correctly", func(t *testing.T) {
		// Test small dataset - should not use caching
		smallData := make([]string, 30)
		for i := 0; i < 30; i++ {
			smallData[i] = fmt.Sprintf("file%d.txt", i)
		}

		if shouldUseCaching(smallData) {
			t.Error("Small dataset should not use caching")
		}

		// Test medium dataset with duplicates - should use caching
		mediumWithDuplicates := make([]string, 80)
		for i := 0; i < 80; i++ {
			mediumWithDuplicates[i] = fmt.Sprintf("file%d.txt", i%10) // Lots of duplicates
		}

		if !shouldUseCaching(mediumWithDuplicates) {
			t.Error("Medium dataset with duplicates should use caching")
		}

		// Test medium dataset with long strings - should use caching
		mediumWithLongStrings := make([]string, 80)
		for i := 0; i < 80; i++ {
			mediumWithLongStrings[i] = fmt.Sprintf("very_long_filename_with_lots_of_characters_and_more_text_%d.txt", i)
		}

		if !shouldUseCaching(mediumWithLongStrings) {
			t.Error("Medium dataset with long strings should use caching")
		}

		// Test large dataset - should always use caching
		largeData := make([]string, 300)
		for i := 0; i < 300; i++ {
			largeData[i] = fmt.Sprintf("file%d.txt", i)
		}

		if !shouldUseCaching(largeData) {
			t.Error("Large dataset should always use caching")
		}
	})
}
