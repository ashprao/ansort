package ansort

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

// Benchmark data generators for performance testing
func generateBenchmarkData(size int, pattern string, seed int64) []string {
	rand.Seed(seed) // Use fixed seed for reproducible benchmarks
	data := make([]string, size)

	switch pattern {
	case "realistic":
		// Mix of common patterns found in real applications
		patterns := []string{
			"file", "document", "image", "video", "BlueBungalow_",
			"Product", "SKU", "Item", "Order", "Invoice",
		}
		suffixes := []string{".txt", ".pdf", ".jpg", ".mp4", "_test", "_prod", ""}

		for i := 0; i < size; i++ {
			prefix := patterns[rand.Intn(len(patterns))]
			number := rand.Intn(10000) + 1
			suffix := suffixes[rand.Intn(len(suffixes))]

			if rand.Float32() < 0.3 { // 30% chance of version numbers
				minor := rand.Intn(100)
				patch := rand.Intn(100)
				data[i] = fmt.Sprintf("%s%d.%d.%d%s", prefix, number, minor, patch, suffix)
			} else {
				data[i] = prefix + strconv.Itoa(number) + suffix
			}
		}
	case "versioning":
		for i := 0; i < size; i++ {
			major := rand.Intn(10) + 1
			minor := rand.Intn(20)
			patch := rand.Intn(100)
			data[i] = "v" + strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
		}
	case "files":
		prefixes := []string{"file", "document", "image", "video"}
		extensions := []string{".txt", ".pdf", ".jpg", ".mp4", ".docx"}
		for i := 0; i < size; i++ {
			prefix := prefixes[rand.Intn(len(prefixes))]
			number := rand.Intn(1000) + 1
			ext := extensions[rand.Intn(len(extensions))]
			data[i] = prefix + strconv.Itoa(number) + ext
		}
	default:
		for i := 0; i < size; i++ {
			data[i] = "item" + strconv.Itoa(rand.Intn(10000))
		}
	}

	return data
}

// Benchmark comparing original vs optimized implementations
func BenchmarkSortComparison(b *testing.B) {
	sizes := []int{100, 1000, 4000} // 4000 matches your Python benchmark
	patterns := []string{"realistic", "versioning", "files"}

	for _, size := range sizes {
		for _, pattern := range patterns {
			testName := fmt.Sprintf("%s_%d", pattern, size)

			// Generate test data once
			originalData := generateBenchmarkData(size, pattern, 42) // Fixed seed

			// Benchmark original implementation
			b.Run("Original_"+testName, func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					data := make([]string, len(originalData))
					copy(data, originalData)
					b.StartTimer()

					SortStrings(data) // Original implementation
				}
			})

			// Benchmark optimized implementation
			b.Run("Optimized_"+testName, func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					data := make([]string, len(originalData))
					copy(data, originalData)
					b.StartTimer()

					SortStringsOptimized(data) // Optimized implementation
				}
			})

			// Benchmark pooled implementation
			b.Run("Pooled_"+testName, func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					data := make([]string, len(originalData))
					copy(data, originalData)
					b.StartTimer()

					SortStringsPooled(data) // Pooled implementation
				}
			})

			// Benchmark high-performance implementation
			b.Run("HighPerf_"+testName, func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					data := make([]string, len(originalData))
					copy(data, originalData)
					b.StartTimer()

					SortStringsHighPerformance(data) // High-performance implementation
				}
			})
		}
	}
}

// Benchmark tokenization performance
func BenchmarkTokenizationComparison(b *testing.B) {
	testStrings := []string{
		"file123.txt",
		"BlueBungalow_105_test",
		"v1.2.10.build456",
		"very_long_filename_with_multiple_123_numbers_456_embedded_789.extension",
	}

	for _, str := range testStrings {
		name := fmt.Sprintf("len_%d", len(str))

		// Original tokenization
		b.Run("Original_"+name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = parseString(str)
			}
		})

		// Optimized tokenization
		b.Run("Optimized_"+name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = parseStringOptimized(str)
			}
		})

		// Pooled tokenization
		b.Run("Pooled_"+name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = parseStringPooled(str)
			}
		})
	}
}

// Benchmark comparison functions
func BenchmarkComparisonComparison(b *testing.B) {
	testPairs := []struct {
		name string
		a, b string
	}{
		{"SimpleNumbers", "file1.txt", "file10.txt"},
		{"ComplexVersions", "v1.2.10", "v1.10.2"},
		{"LongStrings", "BlueBungalow_105_test_v1.2.3", "BlueBungalow_12_prod_v2.1.0"},
		{"SimilarStrings", "item123test", "item123prod"},
	}

	for _, tp := range testPairs {
		// Original comparison
		b.Run("Original_"+tp.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = Compare(tp.a, tp.b)
			}
		})

		// Optimized comparison
		b.Run("Optimized_"+tp.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = CompareOptimized(tp.a, tp.b)
			}
		})
	}
}

// Benchmark cache effectiveness
func BenchmarkCacheEffectiveness(b *testing.B) {
	// Create data with many repeated strings to test cache effectiveness
	size := 1000
	uniqueStrings := 100 // 10% unique strings, 90% repetition

	baseStrings := generateBenchmarkData(uniqueStrings, "realistic", 42)

	// Create test data with repetition
	data := make([]string, size)
	for i := 0; i < size; i++ {
		data[i] = baseStrings[i%uniqueStrings]
	}

	// Benchmark without cache (original)
	b.Run("WithoutCache", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(data))
			copy(testData, data)
			b.StartTimer()

			SortStrings(testData)
		}
	})

	// Benchmark with cache
	b.Run("WithCache", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(data))
			copy(testData, data)
			// Clear cache between runs for fair comparison
			ClearGlobalCache()
			b.StartTimer()

			SortStringsOptimized(testData)
		}
	})

	// Benchmark with high-performance implementation
	b.Run("HighPerformance", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(data))
			copy(testData, data)
			b.StartTimer()

			SortStringsHighPerformance(testData)
		}
	})
}

// Benchmark your specific scenario that was 2x slower than Python
func BenchmarkPythonComparisonScenario(b *testing.B) {
	// Simulate your exact benchmark data pattern
	baseData := []string{
		"10", "a", "A1", "A2", "A3", "A4", "A5", "A6", "A7", "abcd", "ABCD",
		"BlueBungalow_01", "bluebungalow_02", "bluebungalow_03", "BlueBungalow_105",
		"BlueBungalow_118", "BlueBungalow_119", "BlueBungalow_12", "BlueBungalow_120",
		"BlueBungalow_121", "BlueBungalow_13", "BlueBungalow_131", "BlueBungalow_132",
		"BlueBungalow@1", "MediumPouch", "mediumPouch", "Satya", "satya ",
		"Satya1", "Satya1", "SatyaTesting", "satyatesting", "SatyaTesting1",
		"Satyatesting1", "1A2B", "101", "22222",
	}

	// Extend to 4000 items like your benchmark
	extendedData := make([]string, 0, 4000)
	for len(extendedData) < 4000 {
		extendedData = append(extendedData, baseData...)
		if len(extendedData) > 4000 {
			extendedData = extendedData[:4000]
		}
	}

	// Original implementation (your current performance)
	b.Run("Original_4000_items", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(extendedData))
			copy(testData, extendedData)
			b.StartTimer()

			SortStrings(testData)
		}
	})

	// Optimized implementation
	b.Run("Optimized_4000_items", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(extendedData))
			copy(testData, extendedData)
			b.StartTimer()

			SortStringsOptimized(testData)
		}
	})

	// High-performance implementation
	b.Run("HighPerf_4000_items", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(extendedData))
			copy(testData, extendedData)
			b.StartTimer()

			SortStringsHighPerformance(testData)
		}
	})
}

// Memory allocation analysis
func BenchmarkMemoryProfile(b *testing.B) {
	data := generateBenchmarkData(1000, "realistic", 42)

	b.Run("MemoryProfile_Original", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(data))
			copy(testData, data)
			b.StartTimer()

			SortStrings(testData)
		}
	})

	b.Run("MemoryProfile_Optimized", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			testData := make([]string, len(data))
			copy(testData, data)
			b.StartTimer()

			SortStringsHighPerformance(testData)
		}
	})
}
