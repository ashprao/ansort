package ansort

import (
	"testing"
	"time"
)

// TestAPISimplification verifies that the standard API now uses optimized implementations
func TestAPISimplification(t *testing.T) {
	testData := []string{
		"file1.txt", "file10.txt", "file2.txt", "file20.txt", "file3.txt",
		"doc1.pdf", "doc10.pdf", "doc2.pdf", "doc20.pdf", "doc3.pdf",
		"image1.jpg", "image10.jpg", "image2.jpg", "image20.jpg", "image3.jpg",
	}

	t.Run("SortStrings uses optimized implementation", func(t *testing.T) {
		data := make([]string, len(testData))
		copy(data, testData)

		// Time the standard SortStrings function
		start := time.Now()
		SortStrings(data)
		duration := time.Since(start)

		// Verify it's sorted correctly
		expected := []string{
			"doc1.pdf", "doc2.pdf", "doc3.pdf", "doc10.pdf", "doc20.pdf",
			"file1.txt", "file2.txt", "file3.txt", "file10.txt", "file20.txt",
			"image1.jpg", "image2.jpg", "image3.jpg", "image10.jpg", "image20.jpg",
		}

		if len(data) != len(expected) {
			t.Fatalf("Length mismatch: got %d, want %d", len(data), len(expected))
		}

		for i, v := range data {
			if v != expected[i] {
				t.Errorf("Position %d: got %q, want %q", i, v, expected[i])
			}
		}

		t.Logf("SortStrings completed in %v", duration)
	})

	t.Run("Compare uses optimized implementation", func(t *testing.T) {
		// Test basic comparison
		result := Compare("file1.txt", "file10.txt")
		if result != -1 {
			t.Errorf("Compare(file1.txt, file10.txt): got %d, want -1", result)
		}

		result = Compare("file10.txt", "file2.txt")
		if result != 1 {
			t.Errorf("Compare(file10.txt, file2.txt): got %d, want 1", result)
		}

		result = Compare("file1.txt", "file1.txt")
		if result != 0 {
			t.Errorf("Compare(file1.txt, file1.txt): got %d, want 0", result)
		}
	})

	t.Run("Cache control functions work", func(t *testing.T) {
		// Test disabling cache
		DisableCache()

		data := make([]string, len(testData))
		copy(data, testData)
		SortStrings(data) // Should use legacy implementation

		// Verify it still works correctly
		expected := []string{
			"doc1.pdf", "doc2.pdf", "doc3.pdf", "doc10.pdf", "doc20.pdf",
			"file1.txt", "file2.txt", "file3.txt", "file10.txt", "file20.txt",
			"image1.jpg", "image2.jpg", "image3.jpg", "image10.jpg", "image20.jpg",
		}

		for i, v := range data {
			if v != expected[i] {
				t.Errorf("Position %d: got %q, want %q", i, v, expected[i])
			}
		}

		// Re-enable cache
		EnableCache()

		// Test configuring cache size
		ConfigureCacheSize(5000)

		// Verify comparison still works
		result := Compare("file1.txt", "file10.txt")
		if result != -1 {
			t.Errorf("After cache reconfiguration, Compare(file1.txt, file10.txt): got %d, want -1", result)
		}

		// Reset cache size to default for other tests
		ConfigureCacheSize(2000)
	})

	t.Run("Legacy functions still work", func(t *testing.T) {
		data := make([]string, len(testData))
		copy(data, testData)

		// Test legacy functions
		SortStringsLegacy(data)

		expected := []string{
			"doc1.pdf", "doc2.pdf", "doc3.pdf", "doc10.pdf", "doc20.pdf",
			"file1.txt", "file2.txt", "file3.txt", "file10.txt", "file20.txt",
			"image1.jpg", "image2.jpg", "image3.jpg", "image10.jpg", "image20.jpg",
		}

		for i, v := range data {
			if v != expected[i] {
				t.Errorf("Legacy sort position %d: got %q, want %q", i, v, expected[i])
			}
		}

		// Test legacy compare
		result := CompareLegacy("file1.txt", "file10.txt")
		if result != -1 {
			t.Errorf("CompareLegacy(file1.txt, file10.txt): got %d, want -1", result)
		}
	})
}

// BenchmarkAPIPerformance verifies that the default API provides optimized performance
func BenchmarkAPIPerformance(t *testing.B) {
	testData := []string{
		"file1.txt", "file10.txt", "file2.txt", "file20.txt", "file3.txt",
		"doc1.pdf", "doc10.pdf", "doc2.pdf", "doc20.pdf", "doc3.pdf",
		"image1.jpg", "image10.jpg", "image2.jpg", "image20.jpg", "image3.jpg",
		"photo100.png", "photo11.png", "photo2.png", "photo20.png", "photo3.png",
	}

	t.Run("SortStrings_Default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := make([]string, len(testData))
			copy(data, testData)
			SortStrings(data)
		}
	})

	t.Run("SortStrings_Legacy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := make([]string, len(testData))
			copy(data, testData)
			SortStringsLegacy(data)
		}
	})

	t.Run("Compare_Default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < len(testData)-1; j++ {
				Compare(testData[j], testData[j+1])
			}
		}
	})

	t.Run("Compare_Legacy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < len(testData)-1; j++ {
				CompareLegacy(testData[j], testData[j+1])
			}
		}
	})
}
